package gdenticon

import (
	"image/color"
	"testing"
)

func TestParseColor(t *testing.T) {
	tests := []struct {
		input string
		want  color.NRGBA
	}{
		{"#123456", color.NRGBA{0x12, 0x34, 0x56, 0xff}},
		{"#12345678", color.NRGBA{0x12, 0x34, 0x56, 0x78}},
	}
	for _, test := range tests {
		got, err := ParseColor(test.input)
		if err != nil {
			t.Fatal(err)
		}
		if got != test.want {
			t.Fatalf("ParseColor(%q) = %#v, want %#v", test.input, got, test.want)
		}
	}

	for _, input := range []string{"123456", "#123", "#gggggg"} {
		if _, err := ParseColor(input); err == nil {
			t.Fatalf("ParseColor(%q) returned nil error", input)
		}
	}
}

func TestColourHelpers(t *testing.T) {
	if got := HexOf(color.NRGBA{R: 1, G: 2, B: 3, A: 0xff}); got != "#010203" {
		t.Fatalf("opaque HexOf = %q", got)
	}
	if got := HexOf(color.NRGBA{R: 1, G: 2, B: 3, A: 4}); got != "#01020304" {
		t.Fatalf("alpha HexOf = %q", got)
	}
	fixed := color.NRGBA{R: 10, G: 20, B: 30, A: 255}
	if got := FixedColor(fixed)([]byte("ignored")); got != fixed {
		t.Fatalf("FixedColor returned %#v", got)
	}

	a := HueColor([]byte("same"))
	b := HueColor([]byte("same"))
	if a != b {
		t.Fatalf("HueColor is not deterministic: %#v != %#v", a, b)
	}
	ratio := contrastRatio(a, DefaultBackground)
	if ratio < minContrast || ratio > maxContrast+0.1 {
		t.Fatalf("contrast ratio %f is outside expected range", ratio)
	}
}
