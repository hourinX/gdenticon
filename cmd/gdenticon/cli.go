// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// normalizeArgs rewrites long flags spelled with a single dash.
func normalizeArgs(cmd *cobra.Command, args []string) []string {
	normalized := make([]string, 0, len(args))
	for _, arg := range args {
		name, _, _ := strings.Cut(arg, "=")
		if strings.HasPrefix(name, "-") && !strings.HasPrefix(name, "--") &&
			len(name) > 2 && cmd.Flags().Lookup(name[1:]) != nil {
			normalized = append(normalized, "-"+arg)
			continue
		}
		normalized = append(normalized, arg)
	}
	return normalized
}

type usageError struct {
	err error
}

func (e usageError) Error() string {
	return e.err.Error()
}

func (e usageError) Unwrap() error {
	return e.err
}

func exactlyOneID(_ *cobra.Command, args []string) error {
	if len(args) != 1 {
		return usageError{fmt.Errorf("expected exactly one id, got %d", len(args))}
	}
	return nil
}
