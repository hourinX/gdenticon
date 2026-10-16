package gdenticon

import (
	"testing"

	"github.com/hourinX/gdenticon/pkg/gdenticon/model"
)

func TestRasterIsMirrored(t *testing.T) {
	for _, gridSize := range []int{1, 2, 3, 4, 5, 6, 7, 8, 11} {
		for _, size := range []int{1, 2, 5, 8, 16, 17, 32, 64, 100, 255, 256, 257, 420} {
			icon, err := Generate("alice@example.com", model.Image{GridSize: gridSize, Mirror: true})
			if err != nil {
				t.Fatalf("grid %d: Generate: %v", gridSize, err)
			}

			m := raster(icon, size)
			for y := range size {
				for x := range size / 2 {
					r1, g1, b1, a1 := m.At(x, y).RGBA()
					r2, g2, b2, a2 := m.At(size-1-x, y).RGBA()
					if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
						t.Fatalf("grid %d size %d: pixel (%d,%d) differs from its mirror", gridSize, size, x, y)
					}
				}
			}
		}
	}
}

func TestEdgeKeepsSmallImagesUsable(t *testing.T) {
	const gridSize = 5

	for _, size := range []int{8, 12, 16, 20, 24, 32, 48, 64, 128, 256, 420} {
		if content := edge(size, gridSize, gridSize) - edge(size, gridSize, 0); content*4 < size*3 {
			t.Errorf("size %d: pattern is %d px wide, under three quarters of the image", size, content)
		}
		if left, right := edge(size, gridSize, 0), size-edge(size, gridSize, gridSize); left != right {
			t.Errorf("size %d: margins are %d and %d, they must match", size, left, right)
		}
	}
}
