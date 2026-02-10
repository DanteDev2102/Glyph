package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// ListCmd adds the "list" command to the CLI root command.
func (cli *Base) ListCmd() {
	cli.Root.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all registered templates",
		Run: func(cmd *cobra.Command, args []string) {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "NAME\tSUMMARY\tSOURCE")

			for _, command := range cli.Conf.Commmands {
				source := command.Repo
				if command.LocalPath != "" {
					source = command.LocalPath + " (Local)"
				}
				fmt.Fprintf(w, "%s\t%s\t%s\n", command.Key, command.Short, source)
			}
			w.Flush()
		},
	})
}
