package parser

import (
	"fmt"

	"github.com/pelletier/go-toml/v2"
)

// DeleteSection removes a section from the configuration by its name.
func (p *Parser) DeleteSection(name string) {
	if err := ValidateTemplateName(name); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	var config map[string]interface{}

	err := toml.Unmarshal(p.Content, &config)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, ok := config[name]
	if !ok {
		fmt.Println("Not Exist this command")
		return
	}

	delete(config, name)

	data, er := toml.Marshal(config)
	if er != nil {
		fmt.Println(er)
		return
	}

	err = p.safeWrite(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	p.Refresh()
}
