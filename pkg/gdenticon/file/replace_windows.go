// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

//go:build windows

package file

import (
	"io"
	"os"

	"golang.org/x/sys/windows"
)

// replacePath atomically replaces target with source on Windows.
func replacePath(source, target string) error {
	sourcePtr, err := windows.UTF16PtrFromString(source)
	if err != nil {
		return err
	}
	targetPtr, err := windows.UTF16PtrFromString(target)
	if err != nil {
		return err
	}
	if err := windows.MoveFileEx(sourcePtr, targetPtr, windows.MOVEFILE_REPLACE_EXISTING|windows.MOVEFILE_WRITE_THROUGH); err == nil {
		return nil
	}

	// Some restricted Windows environments deny MoveFileEx even within one
	// directory. Preserve functional replacement there, without claiming the
	// fallback is atomic.
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
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
