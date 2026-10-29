// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

package main

import "path/filepath"

type outputRecord struct {
	ID      string
	Path    string
	Created bool
	Removed bool
}

type interactiveSession struct {
	outputs []outputRecord
}

func (s *interactiveSession) add(result renderResult) {
	for i := range s.outputs {
		if samePath(s.outputs[i].Path, result.Path) {
			s.outputs[i].Created = s.outputs[i].Created || result.Created
			s.outputs[i].Removed = false
			return
		}
	}
	s.outputs = append(s.outputs, outputRecord{
		ID:      result.ID,
		Path:    result.Path,
		Created: result.Created,
	})
}

func samePath(a, b string) bool {
	a, errA := filepath.Abs(a)
	b, errB := filepath.Abs(b)
	if errA == nil && errB == nil {
		return filepath.Clean(a) == filepath.Clean(b)
	}
	return filepath.Clean(a) == filepath.Clean(b)
}
