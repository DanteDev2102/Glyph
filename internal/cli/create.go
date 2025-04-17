package cli

import (
	"fmt"

	"github.com/DanteDev2102/Glyph/config"
	"github.com/DanteDev2102/Glyph/internal/parser"
	"github.com/spf13/cobra"
)

// CreateCmd adds the "create" command to the CLI root command.
func (cli *Base) CreateCmd() {
	cli.Root.AddCommand(&cobra.Command{
		Use:   config.CreateUse,
		Short: config.CreateSummary,
		Long:  config.CreateDescription,
		Run: func(cmd *cobra.Command, args []string) {
			if name == "" || repo == "" {
				fmt.Println("name and repo flags required")
				return
			} else if branch != "" && tag != "" {
				fmt.Println("Use only branch or only tag for create new project template")
				return
			}

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
