// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

package gdenticon

const (
	maxAttempts         = 32  // balancedGrid: candidates drawn before the least bad is kept
	minFragmentedBlocks = 4   // degenerate: block count below which "fragmented" means nothing
	minBandSize         = 3   // degenerate: smallest grid the density band applies to
	minDensityPermille  = 320 // degenerate: least acceptable fill, in thousandths of the grid
	maxDensityPermille  = 680 // degenerate: most acceptable fill, in thousandths of the grid
)

// balancedGrid draws a grid, redrawing while the candidate falls outside the
// accepted band. If no candidate fits, the least bad one is kept rather than the
// last one drawn, so the result never depends on where the redraw happened to
// stop. The pickiest grid accepts under two thirds of what it can draw, so all
// maxAttempts draws failing has probability below 1e-14.
func balancedGrid(e *entropy, size int, mirror bool) []bool {
	return balancedGridWithin(e, size, mirror, maxAttempts)
}

// balancedGridWithin is balancedGrid with an explicit attempt budget, so a test
// can drive the fallback that production reaches once in 1e14 ids.
func balancedGridWithin(e *entropy, size int, mirror bool, attempts int) []bool {
	var best []bool
	bestPenalty := 0

	for attempt := range attempts {
		cells := grid(e, size, mirror)
		if !degenerate(cells, size) {
			return cells
		}
		if penalty := visualPenalty(cells, size); attempt == 0 || penalty < bestPenalty {
			best, bestPenalty = cells, penalty
		}
	}
	return best
}

// degenerate reports whether cells falls outside the band of patterns that read
// as an identicon: blank, solid, too sparse, too dense, or split into blocks
// averaging under two cells each.
func degenerate(cells []bool, size int) bool {
	filled := filledCells(cells)
	if filled == 0 || filled == len(cells) {
		return true
	}

	if size >= minBandSize {
		if permille := filled * 1000 / len(cells); permille < minDensityPermille || permille > maxDensityPermille {
			return true
		}
	}

	blocks := blocks(cells, size)
	return blocks >= minFragmentedBlocks && blocks*2 > filled
}

// visualPenalty scores how far cells falls outside the accepted band, in
// thousandths: zero is inside it. Both terms are on the same scale, so neither
// wins by unit choice alone.
func visualPenalty(cells []bool, size int) int {
	filled := filledCells(cells)
	if filled == 0 || filled == len(cells) {
		return 2 * 1000
	}

	penalty := 0
	if size >= minBandSize {
		permille := filled * 1000 / len(cells)
		switch {
		case permille < minDensityPermille:
			penalty += minDensityPermille - permille
		case permille > maxDensityPermille:
			penalty += permille - maxDensityPermille
		}
	}

	if blocks := blocks(cells, size); blocks >= minFragmentedBlocks && blocks*2 > filled {
		penalty += (blocks*2 - filled) * 1000 / filled
	}
	return penalty
}

// filledCells counts the filled cells of a grid.
func filledCells(cells []bool) int {
	filled := 0
	for _, cell := range cells {
		if cell {
			filled++
		}
	}
	return filled
}

// blocks counts the 4-connected blocks of filled cells.
func blocks(cells []bool, size int) int {
	seen := make([]bool, len(cells))
	count := 0

	for start := range cells {
		if !cells[start] || seen[start] {
			continue
		}
		count++

		seen[start] = true
		stack := []int{start}
		for len(stack) > 0 {
			at := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			row, col := at/size, at%size
			for _, step := range [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
				r, c := row+step[0], col+step[1]
				if r < 0 || r >= size || c < 0 || c >= size {
					continue
				}
				next := r*size + c
				if cells[next] && !seen[next] {
					seen[next] = true
					stack = append(stack, next)
				}
			}
		}
	}
	return count
}
