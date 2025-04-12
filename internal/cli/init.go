package cli

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// ChargeTemplates initializes and adds template-related commands to the CLI.
func (cli *Base) ChargeTemplates() {
	commands := cli.Conf.Commmands
	initCmd := &cobra.Command{
		Use:   "init [template]",
		Short: "example",
		Long:  "example",
	}

	for i := range commands {
		command := commands[i]

		initCmd.AddCommand(&cobra.Command{
			Use:   fmt.Sprintf("%s [path]", command.Key),
			Short: command.Short,
			Long:  command.Long,
			Run: func(_ *cobra.Command, args []string) {
				cmd := strings.TrimSpace(command.Cmd)
				if len(cmd) > 0 {
					if err := exec.Command(cmd).Run(); err != nil {
						fmt.Println(err)
					}
					return
				}

				execute := exec.Command("git", "clone", command.Repo)
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

	cli.Root.AddCommand(initCmd)
}
