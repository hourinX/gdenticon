// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

package gdenticon

import (
	"fmt"

	"github.com/hourinX/gdenticon/pkg/gdenticon/model"
)

// Generate derives the icon for id from img. The same id and img always produce
// the same Icon.
func Generate(id string, img model.Image) (model.Icon, error) {
	if img.GridSize < 1 {
		return model.Icon{}, fmt.Errorf("grid size must be positive, got %d", img.GridSize)
	}

	hash := img.Hash
	if hash == nil {
		hash = MD5
	}
	digest := hash(id)

	background := img.Background
	if background == nil {
		background = model.DefaultBackground
	}

	foreground := img.Foreground
	if foreground == nil {
		foreground = HueColor
	}

	return model.Icon{
		GridSize:   img.GridSize,
		Filled:     grid(newEntropy(digest), img.GridSize, img.Mirror),
		Background: background,
		Foreground: foreground(digest),
	}, nil
}
