package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// Template represents the structure for storing template metadata.
type Template struct {
	Name        string
	Repo        string
	LocalPath   string
	Description string
	Summary     string
	Branch      string
	Tag         string
	License     string
	Author      string
}

func (p *Parser) Write(tmpl *Template) {
	if err := ValidateTemplateName(tmpl.Name); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	var config map[string]interface{}
	if err := toml.Unmarshal(p.Content, &config); err != nil || config == nil {
		config = make(map[string]interface{})
	}

	if _, ok := config[tmpl.Name]; ok {
		fmt.Println("This command or section already exists")
		return
	}

	newTmpl := map[string]string{
		"repo":        tmpl.Repo,
		"local_path":  tmpl.LocalPath,
		"branch":      tmpl.Branch,
		"tag":         tmpl.Tag,
		"summary":     tmpl.Summary,
		"description": tmpl.Description,
		"license":     tmpl.License,
		"author":      tmpl.Author,
	}

	// Remove empty strings to keep config clean
	for k, v := range newTmpl {
		if v == "" {
			delete(newTmpl, k)
		}
	}

	config[tmpl.Name] = newTmpl

	data, err := toml.Marshal(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	dir := filepath.Dir(p.File)
	os.MkdirAll(dir, 0755)

	if err := p.safeWrite(data); err != nil {
		fmt.Println(err)
		return
	}

	p.Refresh()
}

// WriteSection updates or adds a section in the TOML configuration based on the provided template and section name.
func (p *Parser) WriteSection(tmpl *Template, name string) {
	if tmpl.Name != "" {
		if err := ValidateTemplateName(tmpl.Name); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
	}

	var config map[string]map[string]string

	err := toml.Unmarshal(p.Content, &config)
	if err != nil {
		fmt.Println(err)
		return
	}

	if strings.EqualFold(name, "config") {
		fmt.Println("Cannot modify reserved 'config' section")
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
	for k, v := range config[name] {
		newValues[k] = v
	}

	if len(tmpl.Repo) > 0 {
		newValues["repo"] = tmpl.Repo
	}

	if len(tmpl.LocalPath) > 0 {
		newValues["local_path"] = tmpl.LocalPath
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

	if len(tmpl.License) > 0 {
		newValues["license"] = tmpl.License
	}

	if len(tmpl.Author) > 0 {
		newValues["author"] = tmpl.Author
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

	dir := filepath.Dir(p.File)
	os.MkdirAll(dir, 0755)

	err = p.safeWrite(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	p.Refresh()
}
