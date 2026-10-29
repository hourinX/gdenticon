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
	return runWithGenerated(args, stdout, stderr, nil)
}

func runProgram(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if interactiveRequested(args) {
		return runInteractive(stdin, stdout, stderr)
	}
	return run(args, stdout, stderr)
}

func runWithGenerated(args []string, stdout, stderr io.Writer, generated func(renderResult)) int {
	return execute(createCommandTree(generated), args, stdout, stderr)
}

func createCommandTree(generated func(renderResult)) *cobra.Command {
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
			result, err := execRender(args[0], f)
			if err != nil {
				return err
			}
			if generated != nil {
				generated(result)
			}
			fmt.Fprintln(cmd.OutOrStdout(), result.Path)
			return nil
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
