package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSessionRemoveByNumberNameAndLast(t *testing.T) {
	root := t.TempDir()
	paths := []string{
		filepath.Join(root, "one", "gdenticon", "one.png"),
		filepath.Join(root, "two", "gdenticon", "two.png"),
		filepath.Join(root, "three", "gdenticon", "three.png"),
	}
	session := &interactiveSession{}
	for i, path := range paths {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("icon"), 0o600); err != nil {
			t.Fatal(err)
		}
		session.add(renderResult{ID: string(rune('a' + i)), Path: path, Created: true})
	}

	if _, err := session.remove("2"); err != nil {
		t.Fatal(err)
	}
	if _, err := session.remove("one.png"); err != nil {
		t.Fatal(err)
	}
	if output, err := session.remove("last"); err != nil {
		t.Fatal(err)
	} else if output.Path != paths[2] {
		t.Fatalf("removed %q, want %q", output.Path, paths[2])
	}
	for _, path := range paths {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("%s still exists: %v", path, err)
		}
	}
}

func TestSessionCleanAcrossDirectoriesPreservesOtherFiles(t *testing.T) {
	root := t.TempDir()
	session := &interactiveSession{}
	for _, name := range []string{"one", "two"} {
		dir := filepath.Join(root, name, "gdenticon")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(dir, name+".png")
		if err := os.WriteFile(path, []byte("icon"), 0o600); err != nil {
			t.Fatal(err)
		}
		session.add(renderResult{ID: name, Path: path, Created: true})
	}

	keptDir := filepath.Dir(session.outputs[0].Path)
	keptPath := filepath.Join(keptDir, "historical.png")
	if err := os.WriteFile(keptPath, []byte("old"), 0o600); err != nil {
		t.Fatal(err)
	}
	cachedPath := filepath.Join(root, "cached", "gdenticon", "cached.png")
	if err := os.MkdirAll(filepath.Dir(cachedPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cachedPath, []byte("cached"), 0o600); err != nil {
		t.Fatal(err)
	}
	session.add(renderResult{ID: "cached", Path: cachedPath, Created: false})

	removed, err := session.clean()
	if err != nil {
		t.Fatal(err)
	}
	if removed != 2 {
		t.Fatalf("removed %d files, want 2", removed)
	}
	for _, path := range []string{keptPath, cachedPath} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("preserved file %s: %v", path, err)
		}
	}
}

func TestSessionRejectsAmbiguousFilename(t *testing.T) {
	session := &interactiveSession{}
	session.add(renderResult{Path: filepath.Join("one", "same.png")})
	session.add(renderResult{Path: filepath.Join("two", "same.png")})
	if _, err := session.remove("same.png"); err == nil || !strings.Contains(err.Error(), "matches outputs") {
		t.Fatalf("ambiguous remove error = %v", err)
	}
}

func TestInteractiveCleanConfirmation(t *testing.T) {
	dir := t.TempDir()
	input := strings.Join([]string{
		`generate first@example.com --output "` + dir + `"`,
		"outputs",
		"clean",
		"y",
		"outputs",
		"quit",
	}, "\n")
	var stdout, stderr bytes.Buffer
	if code := runProgram(nil, strings.NewReader(input), &stdout, &stderr); code != 0 {
		t.Fatalf("exit code = %d, stderr:\n%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Remove 1 files created in this session?") ||
		!strings.Contains(stdout.String(), "Removed 1 files.") {
		t.Fatalf("cleanup output missing:\n%s", stdout.String())
	}
}
