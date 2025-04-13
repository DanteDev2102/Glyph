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
	var builder strings.Builder
	builder.Grow(len(p.Content) + len(tmpl.Name) + len(tmpl.Repo) + 50)

	builder.WriteString("\n")
	builder.WriteString(string(p.Content))
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("[%s]", tmpl.Name))
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("repo = \"%s\"\n", tmpl.Repo))

	if len(tmpl.Branch) > 0 {
		builder.WriteString(fmt.Sprintf("branch = \"%s\"\n", tmpl.Branch))
	}

	if len(tmpl.Tag) > 0 {
		builder.WriteString(fmt.Sprintf("tag = \"%s\"\n", tmpl.Tag))
	}

	if len(tmpl.Summary) > 0 {
		builder.WriteString(fmt.Sprintf("short = \"%s\"\n", tmpl.Summary))
	}

	if len(tmpl.Description) > 0 {
		builder.WriteString(fmt.Sprintf("long = \"%s\"\n", tmpl.Branch))
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
	var config map[string]interface{}

	err := toml.Unmarshal(p.Content, &config)
	if err != nil {
		fmt.Println(err)
		return
	}

	value, ok := config[name]
	if !ok {
		fmt.Println("Not Exist this command")
		return
	}

	var newValues = make(map[string]string)

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

	fmt.Println(config)
	fmt.Println(value)
	fmt.Println(newValues)
	fmt.Println(tmpl)
}
