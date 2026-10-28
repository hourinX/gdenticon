// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/hourinX/gdenticon/v2/pkg/asset"
	"github.com/hourinX/gdenticon/v2/pkg/gdenticon"
	iconfile "github.com/hourinX/gdenticon/v2/pkg/gdenticon/file"
	"github.com/spf13/cobra"
)

var styleNames = map[string]gdenticon.PatternStyle{
	"classic":  gdenticon.PatternClassic,
	"balanced": gdenticon.PatternBalanced,
} // -style: flag name of each pattern style

type renderFlags struct {
	grid   int
	size   int
	mirror bool
	style  string
	hash   string
	format string
	bg     string
	fg     string
	output string
}

func (f *renderFlags) register(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.IntVar(&f.grid, "grid", gdenticon.DefaultGridSize, "side of the logical grid, in cells")
	flags.IntVar(&f.size, "size", gdenticon.DefaultImageSize, "side of the output image")
	flags.BoolVar(&f.mirror, "mirror", gdenticon.DefaultMirror, "mirror the left half of the grid onto the right")
	flags.StringVar(&f.style, "style", styleName(gdenticon.DefaultStyle), "pattern style: classic or balanced")
	flags.StringVar(&f.hash, "hash", gdenticon.DefaultDigest, "digest function")
	flags.StringVar(&f.format, "format", string(gdenticon.DefaultFormat), "image format: png, svg, jpeg, gif, bmp, or tiff")
	flags.StringVar(&f.bg, "bg", gdenticon.HexOf(gdenticon.DefaultBackground), "background colour, as #rrggbb or #rrggbbaa")
	flags.StringVar(&f.fg, "fg", "", "foreground colour, as #rrggbb or #rrggbbaa; empty derives it from the id")
	flags.StringVarP(&f.output, "output", "o", "", "output base directory; empty means the system temp directory")
}

// styleName returns the flag name that selects style.
func styleName(style gdenticon.PatternStyle) string {
	for name, s := range styleNames {
		if s == style {
			return name
		}
	}
	return ""
}

// execRender generates the identicon for id, saves it, and prints the path.
func execRender(cmd *cobra.Command, id string, f *renderFlags) error {
	if f.grid < 1 {
		return usageError{fmt.Errorf("-grid must be positive, got %d", f.grid)}
	}
	if f.size < 1 {
		return usageError{fmt.Errorf("-size must be positive, got %d", f.size)}
	}

	digest, ok := gdenticon.Digest(f.hash)
	if !ok {
		return usageError{fmt.Errorf("unknown -hash %q", f.hash)}
	}

	style, ok := styleNames[f.style]
	if !ok {
		return usageError{fmt.Errorf("unknown -style %q", f.style)}
	}
	format := gdenticon.Format(f.format)
	if format != gdenticon.FormatPNG && format != gdenticon.FormatSVG &&
		format != gdenticon.FormatJPEG && format != gdenticon.FormatGIF &&
		format != gdenticon.FormatBMP && format != gdenticon.FormatTIFF {
		return usageError{fmt.Errorf("unknown -format %q", f.format)}
	}

	background, err := gdenticon.ParseColor(f.bg)
	if err != nil {
		return usageError{fmt.Errorf("-bg: %w", err)}
	}

	opts := gdenticon.GenerateOptions{
		GridSize:   f.grid,
		Mirror:     f.mirror,
		Style:      style,
		Digest:     digest,
		Background: background,
		Format:     format,
	}

	if f.fg != "" {
		c, err := gdenticon.ParseColor(f.fg)
		if err != nil {
			return usageError{fmt.Errorf("-fg: %w", err)}
		}
		opts.Foreground = gdenticon.FixedColor(c)
	}

	icon, err := gdenticon.Generate(id, opts)
	if err != nil {
		return err
	}

	base := f.output
	if base == "" {
		base = os.TempDir()
	}
	dir := filepath.Join(base, asset.FolderName)
	path, err := iconfile.PutCache(dir, icon.Format.Extension(), 0o644, func(w io.Writer) error {
		return gdenticon.Write(w, icon, gdenticon.WriteOptions{Size: f.size})
	})
	if err != nil {
		return err
	}

	fmt.Fprintln(cmd.OutOrStdout(), path)
	return nil
}
