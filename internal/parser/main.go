package parser

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Config holds global configuration for Glyph.
type Config struct {
	Author string `toml:"author"`
}

// Parser is a struct that holds the file information and parsed commands.
type Parser struct {
	File        string
	Config      Config
	Commmands   []Command
	Content     []byte
	ContentRead bool
}

// Command represents a parsed command with its associated metadata.
type Command struct {
	Repo      string `toml:"repo"`
	LocalPath string `toml:"local_path"`
	Key       string `toml:"-"`
	Short     string `toml:"summary"`
	Long      string `toml:"description"`
	Branch    string `toml:"branch"`
	Tag       string `toml:"tag"`
	License   string `toml:"license"`
	Author    string `toml:"author"`
}

// IParser defines the interface for parsing operations.
type IParser interface {
	Read() (map[string]interface{}, error)
	Create()
	ExtractCommands()
	Refresh()
	DeleteSection()
}

// ValidateTemplateName checks if the provided name is a valid and safe template name.
func ValidateTemplateName(name string) error {
	if len(name) == 0 {
		return errors.New("template name cannot be empty")
	}

	if len(name) > 50 {
		return errors.New("template name is too long (max 50 characters)")
	}

	// Alphanumeric, hyphens, and underscores only
	validName := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validName.MatchString(name) {
		return errors.New("template name contains invalid characters (only alphanumeric, hyphens, and underscores allowed)")
	}

	// Reserved names check
	reserved := []string{"config", "help", "init", "create", "rm", "set", "list", "glyph"}
	for _, r := range reserved {
		if strings.ToLower(name) == r {
			return fmt.Errorf("'%s' is a reserved name and cannot be used as a template name", name)
		}
	}

	return nil
}

// safeWrite checks if the target file is a symlink before writing data to it.
func (p *Parser) safeWrite(data []byte) error {
	if info, err := os.Lstat(p.File); err == nil {
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("security error: %s is a symbolic link", p.File)
		}
	}

	return os.WriteFile(p.File, data, 0600)
}
