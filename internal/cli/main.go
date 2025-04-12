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
}

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

	Cli.InitCmd()
}
