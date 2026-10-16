// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

package gdenticon

import (
	"github.com/hourinX/gdenticon/pkg/gdenticon/model"
)

// Generate derives the icon for id from img. The same id and img always produce
// the same Icon.
func Generate(id string, img model.Image) (model.Icon, error) {
	if err := validateGrid(img); err != nil {
		return model.Icon{}, err
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
		Filled:     fill(digest, img),
		Background: background,
		Foreground: foreground(digest),
	}, nil
}

// fill draws the cell grid for img from digest, by whichever style img asks for.
func fill(digest []byte, img model.Image) []bool {
	style := img.Style
	if style == model.PatternDefault {
		style = model.DefaultStyle
	}

	e := newEntropy(digest)
	if style == model.PatternBalanced {
		return balancedGrid(e, img.GridSize, img.Mirror)
	}
	return grid(e, img.GridSize, img.Mirror)
}
