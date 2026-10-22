package scoring

import (
	"fmt"
	"io"

	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/summary"
	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

const skipNote = "SKIP-NOTE SKIP-SENSITIVE lists are complete for PASS cases and a lower bound for SPLIT-FAIL cases"

// Score runs every golden (whole writes, probe twins and every split variant), prints one line per case and one
// SKIP-SENSITIVE line per skipped variant (complete for a PASS case, else a lower bound), then writes the summary.
func Score(out io.Writer, id summary.Summary, root, path string, gs []*runner.Golden, f runner.Factory, pr *runner.Probes) int {
	var rs []runner.Result
	lower := false
	for _, g := range gs {
		r := runner.RunCase(g, f, pr, runner.Options{Splits: true})
		lower = report(out, r) || lower
		rs = append(rs, r)
	}
	if lower {
		fmt.Fprintln(out, skipNote)
	}
	return Write(out, id, root, path, rs)
}

// report prints r's line and its SKIP-SENSITIVE lines; true when they are a lower bound (r did not pass).
func report(out io.Writer, r runner.Result) bool {
	fmt.Fprintln(out, r.Status, r.ID, r.Diff)
	for _, v := range r.Skipped {
		kind := "complete"
		if r.Status != runner.PASS {
			kind = "lowerbound"
		}
		fmt.Fprintln(out, "SKIP-SENSITIVE", r.ID, v, kind)
	}
	return r.Status != runner.PASS && len(r.Skipped) > 0
}
