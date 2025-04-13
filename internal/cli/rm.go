package cli

import (
	"github.com/spf13/cobra"
)

// DeleteCmd adds the "rm" command to the CLI, allowing users to delete a configuration section by name.
func (cli *Base) DeleteCmd() {
	cli.Root.AddCommand(&cobra.Command{
		Use:   "rm [name template]",
		Short: "example",
		Long:  "example",
		Run: func(cmd *cobra.Command, args []string) {
			cli.Conf.DeleteSection(args[0])
		},
	})
}
