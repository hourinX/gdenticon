package gdenticon

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hourinX/gdenticon/pkg/asset"
)

func TestGeneratePNGFiles(t *testing.T) {
	outputDir := filepath.Join(os.TempDir(), asset.FolderName)
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Logf("generated PNG files: %s", outputDir)

	for i := 1; i <= 100; i++ {
		id := fmt.Sprintf("user-%03d:name-%03d:user%03d@example.com", i, i, i)
		icon, err := Generate(id, GenerateOptions{
			GridSize: 5,
			Mirror:   true,
			Digest:   SHA256,
			Format:   FormatPNG,
		})
		if err != nil {
			t.Fatalf("Generate %d: %v", i, err)
		}

		path := filepath.Join(outputDir, fmt.Sprintf("user-%03d.png", i))
		file, err := os.Create(path)
		if err != nil {
			t.Fatalf("Create %s: %v", path, err)
		}
		if err := Write(file, icon, WriteOptions{Size: DefaultImageSize}); err != nil {
			_ = file.Close()
			t.Fatalf("Write %s: %v", path, err)
		}
		if err := file.Close(); err != nil {
			t.Fatalf("Close %s: %v", path, err)
		}
	}
}

func TestGenerateResolvesFormat(t *testing.T) {
	for _, test := range []struct {
		name string
		give Format
		want Format
	}{
		{name: "default", give: FormatDefault, want: DefaultFormat},
		{name: "PNG", give: FormatPNG, want: FormatPNG},
		{name: "SVG", give: FormatSVG, want: FormatSVG},
	} {
		t.Run(test.name, func(t *testing.T) {
			icon, err := Generate("alice@example.com", GenerateOptions{GridSize: 5, Format: test.give})
			if err != nil {
				t.Fatal(err)
			}
			if icon.Format != test.want {
				t.Fatalf("format is %q, want %q", icon.Format, test.want)
			}
		})
	}
}
