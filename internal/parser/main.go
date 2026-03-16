package parser

import (
	"fmt"
	"os"
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

// safeWrite checks if the target file is a symlink before writing data to it.
func (p *Parser) safeWrite(data []byte) error {
	if info, err := os.Lstat(p.File); err == nil {
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("security error: %s is a symbolic link", p.File)
		}
	}

	if err := os.WriteFile(p.File, data, 0600); err != nil {
		return err
	}
	return os.Chmod(p.File, 0600)
}
