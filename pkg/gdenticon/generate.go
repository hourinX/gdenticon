// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

package gdenticon

import "fmt"

// Generate derives an icon from id and opts, resolving the default format when
// none is selected. The same id and options always produce the same Icon.
func Generate(id string, opts GenerateOptions) (Icon, error) {
	if err := validateGrid(opts.GridSize); err != nil {
		return Icon{}, err
	}

	digestFunc := opts.Digest
	if digestFunc == nil {
		digestFunc = MD5
	}
	digest := digestFunc(id)

	background := opts.Background
	if background == nil {
		background = DefaultBackground
	}

	foreground := opts.Foreground
	if foreground == nil {
		foreground = HueColor
	}
	format := opts.Format
	if format == FormatDefault {
		format = DefaultFormat
	}
	if !supportedFormat(format) {
		return Icon{}, fmt.Errorf("unsupported image format %q", format)
	}

	return Icon{
		GridSize:   opts.GridSize,
		Filled:     fill(digest, opts),
		Background: background,
		Foreground: foreground(digest),
		Format:     format,
	}, nil
}

// supportedFormat reports whether Write has an encoder for format.
func supportedFormat(format Format) bool {
	switch format {
	case FormatPNG, FormatSVG, FormatJPEG, FormatGIF, FormatBMP, FormatTIFF:
		return true
	default:
		return false
	}
}

// fill draws the cell grid for opts from digest, using the requested style.
func fill(digest []byte, opts GenerateOptions) []bool {
	style := opts.Style
	if style == PatternDefault {
		style = DefaultStyle
	}

	e := newEntropy(digest)
	if style == PatternBalanced {
		return balancedGrid(e, opts.GridSize, opts.Mirror)
	}
	return grid(e, opts.GridSize, opts.Mirror)
}
