package gdenticon

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"io"
	"strings"
	"testing"
)

func testIcon(format Format) Icon {
	return Icon{
		GridSize:   2,
		Filled:     []bool{true, false, false, true},
		Background: color.White,
		Foreground: color.Black,
		Format:     format,
	}
}

func TestWriteFormats(t *testing.T) {
	tests := []struct {
		format Format
		prefix []byte
	}{
		{FormatPNG, []byte("\x89PNG")},
		{FormatSVG, []byte("<svg")},
		{FormatJPEG, []byte("\xff\xd8")},
		{FormatGIF, []byte("GIF8")},
		{FormatBMP, []byte("BM")},
		{FormatTIFF, []byte("II*")},
	}

	for _, test := range tests {
		t.Run(string(test.format), func(t *testing.T) {
			var output bytes.Buffer
			if err := Write(&output, testIcon(test.format), WriteOptions{Size: 17}); err != nil {
				t.Fatal(err)
			}
			if !bytes.HasPrefix(output.Bytes(), test.prefix) {
				t.Fatalf("output prefix is %x, want %x", output.Bytes()[:len(test.prefix)], test.prefix)
			}
			if test.format != FormatSVG {
				decoded, _, err := image.Decode(bytes.NewReader(output.Bytes()))
				if err != nil {
					t.Fatal(err)
				}
				if decoded.Bounds().Dx() != 17 || decoded.Bounds().Dy() != 17 {
					t.Fatalf("image bounds are %v", decoded.Bounds())
				}
			}
		})
	}
}

func TestWriteDefaultsAndErrors(t *testing.T) {
	var output bytes.Buffer
	if err := Write(&output, testIcon(FormatPNG), WriteOptions{}); err != nil {
		t.Fatal(err)
	}
	decoded, _, err := image.Decode(bytes.NewReader(output.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	if decoded.Bounds().Dx() != DefaultImageSize {
		t.Fatalf("default size is %d, want %d", decoded.Bounds().Dx(), DefaultImageSize)
	}

	tests := []struct {
		name string
		icon Icon
		opts WriteOptions
		want string
	}{
		{"format", testIcon(Format("webp")), WriteOptions{Size: 10}, "unsupported image format"},
		{"quality low", testIcon(FormatJPEG), WriteOptions{Size: 10, JPEGQuality: -1}, "JPEG quality"},
		{"quality high", testIcon(FormatJPEG), WriteOptions{Size: 10, JPEGQuality: 101}, "JPEG quality"},
		{"small raster", testIcon(FormatPNG), WriteOptions{Size: 1}, "at least the grid size"},
		{"large raster", testIcon(FormatPNG), WriteOptions{Size: MaxImageSize + 1}, "at most"},
		{"empty cells", Icon{GridSize: 2, Background: color.White, Foreground: color.Black, Format: FormatPNG}, WriteOptions{Size: 10}, "cells"},
		{"nil colour", Icon{GridSize: 1, Filled: []bool{true}, Format: FormatSVG}, WriteOptions{Size: 10}, "colours"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Write(io.Discard, test.icon, test.opts)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("error is %v, want substring %q", err, test.want)
			}
		})
	}
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, errors.New("write failed") }

func TestWritePropagatesWriterError(t *testing.T) {
	if err := Write(failingWriter{}, testIcon(FormatSVG), WriteOptions{Size: 10}); err == nil {
		t.Fatal("Write returned nil for a failing writer")
	}
}
