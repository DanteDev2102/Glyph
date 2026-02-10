package parser

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
