// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/hourinX/gdenticon/v2/pkg/asset"
	"github.com/spf13/cobra"
)

func main() {
	os.Exit(runProgram(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	return execute(createCommandTree(), args, stdout, stderr)
}

func runProgram(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if interactiveRequested(args) {
		return runInteractive(stdin, stdout, stderr)
	}
	return run(args, stdout, stderr)
}

func createCommandTree() *cobra.Command {
	f := &renderFlags{}

	cmd := &cobra.Command{
		Use:               "gdenticon <id>",
		Short:             "Generate a deterministic identicon for id",
		Version:           asset.Version,
		Args:              exactlyOneID,
		SilenceUsage:      true,
		SilenceErrors:     true,
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		RunE: func(cmd *cobra.Command, args []string) error {
			return execRender(cmd, args[0], f)
		},
	}
	cmd.SetVersionTemplate("gdenticon {{.Version}}\n")
	f.register(cmd)

	return cmd
}

func execute(cmd *cobra.Command, args []string, stdout, stderr io.Writer) int {
	cmd.InitDefaultHelpFlag()
	cmd.InitDefaultVersionFlag()
	cmd.SetArgs(normalizeArgs(cmd, args))
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		return usageError{err}
	})

	err := cmd.Execute()
	if err == nil {
		return 0
	}

	var ue usageError
	if errors.As(err, &ue) {
		fmt.Fprintln(stderr, "gdenticon:", err)
		fmt.Fprint(stderr, cmd.UsageString())
		return 2
	}

	fmt.Fprintln(stderr, "gdenticon:", err)
	return 1
}
