package gitutils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateNonEmpty(t *testing.T) {
	tests := []struct {
		input    string
		wantErr bool
	}{
		{"user", false},
		{"  user  ", false},
		{"", true},
		{"   ", true},
	}

	for _, tt := range tests {
		err := ValidateNonEmpty(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateNonEmpty(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
	}
}

func TestValidateFileExists(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-auth-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	existingFile := filepath.Join(tmpDir, "exists")
	err = os.WriteFile(existingFile, []byte("data"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	nonExistingFile := filepath.Join(tmpDir, "not-exists")

	tests := []struct {
		path     string
		wantErr bool
	}{
		{existingFile, false},
		{nonExistingFile, true},
	}

	for _, tt := range tests {
		err := ValidateFileExists(tt.path)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateFileExists(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
		}
	}
}
