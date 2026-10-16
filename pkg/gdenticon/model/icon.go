// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

package model

import "image/color"

// Icon is a generated picture: a square field of cells, independent of any output format.
type Icon struct {
	GridSize   int
	Filled     []bool // row-major, true = foreground
	Background color.Color
	Foreground color.Color
}
