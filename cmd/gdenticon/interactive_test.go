package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInteractiveRequested(t *testing.T) {
	for _, args := range [][]string{nil, {"shell"}, {"--interactive"}} {
		if !interactiveRequested(args) {
			t.Fatalf("interactiveRequested(%q) = false", args)
		}
	}
	if interactiveRequested([]string{"alice@example.com"}) {
		t.Fatal("identifier requested interactive mode")
	}
}

func TestSplitCommandLine(t *testing.T) {
	got, err := splitCommandLine(`generate "Alice Example" --format svg --output "C:\Icons Here"`)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"generate", "Alice Example", "--format", "svg", "--output", `C:\Icons Here`}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("arguments are %#v, want %#v", got, want)
	}
	if _, err := splitCommandLine(`generate "unfinished`); err == nil {
		t.Fatal("unterminated quote returned nil error")
	}
}

func TestInteractiveGenerateAndContinueAfterError(t *testing.T) {
	dir := t.TempDir()
	input := strings.Join([]string{
		"help",
		"version",
		"hellp",
		"generate",
		`generate alice@example.com --format svg --size 32 --output "` + dir + `"`,
		"quit",
	}, "\n")
	var stdout, stderr bytes.Buffer

	if code := runProgram(nil, strings.NewReader(input), &stdout, &stderr); code != 0 {
		t.Fatalf("exit code = %d", code)
	}
	if !strings.Contains(stdout.String(), "interactive mode") || !strings.Contains(stdout.String(), "gdenticon> ") {
		t.Fatalf("stdout does not contain shell output:\n%s", stdout.String())
	}
	if !strings.Contains(stdout.String(), "Created:\n  ") {
		t.Fatalf("stdout does not label generated output:\n%s", stdout.String())
	}
	if !strings.Contains(stderr.String(), "expected exactly one id") {
		t.Fatalf("stderr does not contain recoverable command error:\n%s", stderr.String())
	}
	if !strings.Contains(stderr.String(), `Error: unknown command "hellp"`) {
		t.Fatalf("stderr does not contain unknown-command error:\n%s", stderr.String())
	}

	entries, err := os.ReadDir(filepath.Join(dir, "gdenticon"))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || filepath.Ext(entries[0].Name()) != ".svg" {
		t.Fatalf("generated files = %v", entries)
	}
}
