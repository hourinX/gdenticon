package gdenticon

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hourinX/gdenticon/pkg/gdenticon/model"
)

// cellsFrom turns rows of '#' and '.' into a cell grid.
func cellsFrom(rows ...string) []bool {
	cells := make([]bool, 0, len(rows)*len(rows[0]))
	for _, row := range rows {
		for _, cell := range row {
			cells = append(cells, cell == '#')
		}
	}
	return cells
}

// seedOf packs 15-bit candidates into a digest the entropy reader hands back one
// grid at a time.
func seedOf(candidates ...uint16) []byte {
	seed := make([]byte, (len(candidates)*15+7)/8)
	for i, bits := range candidates {
		for j := range 15 {
			if bits>>(14-j)&1 == 1 {
				bit := i*15 + j
				seed[bit/8] |= 1 << (7 - bit%8)
			}
		}
	}
	return seed
}

func TestDegenerateFlagsKnownPatterns(t *testing.T) {
	for _, tc := range []struct {
		name string
		rows []string
		want bool
	}{
		{"blank", []string{".....", ".....", ".....", ".....", "....."}, true},
		{"solid", []string{"#####", "#####", "#####", "#####", "#####"}, true},
		{"over the band", []string{"#####", "#####", "#####", "#####", "###.."}, true},
		{"under the band", []string{".....", "..#..", "..#..", ".....", "....."}, true},
		{"nine cells in one block", []string{".....", ".###.", ".###.", ".###.", "....."}, false},
		{"four lone cells", []string{"#...#", ".....", ".....", ".....", "#...#"}, true},
		{"checkerboard", []string{"#.#.#", ".#.#.", "#.#.#", ".#.#.", "#.#.#"}, true},
		{"compact blob inside the band", []string{"###..", "####.", ".####", ".###.", "..##."}, false},
	} {
		if got := degenerate(cellsFrom(tc.rows...), len(tc.rows)); got != tc.want {
			t.Errorf("%s: degenerate = %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestVisualPenaltyRanksPatterns(t *testing.T) {
	// The steps walk out from the middle of the band in both directions, so the
	// penalty must never fall. The last two are the two ways to leave the band
	// sideways: shredded at an acceptable density, then solid.
	rising := [][]string{
		{"#####", "###..", ".....", ".....", "....."}, // 32%: at the floor
		{"#####", "#####", "##...", ".....", "....."},
		{"#####", "#####", "#####", "##...", "....."}, // 68%: at the ceiling
		{"#####", "#####", "#####", "#####", "....."},
		{"#####", ".....", ".....", ".....", "....."},
		{"###..", ".....", ".....", ".....", "....."},
		{"##...", ".....", ".....", ".....", "....."},
		{"#....", ".....", ".....", ".....", "....."},
		{"#.#.#", ".#.#.", "#.#.#", ".#.#.", "#.#.#"},
		{"#####", "#####", "#####", "#####", "#####"},
	}

	last := -1
	for _, rows := range rising {
		cells := cellsFrom(rows...)
		penalty := visualPenalty(cells, len(rows))
		if penalty < last {
			t.Fatalf("penalty fell from %d to %d going from %v", last, penalty, rows)
		}
		last = penalty
	}
}

func TestBalancedGridFallsBackToTheLeastBadCandidate(t *testing.T) {
	// Candidates 1, 4 and 5 onward are blank; candidate 2 is solid; candidate 3
	// is 96% full and so carries the smallest penalty of the eight. Keeping
	// whichever candidate was drawn last would return a blank grid, and keeping
	// the first drawn would return the solid one.
	seed := seedOf(0, 0x7fff, 0x6fff, 0, 0, 0, 0, 0)

	cells := balancedGridWithin(newEntropy(seed), 5, true, 8)
	if filled := filledCells(cells); filled != 24 {
		t.Fatalf("kept a grid with %d filled cells, want the least bad candidate's 24", filled)
	}
}

func TestBalancedGridStaysInsideTheBand(t *testing.T) {
	const ids = 20000

	outside := 0
	for _, gridSize := range []int{3, 4, 5, 6, 7, 8} {
		for i := range ids {
			img := model.Image{GridSize: gridSize, Mirror: true, Style: model.PatternBalanced}
			icon, err := Generate(fmt.Sprintf("user%d", i), img)
			if err != nil {
				t.Fatalf("grid %d: Generate: %v", gridSize, err)
			}

			// Measured here rather than through degenerate, so this test states
			// the intended band and not the rule that produced the grid.
			filled := filledCells(icon.Filled)
			if filled == 0 || filled == len(icon.Filled) {
				t.Fatalf("grid %d id %d: blank or solid grid escaped the filter", gridSize, i)
			}
			if permille := filled * 1000 / len(icon.Filled); permille < minDensityPermille || permille > maxDensityPermille {
				outside++
			}
		}
	}

	if outside != 0 {
		t.Errorf("%d of %d grids fell outside the band", outside, 6*ids)
	}
}

func TestUnsetStyleFollowsDefaultStyle(t *testing.T) {
	differed := false
	for i := range 100 {
		id := fmt.Sprintf("user%d", i)

		unset, err := Generate(id, model.Image{GridSize: 5, Mirror: true})
		if err != nil {
			t.Fatalf("Generate: %v", err)
		}
		tagged, err := Generate(id, model.Image{GridSize: 5, Mirror: true, Style: model.DefaultStyle})
		if err != nil {
			t.Fatalf("Generate: %v", err)
		}
		classic, err := Generate(id, model.Image{GridSize: 5, Mirror: true, Style: model.PatternClassic})
		if err != nil {
			t.Fatalf("Generate: %v", err)
		}

		if !reflect.DeepEqual(unset.Filled, tagged.Filled) {
			t.Fatalf("%s: a zero Style did not follow DefaultStyle", id)
		}
		if !reflect.DeepEqual(unset.Filled, classic.Filled) {
			differed = true
		}
	}

	if !differed {
		t.Fatal("no id drew differently from PatternClassic, so this test proves nothing")
	}
}

func TestBalancedGridStaysMirrored(t *testing.T) {
	for _, gridSize := range []int{1, 2, 3, 4, 5, 7, 8} {
		for _, id := range []string{"alice@example.com", "user1", "user2", "user3"} {
			icon, err := Generate(id, model.Image{GridSize: gridSize, Mirror: true, Style: model.PatternBalanced})
			if err != nil {
				t.Fatalf("grid %d: Generate: %v", gridSize, err)
			}
			if len(icon.Filled) != gridSize*gridSize {
				t.Fatalf("grid %d: got %d cells, want %d", gridSize, len(icon.Filled), gridSize*gridSize)
			}

			for row := range gridSize {
				for col := range gridSize {
					cell := icon.Filled[row*gridSize+col]
					mirror := icon.Filled[row*gridSize+gridSize-1-col]
					if cell != mirror {
						t.Fatalf("grid %d id %q: cell (%d,%d) is %v but its mirror is %v", gridSize, id, col, row, cell, mirror)
					}
				}
			}
		}
	}
}

func TestBalancedGridIsDeterministic(t *testing.T) {
	img := model.Image{GridSize: model.DefaultGridSize, Mirror: true, Style: model.PatternBalanced}

	first, err := Generate("alice@example.com", img)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	second, err := Generate("alice@example.com", img)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("same id twice gave different icons:\n%+v\n%+v", first, second)
	}
}
