// SPDX-License-Identifier: MIT
// Copyright (c) 2026 hourinX

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

func interactiveRequested(args []string) bool {
	if len(args) == 0 {
		return true
	}
	return len(args) == 1 && (args[0] == "shell" || args[0] == "--interactive")
}

func runInteractive(stdin io.Reader, stdout, stderr io.Writer) int {
	fmt.Fprintln(stdout, "gdenticon interactive mode")
	fmt.Fprintln(stdout, `Type "help" for commands and "exit" to quit.`)
	fmt.Fprintln(stdout)

	scanner := bufio.NewScanner(stdin)
	session := &interactiveSession{}
	for {
		fmt.Fprint(stdout, "gdenticon> ")
		if !scanner.Scan() {
			break
		}

		args, err := splitCommandLine(scanner.Text())
		if err != nil {
			fmt.Fprintln(stderr, "Error:", err)
			fmt.Fprintln(stderr)
			continue
		}
		if len(args) == 0 {
			continue
		}

		switch args[0] {
		case "exit", "quit":
			return 0
		case "help":
			printInteractiveHelp(stdout)
			fmt.Fprintln(stdout)
			continue
		case "outputs":
			if len(args) != 1 {
				printInteractiveError(stderr, `outputs does not accept arguments`)
				continue
			}
			session.printOutputs(stdout)
			continue
		case "remove":
			if len(args) != 2 {
				printInteractiveError(stderr, `remove requires an output number, filename, or "last"`)
				continue
			}
			output, err := session.remove(args[1])
			if err != nil {
				printInteractiveError(stderr, err.Error())
				continue
			}
			fmt.Fprintf(stdout, "Removed:\n  %s\n\n", output.Path)
			continue
		case "clean":
			if len(args) != 1 {
				printInteractiveError(stderr, `clean does not accept arguments`)
				continue
			}
			count := session.removableCount()
			if count == 0 {
				fmt.Fprintln(stdout, "No files created in this session are available to clean.")
				fmt.Fprintln(stdout)
				continue
			}
			fmt.Fprintf(stdout, "Remove %d files created in this session? [y/N] ", count)
			if !scanner.Scan() {
				if err := scanner.Err(); err != nil {
					fmt.Fprintln(stderr, "gdenticon: read confirmation:", err)
					return 1
				}
				return 0
			}
			answer := strings.ToLower(strings.TrimSpace(scanner.Text()))
			if answer != "y" && answer != "yes" {
				fmt.Fprintln(stdout, "Cancelled.")
				fmt.Fprintln(stdout)
				continue
			}
			removed, err := session.clean()
			if err != nil {
				printInteractiveError(stderr, err.Error())
				continue
			}
			fmt.Fprintf(stdout, "Removed %d files.\n\n", removed)
			continue
		case "version":
			args = []string{"--version"}
		case "generate":
			args = args[1:]
		default:
			fmt.Fprintf(stderr, "Error: unknown command %q\n", args[0])
			fmt.Fprintln(stderr, `Type "help" for available commands.`)
			fmt.Fprintln(stderr)
			continue
		}

		var commandOut, commandErr bytes.Buffer
		var generated *renderResult
		code := runWithGenerated(args, &commandOut, &commandErr, func(result renderResult) {
			session.add(result)
			generated = &result
		})
		printInteractiveResult(stdout, stderr, args, code, commandOut.String(), commandErr.String(), generated)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(stderr, "gdenticon: read command:", err)
		return 1
	}
	return 0
}

func printInteractiveResult(stdout, stderr io.Writer, args []string, code int, output, failure string, generated *renderResult) {
	if code != 0 {
		message, _, _ := strings.Cut(strings.TrimSpace(failure), "\n")
		message = strings.TrimPrefix(message, "gdenticon: ")
		fmt.Fprintln(stderr, "Error:", message)
		fmt.Fprintln(stderr, `Run "generate --help" to view generation options.`)
		fmt.Fprintln(stderr)
		return
	}

	if hasHelpFlag(args) || len(args) == 1 && args[0] == "--version" {
		fmt.Fprintln(stdout, strings.TrimSpace(output))
		fmt.Fprintln(stdout)
		return
	}
	label := "Created:"
	if generated != nil && !generated.Created {
		label = "Cached:"
	}
	fmt.Fprintln(stdout, label)
	fmt.Fprintf(stdout, "  %s\n\n", strings.TrimSpace(output))
}

func printInteractiveError(w io.Writer, message string) {
	fmt.Fprintln(w, "Error:", message)
	fmt.Fprintln(w, `Type "help" for available commands.`)
	fmt.Fprintln(w)
}

func hasHelpFlag(args []string) bool {
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			return true
		}
	}
	return false
}

func printInteractiveHelp(w io.Writer) {
	fmt.Fprintln(w, "Commands:")
	fmt.Fprintln(w, "  generate <id> [flags]  generate an identicon")
	fmt.Fprintln(w, "  outputs                list outputs from this session")
	fmt.Fprintln(w, "  remove <target>        remove by number, filename, or last")
	fmt.Fprintln(w, "  clean                  remove files created in this session")
	fmt.Fprintln(w, "  version                print the version")
	fmt.Fprintln(w, "  help                   show this help")
	fmt.Fprintln(w, "  exit, quit             leave interactive mode")
	fmt.Fprintln(w, "Run generate --help to list generation flags.")
}

func splitCommandLine(line string) ([]string, error) {
	var args []string
	var token strings.Builder
	var quote rune
	haveToken := false

	flush := func() {
		if haveToken {
			args = append(args, token.String())
			token.Reset()
			haveToken = false
		}
	}

	for _, r := range line {
		if quote != 0 {
			if r == quote {
				quote = 0
				continue
			}
			token.WriteRune(r)
			haveToken = true
			continue
		}

		switch {
		case r == '\'' || r == '"':
			quote = r
			haveToken = true
		case unicode.IsSpace(r):
			flush()
		default:
			token.WriteRune(r)
			haveToken = true
		}
	}
	if quote != 0 {
		return nil, fmt.Errorf("unterminated quoted value")
	}
	flush()
	return args, nil
}
