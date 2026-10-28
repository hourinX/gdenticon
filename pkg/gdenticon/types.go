// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

package gdenticon

import "image/color"

// DigestFunc derives a digest from an identifier.
type DigestFunc func(string) []byte

// ColorFunc derives a colour from the digest of an identifier.
type ColorFunc func([]byte) color.Color

// PatternStyle selects how the cell grid is drawn from the digest.
type PatternStyle int

const (
	PatternDefault  PatternStyle = iota // generation setting: whichever style DefaultStyle names
	PatternClassic                      // generation setting: every digest bit lands as drawn
	PatternBalanced                     // generation setting: redraws patterns outside the accepted band
)

// Format selects the encoded image type written for an icon.
type Format string

const (
	FormatDefault Format = ""     // generation setting: whichever format DefaultFormat names
	FormatPNG     Format = "png"  // generation setting: PNG image
	FormatSVG     Format = "svg"  // generation setting: SVG image
	FormatJPEG    Format = "jpeg" // generation setting: JPEG image
	FormatGIF     Format = "gif"  // generation setting: GIF image
	FormatBMP     Format = "bmp"  // generation setting: BMP image
	FormatTIFF    Format = "tiff" // generation setting: TIFF image
)

// GenerateOptions holds the values that generate an icon and select its format.
type GenerateOptions struct {
	GridSize   int
	Mirror     bool
	Style      PatternStyle
	Digest     DigestFunc
	Foreground ColorFunc
	Background color.Color
	Format     Format
}

const (
	DefaultGridSize  = 5               // default generation setting: logical grid, in cells
	DefaultMirror    = true            // default generation setting: mirrored left and right
	DefaultStyle     = PatternBalanced // default generation setting: pattern style
	DefaultFormat    = FormatPNG       // default generation setting: encoded image type
	MaxGridSize      = 64              // generation limit: largest logical grid side, in cells
	DefaultImageSize = 256             // default write setting: output side length
	MaxImageSize     = 4096            // write limit: largest raster output side, in pixels
)

var DefaultBackground color.Color = color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff} // default generation setting: background, white

// Icon is a generated square field of cells and its requested output format.
type Icon struct {
	GridSize   int
	Filled     []bool // row-major, true means foreground
	Background color.Color
	Foreground color.Color
	Format     Format
}

// WriteOptions holds settings shared by image writers.
type WriteOptions struct {
	Size        int // output side length; zero uses DefaultImageSize
	JPEGQuality int // JPEG quality from 1 to 100; zero uses 90
}

// Extension returns the conventional file extension for f.
func (f Format) Extension() string {
	return "." + string(f)
}
