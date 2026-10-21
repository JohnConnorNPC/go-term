// Command conformance scores a terminal engine against the xterm.js headless goldens.
//
//	go run ./conformance/cmd/conformance -engine stub -ref x53
//	go run ./conformance/cmd/conformance -count
//
// Exit status: 0 every case passes (or the counts meet the minimums), 1 a case fails, 2 something is unmeasured.
// The body is runner.CLI, so other modules can run it with their own engines.
package main

import (
	"os"

	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

func main() { os.Exit(runner.CLI(os.Args[1:], os.Stdout, nil)) }
