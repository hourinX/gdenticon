package gdenticon

import "testing"

func TestGridMirrorsAndExtendsEntropy(t *testing.T) {
	cells := grid(newEntropy([]byte{0xff}), 9, true)
	for row := range 9 {
		for col := range 9 {
			if cells[row*9+col] != cells[row*9+8-col] {
				t.Fatalf("row %d is not mirrored", row)
			}
		}
	}

	e := newEntropy([]byte{0x80})
	if !e.Bool() {
		t.Fatal("first entropy bit should be true")
	}
	for range 20 {
		_ = e.Bool()
	}
}

func TestPatternScoring(t *testing.T) {
	if !degenerate(make([]bool, 9), 3) {
		t.Fatal("blank grid should be degenerate")
	}
	solid := []bool{true, true, true, true}
	if !degenerate(solid, 2) || visualPenalty(solid, 2) == 0 {
		t.Fatal("solid grid should have a penalty")
	}
	connected := []bool{
		true, true, false,
		false, true, false,
		false, false, true,
	}
	if got := blocks(connected, 3); got != 2 {
		t.Fatalf("blocks = %d, want 2", got)
	}
	if got := filledCells(connected); got != 4 {
		t.Fatalf("filledCells = %d, want 4", got)
	}
	if got := balancedGridWithin(newEntropy([]byte{0}), 3, false, 1); len(got) != 9 {
		t.Fatalf("balanced grid length = %d", len(got))
	}
}

func TestGeometry(t *testing.T) {
	if got := edge(17, 5, 5); got != 17-edge(17, 5, 0) {
		t.Fatalf("opposite edge = %d", got)
	}
	r := cell(17, 5, 2, 3)
	if r.Empty() {
		t.Fatal("cell rectangle is empty")
	}
}
