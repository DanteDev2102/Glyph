package cli

import (
	"github.com/DanteDev2102/Glyph/config"
	"github.com/DanteDev2102/Glyph/internal/parser"
	"github.com/spf13/cobra"
)

// UpdateCmd adds the "set" command to the CLI, allowing users to configure templates.
func (cli *Base) UpdateCmd() {
	cli.Root.AddCommand(&cobra.Command{
		Use:   config.SetUse,
		Short: config.SetSummary,
		Long:  config.SetDescription,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cli.Conf.WriteSection(&parser.Template{
				Name:        name,
				Repo:        repo,
				Description: description,
				Summary:     summary,
				Branch:      branch,
				Tag:         tag,
			}, args[0])
		},
	})
}
