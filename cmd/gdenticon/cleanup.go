// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (s *interactiveSession) printOutputs(w io.Writer) {
	found := false
	for i, output := range s.outputs {
		if output.Removed {
			continue
		}
		if !found {
			fmt.Fprintln(w, "Session outputs:")
			found = true
		}
		status := "cached"
		if output.Created {
			status = "created"
		}
		fmt.Fprintf(w, "  [%d] %-7s %s\n", i+1, status, output.Path)
	}
	if !found {
		fmt.Fprintln(w, "No outputs have been generated in this session.")
	}
	fmt.Fprintln(w)
}

func (s *interactiveSession) removableCount() int {
	count := 0
	for _, output := range s.outputs {
		if output.Created && !output.Removed {
			count++
		}
	}
	return count
}

func (s *interactiveSession) remove(selector string) (outputRecord, error) {
	index, err := s.resolve(selector)
	if err != nil {
		return outputRecord{}, err
	}
	output := s.outputs[index]
	if err := removeGeneratedFile(output.Path); err != nil {
		return outputRecord{}, err
	}
	s.outputs[index].Removed = true
	return output, nil
}

func (s *interactiveSession) clean() (int, error) {
	removed := 0
	var failures []error
	for i, output := range s.outputs {
		if !output.Created || output.Removed {
			continue
		}
		if err := removeGeneratedFile(output.Path); err != nil {
			failures = append(failures, err)
			continue
		}
		s.outputs[i].Removed = true
		removed++
	}
	return removed, errors.Join(failures...)
}

func (s *interactiveSession) resolve(selector string) (int, error) {
	if selector == "last" {
		for i := len(s.outputs) - 1; i >= 0; i-- {
			if !s.outputs[i].Removed {
				return i, nil
			}
		}
		return 0, fmt.Errorf("no session output is available")
	}

	if number, err := strconv.Atoi(selector); err == nil {
		index := number - 1
		if index < 0 || index >= len(s.outputs) || s.outputs[index].Removed {
			return 0, fmt.Errorf("output number %d is not available", number)
		}
		return index, nil
	}

	matches := make([]int, 0, 1)
	for i, output := range s.outputs {
		if output.Removed {
			continue
		}
		if selector == output.Path || selector == filepath.Base(output.Path) {
			matches = append(matches, i)
		}
	}
	switch len(matches) {
	case 0:
		return 0, fmt.Errorf("%q is not an output from this session", selector)
	case 1:
		return matches[0], nil
	default:
		numbers := make([]string, len(matches))
		for i, match := range matches {
			numbers[i] = strconv.Itoa(match + 1)
		}
		return 0, fmt.Errorf("%q matches outputs [%s]; remove one by number", selector, strings.Join(numbers, ", "))
	}
}

func removeGeneratedFile(path string) error {
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove %s: %w", path, err)
	}

	// Best-effort only: os.Remove cannot remove a non-empty directory, so
	// unrelated and historical files keep both themselves and their directory.
	_ = os.Remove(filepath.Dir(path))
	return nil
}
