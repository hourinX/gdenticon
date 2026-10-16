// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

package gdenticon

import (
	"os"
	"path/filepath"

	"github.com/hourinX/gdenticon/pkg/asset"
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
