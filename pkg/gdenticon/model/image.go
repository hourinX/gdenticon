// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

package model

import "image/color"

// ColorFunc derives a colour from the digest of an id.
type ColorFunc func(hash []byte) color.Color

// Image holds the values that shape a generated image.
type Image struct {
	GridSize   int
	ImageSize  int
	Mirror     bool
	Hash       func(string) []byte
	Foreground ColorFunc
	Background color.Color
}

const (
	DefaultGridSize  = 5    // default image setting: logical grid, in cells
	DefaultImageSize = 256  // default image setting: output side length, in pixels
	DefaultMirror    = true // default image setting: mirrored left and right
)

var DefaultBackground color.Color = color.RGBA{0xff, 0xff, 0xff, 0xff} // default image setting: background, white
