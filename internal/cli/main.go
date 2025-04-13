package cli

import (
	parser "github.com/DanteDev2102/Glyph/internal/parser"
	cli "github.com/spf13/cobra"
)

// Base is a struct that holds the configuration parser and the root command.
type Base struct {
	Conf *parser.Parser
	Root *cli.Command
}

// Cli is an instance of Base that holds the configuration parser and root command.
var Cli *Base

// IBase defines the interface for initializing commands.
type IBase interface {
	InitCmd()
	CreateCmd()
	UpdateCmd()
	DeleteCmd()
}

var (
	repo        string
	branch      string
	tag         string
	summary     string
	name        string
	description string
)

func init() {
	conf := &parser.Parser{File: "/app/test.toml"}

	err := conf.ExtractCommands()
	if err != nil {
		panic("END")
	}

	Cli = &Base{
		Conf: conf,
		Root: &cli.Command{
			Use:   "glyph",
			Short: "example",
			Long:  "example",
			CompletionOptions: cli.CompletionOptions{
				HiddenDefaultCmd: true,
			},
		},
	}

	Cli.Root.PersistentFlags().StringVarP(&name, "name", "n", "", "--name or -n [name template]")
	Cli.Root.PersistentFlags().StringVarP(&repo, "repo", "r", "", "--repo or -r [repository url]")
	Cli.Root.PersistentFlags().StringVarP(&branch, "branch", "b", "", "--branch or -b [branch]")
	Cli.Root.PersistentFlags().StringVarP(&tag, "tag", "t", "", "--tag or -t [tag]")
	Cli.Root.PersistentFlags().StringVarP(&summary, "summary", "s", "", "--summary or -s [summary]")
	Cli.Root.PersistentFlags().StringVarP(&description, "description", "d", "", "--description or -d [description]")

	Cli.CreateCmd()
	Cli.InitCmd()
	Cli.UpdateCmd()
	Cli.DeleteCmd()

}
