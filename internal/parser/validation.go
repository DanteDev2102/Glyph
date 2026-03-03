package parser

import (
	"errors"
	"fmt"
	"regexp"
)

var nameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

var reservedWords = map[string]bool{
	"config":     true,
	"help":       true,
	"init":       true,
	"create":     true,
	"rm":         true,
	"set":        true,
	"list":       true,
	"glyph":      true,
	"completion": true,
}

// ValidateTemplateName ensures that the template name is valid and not a reserved word.
func ValidateTemplateName(name string) error {
	if len(name) == 0 {
		return errors.New("template name cannot be empty")
	}
	if len(name) > 50 {
		return fmt.Errorf("template name %q is too long (max 50 characters)", name)
	}
	if !nameRegex.MatchString(name) {
		return fmt.Errorf("template name %q contains invalid characters (only alphanumeric, hyphens, and underscores allowed)", name)
	}
	if reservedWords[name] {
		return fmt.Errorf("template name %q is a reserved word", name)
	}
	return nil
}
