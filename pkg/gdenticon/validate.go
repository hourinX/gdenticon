// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

package gdenticon

import (
	"fmt"

	"github.com/hourinX/gdenticon/pkg/gdenticon/model"
)

// validateGrid checks the constraints on the grid an Image asks for. Any entry
// point that draws from an Image needs them, whatever it renders to.
func validateGrid(img model.Image) error {
	if img.GridSize < 1 {
		return fmt.Errorf("grid size must be positive, got %d", img.GridSize)
	}
	if img.GridSize > model.MaxGridSize {
		return fmt.Errorf("grid size must be at most %d, got %d", model.MaxGridSize, img.GridSize)
	}
	return nil
}

// validateImage checks the constraints that only rasterising imposes, which is
// why Generate does not apply them: they say nothing about a vector icon, and
// everything about a grid whose cells have to land on whole pixels.
func validateImage(img model.Image) error {
	if img.ImageSize < 1 {
		return fmt.Errorf("image size must be positive, got %d", img.ImageSize)
	}
	if img.ImageSize > model.MaxImageSize {
		return fmt.Errorf("image size must be at most %d, got %d", model.MaxImageSize, img.ImageSize)
	}
	if img.ImageSize < img.GridSize {
		return fmt.Errorf("image size must be at least the grid size %d, got %d", img.GridSize, img.ImageSize)
	}
	return nil
}

// validate checks everything an Image must satisfy to be drawn on pixels. It
// takes the grid first, so an image that fails both is reported for the reason
// that is actually the caller's mistake rather than the consequence of it.
func validate(img model.Image) error {
	if err := validateGrid(img); err != nil {
		return err
	}
	return validateImage(img)
}
