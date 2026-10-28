// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

package file

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// WriteFunc writes encoded content to w without closing it.
type WriteFunc func(w io.Writer) error

// Create writes a new file at path and fails when the path already exists.
func Create(path string, perm fs.FileMode, write WriteFunc) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, perm)
	if err != nil {
		return err
	}
	return finish(f, write)
}

// Replace writes path through a temporary sibling and then replaces path.
func Replace(path string, perm fs.FileMode, write WriteFunc) error {
	dir := filepath.Dir(path)
	f, err := os.CreateTemp(dir, ".gdenticon-*")
	if err != nil {
		return err
	}
	temp := f.Name()
	ok := false
	defer func() {
		if !ok {
			_ = os.Remove(temp)
		}
	}()

	if err := f.Chmod(perm); err != nil {
		_ = f.Close()
		return err
	}
	if err := finish(f, write); err != nil {
		return err
	}
	if err := replacePath(temp, path); err != nil {
		return fmt.Errorf("replace %s: %w", path, err)
	}
	ok = true
	return nil
}

// CreateTemp writes a uniquely named file in dir and returns its path.
func CreateTemp(dir, pattern string, perm fs.FileMode, write WriteFunc) (string, error) {
	f, err := os.CreateTemp(dir, pattern)
	if err != nil {
		return "", err
	}
	path := f.Name()
	ok := false
	defer func() {
		if !ok {
			_ = os.Remove(path)
		}
	}()

	if err := f.Chmod(perm); err != nil {
		_ = f.Close()
		return "", err
	}
	if err := finish(f, write); err != nil {
		return "", err
	}
	ok = true
	return path, nil
}

// finish writes and closes f.
func finish(f *os.File, write WriteFunc) error {
	if err := write(f); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}
