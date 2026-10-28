package gdenticon

import (
	"encoding/hex"
	"testing"
)

func TestDigests(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"md5", 16}, {"sha1", 20}, {"sha224", 28}, {"sha256", 32},
		{"sha384", 48}, {"sha512", 64}, {"sha512-224", 28}, {"sha512-256", 32},
		{"sha3-224", 28}, {"sha3-256", 32}, {"sha3-384", 48}, {"sha3-512", 64},
	}
	for _, test := range tests {
		fn, ok := Digest(test.name)
		if !ok {
			t.Fatalf("Digest(%q) was not found", test.name)
		}
		if got := len(fn("gdenticon")); got != test.length {
			t.Fatalf("%s length = %d, want %d", test.name, got, test.length)
		}
	}
	if _, ok := Digest("unknown"); ok {
		t.Fatal("unknown digest was found")
	}
	if got := hex.EncodeToString(MD5("abc")); got != "900150983cd24fb0d6963f7d28e17f72" {
		t.Fatalf("MD5 known vector = %s", got)
	}
}
