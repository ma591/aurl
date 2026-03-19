package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/shawnpana/aurl/internal/config"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered APIs",
	RunE:  runList,
}

func runList(cmd *cobra.Command, args []string) error {
	apis, err := config.ListAPIs()
	if err != nil {
		return err
	}
	gqls, err := config.ListGraphQL()
	if err != nil {
		return err
	}

	if len(apis) == 0 && len(gqls) == 0 {
		fmt.Println("No APIs registered. Use 'aurl add [name] [spec]' to register one.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tTYPE\tTITLE\tVERSION\tENDPOINT")
	for _, api := range apis {
		fmt.Fprintf(w, "%s\tapi\t%s\t%s\t%s\n", api.Name, api.Title, api.Version, api.BaseURL)
	}
	for _, gql := range gqls {
		fmt.Fprintf(w, "%s\tgraphql\t\t\t%s\n", gql.Name, gql.Endpoint)
	}
	w.Flush()

	return nil
}
