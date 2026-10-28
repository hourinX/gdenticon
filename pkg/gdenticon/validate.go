// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

package gdenticon

import "fmt"

// validateGrid checks the constraints on a logical icon grid.
func validateGrid(size int) error {
	if size < 1 {
		return fmt.Errorf("grid size must be positive, got %d", size)
	}
	if size > MaxGridSize {
		return fmt.Errorf("grid size must be at most %d, got %d", MaxGridSize, size)
	}
	return nil
}

// validateRaster checks the constraints shared by raster encoders.
func validateRaster(icon Icon, size int, format Format) error {
	if err := validateIcon(icon); err != nil {
		return err
	}
	if size < 1 {
		return fmt.Errorf("%s size must be positive, got %d", format, size)
	}
	if size > MaxImageSize {
		return fmt.Errorf("%s size must be at most %d, got %d", format, MaxImageSize, size)
	}
	if size < icon.GridSize {
		return fmt.Errorf("%s size must be at least the grid size %d, got %d", format, icon.GridSize, size)
	}
	return nil
}

// validateSVG checks the constraints imposed by writing an icon as SVG.
func validateSVG(icon Icon, size int) error {
	if err := validateIcon(icon); err != nil {
		return err
	}
	if size < 1 {
		return fmt.Errorf("SVG size must be positive, got %d", size)
	}
	return nil
}

// validateIcon checks the invariants expected by every encoder.
func validateIcon(icon Icon) error {
	if err := validateGrid(icon.GridSize); err != nil {
		return err
	}
	if len(icon.Filled) != icon.GridSize*icon.GridSize {
		return fmt.Errorf("icon has %d cells, want %d", len(icon.Filled), icon.GridSize*icon.GridSize)
	}
	if icon.Background == nil || icon.Foreground == nil {
		return fmt.Errorf("icon colours must not be nil")
	}
	return nil
}
