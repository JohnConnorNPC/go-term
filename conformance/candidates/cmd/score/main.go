// Command score scores one candidate on one golden set under one equivalence mode, prints the runner's lines and
// writes the E6 summary:
//
//	go run ./cmd/score -engine gitpod -ref x53 -equiv full -out scores/gitpod-x53-full.json
//
// Exit status: 0 every case passes, 1 a case fails, 2 something is UNMEASURED (no file is written then): a STALE
// golden, a golden or category missing, an unknown engine or equiv mode. An unobserved path is a FAIL line, never
// UNMEASURED.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/equiv"
	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/scoring"
	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/summary"
	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

func main() { os.Exit(run(os.Args[1:], os.Stdout)) }

func run(args []string, out io.Writer) int {
	fs := flag.NewFlagSet("score", flag.ContinueOnError)
	engine := fs.String("engine", "", "gitpod, magnet, xvt (or the anchors stub, oracle)")
	ref := fs.String("ref", "x53", "golden set: x53 or xcur")
	mode := fs.String("equiv", "full", "equivalence table: strict, case, default or full (the gate)")
	root := fs.String("root", "..", "the conformance directory (cases/, golden/, probes.json)")
	table := fs.String("table", "equiv/table.json", "the equivalence table")
	outPath := fs.String("out", "", "the summary file")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	pr, err := runner.LoadProbes(filepath.Join(*root, "probes.json"))
	f, ok := scoring.Lookup(*engine, pr)
	gs, gerr := runner.LoadAll(*root, *ref)
	lines, counted := runner.CheckCounts(*root, []string{*ref})
	if eerr := equiv.Load(*table, *mode); !ok || err != nil || gerr != nil || eerr != nil || !counted || *outPath == "" {
		fmt.Fprintf(out, "UNMEASURED engine=%s found=%v probes=%v goldens=%v equiv=%v counts=%v out=%q\n", *engine, ok, err, gerr, eerr, lines, *outPath)
		return 2
	}
	for _, g := range gs {
		if g.Stale() {
			fmt.Fprintln(out, "UNMEASURED STALE", g.ID)
			return 2
		}
	}
	id := summary.Summary{Engine: *engine, Ref: *ref, Equiv: *mode, Pin: scoring.Pin(*engine)}
	return scoring.Score(out, id, *root, *outPath, gs, f, pr)
}
