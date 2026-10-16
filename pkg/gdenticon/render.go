// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

package gdenticon

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"

	"github.com/hourinX/gdenticon/pkg/gdenticon/model"
)

// WritePNG writes the PNG icon for id to w.
func WritePNG(w io.Writer, id string, img model.Image) error {
	if img.ImageSize < 1 {
		return fmt.Errorf("image size must be positive, got %d", img.ImageSize)
	}

	icon, err := Generate(id, img)
	if err != nil {
		return err
	}
	return png.Encode(w, raster(icon, img.ImageSize))
}

// raster paints icon into a size×size image.
func raster(icon model.Icon, size int) *image.NRGBA {
	m := image.NewNRGBA(image.Rect(0, 0, size, size))
	draw.Draw(m, m.Bounds(), image.NewUniform(icon.Background), image.Point{}, draw.Src)

	for row := range icon.GridSize {
		for col := range icon.GridSize {
			if icon.Filled[row*icon.GridSize+col] {
				draw.Draw(m, cell(size, icon.GridSize, col, row), image.NewUniform(icon.Foreground), image.Point{}, draw.Src)
			}
		}
	}
	return m
}

// cell returns the pixel rectangle of the cell at col, row.
func cell(size, gridSize, col, row int) image.Rectangle {
	return image.Rect(
		edge(size, gridSize, col), edge(size, gridSize, row),
		edge(size, gridSize, col+1), edge(size, gridSize, row+1),
	)
}

// edge returns the pixel coordinate of grid line i. The grid is inset by half a
// cell so that it floats in the image rather than touching its edges, and lines
// past the middle are derived from their opposite by subtraction rather than by
// dividing again, so that a cell and the cell it mirrors onto occupy matching
// pixel ranges even when size is not a multiple of gridSize.
func edge(size, gridSize, i int) int {
	if i > gridSize/2 {
		return size - edge(size, gridSize, gridSize-i)
	}
	cell := size / (gridSize + 1)
	return (size-gridSize*cell)/2 + i*cell
}
