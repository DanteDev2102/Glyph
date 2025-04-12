package cli

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/DanteDev2102/Glyph/internal/parser"
	"github.com/spf13/cobra"
)

func chargeTemplates(initCmd *cobra.Command, commands *[]parser.Command) {
	for i := range *commands {
		command := (*commands)[i]

		initCmd.AddCommand(&cobra.Command{
			Use:   fmt.Sprintf("%s [path]", command.Key),
			Short: command.Short,
			Long:  command.Long,
			Run: func(_ *cobra.Command, args []string) {
				execute := exec.Command("git", "clone", strings.TrimSpace(command.Repo))
				if len(args) > 0 {
					execute = exec.Command("git", "clone", command.Repo, args[0])
				}

				err := execute.Run()
				if err != nil {
					fmt.Println(err)
				}
			},
		})
	}
}

// InitCmd initializes the CLI with the "init" command.
func (cli *Base) InitCmd() {
	initCmd := &cobra.Command{
		Use:   "init [template]",
		Short: "example",
		Long:  "example",
	}

	chargeTemplates(initCmd, cli.Conf.Commmands)

	cli.Root.AddCommand(initCmd)
}
