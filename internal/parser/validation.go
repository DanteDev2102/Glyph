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

var nameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// ValidateTemplateName checks if the template name is valid.
func ValidateTemplateName(name string) error {
	if len(name) == 0 {
		return errors.New("template name cannot be empty")
	}

	if len(name) > 50 {
		return fmt.Errorf("template name too long (max 50 characters): %d", len(name))
	}

	if !nameRegex.MatchString(name) {
		return errors.New("template name contains invalid characters (only alphanumeric, hyphens, and underscores allowed)")
	}

	if reservedWords[strings.ToLower(name)] {
		return fmt.Errorf("template name '%s' is a reserved word and cannot be used", name)
	}

	return nil
}
