package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shawnpana/arc/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "arc",
	Short: "Turn any API spec into a CLI command",
	Long:  "Register OpenAPI specs or GraphQL endpoints as named subcommands and invoke them from the terminal.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// GetRootCmd returns the root command for doc generation.
func GetRootCmd() *cobra.Command {
	return rootCmd
}

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(authCmd)

	// Register dynamic commands from stored specs
	apisDir := config.APIsDir()
	entries, err := os.ReadDir(apisDir)
	if err != nil {
		return
	}

	reserved := map[string]bool{
		"add": true, "list": true, "remove": true,
		"auth": true, "help": true, "completion": true,
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		name := strings.TrimSuffix(entry.Name(), ".json")
		if reserved[name] {
			continue
		}
		specPath := filepath.Join(apisDir, entry.Name())
		authPath := filepath.Join(config.AuthDir(), name+".json")
		rootCmd.AddCommand(NewAPICommand(name, specPath, authPath))
	}

	// Register dynamic commands from stored GraphQL schemas
	gqlDir := config.GraphQLDir()
	gqlEntries, err := os.ReadDir(gqlDir)
	if err == nil {
		for _, entry := range gqlEntries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
				continue
			}
			name := strings.TrimSuffix(entry.Name(), ".json")
			if reserved[name] {
				continue
			}
			configPath := filepath.Join(gqlDir, entry.Name())
			authPath := filepath.Join(config.AuthDir(), name+".json")
			rootCmd.AddCommand(NewGraphQLCommand(name, configPath, authPath))
		}
	}
}

func isReservedName(name string) bool {
	reserved := []string{"add", "list", "remove", "auth", "help", "completion"}
	for _, r := range reserved {
		if name == r {
			return true
		}
	}
	return false
}

func exitWithError(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
