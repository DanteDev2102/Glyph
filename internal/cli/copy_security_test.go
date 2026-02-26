package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyDir_SymlinkSource(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	srcDir := filepath.Join(tmpDir, "src")
	err = os.Mkdir(srcDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	secretFile := filepath.Join(tmpDir, "secret")
	err = os.WriteFile(secretFile, []byte("TOP SECRET"), 0600)
	if err != nil {
		t.Fatal(err)
	}

	symlinkFile := filepath.Join(srcDir, "link")
	err = os.Symlink(secretFile, symlinkFile)
	if err != nil {
		t.Fatal(err)
	}

	dstDir := filepath.Join(tmpDir, "dst")
	err = copyDir(srcDir, dstDir)
	if err != nil {
		t.Fatal(err)
	}

	copiedFile := filepath.Join(dstDir, "link")
	if _, err := os.Stat(copiedFile); !os.IsNotExist(err) {
		t.Errorf("copyDir should have skipped symlink at source, but file exists at destination")
	}
}

func TestCopyDir_SymlinkDestination(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	srcDir := filepath.Join(tmpDir, "src")
	err = os.Mkdir(srcDir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(filepath.Join(srcDir, "file"), []byte("new content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	targetFile := filepath.Join(tmpDir, "target")
	err = os.WriteFile(targetFile, []byte("original content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	dstDir := filepath.Join(tmpDir, "dst")
	err = os.Mkdir(dstDir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.Symlink(targetFile, filepath.Join(dstDir, "file"))
	if err != nil {
		t.Fatal(err)
	}

	err = copyDir(srcDir, dstDir)
	if err == nil {
		t.Errorf("copyDir should have failed when destination is a symlink")
	}

	content, err := os.ReadFile(targetFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "original content" {
		t.Errorf("copyDir followed symlink at destination and overwrote target file!")
	}
}

func TestCopyFile_SymlinkDestination(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	srcFile := filepath.Join(tmpDir, "src")
	err = os.WriteFile(srcFile, []byte("new license"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	targetFile := filepath.Join(tmpDir, "target")
	err = os.WriteFile(targetFile, []byte("original content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	dstFile := filepath.Join(tmpDir, "dst_link")
	err = os.Symlink(targetFile, dstFile)
	if err != nil {
		t.Fatal(err)
	}

	err = copyFile(srcFile, dstFile)
	if err == nil {
		t.Errorf("copyFile should have failed when destination is a symlink")
	}

	content, err := os.ReadFile(targetFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "original content" {
		t.Errorf("copyFile followed symlink at destination and overwrote target file!")
	}
}

func TestCopyFile_SymlinkSource(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "glyph-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	secretFile := filepath.Join(tmpDir, "secret")
	err = os.WriteFile(secretFile, []byte("TOP SECRET"), 0600)
	if err != nil {
		t.Fatal(err)
	}

	srcFile := filepath.Join(tmpDir, "src_link")
	err = os.Symlink(secretFile, srcFile)
	if err != nil {
		t.Fatal(err)
	}

	dstFile := filepath.Join(tmpDir, "dst")

	err = copyFile(srcFile, dstFile)
	if err == nil {
		t.Errorf("copyFile should have failed when source is a symlink")
	}

	if _, err := os.Stat(dstFile); !os.IsNotExist(err) {
		t.Errorf("copyFile should not have created destination file when source is a symlink")
	}
}
