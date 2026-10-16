// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

package gdenticon

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"

	"github.com/hourinX/gdenticon/pkg/asset"
	"github.com/hourinX/gdenticon/pkg/gdenticon/model"
)

// Dir returns the directory that generated icons are written to. An empty base
// selects the system temp directory.
func Dir(base string) string {
	if base == "" {
		base = os.TempDir()
	}
	return filepath.Join(base, asset.FolderName)
}

// Prepare returns Dir(base), creating it if it does not already exist.
func Prepare(base string) (string, error) {
	dir := Dir(base)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

// Save renders the icon for id, writes it under base, and returns the path written.
func Save(id string, img model.Image, base string) (string, error) {
	var buf bytes.Buffer
	if err := WritePNG(&buf, id, img); err != nil {
		return "", err
	}

	dir, err := Prepare(base)
	if err != nil {
		return "", err
	}

	path := filepath.Join(dir, fileName(buf.Bytes()))
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		return "", err
	}
	return path, nil
}

// fileName names a rendered image after its contents.
func fileName(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])[:16] + ".png"
}
