package parser

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

func (p *Parser) ensureContentRead() error {
	if !p.ContentRead {
		if info, err := os.Lstat(p.File); err == nil {
			if info.Mode()&os.ModeSymlink != 0 {
				return fmt.Errorf("security error: %s is a symbolic link", p.File)
			}
		}

		templates, err := os.ReadFile(p.File)
		if err != nil {
			if os.IsNotExist(err) {
				p.ContentRead = true
				p.Content = []byte("")
				return nil
			}
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
func (p *Parser) ExtractCommands() error {
	if err := p.ensureContentRead(); err != nil {
		return err
	}

	var data map[string]interface{}
	err := toml.Unmarshal(p.Content, &data)
	if err != nil {
		return err
	}

	var commands []Command
	for key, value := range data {
		if key == "config" {
			configBytes, _ := toml.Marshal(value)
			toml.Unmarshal(configBytes, &p.Config)
			continue
		}

		// Convert value back to bytes to unmarshal into Command struct
		cmdBytes, _ := toml.Marshal(value)
		var cmd Command
		toml.Unmarshal(cmdBytes, &cmd)
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
