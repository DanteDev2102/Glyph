package cli

import (
	"fmt"

	"github.com/DanteDev2102/Glyph/config"
	"github.com/DanteDev2102/Glyph/internal/gitutils"
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
			if repo != "" {
				_, err := gitutils.ValidateRepo(repo)
				if err != nil {
					fmt.Printf("Error validating repository: %v\n", err)
					return
				}
			}

			cli.Conf.WriteSection(&parser.Template{
				Name:        name,
				Repo:        repo,
				LocalPath:   localPath,
				Description: description,
				Summary:     summary,
				Branch:      branch,
				Tag:         tag,
				License:     license,
				Author:      author,
			}, args[0])
		},
	})
}
