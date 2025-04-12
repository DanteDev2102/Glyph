package parser

import (
	"fmt"
	"os"
	"strings"
)

type Template struct {
	Name        string
	Repo        string
	Description string
	Summary     string
	Branch      string
	Tag         string
}

func (p *Parser) Write(tmpl *Template) {
	content := string(p.Content)

	builder := strings.Builder{}

	builder.WriteString("\n")
	builder.WriteString(content)
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("[%s]", tmpl.Name))
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("repo = \"%s\"", tmpl.Repo))

	if len(tmpl.Branch) > 0 {
		builder.WriteString("\n")
		builder.WriteString(fmt.Sprintf("branch = \"%s\"", tmpl.Branch))
	}

	if len(tmpl.Tag) > 0 {
		builder.WriteString("\n")
		builder.WriteString(fmt.Sprintf("tag = \"%s\"", tmpl.Tag))
	}

	if len(tmpl.Summary) > 0 {
		builder.WriteString("\n")
		builder.WriteString(fmt.Sprintf("summary = \"%s\"", tmpl.Summary))
	}

	if len(tmpl.Description) > 0 {
		builder.WriteString("\n")
		builder.WriteString(fmt.Sprintf("description = \"%s\"", tmpl.Branch))
	}

	err := os.WriteFile(p.File, []byte(builder.String()), 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	p.Refresh()
}
