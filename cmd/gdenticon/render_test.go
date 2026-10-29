package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/hourinX/gdenticon/v2/pkg/gdenticon"
)

func validRenderFlags(output string) renderFlags {
	return renderFlags{
		grid:   gdenticon.DefaultGridSize,
		size:   gdenticon.DefaultImageSize,
		mirror: gdenticon.DefaultMirror,
		style:  styleName(gdenticon.DefaultStyle),
		hash:   gdenticon.DefaultDigest,
		format: string(gdenticon.DefaultFormat),
		bg:     gdenticon.HexOf(gdenticon.DefaultBackground),
		output: output,
	}
}

func TestStyleName(t *testing.T) {
	if got := styleName(gdenticon.PatternClassic); got != "classic" {
		t.Fatalf("styleName(classic) = %q", got)
	}
	if got := styleName(gdenticon.PatternStyle(255)); got != "" {
		t.Fatalf("styleName(unknown) = %q", got)
	}
}

func TestExecRenderValidatesFlags(t *testing.T) {
	tests := []struct {
		name   string
		change func(*renderFlags)
	}{
		{"grid", func(f *renderFlags) { f.grid = 0 }},
		{"size", func(f *renderFlags) { f.size = 0 }},
		{"hash", func(f *renderFlags) { f.hash = "unknown" }},
		{"style", func(f *renderFlags) { f.style = "unknown" }},
		{"format", func(f *renderFlags) { f.format = "unknown" }},
		{"background", func(f *renderFlags) { f.bg = "invalid" }},
		{"foreground", func(f *renderFlags) { f.fg = "invalid" }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			flags := validRenderFlags(t.TempDir())
			test.change(&flags)
			_, err := execRender("first@example.com", &flags)
			if err == nil {
				t.Fatal("execRender returned nil error")
			}
			var usage usageError
			if !errors.As(err, &usage) {
				t.Fatalf("error type = %T, want usageError", err)
			}
		})
	}
}

func TestExecRenderCreatesAndReusesImage(t *testing.T) {
	flags := validRenderFlags(t.TempDir())
	flags.fg = "#2563eb"
	flags.bg = "#f5f8ff"

	first, err := execRender("first@example.com", &flags)
	if err != nil {
		t.Fatal(err)
	}
	if !first.Created {
		t.Fatal("first render was not reported as created")
	}
	if first.ID != "first@example.com" {
		t.Fatalf("result ID = %q", first.ID)
	}
	if filepath.Ext(first.Path) != ".png" {
		t.Fatalf("result path = %q", first.Path)
	}
	if info, err := os.Stat(first.Path); err != nil || info.Size() == 0 {
		t.Fatalf("generated file is invalid: info=%v err=%v", info, err)
	}

	second, err := execRender("first@example.com", &flags)
	if err != nil {
		t.Fatal(err)
	}
	if second.Created {
		t.Fatal("cached render was reported as created")
	}
	if second.Path != first.Path {
		t.Fatalf("cache paths differ: %q != %q", second.Path, first.Path)
	}
}
