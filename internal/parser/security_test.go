package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWrite_Symlink(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-parser-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "repositories.toml")
	targetFile := filepath.Join(tmpDir, "sensitive_file")
	err = os.WriteFile(targetFile, []byte("SENSITIVE DATA"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Create a symlink from configPath to targetFile
	err = os.Symlink(targetFile, configPath)
	if err != nil {
		t.Fatal(err)
	}

	p := &Parser{
		File: configPath,
	}

	tmpl := &Template{
		Name: "test-template",
		Repo: "https://github.com/test/repo",
	}

	// This should NOT overwrite targetFile if we have symlink protection
	p.Write(tmpl)

	// Check if targetFile was modified
	content, err := os.ReadFile(targetFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != "SENSITIVE DATA" {
		t.Errorf("Target file was modified through symlink! Content: %s", string(content))
	}
}

func TestValidateTemplateName(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"valid-name", false},
		{"ValidName_123", false},
		{"", true},
		{"config", true},
		{"help", true},
		{"init", true},
		{"create", true},
		{"rm", true},
		{"set", true},
		{"list", true},
		{"glyph", true},
		{"CONFIG", true},
		{"invalid name", true},
		{"invalid@name", true},
		{"invalid/name", true},
		{"this-name-is-way-too-long-for-the-validation-to-pass-successfully", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateTemplateName(tt.name); (err != nil) != tt.wantErr {
				t.Errorf("ValidateTemplateName(%q) error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}
