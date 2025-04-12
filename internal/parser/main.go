package parser

// Parser is a struct that holds the file information and parsed commands.
type Parser struct {
	File        string
	Commmands   []Command
	Content     []byte
	ContentRead bool
}

// Command represents a parsed command with its associated metadata.
type Command struct {
	Repo  string
	Key   string
	Short string
	Long  string
}

// IParser defines the interface for parsing operations.
type IParser interface {
	Read() (map[string]interface{}, error)
	Create()
	ExtractCommands()
	Refresh()
}
