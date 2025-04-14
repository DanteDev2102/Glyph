package cli

import (
	"github.com/DanteDev2102/Glyph/internal/parser"
	"github.com/spf13/cobra"
)

// UpdateCmd adds the "set" command to the CLI, allowing users to configure templates.
func (cli *Base) UpdateCmd() {
	cli.Root.AddCommand(&cobra.Command{
		Use:   "set [name template] --repo [repository url] --name [name template] --branch [branch] --tag [tag] --summary [summary] --description [description]",
		Short: "example",
		Long:  "example",
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
