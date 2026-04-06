package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pelletier/go-toml/v2"
)

func TestWriteSection_ConfigProtection(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-parser-write-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "repositories.toml")
	initialContent := []byte("[config]\nauthor = \"original-author\"\n")
	err = os.WriteFile(configPath, initialContent, 0600)
	if err != nil {
		t.Fatal(err)
	}

	p := &Parser{
		File:    configPath,
		Content: initialContent,
	}

	tmpl := &Template{
		Author: "malicious-author",
	}

	// Try to modify 'config' using WriteSection
	p.WriteSection(tmpl, "config")

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}

	var config map[string]map[string]string
	err = toml.Unmarshal(data, &config)
	if err != nil {
		t.Fatal(err)
	}

	if author, ok := config["config"]["author"]; ok && author == "malicious-author" {
		t.Errorf("Config section was modified! Protection for WriteSection is missing.")
	}
}

func TestWriteSection_CaseInsensitiveConfigProtection(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-parser-write-case-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "repositories.toml")
	initialContent := []byte("[config]\nauthor = \"original-author\"\n")
	err = os.WriteFile(configPath, initialContent, 0600)
	if err != nil {
		t.Fatal(err)
	}

	p := &Parser{
		File:    configPath,
		Content: initialContent,
	}

	tmpl := &Template{
		Author: "malicious-author",
	}

	// Try to modify 'CONFIG' using WriteSection
	p.WriteSection(tmpl, "CONFIG")

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}

	var config map[string]map[string]string
	toml.Unmarshal(data, &config)

	if author, ok := config["config"]["author"]; ok && author == "malicious-author" {
		t.Errorf("Config section was modified using uppercase 'CONFIG'! Case-insensitive protection for WriteSection is missing.")
	}
}
