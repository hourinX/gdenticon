package gdenticon

import (
	"fmt"
	"testing"

	"github.com/hourinX/gdenticon/pkg/gdenticon/model"
)

func TestHueColorKeepsContrastInsideTheBand(t *testing.T) {
	const ids = 20000

	lightest := 99.0
	for i := range ids {
		c := HueColor(SHA256(fmt.Sprintf("user%d", i)))
		ratio := contrastRatio(c, model.DefaultBackground)

		if ratio < minContrast {
			t.Fatalf("id %d: contrast %.4f is under the floor of %.1f", i, ratio, minContrast)
		}
		if ratio > maxContrast {
			t.Fatalf("id %d: contrast %.4f is over the ceiling of %.1f", i, ratio, maxContrast)
		}
		if ratio < lightest {
			lightest = ratio
		}
	}
	t.Logf("lightest of %d colours sits at %.4f", ids, lightest)
}

func TestLuminanceColorHoldsTheBandAcrossItsWholeRange(t *testing.T) {
	// Swept over the parameters the solver takes rather than over digests, so
	// the cases the digest distribution happens to miss are covered too.
	checked := 0
	for hue := 0.0; hue < 1; hue += 0.01 {
		for saturation := 0.45; saturation <= 0.64; saturation += 0.005 {
			for want := minContrast; want <= maxContrast; want += 0.025 {
				ratio := contrastRatio(luminanceColor(hue, saturation, want), model.DefaultBackground)
				if ratio < minContrast || ratio > maxContrast {
					t.Fatalf("hue %.2f saturation %.3f want %.3f: contrast %.4f left the band", hue, saturation, want, ratio)
				}
				checked++
			}
		}
	}
	t.Logf("%d combinations stayed inside the band", checked)
}

func TestHueColorStillVaries(t *testing.T) {
	// The band would hold trivially if every id got the same colour, so the
	// variety the digest is supposed to drive is asserted alongside it.
	seen := map[string]bool{}
	for i := range 1000 {
		r, g, b, _ := HueColor(SHA256(fmt.Sprintf("user%d", i))).RGBA()
		seen[fmt.Sprintf("%d-%d-%d", uint8(r>>8), uint8(g>>8), uint8(b>>8))] = true
	}

	if len(seen) < 500 {
		t.Errorf("1000 ids produced only %d colours", len(seen))
	}
}

func TestContrastRatioIsTheWcagOne(t *testing.T) {
	// The two anchors WCAG names, which pin the formula down on its own.
	if got := contrastRatio(model.DefaultBackground, model.DefaultBackground); got != 1 {
		t.Errorf("white on white is %.4f, want 1", got)
	}
	if got := contrastRatio(hsl(0, 0, 0), model.DefaultBackground); got < 20.9 || got > 21.01 {
		t.Errorf("black on white is %.4f, want 21", got)
	}
}
