package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/summary"
)

// TestSplits: skips split by complete; SPLITS counts the SPLIT cases, not the FAILs (3 here; whole-pass adds every
// case with a variant to PASS); DIAG counts the FAILs whose first path is unobserved (rulings-8), not every FAIL.
func TestSplits(t *testing.T) {
	s := summary.Summary{Total: 6, Passed: 1, Skipped: []summary.Skip{{Complete: true}, {Complete: false}},
		Cases: []summary.Case{{Result: "PASS"}, {Result: "SPLIT", Variant: "s0k1"}, {Result: "SPLIT", Variant: "bytes"},
			{Result: "FAIL", Detail: "unobserved step[0].P.attrs"}, {Result: "FAIL", Detail: "step[0].P.cols: want 80 got 0"},
			{Result: "FAIL", Detail: "step[0].P.rows: want 24 got 0"}}}
	var out bytes.Buffer
	Splits(&out, "e", "x53", s)
	for _, want := range []string{"SKIP e x53 pass=1 nonpass>=1\n", "SPLITS e x53 whole-pass=3 split-fail=2\n", "DIAG e x53 unobserved-first 1/6\n"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("no %q in\n%s", want, out.String())
		}
	}
	out.Reset()
	Splits(&out, "e", "x53", summary.Summary{Total: 3, Passed: 1, Cases: []summary.Case{{Result: "PASS"},
		{Result: "FAIL", Detail: "unobserved step[0].P.lineFeeds"}, {Result: "FAIL", Detail: "step[0].P.cols: want 80 got 0"}}})
	if !strings.Contains(out.String(), "DIAG e x53 unobserved-first 1/3\n") {
		t.Errorf("rulings-8 vector:\n%s", out.String())
	}
}
