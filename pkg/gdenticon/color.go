// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

package gdenticon

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
)

func ParseColor(s string) (color.Color, error) {
	digits, ok := strings.CutPrefix(s, "#")
	if !ok {
		return nil, fmt.Errorf("want a leading #, got %q", s)
	}
	if len(digits) != 6 && len(digits) != 8 {
		return nil, fmt.Errorf("want 6 or 8 digits, got %d in %q", len(digits), s)
	}

	value, err := strconv.ParseUint(digits, 16, 32)
	if err != nil {
		return nil, fmt.Errorf("want hex digits, got %q", s)
	}

	if len(digits) == 6 {
		return color.NRGBA{R: uint8(value >> 16), G: uint8(value >> 8), B: uint8(value), A: 0xff}, nil
	}
	return color.NRGBA{R: uint8(value >> 24), G: uint8(value >> 16), B: uint8(value >> 8), A: uint8(value)}, nil
}

// HexOf returns c as "#rrggbb", or as "#rrggbbaa" when it is not opaque.
func HexOf(c color.Color) string {
	n := color.NRGBAModel.Convert(c).(color.NRGBA)
	if n.A == 0xff {
		return fmt.Sprintf("#%02x%02x%02x", n.R, n.G, n.B)
	}
	return fmt.Sprintf("#%02x%02x%02x%02x", n.R, n.G, n.B, n.A)
}

// FixedColor returns a ColorFunc that ignores the hash and always returns c.
func FixedColor(c color.Color) ColorFunc {
	return func([]byte) color.Color { return c }
}
