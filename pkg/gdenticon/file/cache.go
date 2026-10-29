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

// CacheEntry describes the result of a content-addressed cache write.
type CacheEntry struct {
	Path    string
	Created bool
}

// PutCache writes content to a content-addressed file in dir. Equal encoded
// content and extension resolve to the same path.
func PutCache(dir, extension string, perm fs.FileMode, write WriteFunc) (string, error) {
	entry, err := PutCacheEntry(dir, extension, perm, write)
	return entry.Path, err
}

// PutCacheEntry writes content like PutCache and reports whether this call
// created the returned file instead of resolving an existing cache entry.
func PutCacheEntry(dir, extension string, perm fs.FileMode, write WriteFunc) (CacheEntry, error) {
	if err := EnsureDir(dir, 0o755); err != nil {
		return CacheEntry{}, err
	}

	f, err := os.CreateTemp(dir, ".gdenticon-cache-*")
	if err != nil {
		return CacheEntry{}, err
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
		return CacheEntry{}, err
	}
	if err := f.Chmod(perm); err != nil {
		_ = f.Close()
		return CacheEntry{}, err
	}
	if err := f.Close(); err != nil {
		return CacheEntry{}, err
	}

	name := hex.EncodeToString(hash.Sum(nil))[:16]
	path := filepath.Join(dir, name+extension)
	if _, err := os.Stat(path); err == nil {
		return CacheEntry{Path: path}, nil
	} else if !os.IsNotExist(err) {
		return CacheEntry{}, err
	}
	if err := moveNew(temp, path, perm); err != nil {
		// A concurrent writer may have won after the Stat above.
		if _, statErr := os.Stat(path); statErr == nil {
			return CacheEntry{Path: path}, nil
		}
		return CacheEntry{}, err
	}
	keep = true
	return CacheEntry{Path: path, Created: true}, nil
}

// moveNew installs source without replacing an existing target. Copying is a
// fallback for filesystems that do not support hard links.
func moveNew(source, target string, perm fs.FileMode) error {
	return moveNewWithLink(source, target, perm, os.Link)
}

func moveNewWithLink(source, target string, perm fs.FileMode, link func(string, string) error) error {
	if err := link(source, target); err == nil {
		if err := os.Remove(source); err != nil {
			_ = os.Remove(target)
			return err
		}
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
	keep := false
	defer func() {
		if !keep {
			_ = os.Remove(target)
		}
	}()
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
	if err := os.Remove(source); err != nil {
		return err
	}
	keep = true
	return nil
}
