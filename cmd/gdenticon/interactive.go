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
	if _, err := fmt.Fprintln(stdout, "gdenticon interactive mode"); err != nil {
		return 1
	}
	if _, err := fmt.Fprintln(stdout, `Type "help" for commands and "exit" to quit.`); err != nil {
		return 1
	}
	if _, err := fmt.Fprintln(stdout); err != nil {
		return 1
	}

	scanner := bufio.NewScanner(stdin)
	session := &interactiveSession{}
	for {
		if _, err := fmt.Fprint(stdout, "gdenticon> "); err != nil {
			return 1
		}
		if !scanner.Scan() {
			break
		}

		args, err := splitCommandLine(scanner.Text())
		if err != nil {
			if _, writeErr := fmt.Fprintln(stderr, "Error:", err); writeErr != nil {
				return 1
			}
			if _, writeErr := fmt.Fprintln(stderr); writeErr != nil {
				return 1
			}
			continue
		}
		if len(args) == 0 {
			continue
		}

		switch args[0] {
		case "exit", "quit":
			return 0
		case "help":
			if err := printInteractiveHelp(stdout); err != nil {
				return 1
			}
			if _, err := fmt.Fprintln(stdout); err != nil {
				return 1
			}
			continue
		case "outputs":
			if len(args) != 1 {
				if err := printInteractiveError(stderr, `outputs does not accept arguments`); err != nil {
					return 1
				}
				continue
			}
			if err := session.printOutputs(stdout); err != nil {
				return 1
			}
			continue
		case "remove":
			if len(args) != 2 {
				if err := printInteractiveError(stderr, `remove requires an output number, filename, or "last"`); err != nil {
					return 1
				}
				continue
			}
			output, err := session.remove(args[1])
			if err != nil {
				if writeErr := printInteractiveError(stderr, err.Error()); writeErr != nil {
					return 1
				}
				continue
			}
			if _, err := fmt.Fprintf(stdout, "Removed:\n  %s\n\n", output.Path); err != nil {
				return 1
			}
			continue
		case "clean":
			if len(args) != 1 {
				if err := printInteractiveError(stderr, `clean does not accept arguments`); err != nil {
					return 1
				}
				continue
			}
			count := session.removableCount()
			if count == 0 {
				if _, err := fmt.Fprintln(stdout, "No files created in this session are available to clean."); err != nil {
					return 1
				}
				if _, err := fmt.Fprintln(stdout); err != nil {
					return 1
				}
				continue
			}
			if _, err := fmt.Fprintf(stdout, "Remove %d files created in this session? [y/N] ", count); err != nil {
				return 1
			}
			if !scanner.Scan() {
				if err := scanner.Err(); err != nil {
					if _, writeErr := fmt.Fprintln(stderr, "gdenticon: read confirmation:", err); writeErr != nil {
						return 1
					}
					return 1
				}
				return 0
			}
			answer := strings.ToLower(strings.TrimSpace(scanner.Text()))
			if answer != "y" && answer != "yes" {
				if _, err := fmt.Fprintln(stdout, "Cancelled."); err != nil {
					return 1
				}
				if _, err := fmt.Fprintln(stdout); err != nil {
					return 1
				}
				continue
			}
			removed, err := session.clean()
			if err != nil {
				if writeErr := printInteractiveError(stderr, err.Error()); writeErr != nil {
					return 1
				}
				continue
			}
			if _, err := fmt.Fprintf(stdout, "Removed %d files.\n\n", removed); err != nil {
				return 1
			}
			continue
		case "version":
			args = []string{"--version"}
		case "generate":
			args = args[1:]
		default:
			if _, err := fmt.Fprintf(stderr, "Error: unknown command %q\n", args[0]); err != nil {
				return 1
			}
			if _, err := fmt.Fprintln(stderr, `Type "help" for available commands.`); err != nil {
				return 1
			}
			if _, err := fmt.Fprintln(stderr); err != nil {
				return 1
			}
			continue
		}

		var commandOut, commandErr bytes.Buffer
		var generated *renderResult
		code := runWithGenerated(args, &commandOut, &commandErr, func(result renderResult) {
			session.add(result)
			generated = &result
		})
		if err := printInteractiveResult(stdout, stderr, args, code, commandOut.String(), commandErr.String(), generated); err != nil {
			return 1
		}
	}
	if err := scanner.Err(); err != nil {
		if _, writeErr := fmt.Fprintln(stderr, "gdenticon: read command:", err); writeErr != nil {
			return 1
		}
		return 1
	}
	return 0
}

func printInteractiveResult(stdout, stderr io.Writer, args []string, code int, output, failure string, generated *renderResult) error {
	if code != 0 {
		message, _, _ := strings.Cut(strings.TrimSpace(failure), "\n")
		message = strings.TrimPrefix(message, "gdenticon: ")
		if _, err := fmt.Fprintln(stderr, "Error:", message); err != nil {
			return err
		}
		return writeLines(stderr, `Run "generate --help" to view generation options.`, "")
	}

	if hasHelpFlag(args) || len(args) == 1 && args[0] == "--version" {
		return writeLines(stdout, strings.TrimSpace(output), "")
	}
	label := "Created:"
	if generated != nil && !generated.Created {
		label = "Cached:"
	}
	if _, err := fmt.Fprintln(stdout, label); err != nil {
		return err
	}
	_, err := fmt.Fprintf(stdout, "  %s\n\n", strings.TrimSpace(output))
	return err
}

func printInteractiveError(w io.Writer, message string) error {
	if _, err := fmt.Fprintln(w, "Error:", message); err != nil {
		return err
	}
	return writeLines(w, `Type "help" for available commands.`, "")
}

func hasHelpFlag(args []string) bool {
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			return true
		}
	}
	return false
}

func printInteractiveHelp(w io.Writer) error {
	return writeLines(w,
		"Commands:",
		"  generate <id> [flags]  generate an identicon",
		"  outputs                list outputs from this session",
		"  remove <target>        remove by number, filename, or last",
		"  clean                  remove files created in this session",
		"  version                print the version",
		"  help                   show this help",
		"  exit, quit             leave interactive mode",
		"Run generate --help to list generation flags.",
	)
}

func writeLines(w io.Writer, lines ...string) error {
	for _, line := range lines {
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
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
