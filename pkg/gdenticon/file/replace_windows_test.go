//go:build windows

package file

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReplacePathRejectsInvalidNames(t *testing.T) {
	if err := replacePath("source\x00", "target"); err == nil {
		t.Fatal("replacePath accepted a source containing NUL")
	}
	if err := replacePath("source", "target\x00"); err == nil {
		t.Fatal("replacePath accepted a target containing NUL")
	}
}

func TestReplacePathFailures(t *testing.T) {
	dir := t.TempDir()
	if err := replacePath(filepath.Join(dir, "missing"), filepath.Join(dir, "target")); err == nil {
		t.Fatal("replacePath returned nil for missing source")
	}

	source := filepath.Join(dir, "source")
	target := filepath.Join(dir, "directory")
	if err := os.WriteFile(source, []byte("content"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(target, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := replacePath(source, target); err == nil {
		t.Fatal("replacePath replaced a directory")
	}
}
