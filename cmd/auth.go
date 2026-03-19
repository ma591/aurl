package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/shawnpana/arc/internal/config"
	"github.com/shawnpana/arc/internal/openapi"
	"github.com/spf13/cobra"
)

var authHeaders []string

var authCmd = &cobra.Command{
	Use:   "auth [name]",
	Short: "Configure auth for a registered API",
	Args:  cobra.ExactArgs(1),
	RunE:  runAuth,
}

func init() {
	authCmd.Flags().StringArrayVar(&authHeaders, "header", nil, "Auth header (e.g., \"Authorization: Bearer token\")")
}

func runAuth(cmd *cobra.Command, args []string) error {
	name := args[0]

	// Check if this is a registered API or GraphQL endpoint
	data, err := config.LoadSpec(name)
	isGraphQL := false
	if err != nil {
		// Try GraphQL
		_, err = config.LoadGraphQL(name)
		if err != nil {
			return fmt.Errorf("%q not registered. Run 'arc list' to see registered commands.", name)
		}
		isGraphQL = true
	}

	// Load existing auth to preserve base_url_override
	existingAuth, _ := config.LoadAuth(name)
	auth := &config.AuthConfig{
		Headers:     make(map[string]string),
		QueryParams: make(map[string]string),
	}
	if existingAuth != nil {
		auth.BaseURLOverride = existingAuth.BaseURLOverride
	}

	// Add manual headers from flags
	for _, h := range authHeaders {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) == 2 {
			auth.Headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	// If no manual headers, run interactive detection (OpenAPI only)
	if len(authHeaders) == 0 {
		if isGraphQL {
			fmt.Println("Use --header to set auth for GraphQL APIs.")
			fmt.Printf("  Example: arc auth %s --header \"Authorization: Bearer your-token\"\n", name)
			return nil
		}

		spec, err := openapi.Parse(data)
		if err != nil {
			return fmt.Errorf("failed to parse spec: %w", err)
		}

		schemes := openapi.DetectAuth(spec)
		if len(schemes) > 0 {
			fmt.Println("Auth detected:")
			for i, scheme := range schemes {
				fmt.Printf("  [%d] %s (%s)\n", i+1, scheme.Name, scheme.Description)
			}
			fmt.Println()

			reader := bufio.NewReader(os.Stdin)
			for _, scheme := range schemes {
				promptAndStoreAuth(reader, scheme, auth)
			}
		} else {
			fmt.Println("No securitySchemes found in spec. Use --header to set auth manually.")
			fmt.Printf("  Example: arc auth %s --header \"Authorization: Bearer your-token\"\n", name)
			return nil
		}
	}

	if err := config.SaveAuth(name, auth); err != nil {
		return fmt.Errorf("failed to save auth: %w", err)
	}

	fmt.Println("Auth updated.")
	return nil
}
