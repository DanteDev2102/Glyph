package cli

import (
	"github.com/DanteDev2102/Glyph/config"
	"github.com/spf13/cobra"
)

// DeleteCmd adds the "rm" command to the CLI, allowing users to delete a configuration section by name.
func (cli *Base) DeleteCmd() {
	cli.Root.AddCommand(&cobra.Command{
		Use:   config.RmUse,
		Short: config.RmSummary,
		Long:  config.RmDescription,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cli.Conf.DeleteSection(args[0])
		},
	})
}
