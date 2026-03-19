package cmd

import (
	"fmt"

	"github.com/shawnpana/arc/internal/config"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Unregister an API",
	Args:  cobra.ExactArgs(1),
	RunE:  runRemove,
}

func runRemove(cmd *cobra.Command, args []string) error {
	name := args[0]

	// Try API first, then GraphQL
	errAPI := config.DeleteAPI(name)
	errGQL := config.DeleteGraphQL(name)

	if errAPI != nil && errGQL != nil {
		return fmt.Errorf("%q not found. Run 'arc list' to see registered commands.", name)
	}

	fmt.Printf("Removed %q.\n", name)
	return nil
}
