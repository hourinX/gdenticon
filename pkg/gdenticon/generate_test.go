package gdenticon

import (
	"fmt"
	"testing"

	"github.com/hourinX/gdenticon/pkg/gdenticon/model"
)

func TestGenerateAndSave(t *testing.T) {
	img := model.Image{GridSize: 5, ImageSize: 256, Mirror: true, Hash: SHA256}

	for i := 0; i < 100; i++ {
		if _, err := Generate(fmt.Sprintf("alice@example%d.com", i), img); err != nil {
			t.Fatalf("Generate %d: %v", i, err)
		}
		_, err := Save(fmt.Sprintf("alice@example%d.com", i), img, "")
		if err != nil {
			t.Fatalf("Save %d: %v", i, err)
		}
	}
}
