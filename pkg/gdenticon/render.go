// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

package gdenticon

import (
	"image"
	"image/draw"
)

// raster paints icon into a size×size image.
func raster(icon Icon, size int) *image.NRGBA {
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

// edge returns the pixel coordinate of grid line i. The cell size is rounded
// rather than truncated, because truncating throws away up to a whole cell: at
// 16 pixels that leaves the pattern filling barely half the image. Lines past
// the middle are derived from their opposite by subtraction rather than by
// dividing again, so a cell and the cell it mirrors onto cover matching pixels.
func edge(size, gridSize, i int) int {
	if i > gridSize/2 {
		return size - edge(size, gridSize, gridSize-i)
	}
	cell := (2*size + gridSize + 1) / (2 * (gridSize + 1))
	if room := size / gridSize; cell > room {
		cell = room
	}
	return (size-gridSize*cell)/2 + i*cell
}
