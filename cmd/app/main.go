package main

import (
	"fmt"
	"os"

	cli "github.com/DanteDev2102/Glyph/internal/cli"
)

func main() {
	cli := cli.Cli.Root

	if err := cli.Execute(); err != nil {
		fmt.Println("Error")
		os.Exit(1)
	}
}
