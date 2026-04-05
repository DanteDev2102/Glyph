package parser

import (
	"os"
	"path/filepath"
	"testing"
	"strings"
)

func TestWriteSection_Config(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-parser-write-test-*")
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

	// Attempt to modify the [config] section using WriteSection
	p.WriteSection(&Template{
		Repo: "https://evil.com/repo",
	}, "config")

	// Verify the [config] section was not corrupted or overwritten with template fields
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}

	// Currently, WriteSection will fail to unmarshal if [config] has different structure
	// Let's see what happened
	t.Logf("Config content after WriteSection:\n%s", string(content))

	if strings.Contains(string(content), "repo = \"https://evil.com/repo\"") && strings.Contains(string(content), "[config]") {
		t.Error("[config] section was overwritten with template fields!")
	}
}
