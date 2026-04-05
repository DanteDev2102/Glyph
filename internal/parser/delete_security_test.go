package parser

import (
	"os"
	"path/filepath"
	"testing"
	"strings"
)

func TestDeleteSection_Config(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-parser-delete-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "repositories.toml")
	initialContent := `[config]
author = "Original Author"

[test-template]
repo = "https://github.com/test/repo"
`
	err = os.WriteFile(configPath, []byte(initialContent), 0600)
	if err != nil {
		t.Fatal(err)
	}

	p := &Parser{
		File:    configPath,
		Content: []byte(initialContent),
	}

	// Attempt to delete the [config] section
	p.DeleteSection("config")

	// Verify the [config] section still exists
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(content), "[config]") {
		t.Error("[config] section was deleted!")
	}
}
