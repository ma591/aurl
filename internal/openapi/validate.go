package openapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidateRequest validates a request against the spec.
// Returns enum violations (hard errors) and missing required fields (warnings) separately.
func ValidateRequest(spec *Spec, method, path string, body []byte) (enumErrors []ValidationError, missingFields []ValidationError) {
	op, _ := FindOperation(spec, method, path)
	if op == nil {
		return nil, nil
	}

	// Parse query params from the path
	queryParams := extractQueryParams(path)

	// Validate enum constraints on parameters
	for _, param := range op.Parameters {
		if len(param.Enum) == 0 {
			continue
		}

		var value string
		var hasValue bool

		switch param.In {
		case "query":
			value, hasValue = queryParams[param.Name]
		case "path":
			value, hasValue = extractPathParam(spec, method, path, param.Name)
		}

		if !hasValue {
			continue
		}

		if !enumContains(param.Enum, value) {
			vals := make([]string, len(param.Enum))
			for i, v := range param.Enum {
				vals[i] = fmt.Sprintf("%v", v)
			}
			enumErrors = append(enumErrors, ValidationError{
				Field:   param.Name,
				Message: fmt.Sprintf("invalid value %q\n  allowed: %s", value, strings.Join(vals, ", ")),
			})
		}
	}

	// Validate required body fields
	if op.RequestBody != nil && op.RequestBody.Schema != nil && len(body) > 0 {
		var bodyMap map[string]any
		if err := json.Unmarshal(body, &bodyMap); err == nil {
			required := SchemaRequiredFields(op.RequestBody.Schema)
			for _, field := range required {
				if _, ok := bodyMap[field]; !ok {
					missingFields = append(missingFields, ValidationError{
						Field:   field,
						Message: "required field missing from request body",
					})
				}
			}
		}
	}

	return enumErrors, missingFields
}

func extractQueryParams(path string) map[string]string {
	params := make(map[string]string)
	idx := strings.Index(path, "?")
	if idx == -1 {
		return params
	}
	query := path[idx+1:]
	values, err := url.ParseQuery(query)
	if err != nil {
		return params
	}
	for k, v := range values {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}
	return params
}

func extractPathParam(spec *Spec, method, reqPath, paramName string) (string, bool) {
	// Strip query string
	if idx := strings.Index(reqPath, "?"); idx != -1 {
		reqPath = reqPath[:idx]
	}

	method = strings.ToLower(method)
	for specPath, item := range spec.Paths {
		if _, ok := item.Operations[method]; !ok {
			continue
		}
		if !MatchPath(specPath, reqPath) {
			continue
		}

		specParts := strings.Split(strings.Trim(specPath, "/"), "/")
		reqParts := strings.Split(strings.Trim(reqPath, "/"), "/")

		for i, sp := range specParts {
			if sp == "{"+paramName+"}" && i < len(reqParts) {
				return reqParts[i], true
			}
		}
	}
	return "", false
}

func enumContains(enum []any, value string) bool {
	for _, v := range enum {
		if fmt.Sprintf("%v", v) == value {
			return true
		}
	}
	return false
}
