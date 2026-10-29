// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

package gdenticon

import (
	"image/color"
	"strings"
	"testing"
)

func TestValidateGridBounds(t *testing.T) {
	for _, test := range []struct {
		name string
		size int
		want string
	}{
		{name: "zero", size: 0, want: "positive"},
		{name: "too large", size: MaxGridSize + 1, want: "at most"},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := validateGrid(test.size)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("validateGrid(%d) = %v, want %q", test.size, err, test.want)
			}
		})
	}
	if err := validateGrid(MaxGridSize); err != nil {
		t.Fatalf("validateGrid(MaxGridSize): %v", err)
	}
}

func TestValidateRasterBranches(t *testing.T) {
	icon := testIcon(FormatPNG)
	for _, test := range []struct {
		name string
		icon Icon
		size int
		want string
	}{
		{name: "invalid icon", icon: Icon{}, size: 10, want: "grid size"},
		{name: "zero size", icon: icon, size: 0, want: "positive"},
		{name: "too large", icon: icon, size: MaxImageSize + 1, want: "at most"},
		{name: "smaller than grid", icon: icon, size: icon.GridSize - 1, want: "at least"},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := validateRaster(test.icon, test.size, FormatPNG)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("validateRaster() = %v, want %q", err, test.want)
			}
		})
	}
	if err := validateRaster(icon, icon.GridSize, FormatPNG); err != nil {
		t.Fatalf("valid raster: %v", err)
	}
}

func TestValidateSVGAndIconBranches(t *testing.T) {
	icon := testIcon(FormatSVG)
	if err := validateSVG(Icon{}, 10); err == nil {
		t.Fatal("validateSVG accepted an invalid icon")
	}
	if err := validateSVG(icon, 0); err == nil || !strings.Contains(err.Error(), "positive") {
		t.Fatalf("validateSVG zero size = %v", err)
	}
	if err := validateSVG(icon, 1); err != nil {
		t.Fatalf("valid SVG: %v", err)
	}

	tests := []Icon{
		{GridSize: 2, Filled: []bool{true}, Background: color.White, Foreground: color.Black},
		{GridSize: 1, Filled: []bool{true}, Foreground: color.Black},
		{GridSize: 1, Filled: []bool{true}, Background: color.White},
	}
	for i, invalid := range tests {
		if err := validateIcon(invalid); err == nil {
			t.Errorf("validateIcon case %d returned nil", i)
		}
	}
}
