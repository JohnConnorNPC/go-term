// Command verdict prints the E7 evidence lines of the score summaries, in a fixed order, then the E8 verdict:
//
//	go run ./cmd/verdict > scores/verdict.txt
//
// Exit 2 (UNMEASURED) when a summary, golden/diff.json or a fact file is missing, or a summary is not its run's (engine,
// ref, mode, total, pin: report.Load); no verdict is printed then.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/report"
)

const skipNote = "SKIP-NOTE SKIP-SENSITIVE lists are complete for PASS cases and a lower bound for SPLIT-FAIL cases"

const diagNote = "DIAG-NOTE n cases fail first on a path the engine cannot observe; their later paths were not compared"

func main() { os.Exit(run(os.Args[1:], os.Stdout)) }

func run(args []string, out io.Writer) int {
	fs := flag.NewFlagSet("verdict", flag.ContinueOnError)
	scores := fs.String("scores", "scores", "the summaries and lic.txt, pin.txt")
	diff := fs.String("diff", "../golden/diff.json", "the x53-vs-xcur golden diff")
	root := fs.String("root", "..", "the conformance directory: a summary's total must equal its ref's goldens there")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	runs, err := report.Load(*scores, *root)
	note, err2 := report.XcurNote(*diff)
	var verdict string
	err3 := err
	if err == nil {
		verdict, err3 = report.Rule(runs, *scores+"/lic.txt", *scores+"/pin.txt")
	}
	for _, e := range []error{err, err2, err3} {
		if e != nil {
			fmt.Fprintln(out, e)
			return 2
		}
	}
	report.Scores(out, runs)
	report.PerEngine(out, runs)
	fmt.Fprintln(out, skipNote)
	fmt.Fprintln(out, diagNote)
	fmt.Fprintln(out, note)
	fmt.Fprintln(out, verdict)
	return 0
}
