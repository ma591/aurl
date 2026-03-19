package graphql

import (
	"fmt"
	"strings"
)

type Schema struct {
	Raw            map[string]any
	Endpoint       string
	QueryType      string // usually "Query"
	MutationType   string // usually "Mutation"
	Types          map[string]*GQLType
}

type GQLType struct {
	Name        string
	Kind        string // OBJECT, INPUT_OBJECT, ENUM, SCALAR, etc.
	Description string
	Fields      []GQLField
	EnumValues  []string
}

type GQLField struct {
	Name        string
	Description string
	Args        []GQLArg
	Type        string // formatted, e.g. "[User!]!"
}

type GQLArg struct {
	Name         string
	Description  string
	Type         string // e.g. "String!", "ID!"
	Required     bool
	DefaultValue string
}

// ParseSchema extracts types from an introspection result.
func ParseSchema(raw map[string]any, endpoint string) (*Schema, error) {
	data, ok := raw["data"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("introspection result missing 'data' field")
	}

	schemaData, ok := data["__schema"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("introspection result missing '__schema' field")
	}

	schema := &Schema{
		Raw:      raw,
		Endpoint: endpoint,
		Types:    make(map[string]*GQLType),
	}

	// Root type names
	if qt, ok := schemaData["queryType"].(map[string]any); ok {
		schema.QueryType, _ = qt["name"].(string)
	}
	if mt, ok := schemaData["mutationType"].(map[string]any); ok {
		schema.MutationType, _ = mt["name"].(string)
	}

	// Parse types
	types, ok := schemaData["types"].([]any)
	if !ok {
		return schema, nil
	}

	for _, t := range types {
		typeMap, ok := t.(map[string]any)
		if !ok {
			continue
		}

		name, _ := typeMap["name"].(string)
		if name == "" || strings.HasPrefix(name, "__") {
			continue // skip introspection types
		}

		gqlType := &GQLType{
			Name: name,
		}
		gqlType.Kind, _ = typeMap["kind"].(string)
		gqlType.Description, _ = typeMap["description"].(string)

		// Parse fields
		if fields, ok := typeMap["fields"].([]any); ok {
			for _, f := range fields {
				fieldMap, ok := f.(map[string]any)
				if !ok {
					continue
				}
				field := GQLField{}
				field.Name, _ = fieldMap["name"].(string)
				field.Description, _ = fieldMap["description"].(string)

				if typeRef, ok := fieldMap["type"].(map[string]any); ok {
					field.Type = formatTypeRef(typeRef)
				}

				// Parse args
				if args, ok := fieldMap["args"].([]any); ok {
					for _, a := range args {
						argMap, ok := a.(map[string]any)
						if !ok {
							continue
						}
						arg := GQLArg{}
						arg.Name, _ = argMap["name"].(string)
						arg.Description, _ = argMap["description"].(string)
						if dv, ok := argMap["defaultValue"].(string); ok {
							arg.DefaultValue = dv
						}

						if typeRef, ok := argMap["type"].(map[string]any); ok {
							arg.Type = formatTypeRef(typeRef)
							arg.Required = isNonNull(typeRef)
						}

						field.Args = append(field.Args, arg)
					}
				}

				gqlType.Fields = append(gqlType.Fields, field)
			}
		}

		// Parse input fields (for INPUT_OBJECT types)
		if inputFields, ok := typeMap["inputFields"].([]any); ok {
			for _, f := range inputFields {
				fieldMap, ok := f.(map[string]any)
				if !ok {
					continue
				}
				field := GQLField{}
				field.Name, _ = fieldMap["name"].(string)
				field.Description, _ = fieldMap["description"].(string)

				if typeRef, ok := fieldMap["type"].(map[string]any); ok {
					field.Type = formatTypeRef(typeRef)
				}

				gqlType.Fields = append(gqlType.Fields, field)
			}
		}

		// Parse enum values
		if enumValues, ok := typeMap["enumValues"].([]any); ok {
			for _, ev := range enumValues {
				evMap, ok := ev.(map[string]any)
				if !ok {
					continue
				}
				if name, ok := evMap["name"].(string); ok {
					gqlType.EnumValues = append(gqlType.EnumValues, name)
				}
			}
		}

		schema.Types[name] = gqlType
	}

	return schema, nil
}

// formatTypeRef converts a GraphQL type reference to a readable string.
// e.g., NON_NULL(LIST(NON_NULL(OBJECT("User")))) → "[User!]!"
func formatTypeRef(ref map[string]any) string {
	kind, _ := ref["kind"].(string)
	name, _ := ref["name"].(string)

	switch kind {
	case "NON_NULL":
		if ofType, ok := ref["ofType"].(map[string]any); ok {
			inner := formatTypeRef(ofType)
			if inner != "" && inner != "?" {
				return inner + "!"
			}
			// Truncated — just show what we have without "!"
			return inner
		}
		return ""
	case "LIST":
		if ofType, ok := ref["ofType"].(map[string]any); ok {
			inner := formatTypeRef(ofType)
			if inner == "" || inner == "?" {
				inner = "..."
			}
			return "[" + inner + "]"
		}
		return "[...]"
	default:
		if name != "" {
			return name
		}
		// Bottom of truncated type ref — try to find the deepest name
		if ofType, ok := ref["ofType"].(map[string]any); ok {
			return formatTypeRef(ofType)
		}
		return "?"
	}
}

// isNonNull checks if a type reference is wrapped in NON_NULL.
func isNonNull(ref map[string]any) bool {
	kind, _ := ref["kind"].(string)
	return kind == "NON_NULL"
}

// IsUserType returns true if this is a user-defined type worth showing.
func (t *GQLType) IsUserType() bool {
	if strings.HasPrefix(t.Name, "__") {
		return false
	}
	// Skip built-in scalars
	builtins := map[string]bool{
		"String": true, "Int": true, "Float": true,
		"Boolean": true, "ID": true,
	}
	if t.Kind == "SCALAR" && builtins[t.Name] {
		return false
	}
	return true
}
