package openapi

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Spec struct {
	Raw          map[string]any
	Title        string
	Version      string
	Description  string
	BaseURL      string
	ExternalDocs *ExternalDocs
	Tags         map[string]TagInfo // tag name → info
	Paths        map[string]*PathItem
}

type ExternalDocs struct {
	Description string
	URL         string
}

type TagInfo struct {
	Description  string
	ExternalDocs *ExternalDocs
}

type PathItem struct {
	Operations map[string]*Operation // key = lowercase method (get, post, etc.)
}

type Operation struct {
	Tags         []string
	Summary      string
	Description  string
	OperationID  string
	ExternalDocs *ExternalDocs
	Parameters   []Parameter
	RequestBody  *RequestBody
	Responses    map[string]*Response // key = status code string
}

type Parameter struct {
	Name        string
	In          string // "path", "query", "header", "cookie"
	Required    bool
	Description string
	Schema      map[string]any
	Enum        []any
}

type RequestBody struct {
	Required    bool
	Description string
	Schema      map[string]any // resolved schema
}

type Response struct {
	Description string
	SchemaRef   string // short name for display, e.g. "Pet"
}

func Parse(data []byte) (*Spec, error) {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	// Validate it looks like an OpenAPI spec
	if _, ok := raw["openapi"]; !ok {
		if _, ok := raw["swagger"]; !ok {
			return nil, fmt.Errorf("not an OpenAPI spec (missing 'openapi' or 'swagger' field)")
		}
	}

	spec := &Spec{
		Raw:  raw,
		Tags: make(map[string]TagInfo),
		Paths: make(map[string]*PathItem),
	}

	// info
	if info, ok := raw["info"].(map[string]any); ok {
		spec.Title, _ = info["title"].(string)
		spec.Version, _ = info["version"].(string)
		spec.Description, _ = info["description"].(string)
	}

	// externalDocs
	spec.ExternalDocs = parseExternalDocs(raw)

	// tags
	if tags, ok := raw["tags"].([]any); ok {
		for _, t := range tags {
			tagMap, ok := t.(map[string]any)
			if !ok {
				continue
			}
			name, _ := tagMap["name"].(string)
			if name == "" {
				continue
			}
			info := TagInfo{}
			info.Description, _ = tagMap["description"].(string)
			info.ExternalDocs = parseExternalDocs(tagMap)
			spec.Tags[name] = info
		}
	}

	// servers[0].url
	if servers, ok := raw["servers"].([]any); ok && len(servers) > 0 {
		if s, ok := servers[0].(map[string]any); ok {
			spec.BaseURL, _ = s["url"].(string)
		}
	}

	// paths
	paths, ok := raw["paths"].(map[string]any)
	if !ok {
		return spec, nil
	}

	methods := []string{"get", "post", "put", "delete", "patch", "head", "options"}

	for path, pathData := range paths {
		pathMap, ok := pathData.(map[string]any)
		if !ok {
			continue
		}

		item := &PathItem{Operations: make(map[string]*Operation)}

		// Path-level parameters (shared by all operations)
		var pathParams []Parameter
		if params, ok := pathMap["parameters"].([]any); ok {
			pathParams = parseParameters(params, raw)
		}

		for _, method := range methods {
			opData, ok := pathMap[method].(map[string]any)
			if !ok {
				continue
			}

			op := parseOperation(opData, raw)

			// Merge path-level params (operation params take precedence)
			op.Parameters = mergeParameters(pathParams, op.Parameters)

			item.Operations[method] = op
		}

		if len(item.Operations) > 0 {
			spec.Paths[path] = item
		}
	}

	return spec, nil
}

func parseOperation(data map[string]any, root map[string]any) *Operation {
	op := &Operation{
		Responses: make(map[string]*Response),
	}

	// tags
	if tags, ok := data["tags"].([]any); ok {
		for _, t := range tags {
			if s, ok := t.(string); ok {
				op.Tags = append(op.Tags, s)
			}
		}
	}

	op.Summary, _ = data["summary"].(string)
	op.Description, _ = data["description"].(string)
	op.OperationID, _ = data["operationId"].(string)
	op.ExternalDocs = parseExternalDocs(data)

	// parameters
	if params, ok := data["parameters"].([]any); ok {
		op.Parameters = parseParameters(params, root)
	}

	// requestBody
	if body, ok := data["requestBody"].(map[string]any); ok {
		op.RequestBody = parseRequestBody(body, root)
	}

	// responses
	if responses, ok := data["responses"].(map[string]any); ok {
		for code, respData := range responses {
			respMap, ok := respData.(map[string]any)
			if !ok {
				continue
			}
			resp := &Response{}
			resp.Description, _ = respMap["description"].(string)

			// Try to get a short schema name for display
			if content, ok := respMap["content"].(map[string]any); ok {
				if jsonContent, ok := content["application/json"].(map[string]any); ok {
					if schema, ok := jsonContent["schema"].(map[string]any); ok {
						resp.SchemaRef = schemaDisplayName(schema)
					}
				}
			}

			op.Responses[code] = resp
		}
	}

	return op
}

func parseParameters(params []any, root map[string]any) []Parameter {
	var result []Parameter
	for _, p := range params {
		pMap, ok := p.(map[string]any)
		if !ok {
			continue
		}

		// Resolve $ref on the parameter itself
		if ref, ok := pMap["$ref"].(string); ok {
			resolved, err := ResolveRef(ref, root, make(map[string]bool))
			if err != nil {
				continue
			}
			pMap = resolved
		}

		param := Parameter{}
		param.Name, _ = pMap["name"].(string)
		param.In, _ = pMap["in"].(string)
		param.Required, _ = pMap["required"].(bool)
		param.Description, _ = pMap["description"].(string)

		if schema, ok := pMap["schema"].(map[string]any); ok {
			param.Schema = ResolveSchema(schema, root)
			// Extract enum from schema
			if enum, ok := param.Schema["enum"].([]any); ok {
				param.Enum = enum
			}
		}

		result = append(result, param)
	}
	return result
}

func parseRequestBody(body map[string]any, root map[string]any) *RequestBody {
	rb := &RequestBody{}
	rb.Required, _ = body["required"].(bool)
	rb.Description, _ = body["description"].(string)

	// Resolve $ref on requestBody itself
	if ref, ok := body["$ref"].(string); ok {
		resolved, err := ResolveRef(ref, root, make(map[string]bool))
		if err == nil {
			body = resolved
		}
	}

	if content, ok := body["content"].(map[string]any); ok {
		if jsonContent, ok := content["application/json"].(map[string]any); ok {
			if schema, ok := jsonContent["schema"].(map[string]any); ok {
				rb.Schema = ResolveSchema(schema, root)
			}
		}
	}

	return rb
}

func mergeParameters(pathParams, opParams []Parameter) []Parameter {
	// Operation params override path params with the same name+in
	seen := make(map[string]bool)
	for _, p := range opParams {
		seen[p.Name+"|"+p.In] = true
	}

	merged := make([]Parameter, len(opParams))
	copy(merged, opParams)
	for _, p := range pathParams {
		if !seen[p.Name+"|"+p.In] {
			merged = append(merged, p)
		}
	}
	return merged
}

// schemaDisplayName returns a short display name for a schema.
func schemaDisplayName(schema map[string]any) string {
	if ref, ok := schema["$ref"].(string); ok {
		return RefName(ref)
	}
	if schema["type"] == "array" {
		if items, ok := schema["items"].(map[string]any); ok {
			if ref, ok := items["$ref"].(string); ok {
				return RefName(ref) + "[]"
			}
			if t, ok := items["type"].(string); ok {
				return t + "[]"
			}
		}
		return "array"
	}
	if t, ok := schema["type"].(string); ok {
		return t
	}
	return ""
}

// MatchPath checks if a spec path template matches a request path.
// e.g., "/pets/{petId}" matches "/pets/123"
func MatchPath(specPath, reqPath string) bool {
	// Strip query string from request path
	if idx := strings.Index(reqPath, "?"); idx != -1 {
		reqPath = reqPath[:idx]
	}

	specParts := strings.Split(strings.Trim(specPath, "/"), "/")
	reqParts := strings.Split(strings.Trim(reqPath, "/"), "/")

	if len(specParts) != len(reqParts) {
		return false
	}

	for i, sp := range specParts {
		if strings.HasPrefix(sp, "{") && strings.HasSuffix(sp, "}") {
			continue
		}
		if sp != reqParts[i] {
			return false
		}
	}
	return true
}

func parseExternalDocs(data map[string]any) *ExternalDocs {
	ed, ok := data["externalDocs"].(map[string]any)
	if !ok {
		return nil
	}
	docs := &ExternalDocs{}
	docs.Description, _ = ed["description"].(string)
	docs.URL, _ = ed["url"].(string)
	if docs.URL == "" && docs.Description == "" {
		return nil
	}
	return docs
}

// FindOperation finds the operation matching a method and path in the spec.
func FindOperation(spec *Spec, method, path string) (*Operation, string) {
	method = strings.ToLower(method)
	for specPath, item := range spec.Paths {
		if MatchPath(specPath, path) {
			if op, ok := item.Operations[method]; ok {
				return op, specPath
			}
		}
	}
	return nil, ""
}
