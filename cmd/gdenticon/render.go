// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

package main

import (
	"fmt"

	"github.com/hourinX/gdenticon/pkg/gdenticon"
	"github.com/hourinX/gdenticon/pkg/gdenticon/model"
	"github.com/spf13/cobra"
)

var styleNames = map[string]model.PatternStyle{
	"classic":  model.PatternClassic,
	"balanced": model.PatternBalanced,
} // -style: flag name of each pattern style

type renderFlags struct {
	grid   int
	size   int
	mirror bool
	style  string
	hash   string
	bg     string
	fg     string
	output string
}

func (f *renderFlags) register(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.IntVar(&f.grid, "grid", model.DefaultGridSize, "side of the logical grid, in cells")
	flags.IntVar(&f.size, "size", model.DefaultImageSize, "side of the output image, in pixels")
	flags.BoolVar(&f.mirror, "mirror", model.DefaultMirror, "mirror the left half of the grid onto the right")
	flags.StringVar(&f.style, "style", styleName(model.DefaultStyle), "pattern style: classic or balanced")
	flags.StringVar(&f.hash, "hash", gdenticon.DefaultDigest, "digest function")
	flags.StringVar(&f.bg, "bg", gdenticon.HexOf(model.DefaultBackground), "background colour, as #rrggbb or #rrggbbaa")
	flags.StringVar(&f.fg, "fg", "", "foreground colour, as #rrggbb or #rrggbbaa; empty derives it from the id")
	flags.StringVarP(&f.output, "output", "o", "", "output base directory; empty means the system temp directory")
}

// styleName returns the flag name that selects style.
func styleName(style model.PatternStyle) string {
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

	background, err := gdenticon.ParseColor(f.bg)
	if err != nil {
		return usageError{fmt.Errorf("-bg: %w", err)}
	}

	img := model.Image{
		GridSize:   f.grid,
		ImageSize:  f.size,
		Mirror:     f.mirror,
		Style:      style,
		Hash:       digest,
		Background: background,
	}

	if f.fg != "" {
		c, err := gdenticon.ParseColor(f.fg)
		if err != nil {
			return usageError{fmt.Errorf("-fg: %w", err)}
		}
		img.Foreground = gdenticon.FixedColor(c)
	}

	path, err := gdenticon.Save(id, img, f.output)
	if err != nil {
		return err
	}

	fmt.Fprintln(cmd.OutOrStdout(), path)
	return nil
}
