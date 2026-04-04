package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pelletier/go-toml/v2"
)

func TestConfigSectionProtection(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-config-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "repositories.toml")
	initialConfig := `[config]
author = "original-author"

[test-template]
repo = "https://github.com/test/repo"
`
	err = os.WriteFile(configPath, []byte(initialConfig), 0600)
	if err != nil {
		t.Fatal(err)
	}

	p := &Parser{File: configPath}
	_ = p.ExtractCommands()

	t.Run("WriteSection Protection", func(t *testing.T) {
		p.WriteSection(&Template{Author: "malicious"}, "config")
		p.Refresh()
		var data map[string]interface{}
		content, _ := os.ReadFile(configPath)
		_ = toml.Unmarshal(content, &data)
		configSection := data["config"].(map[string]interface{})
		if configSection["author"] == "malicious" {
			t.Error("Vulnerability: config section modified")
		}
	})

	t.Run("DeleteSection Protection", func(t *testing.T) {
		p.DeleteSection("config")
		p.Refresh()
		var data map[string]interface{}
		content, _ := os.ReadFile(configPath)
		_ = toml.Unmarshal(content, &data)
		if _, ok := data["config"]; !ok {
			t.Error("Vulnerability: config section deleted")
		}
	})
}
