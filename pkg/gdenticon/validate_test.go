package gdenticon

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/hourinX/gdenticon/pkg/gdenticon/model"
)

func TestGenerateRejectsOutOfRangeGridSize(t *testing.T) {
	for _, gridSize := range []int{0, -1, model.MaxGridSize + 1, 1 << 20} {
		if _, err := Generate("alice@example.com", model.Image{GridSize: gridSize}); err == nil {
			t.Errorf("grid %d: Generate returned no error", gridSize)
		}
	}

	for _, gridSize := range []int{1, model.DefaultGridSize, model.MaxGridSize} {
		if _, err := Generate("alice@example.com", model.Image{GridSize: gridSize}); err != nil {
			t.Errorf("grid %d: Generate: %v", gridSize, err)
		}
	}
}

func TestWritePNGRejectsOutOfRangeImageSize(t *testing.T) {
	for _, imageSize := range []int{0, -1, model.DefaultGridSize - 1, model.MaxImageSize + 1} {
		img := model.Image{GridSize: model.DefaultGridSize, ImageSize: imageSize}
		if err := WritePNG(io.Discard, "alice@example.com", img); err == nil {
			t.Errorf("image size %d: WritePNG returned no error", imageSize)
		}
	}

	for _, imageSize := range []int{model.DefaultGridSize, model.DefaultImageSize, model.MaxImageSize} {
		img := model.Image{GridSize: model.DefaultGridSize, ImageSize: imageSize}
		if err := WritePNG(io.Discard, "alice@example.com", img); err != nil {
			t.Errorf("image size %d: WritePNG: %v", imageSize, err)
		}
	}
}

func TestEveryRasterEntryRejectsAnUndersizedImage(t *testing.T) {
	// A grid too big for the image used to be accepted and drawn as nothing at
	// all, so every entry point that takes both has to say so rather than leave
	// the caller with a blank file.
	img := model.Image{GridSize: model.DefaultGridSize, ImageSize: model.DefaultGridSize - 1}

	if err := WritePNG(io.Discard, "alice@example.com", img); err == nil {
		t.Error("WritePNG returned no error")
	}

	path, err := Save("alice@example.com", img, t.TempDir())
	if err == nil {
		t.Errorf("Save wrote %s instead of returning an error", path)
	}
}

func TestUndersizedImageDoesNotBlameTheGrid(t *testing.T) {
	// Both sizes are wrong here, so the message has to name the grid: reporting
	// the image instead would send the caller after the wrong field.
	img := model.Image{GridSize: model.MaxGridSize + 1, ImageSize: 1}

	err := WritePNG(io.Discard, "alice@example.com", img)
	if err == nil {
		t.Fatal("WritePNG returned no error")
	}
	if want := fmt.Sprintf("grid size must be at most %d", model.MaxGridSize); !strings.Contains(err.Error(), want) {
		t.Errorf("error is %q, want it to contain %q", err, want)
	}
}
