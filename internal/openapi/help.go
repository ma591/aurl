package openapi

import (
	"fmt"
	"sort"
	"strings"
)

// FormatHelp generates the full help text for a spec.
func FormatHelp(spec *Spec) string {
	var b strings.Builder

	// Header
	if spec.Title != "" {
		b.WriteString(spec.Title)
		if spec.Version != "" {
			b.WriteString(" v" + spec.Version)
		}
		b.WriteString("\n")
	}
	if spec.BaseURL != "" {
		b.WriteString("Base URL: " + spec.BaseURL + "\n")
	}

	// API description (first 2 lines)
	if spec.Description != "" {
		desc := truncateLines(spec.Description, 2)
		b.WriteString("\n" + desc + "\n")
	}

	// External docs
	if spec.ExternalDocs != nil && spec.ExternalDocs.URL != "" {
		b.WriteString("Docs: " + spec.ExternalDocs.URL + "\n")
	}

	// Collect endpoints grouped by tag
	type endpoint struct {
		method string
		path   string
		op     *Operation
	}

	tagEndpoints := make(map[string][]endpoint)

	for path, item := range spec.Paths {
		methodOrder := []string{"get", "post", "put", "patch", "delete", "head", "options"}
		for _, method := range methodOrder {
			op, ok := item.Operations[method]
			if !ok {
				continue
			}
			tag := "other"
			if len(op.Tags) > 0 {
				tag = op.Tags[0]
			}
			tagEndpoints[tag] = append(tagEndpoints[tag], endpoint{method, path, op})
		}
	}

	// Sort tags
	tags := make([]string, 0, len(tagEndpoints))
	for tag := range tagEndpoints {
		tags = append(tags, tag)
	}
	sort.Strings(tags)

	for _, tag := range tags {
		endpoints := tagEndpoints[tag]

		// Sort endpoints within tag by path, then method
		sort.Slice(endpoints, func(i, j int) bool {
			if endpoints[i].path != endpoints[j].path {
				return endpoints[i].path < endpoints[j].path
			}
			return endpoints[i].method < endpoints[j].method
		})

		// Tag header with optional description
		tagLine := "\n[" + tag + "]"
		if info, ok := spec.Tags[tag]; ok && info.Description != "" {
			desc := info.Description
			if len(desc) > 60 {
				desc = desc[:57] + "..."
			}
			tagLine += " " + desc
		}
		b.WriteString(tagLine + "\n")

		for _, ep := range endpoints {
			// Method and path
			methodStr := fmt.Sprintf("%-6s", strings.ToUpper(ep.method))
			line := fmt.Sprintf("  %s %s", methodStr, ep.path)

			// Summary
			if ep.op.Summary != "" {
				padding := 40 - len(line)
				if padding < 2 {
					padding = 2
				}
				line += strings.Repeat(" ", padding) + "# " + ep.op.Summary
			}
			b.WriteString(line + "\n")

			// Parameters (only path and query, skip header/cookie)
			var params []string
			for _, p := range ep.op.Parameters {
				if p.In != "query" && p.In != "path" {
					continue
				}
				s := p.Name
				if p.Required {
					s += "*"
				}
				if len(p.Enum) > 0 {
					vals := make([]string, len(p.Enum))
					for i, v := range p.Enum {
						vals[i] = fmt.Sprintf("%v", v)
					}
					s += " (" + strings.Join(vals, "|") + ")"
				}
				if p.Description != "" {
					desc := p.Description
					if len(desc) > 50 {
						desc = desc[:47] + "..."
					}
					s += " - " + desc
				}
				params = append(params, s)
			}
			if len(params) > 0 {
				b.WriteString("    params: " + strings.Join(params, ", ") + "\n")
			}

			// Request body description + example
			if ep.op.RequestBody != nil {
				if ep.op.RequestBody.Description != "" {
					desc := ep.op.RequestBody.Description
					if len(desc) > 80 {
						desc = desc[:77] + "..."
					}
					b.WriteString("    " + desc + "\n")
				}
				if ep.op.RequestBody.Schema != nil {
					example := FormatExampleCompact(ep.op.RequestBody.Schema, spec.Raw)
					if example != "" && example != "{}" && example != "null" {
						b.WriteString("    body: '" + example + "'\n")
					}
				}
			}

			// Response codes
			if len(ep.op.Responses) > 0 {
				var respParts []string
				codes := make([]string, 0, len(ep.op.Responses))
				for code := range ep.op.Responses {
					codes = append(codes, code)
				}
				sort.Strings(codes)

				for _, code := range codes {
					resp := ep.op.Responses[code]
					part := code + ":"
					if resp.SchemaRef != "" {
						part += " " + resp.SchemaRef
					} else if resp.Description != "" {
						desc := resp.Description
						if len(desc) > 30 {
							desc = desc[:27] + "..."
						}
						part += " " + desc
					}
					respParts = append(respParts, part)
				}
				b.WriteString("    " + strings.Join(respParts, "  ") + "\n")
			}
		}
	}

	b.WriteString("\n(*) = required\n")
	b.WriteString("Tip: aurl [name] describe METHOD /path for detailed endpoint docs\n")
	return b.String()
}

// FormatDescribe generates detailed documentation for a single operation.
func FormatDescribe(spec *Spec, method, path string) string {
	op, specPath := FindOperation(spec, method, path)
	if op == nil {
		return fmt.Sprintf("No operation found for %s %s\n", strings.ToUpper(method), path)
	}

	var b strings.Builder

	// Header
	b.WriteString(fmt.Sprintf("%s %s\n", strings.ToUpper(method), specPath))

	// Summary
	if op.Summary != "" {
		b.WriteString(op.Summary + "\n")
	}

	// Full description
	if op.Description != "" {
		b.WriteString("\n" + op.Description + "\n")
	}

	// External docs
	if op.ExternalDocs != nil && op.ExternalDocs.URL != "" {
		label := "Docs"
		if op.ExternalDocs.Description != "" {
			label = op.ExternalDocs.Description
		}
		b.WriteString(label + ": " + op.ExternalDocs.URL + "\n")
	}

	// Parameters
	if len(op.Parameters) > 0 {
		b.WriteString("\nParameters:\n")
		for _, p := range op.Parameters {
			marker := ""
			if p.Required {
				marker = "*"
			}
			typ := ""
			if p.Schema != nil {
				if t := schemaType(p.Schema); t != "" {
					typ = ", " + t
				}
			}
			b.WriteString(fmt.Sprintf("  %s%s (%s%s)\n", p.Name, marker, p.In, typ))
			if p.Description != "" {
				b.WriteString("    " + p.Description + "\n")
			}
			if len(p.Enum) > 0 {
				vals := make([]string, len(p.Enum))
				for i, v := range p.Enum {
					vals[i] = fmt.Sprintf("%v", v)
				}
				b.WriteString("    allowed: " + strings.Join(vals, ", ") + "\n")
			}
		}
	}

	// Request body
	if op.RequestBody != nil {
		b.WriteString("\nRequest Body")
		if op.RequestBody.Required {
			b.WriteString(" (required)")
		}
		b.WriteString(":\n")
		if op.RequestBody.Description != "" {
			b.WriteString("  " + op.RequestBody.Description + "\n")
		}
		if op.RequestBody.Schema != nil {
			example := GenerateExampleJSON(op.RequestBody.Schema, spec.Raw)
			if example != "" && example != "{}" {
				b.WriteString("  example: '" + example + "'\n")
			}
			// Show required fields
			required := SchemaRequiredFields(op.RequestBody.Schema)
			if len(required) > 0 {
				b.WriteString("  required fields: " + strings.Join(required, ", ") + "\n")
			}
		}
	}

	// Responses
	if len(op.Responses) > 0 {
		b.WriteString("\nResponses:\n")
		codes := make([]string, 0, len(op.Responses))
		for code := range op.Responses {
			codes = append(codes, code)
		}
		sort.Strings(codes)

		for _, code := range codes {
			resp := op.Responses[code]
			line := "  " + code
			if resp.Description != "" {
				line += " - " + resp.Description
			}
			b.WriteString(line + "\n")
			if resp.SchemaRef != "" {
				b.WriteString("    Schema: " + resp.SchemaRef + "\n")
			}
		}
	}

	return b.String()
}

// truncateLines returns the first n lines of text, adding "..." if truncated.
func truncateLines(text string, n int) string {
	lines := strings.SplitN(text, "\n", n+1)
	if len(lines) <= n {
		return strings.TrimSpace(text)
	}
	return strings.TrimSpace(strings.Join(lines[:n], "\n")) + "..."
}

