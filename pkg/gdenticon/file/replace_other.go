// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

//go:build !windows

package file

import "os"

// replacePath atomically replaces target with source on Unix-like systems.
func replacePath(source, target string) error {
	return os.Rename(source, target)
}
