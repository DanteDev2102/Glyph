package parser

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

func extractValue(data *map[string]any, key, defaultValue string) string {
	if data == nil {
		return defaultValue
	}

	if val, ok := (*data)[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}

	return defaultValue
}

func (p *Parser) Read() (map[string]interface{}, error) {
	var config map[string]interface{}

	if !p.contentRead {
		templates, err := os.ReadFile(p.File)
		if err != nil {
			return config, err
		}

		p.contentRead = true
		p.content = templates
	}

	err := toml.Unmarshal(p.content, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

// ExtractCommands reads and processes commands from the configuration file.
func (p *Parser) ExtractCommands() error {
	data, err := p.Read()
	commands := []Command{}

	if err != nil {
		return err
	}

	for key, values := range data {
		_, isMap := values.(map[string]interface{})

		if !isMap {
			continue
		}

		valueMap, ok := values.(map[string]interface{})

		if !ok {
			continue
		}

		commands = append(commands, Command{
			Key:   key,
			Repo:  extractValue(&valueMap, "repo", ""),
			Long:  extractValue(&valueMap, "long", ""),
			Short: extractValue(&valueMap, "short", ""),
		})
	}

	p.Commmands = &commands

	return nil
}

// Refresh reloads the commands from the configuration file.
func (p *Parser) Refresh() {
	err := p.ExtractCommands()

	if err != nil {
		fmt.Println(err)
	}
}
