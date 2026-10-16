// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

package gdenticon

// grid draws a size×size field of cells from e. With mirror set, only the left
// half is drawn and each cell is reflected onto its opposite column; an odd size
// therefore spends bits on the centre column once, not twice.
func grid(e *entropy, size int, mirror bool) []bool {
	cols := size
	if mirror {
		cols = (size + 1) / 2
	}

	cells := make([]bool, size*size)
	for row := range size {
		for col := range cols {
			if !e.Bool() {
				continue
			}
			cells[row*size+col] = true
			if mirror {
				cells[row*size+size-1-col] = true
			}
		}
	}
	return cells
}
