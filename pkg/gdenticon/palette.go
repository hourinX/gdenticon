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

// palette setting
const (
	minContrast = 3.0 // foreground: least contrast it may have against the background, as a WCAG ratio
	maxContrast = 4.5 // foreground: most contrast it may have against the background, as a WCAG ratio
)

// HueColor derives hue and saturation from digest, then solves lightness so the
// result remains legible against the default background.
func HueColor(digest []byte) color.Color {
	seed := derive(digest, tagColor)
	hue := float64(binary.BigEndian.Uint16(seed)) / 65536
	saturation := 0.45 + 0.18*float64(seed[2])/255
	want := minContrast + (maxContrast-minContrast)*float64(seed[3])/255
	return luminanceColor(hue, saturation, want)
}

// luminanceColor finds the quantised colour nearest want while keeping contrast
// inside the accepted band.
func luminanceColor(hue, saturation, want float64) color.Color {
	lightness := lightnessAt(hue, saturation, want, true)

	if contrastOf(colourLuminance(hue, saturation, lightness), luminance(DefaultBackground)) > maxContrast {
		lightness = lightnessAt(hue, saturation, maxContrast, false)
	}
	return hsl(hue, saturation, lightness)
}

// lightnessAt returns the lightness at which the contrast of a ray crosses
// target: the lightest one still above it, or the darkest one no longer above it.
func lightnessAt(hue, saturation, target float64, above bool) float64 {
	background := luminance(DefaultBackground)
	lo, hi := 0.0, 1.0 // the contrast at lo is above target, the one at hi is not

	for range 12 {
		mid := (lo + hi) / 2
		if contrastOf(colourLuminance(hue, saturation, mid), background) > target {
			lo = mid
		} else {
			hi = mid
		}
	}
	if above {
		return lo
	}
	return hi
}

// colourLuminance measures the same quantised channels hsl writes.
func colourLuminance(hue, saturation, lightness float64) float64 {
	r, g, b := hslChannels(hue, saturation, lightness)
	return 0.2126*linearTable[r] + 0.7152*linearTable[g] + 0.0722*linearTable[b]
}

// contrastOf returns the WCAG contrast ratio of two relative luminances.
func contrastOf(a, b float64) float64 {
	if a < b {
		a, b = b, a
	}
	return (a + 0.05) / (b + 0.05)
}

// contrastRatio returns the WCAG contrast ratio of two opaque colours.
func contrastRatio(a, b color.Color) float64 {
	return contrastOf(luminance(a), luminance(b))
}

// luminance returns WCAG relative luminance from 8-bit colour channels.
func luminance(c color.Color) float64 {
	r, g, b, _ := c.RGBA()
	return 0.2126*linearTable[r>>8] + 0.7152*linearTable[g>>8] + 0.0722*linearTable[b>>8]
}

// linearTable caches the sRGB decode of every channel byte.
var linearTable = func() (table [256]float64) {
	for i := range table {
		c := float64(i) / 255
		if c <= 0.04045 {
			table[i] = c / 12.92
			continue
		}
		table[i] = math.Pow((c+0.055)/1.055, 2.4)
	}
	return
}()

// derive domain-separates a digest for one purpose.
func derive(seed []byte, tag byte) []byte {
	h := sha256.New()
	h.Write(seed)
	h.Write([]byte{tag})
	return h.Sum(nil)
}

// hsl returns the opaque sRGB colour for a hue, saturation and lightness, each
// in [0,1].
func hsl(hue, saturation, lightness float64) color.Color {
	r, g, b := hslChannels(hue, saturation, lightness)
	return color.NRGBA{R: r, G: g, B: b, A: 0xff}
}

// hslChannels returns the 8-bit sRGB channels of a hue, saturation and
// lightness, each in [0,1].
func hslChannels(hue, saturation, lightness float64) (r, g, b uint8) {
	chroma := saturation * math.Min(lightness, 1-lightness)
	channel := func(n float64) uint8 {
		k := math.Mod(n+hue*12, 12)
		return uint8(math.Round(255 * (lightness - chroma*math.Max(-1, math.Min(math.Min(k-3, 9-k), 1)))))
	}
	return channel(0), channel(8), channel(4)
}
