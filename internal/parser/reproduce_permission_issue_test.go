package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSafeWrite_EnforcePermissions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-parser-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	confFile := filepath.Join(tmpDir, "repositories.toml")
	// Pre-create with insecure permissions
	err = os.WriteFile(confFile, []byte(""), 0644)
	if err != nil {
		t.Fatal(err)
	}

	p := &Parser{File: confFile}
	err = p.safeWrite([]byte("test = 1"))
	if err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(confFile)
	if err != nil {
		t.Fatal(err)
	}

	if info.Mode().Perm() != 0600 {
		t.Errorf("safeWrite failed to enforce 0600: expected %o, got %o", 0600, info.Mode().Perm())
	}
}
