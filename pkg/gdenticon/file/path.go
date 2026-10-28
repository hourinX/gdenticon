// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

package file

import (
	"io/fs"
	"os"
	"path/filepath"
)

// Path joins dir with a generated base name and extension.
func Path(dir, name, extension string) string {
	return filepath.Join(dir, name+extension)
}

// EnsureDir creates dir and its parents when necessary.
func EnsureDir(dir string, perm fs.FileMode) error {
	return os.MkdirAll(dir, perm)
}
