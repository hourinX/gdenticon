package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/iotest"
)

type failAtWriter struct {
	failAt int
	writes int
	failed bool
}

func (w *failAtWriter) Write(p []byte) (int, error) {
	w.writes++
	if w.writes == w.failAt {
		w.failed = true
		return 0, errors.New("write failed")
	}
	return len(p), nil
}

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

func TestInteractiveCommandValidation(t *testing.T) {
	input := strings.Join([]string{
		`generate "unfinished`,
		"outputs extra",
		"remove",
		"remove last",
		"clean extra",
		"clean",
		"outputs",
		"quit",
	}, "\n")
	var stdout, stderr bytes.Buffer

	if code := runInteractive(strings.NewReader(input), &stdout, &stderr); code != 0 {
		t.Fatalf("exit code = %d", code)
	}
	for _, want := range []string{
		"unterminated quoted value",
		"outputs does not accept arguments",
		"remove requires an output number",
		"no session output is available",
		"clean does not accept arguments",
	} {
		if !strings.Contains(stderr.String(), want) {
			t.Errorf("stderr does not contain %q:\n%s", want, stderr.String())
		}
	}
	if !strings.Contains(stdout.String(), "No files created in this session") ||
		!strings.Contains(stdout.String(), "No outputs have been generated") {
		t.Fatalf("stdout does not contain empty-session messages:\n%s", stdout.String())
	}
}

func TestInteractivePropagatesOutputErrors(t *testing.T) {
	input := strings.Join([]string{
		"help",
		"outputs",
		"clean",
		"unknown",
		"version",
		"quit",
	}, "\n")

	coveredFailures := 0
	for failAt := 1; failAt <= 40; failAt++ {
		writer := &failAtWriter{failAt: failAt}
		code := runInteractive(strings.NewReader(input), writer, writer)
		if writer.failed {
			coveredFailures++
			if code != 1 {
				t.Fatalf("write %d failed but exit code = %d", failAt, code)
			}
		}
	}
	if coveredFailures < 20 {
		t.Fatalf("covered only %d output failures", coveredFailures)
	}
}

func TestRunReportsWriterErrors(t *testing.T) {
	writer := &failAtWriter{failAt: 1}
	if code := run([]string{"--unknown"}, io.Discard, writer); code != 1 || !writer.failed {
		t.Fatalf("usage write failure: code=%d failed=%v", code, writer.failed)
	}

	writer = &failAtWriter{failAt: 1}
	if code := run([]string{"first@example.com", "--output", t.TempDir()}, writer, io.Discard); code != 1 || !writer.failed {
		t.Fatalf("result write failure: code=%d failed=%v", code, writer.failed)
	}
}

func TestInteractiveRemoveAndCancelClean(t *testing.T) {
	for _, test := range []struct {
		name  string
		lines func(string) []string
		want  string
	}{
		{
			name: "remove last",
			lines: func(dir string) []string {
				return []string{
					`generate remove@example.com --output "` + dir + `"`,
					"remove last",
					"outputs",
					"quit",
				}
			},
			want: "Removed:\n  ",
		},
		{
			name: "cancel clean",
			lines: func(dir string) []string {
				return []string{
					`generate cancel@example.com --output "` + dir + `"`,
					"clean",
					"no",
					"quit",
				}
			},
			want: "Cancelled.",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			input := strings.Join(test.lines(t.TempDir()), "\n")
			if code := runInteractive(strings.NewReader(input), &stdout, &stderr); code != 0 {
				t.Fatalf("exit code = %d, stderr:\n%s", code, stderr.String())
			}
			if !strings.Contains(stdout.String(), test.want) {
				t.Fatalf("stdout does not contain %q:\n%s", test.want, stdout.String())
			}
		})
	}
}

func TestInteractiveCleanConfirmationEOF(t *testing.T) {
	dir := t.TempDir()
	input := strings.Join([]string{
		`generate eof@example.com --output "` + dir + `"`,
		"clean",
	}, "\n")
	var stdout, stderr bytes.Buffer
	if code := runInteractive(strings.NewReader(input), &stdout, &stderr); code != 0 {
		t.Fatalf("exit code = %d, stderr:\n%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "[y/N]") {
		t.Fatalf("confirmation prompt missing:\n%s", stdout.String())
	}
}

func TestInteractiveReportsInputError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := runInteractive(iotest.ErrReader(errors.New("read failed")), &stdout, &stderr); code != 1 {
		t.Fatalf("exit code = %d", code)
	}
	if !strings.Contains(stderr.String(), "read command") {
		t.Fatalf("stderr does not report input failure: %q", stderr.String())
	}
}

func TestRunProgramDirectMode(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := runProgram([]string{"--version"}, strings.NewReader(""), &stdout, &stderr); code != 0 {
		t.Fatalf("exit code = %d, stderr: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "gdenticon") {
		t.Fatalf("version output = %q", stdout.String())
	}
}
