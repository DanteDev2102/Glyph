package parser

import (
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// Template represents the structure for storing template metadata.
type Template struct {
	Name        string
	Repo        string
	Description string
	Summary     string
	Branch      string
	Tag         string
}

func (p *Parser) Write(tmpl *Template) {
	for _, value := range p.Commmands {
		if value.Key == tmpl.Name {
			fmt.Println("This command already exist")
			return
		}
	}

	var builder strings.Builder
	builder.Grow(len(p.Content) + len(tmpl.Name) + len(tmpl.Repo) + 50)

	builder.WriteString(
		fmt.Sprintf("%s\n[%s]\nrepo = \"%s\"\n", string(p.Content), tmpl.Name, tmpl.Repo),
	)

	fields := []struct {
		key   string
		value string
	}{
		{"branch", tmpl.Branch},
		{"tag", tmpl.Tag},
		{"summary", tmpl.Summary},
		{"description", tmpl.Description},
	}

	for _, field := range fields {
		if field.value != "" {
			builder.WriteString(fmt.Sprintf("%s = \"%s\"\n", field.key, field.value))
		}
	}

	builder.WriteString("\n")

	err := os.WriteFile(p.File, []byte(builder.String()), 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	p.Refresh()
}

// WriteSection updates or adds a section in the TOML configuration based on the provided template and section name.
func (p *Parser) WriteSection(tmpl *Template, name string) {
	var config map[string]map[string]string

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

	if tmpl.Branch != "" && tmpl.Tag != "" {
		fmt.Println("Only branch or tag flag")
		return
	}

	var newValues = make(map[string]string)
	newValues = map[string]string{
		"repo":        config[name]["repo"],
		"summary":     config[name]["summary"],
		"description": config[name]["description"],
	}

	if len(tmpl.Repo) > 0 {
		newValues["repo"] = tmpl.Repo
	}

	if len(tmpl.Branch) > 0 {
		newValues["branch"] = tmpl.Branch
	}

	if len(tmpl.Description) > 0 {
		newValues["description"] = tmpl.Description
	}

	if len(tmpl.Summary) > 0 {
		newValues["summary"] = tmpl.Summary
	}

	if len(tmpl.Tag) > 0 {
		newValues["tag"] = tmpl.Tag
	}

	if len(tmpl.Name) > 0 {
		delete(config, name)
		config[tmpl.Name] = newValues
	} else {
		config[name] = newValues
	}

	data, er := toml.Marshal(config)
	if er != nil {
		fmt.Println(er)
		return
	}

	err = os.WriteFile(p.File, data, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	p.Refresh()
}
