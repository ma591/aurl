package graphql

import (
	"fmt"
	"sort"
	"strings"
)

// FormatHelp generates help text showing queries, mutations, and types.
func FormatHelp(schema *Schema) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("GraphQL API\nEndpoint: %s\n", schema.Endpoint))

	// Queries
	if qt, ok := schema.Types[schema.QueryType]; ok && len(qt.Fields) > 0 {
		b.WriteString("\n[Queries]\n")
		for _, f := range qt.Fields {
			writeField(&b, f)
		}
	}

	// Mutations
	if schema.MutationType != "" {
		if mt, ok := schema.Types[schema.MutationType]; ok && len(mt.Fields) > 0 {
			b.WriteString("\n[Mutations]\n")
			for _, f := range mt.Fields {
				writeField(&b, f)
			}
		}
	}

	// User-defined types (OBJECT and ENUM, skip root types and scalars)
	var typeNames []string
	for name, t := range schema.Types {
		if !t.IsUserType() {
			continue
		}
		if name == schema.QueryType || name == schema.MutationType {
			continue
		}
		if t.Kind == "OBJECT" || t.Kind == "ENUM" || t.Kind == "INPUT_OBJECT" {
			typeNames = append(typeNames, name)
		}
	}
	sort.Strings(typeNames)

	if len(typeNames) > 0 {
		b.WriteString("\n[Types]\n")
		for _, name := range typeNames {
			t := schema.Types[name]
			if t.Kind == "ENUM" {
				vals := strings.Join(t.EnumValues, " | ")
				b.WriteString(fmt.Sprintf("  %s: %s\n", name, vals))
			} else {
				var fields []string
				for _, f := range t.Fields {
					fields = append(fields, f.Name+" ("+f.Type+")")
				}
				line := "  " + name + ": " + strings.Join(fields, ", ")
				if len(line) > 120 {
					line = line[:117] + "..."
				}
				b.WriteString(line + "\n")
			}
		}
	}

	b.WriteString("\n(*) = required\n")
	b.WriteString("Tip: arc [name] describe [field] for detailed docs\n")
	return b.String()
}

// FormatDescribe generates detailed docs for a query, mutation, or type.
func FormatDescribe(schema *Schema, name string) string {
	// Check if it's a field on Query or Mutation
	if qt, ok := schema.Types[schema.QueryType]; ok {
		for _, f := range qt.Fields {
			if f.Name == name {
				return formatFieldDetail("Query", f)
			}
		}
	}
	if schema.MutationType != "" {
		if mt, ok := schema.Types[schema.MutationType]; ok {
			for _, f := range mt.Fields {
				if f.Name == name {
					return formatFieldDetail("Mutation", f)
				}
			}
		}
	}

	// Check if it's a type name
	if t, ok := schema.Types[name]; ok {
		return formatTypeDetail(t)
	}

	return fmt.Sprintf("No query, mutation, or type found for %q\n", name)
}

func writeField(b *strings.Builder, f GQLField) {
	line := fmt.Sprintf("  %-28s", f.Name)
	if f.Description != "" {
		desc := f.Description
		if len(desc) > 50 {
			desc = desc[:47] + "..."
		}
		line += "# " + desc
	}
	b.WriteString(line + "\n")

	// Args
	if len(f.Args) > 0 {
		var args []string
		for _, a := range f.Args {
			s := a.Name
			if a.Required {
				s += "*"
			}
			s += " (" + a.Type + ")"
			args = append(args, s)
		}
		b.WriteString("    args: " + strings.Join(args, ", ") + "\n")
	}

	// Return type
	if f.Type != "" {
		b.WriteString("    returns: " + f.Type + "\n")
	}
}

func formatFieldDetail(rootType string, f GQLField) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("%s: %s\n", rootType, f.Name))
	if f.Description != "" {
		b.WriteString(f.Description + "\n")
	}
	b.WriteString(fmt.Sprintf("\nReturns: %s\n", f.Type))

	if len(f.Args) > 0 {
		b.WriteString("\nArguments:\n")
		for _, a := range f.Args {
			marker := ""
			if a.Required {
				marker = "*"
			}
			b.WriteString(fmt.Sprintf("  %s%s (%s)\n", a.Name, marker, a.Type))
			if a.Description != "" {
				b.WriteString("    " + a.Description + "\n")
			}
			if a.DefaultValue != "" {
				b.WriteString("    default: " + a.DefaultValue + "\n")
			}
		}
	}

	return b.String()
}

func formatTypeDetail(t *GQLType) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Type: %s (%s)\n", t.Name, t.Kind))
	if t.Description != "" {
		b.WriteString(t.Description + "\n")
	}

	if t.Kind == "ENUM" && len(t.EnumValues) > 0 {
		b.WriteString("\nValues:\n")
		for _, v := range t.EnumValues {
			b.WriteString("  " + v + "\n")
		}
	}

	if len(t.Fields) > 0 {
		b.WriteString("\nFields:\n")
		for _, f := range t.Fields {
			b.WriteString(fmt.Sprintf("  %s: %s\n", f.Name, f.Type))
			if f.Description != "" {
				b.WriteString("    " + f.Description + "\n")
			}
		}
	}

	return b.String()
}
