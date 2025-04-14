package cli

import (
	"fmt"
	"os/exec"
	"path/filepath"
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
				if len(args) == 0 {
					fmt.Println("Directory path is required")
					return
				}

				if len(strings.TrimSpace(args[0])) <= 3 {
					fmt.Println("Not valid path")
					return
				}

				arg := []string{"clone", "--depth=1", command.Repo, args[0]}

				detail := command.Tag

				if len(detail) == 0 {
					detail = command.Branch
				}

				if len(detail) == 0 {
					detail = branch
				}

				if len(detail) == 0 {
					detail = tag
				}

				if len(detail) > 0 {
					arg = append(arg, "-b", detail, "--single-branch")
				}

				execute := exec.Command("git", arg...)
				err := execute.Run()
				if err != nil {
					return
				}

				gitDir := filepath.Join(args[0], ".git")
				execute = exec.Command("rm", "-rf", gitDir)
				err = execute.Run()
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
		Use:   "init [template name]",
		Short: "example",
		Long:  "example",
	}

	chargeTemplates(initCmd, &cli.Conf.Commmands)

	cli.Root.AddCommand(initCmd)
}
