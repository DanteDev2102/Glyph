package cli

import (
	"github.com/DanteDev2102/Glyph/internal/parser"
	"github.com/spf13/cobra"
)

func (cli *Base) CreateCmd() {
	cli.Root.AddCommand(&cobra.Command{
		Use:   "create --repo [repository url] --name [template name] --short [short description your command] --long [long description your command]",
		Short: "example",
		Long:  "example",
		Run: func(cmd *cobra.Command, args []string) {
			flags := parser.Template{
				Name:        name,
				Summary:     summary,
				Description: description,
				Branch:      branch,
				Tag:         tag,
				Repo:        repo,
			}

			cli.Conf.Write(&flags)
		},
	})
}
