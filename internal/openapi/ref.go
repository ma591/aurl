package openapi

import (
	"fmt"
	"strings"
)

// ResolveRef resolves a JSON $ref pointer like "#/components/schemas/Pet"
// within the root spec map. Returns the resolved object.
// The visited set prevents infinite recursion on circular references.
func ResolveRef(ref string, root map[string]any, visited map[string]bool) (map[string]any, error) {
	if visited[ref] {
		return map[string]any{}, nil // circular ref, return empty
	}
	visited[ref] = true

	if !strings.HasPrefix(ref, "#/") {
		return nil, fmt.Errorf("unsupported $ref: %s (only local refs supported)", ref)
	}

	parts := strings.Split(strings.TrimPrefix(ref, "#/"), "/")
	var current any = root
	for _, part := range parts {
		m, ok := current.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("cannot resolve $ref %s: expected object at %s", ref, part)
		}
		current, ok = m[part]
		if !ok {
			return nil, fmt.Errorf("cannot resolve $ref %s: key %q not found", ref, part)
		}
	}

	result, ok := current.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("$ref %s resolved to non-object", ref)
	}

	// If the resolved object itself has a $ref, resolve it too
	if innerRef, ok := result["$ref"].(string); ok {
		return ResolveRef(innerRef, root, visited)
	}

	return result, nil
}

// ResolveSchema resolves a schema that may contain a $ref.
// If schema has "$ref", resolves it. Otherwise returns schema as-is.
func ResolveSchema(schema map[string]any, root map[string]any) map[string]any {
	if ref, ok := schema["$ref"].(string); ok {
		resolved, err := ResolveRef(ref, root, make(map[string]bool))
		if err != nil {
			return schema
		}
		return resolved
	}
	return schema
}

// RefName extracts the short name from a $ref string.
// e.g., "#/components/schemas/Pet" → "Pet"
func RefName(ref string) string {
	parts := strings.Split(ref, "/")
	return parts[len(parts)-1]
}
