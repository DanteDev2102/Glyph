package parser

type Parser struct {
	File        string
	Commmands   []Command
	content     []byte
	contentRead bool
}

type Command struct {
	Path  string
	Cmd   string
	Repo  string
	Key   string
	Short string
	Long  string
}

type IParser interface {
	Read() (map[string]interface{}, error)
	Create() string
	ExtractCommands() error
	Refresh()
}
