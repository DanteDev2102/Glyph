package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

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

var validNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// ValidateTemplateName ensures the template name is safe and valid.
func ValidateTemplateName(name string) error {
	if len(name) == 0 {
		return errors.New("template name cannot be empty")
	}

	if len(name) > 50 {
		return fmt.Errorf("template name is too long (max 50 characters): %d", len(name))
	}

	if !validNameRegex.MatchString(name) {
		return fmt.Errorf("template name contains invalid characters: %s (only alphanumeric, hyphens, and underscores allowed)", name)
	}

	if reservedWords[strings.ToLower(name)] {
		return fmt.Errorf("template name '%s' is a reserved word", name)
	}

	return nil
}
