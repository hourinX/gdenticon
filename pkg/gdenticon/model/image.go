// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

package model

import "image/color"

// ColorFunc derives a colour from the digest of an id.
type ColorFunc func(hash []byte) color.Color

// Extension enumerates the selectable image types.
type Extension string

const (
	SVG  Extension = "svg"  // vector, needs no ImageSize
	PNG  Extension = "png"  // raster, needs ImageSize
	JPEG Extension = "jpeg" // raster, needs ImageSize, no alpha
	GIF  Extension = "gif"  // raster, needs ImageSize, 256 colours, 1-bit alpha
	BMP  Extension = "bmp"  // raster, needs ImageSize, 32-bit alpha
	TIFF Extension = "tiff" // raster, needs ImageSize, alpha
	WEBP Extension = "webp" // raster, needs ImageSize, alpha
)

// Image holds the values that shape a generated image.
type Image struct {
	GridSize   int
	ImageSize  int
	Mirror     bool
	Hash       func(string) []byte
	Foreground ColorFunc
	Background color.Color
	Format     Extension
}

// Default image setting
const (
	DefaultGridSize  = 5    // logical grid, in cells
	DefaultImageSize = 256  // output side length, in pixels
	DefaultMirror    = true // mirrored left and right
	DefaultFormat    = SVG  // output file extension
)

var DefaultBackground color.Color = color.RGBA{0xff, 0xff, 0xff, 0xff} // background, white
