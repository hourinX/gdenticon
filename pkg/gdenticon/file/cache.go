// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

package file

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// PutCache writes content to a content-addressed file in dir. Equal encoded
// content and extension resolve to the same path.
func PutCache(dir, extension string, perm fs.FileMode, write WriteFunc) (string, error) {
	if err := EnsureDir(dir, 0o755); err != nil {
		return "", err
	}

	f, err := os.CreateTemp(dir, ".gdenticon-cache-*")
	if err != nil {
		return "", err
	}
	temp := f.Name()
	keep := false
	defer func() {
		if !keep {
			_ = os.Remove(temp)
		}
	}()

	hash := sha256.New()
	if err := write(io.MultiWriter(f, hash)); err != nil {
		_ = f.Close()
		return "", err
	}
	if err := f.Chmod(perm); err != nil {
		_ = f.Close()
		return "", err
	}
	if err := f.Close(); err != nil {
		return "", err
	}

	name := hex.EncodeToString(hash.Sum(nil))[:16]
	path := filepath.Join(dir, name+extension)
	if _, err := os.Stat(path); err == nil {
		return path, nil
	} else if !os.IsNotExist(err) {
		return "", err
	}
	if err := moveNew(temp, path, perm); err != nil {
		// A concurrent writer may have won after the Stat above.
		if _, statErr := os.Stat(path); statErr == nil {
			return path, nil
		}
		return "", err
	}
	keep = true
	return path, nil
}

// moveNew installs source at a path that must not already exist. Copying is a
// fallback for restricted filesystems that deny renames of temporary files.
func moveNew(source, target string, perm fs.FileMode) error {
	if err := os.Rename(source, target); err == nil {
		return nil
	}
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, perm)
	if err != nil {
		_ = in.Close()
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = in.Close()
		_ = out.Close()
		return err
	}
	if err := in.Close(); err != nil {
		_ = out.Close()
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}
	return os.Remove(source)
}
