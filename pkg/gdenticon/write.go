// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

package gdenticon

import (
	"fmt"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

const defaultJPEGQuality = 90 // default JPEG setting: encoder quality

// Write encodes icon into w using the format selected when it was generated.
// It does not create or close w.
func Write(w io.Writer, icon Icon, opts WriteOptions) error {
	if opts.Size == 0 {
		opts.Size = DefaultImageSize
	}
	switch icon.Format {
	case FormatPNG:
		return writePNG(w, icon, opts)
	case FormatSVG:
		return writeSVG(w, icon, opts)
	case FormatJPEG:
		return writeJPEG(w, icon, opts)
	case FormatGIF:
		return writeGIF(w, icon, opts)
	case FormatBMP:
		return writeBMP(w, icon, opts)
	case FormatTIFF:
		return writeTIFF(w, icon, opts)
	default:
		return fmt.Errorf("unsupported image format %q", icon.Format)
	}
}

// writePNG encodes icon as PNG.
func writePNG(w io.Writer, icon Icon, opts WriteOptions) error {
	if err := validateRaster(icon, opts.Size, FormatPNG); err != nil {
		return err
	}
	return png.Encode(w, raster(icon, opts.Size))
}

// writeJPEG encodes icon as JPEG.
func writeJPEG(w io.Writer, icon Icon, opts WriteOptions) error {
	if err := validateRaster(icon, opts.Size, FormatJPEG); err != nil {
		return err
	}
	quality := opts.JPEGQuality
	if quality == 0 {
		quality = defaultJPEGQuality
	}
	if quality < 1 || quality > 100 {
		return fmt.Errorf("JPEG quality must be between 1 and 100, got %d", quality)
	}
	return jpeg.Encode(w, raster(icon, opts.Size), &jpeg.Options{Quality: quality})
}

// writeGIF encodes icon as GIF.
func writeGIF(w io.Writer, icon Icon, opts WriteOptions) error {
	if err := validateRaster(icon, opts.Size, FormatGIF); err != nil {
		return err
	}
	return gif.Encode(w, raster(icon, opts.Size), nil)
}

// writeBMP encodes icon as BMP.
func writeBMP(w io.Writer, icon Icon, opts WriteOptions) error {
	if err := validateRaster(icon, opts.Size, FormatBMP); err != nil {
		return err
	}
	return bmp.Encode(w, raster(icon, opts.Size))
}

// writeTIFF encodes icon as TIFF.
func writeTIFF(w io.Writer, icon Icon, opts WriteOptions) error {
	if err := validateRaster(icon, opts.Size, FormatTIFF); err != nil {
		return err
	}
	return tiff.Encode(w, raster(icon, opts.Size), nil)
}

// writeSVG encodes icon as SVG.
func writeSVG(w io.Writer, icon Icon, opts WriteOptions) error {
	if err := validateSVG(icon, opts.Size); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(w,
		`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`,
		opts.Size, opts.Size, opts.Size, opts.Size,
	); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, `<rect width="100%%" height="100%%" fill="%s"/>`, HexOf(icon.Background)); err != nil {
		return err
	}

	foreground := HexOf(icon.Foreground)
	for row := range icon.GridSize {
		for col := range icon.GridSize {
			if !icon.Filled[row*icon.GridSize+col] {
				continue
			}
			x1, x2 := edge(opts.Size, icon.GridSize, col), edge(opts.Size, icon.GridSize, col+1)
			y1, y2 := edge(opts.Size, icon.GridSize, row), edge(opts.Size, icon.GridSize, row+1)
			if _, err := fmt.Fprintf(w, `<rect x="%d" y="%d" width="%d" height="%d" fill="%s"/>`, x1, y1, x2-x1, y2-y1, foreground); err != nil {
				return err
			}
		}
	}
	_, err := io.WriteString(w, "</svg>\n")
	return err
}
