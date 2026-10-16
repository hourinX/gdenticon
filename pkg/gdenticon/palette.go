// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

package gdenticon

import (
	"crypto/sha256"
	"encoding/binary"
	"image/color"
	"math"
)

// tagColor domain-separates the colour stream from the one the grid draws on.
const tagColor = 0x50

// HueColor returns the dark saturated colour that digest implies. It is what
// Generate falls back to when Image.Foreground is nil.
func HueColor(digest []byte) color.Color {
	seed := derive(digest, tagColor)
	hue := float64(binary.BigEndian.Uint16(seed)) / 65536
	saturation := 0.45 + 0.18*float64(seed[2])/255
	lightness := 0.58 + 0.14*float64(seed[3])/255
	return hsl(hue, saturation, lightness)
}

// derive returns seed marked as belonging to one purpose, so that two streams
// taken from the same digest cannot share bits.
func derive(seed []byte, tag byte) []byte {
	h := sha256.New()
	h.Write(seed)
	h.Write([]byte{tag})
	return h.Sum(nil)
}

// hsl returns the opaque sRGB colour for a hue, saturation and lightness, each
// in [0,1].
func hsl(hue, saturation, lightness float64) color.Color {
	chroma := saturation * math.Min(lightness, 1-lightness)
	channel := func(n float64) uint8 {
		k := math.Mod(n+hue*12, 12)
		return uint8(math.Round(255 * (lightness - chroma*math.Max(-1, math.Min(math.Min(k-3, 9-k), 1)))))
	}
	return color.NRGBA{R: channel(0), G: channel(8), B: channel(4), A: 0xff}
}
