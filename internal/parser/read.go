package parser

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

func (p *Parser) ensureContentRead() error {
	if !p.ContentRead {
		templates, err := os.ReadFile(p.File)
		if err != nil {
			return err
		}

		p.ContentRead = true
		p.Content = templates
	}
	return nil
}

func (p *Parser) Read() (map[string]interface{}, error) {
	var config map[string]interface{}

	if err := p.ensureContentRead(); err != nil {
		return config, err
	}

	err := toml.Unmarshal(p.Content, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

// ExtractCommands reads and processes commands from the configuration file.
// Optimization: Unmarshals directly into map[string]Command to avoid generic map allocations.
func (p *Parser) ExtractCommands() error {
	if err := p.ensureContentRead(); err != nil {
		return err
	}

	var data map[string]Command
	// toml.Unmarshal into a structured map is faster than map[string]interface{}
	// as it avoids interface boxing and manual type assertion.
	err := toml.Unmarshal(p.Content, &data)
	if err != nil {
		return err
	}

	// Pre-allocate slice capacity to avoid multiple re-allocations during the loop.
	commands := make([]Command, 0, len(data))

	for key, cmd := range data {
		cmd.Key = key
		commands = append(commands, cmd)
	}

	p.Commmands = commands

	return nil
}

// Refresh reloads the commands from the configuration file.
func (p *Parser) Refresh() {
	err := p.ExtractCommands()

	if err != nil {
		fmt.Println(err)
	}
}
