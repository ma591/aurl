package openapi

import (
	"encoding/json"
	"fmt"
	"strings"
)

// GenerateExample produces an example value from a JSON Schema object.
// Uses the root spec for $ref resolution. The visited set prevents circular refs.
func GenerateExample(schema map[string]any, root map[string]any, visited map[string]bool) any {
	if visited == nil {
		visited = make(map[string]bool)
	}

	// Handle $ref
	if ref, ok := schema["$ref"].(string); ok {
		if visited[ref] {
			return map[string]any{}
		}
		visited[ref] = true
		resolved, err := ResolveRef(ref, root, make(map[string]bool))
		if err != nil {
			return nil
		}
		return GenerateExample(resolved, root, visited)
	}

	// Handle allOf — merge all sub-schemas
	if allOf, ok := schema["allOf"].([]any); ok {
		merged := map[string]any{}
		for _, sub := range allOf {
			subMap, ok := sub.(map[string]any)
			if !ok {
				continue
			}
			result := GenerateExample(subMap, root, visited)
			if m, ok := result.(map[string]any); ok {
				for k, v := range m {
					merged[k] = v
				}
			}
		}
		return merged
	}

	// Handle oneOf / anyOf — pick the first
	for _, key := range []string{"oneOf", "anyOf"} {
		if choices, ok := schema[key].([]any); ok && len(choices) > 0 {
			if first, ok := choices[0].(map[string]any); ok {
				return GenerateExample(first, root, visited)
			}
		}
	}

	// Get type — handle both string and array forms (3.0 vs 3.1)
	typ := schemaType(schema)

	switch typ {
	case "object":
		return generateObjectExample(schema, root, visited)
	case "array":
		return generateArrayExample(schema, root, visited)
	case "string":
		if enum, ok := schema["enum"].([]any); ok && len(enum) > 0 {
			return fmt.Sprintf("%v", enum[0])
		}
		if format, ok := schema["format"].(string); ok {
			switch format {
			case "date-time":
				return "2025-01-01T00:00:00Z"
			case "date":
				return "2025-01-01"
			case "email":
				return "user@example.com"
			case "uri", "url":
				return "https://example.com"
			case "uuid":
				return "550e8400-e29b-41d4-a716-446655440000"
			}
		}
		return "..."
	case "integer":
		if enum, ok := schema["enum"].([]any); ok && len(enum) > 0 {
			return enum[0]
		}
		return 0
	case "number":
		if enum, ok := schema["enum"].([]any); ok && len(enum) > 0 {
			return enum[0]
		}
		return 0.0
	case "boolean":
		return false
	default:
		// If there are properties, treat as object
		if _, ok := schema["properties"]; ok {
			return generateObjectExample(schema, root, visited)
		}
		return nil
	}
}

func generateObjectExample(schema map[string]any, root map[string]any, visited map[string]bool) map[string]any {
	result := map[string]any{}

	props, ok := schema["properties"].(map[string]any)
	if !ok {
		return result
	}

	for name, propData := range props {
		propSchema, ok := propData.(map[string]any)
		if !ok {
			continue
		}
		result[name] = GenerateExample(propSchema, root, visited)
	}

	return result
}

func generateArrayExample(schema map[string]any, root map[string]any, visited map[string]bool) []any {
	if items, ok := schema["items"].(map[string]any); ok {
		example := GenerateExample(items, root, visited)
		return []any{example}
	}
	return []any{}
}

// schemaType extracts the type from a schema, handling both
// OpenAPI 3.0 (type is string) and 3.1 (type can be array).
func schemaType(schema map[string]any) string {
	switch t := schema["type"].(type) {
	case string:
		return t
	case []any:
		// 3.1: type: ["string", "null"] → return first non-null
		for _, v := range t {
			if s, ok := v.(string); ok && s != "null" {
				return s
			}
		}
		if len(t) > 0 {
			if s, ok := t[0].(string); ok {
				return s
			}
		}
	}
	return ""
}

// GenerateExampleJSON produces a compact JSON string from a schema.
func GenerateExampleJSON(schema map[string]any, root map[string]any) string {
	example := GenerateExample(schema, root, nil)
	if example == nil {
		return "{}"
	}
	data, err := json.Marshal(example)
	if err != nil {
		return "{}"
	}
	return string(data)
}

// FindBodyExample finds and generates an example body for an operation.
func FindBodyExample(spec *Spec, method, path string) string {
	op, _ := FindOperation(spec, method, path)
	if op == nil || op.RequestBody == nil || op.RequestBody.Schema == nil {
		return ""
	}
	return GenerateExampleJSON(op.RequestBody.Schema, spec.Raw)
}

// SchemaRequiredFields returns the required field names from a schema.
func SchemaRequiredFields(schema map[string]any) []string {
	req, ok := schema["required"].([]any)
	if !ok {
		return nil
	}
	var fields []string
	for _, r := range req {
		if s, ok := r.(string); ok {
			fields = append(fields, s)
		}
	}
	return fields
}

// FormatExampleCompact formats an example as a single-line JSON for display.
func FormatExampleCompact(schema map[string]any, root map[string]any) string {
	s := GenerateExampleJSON(schema, root)
	// Collapse whitespace for display
	s = strings.ReplaceAll(s, "\n", "")
	return s
}
