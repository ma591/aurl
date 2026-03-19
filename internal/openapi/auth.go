package openapi

// AuthScheme represents a detected authentication scheme from the spec.
type AuthScheme struct {
	Name        string // key in securitySchemes
	Type        string // "apiKey", "http", "oauth2", "openIdConnect"
	Scheme      string // for http type: "bearer", "basic"
	In          string // for apiKey type: "header", "query", "cookie"
	HeaderName  string // for apiKey: the header/query param name
	Description string // human-readable description for prompting
}

// DetectAuth parses components.securitySchemes and returns structured auth info.
func DetectAuth(spec *Spec) []AuthScheme {
	components, ok := spec.Raw["components"].(map[string]any)
	if !ok {
		return nil
	}

	schemes, ok := components["securitySchemes"].(map[string]any)
	if !ok {
		return nil
	}

	var result []AuthScheme

	for name, schemeData := range schemes {
		schemeMap, ok := schemeData.(map[string]any)
		if !ok {
			continue
		}

		// Resolve $ref if present
		if ref, ok := schemeMap["$ref"].(string); ok {
			resolved, err := ResolveRef(ref, spec.Raw, make(map[string]bool))
			if err != nil {
				continue
			}
			schemeMap = resolved
		}

		scheme := AuthScheme{Name: name}
		scheme.Type, _ = schemeMap["type"].(string)

		switch scheme.Type {
		case "apiKey":
			scheme.In, _ = schemeMap["in"].(string)
			scheme.HeaderName, _ = schemeMap["name"].(string)
			scheme.Description = "API Key"
			if scheme.In == "header" {
				scheme.Description += " in header \"" + scheme.HeaderName + "\""
			} else if scheme.In == "query" {
				scheme.Description += " in query param \"" + scheme.HeaderName + "\""
			}

		case "http":
			scheme.Scheme, _ = schemeMap["scheme"].(string)
			switch scheme.Scheme {
			case "bearer":
				scheme.Description = "Bearer token"
			case "basic":
				scheme.Description = "Basic auth (username:password)"
			default:
				scheme.Description = "HTTP " + scheme.Scheme + " auth"
			}

		case "oauth2":
			scheme.Description = "OAuth2 (provide token manually)"

		case "openIdConnect":
			scheme.Description = "OpenID Connect (provide token manually)"

		default:
			scheme.Description = "Unknown auth type: " + scheme.Type
		}

		// Prefer spec-provided description if available
		if specDesc, ok := schemeMap["description"].(string); ok && specDesc != "" {
			scheme.Description = specDesc
		}

		result = append(result, scheme)
	}

	return result
}
