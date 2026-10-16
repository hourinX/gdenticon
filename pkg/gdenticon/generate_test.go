package gdenticon

import (
	"testing"

	"github.com/hourinX/gdenticon/pkg/gdenticon/model"
)

func TestGenerateAndSave(t *testing.T) {
	img := model.Image{GridSize: 5, ImageSize: 256, Mirror: true, Hash: SHA256}

	if _, err := Generate("alice@example.com", img); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	path, err := Save("alice@example.com", img, "")
	if err != nil {
		t.Fatalf("Save: %v", err)
	}

	t.Log(path)
}
