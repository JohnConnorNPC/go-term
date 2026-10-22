package scoring

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/summary"
	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

// TestScore: every golden of a ref is run with its split variants; the oracle passes them all, with wide/zwsp's
// sensitive variants reported as complete; the stub fails every case at its first step, so it has no skips.
func TestScore(t *testing.T) {
	pr, err := runner.LoadProbes("../../../probes.json")
	gs, err2 := runner.LoadAll("../../..", "x53")
	if err != nil || err2 != nil {
		t.Fatal(err, err2)
	}
	for name, want := range map[string]int{"oracle": 0, "stub": 1} {
		f, ok := Lookup(name, pr)
		if !ok {
			t.Fatalf("%s not found", name)
		}
		var out bytes.Buffer
		code := Score(&out, summary.Summary{Engine: name}, "../../..", filepath.Join(t.TempDir(), "s.json"), gs, f, pr)
		skips := strings.Contains(out.String(), "SKIP-SENSITIVE wide/zwsp ")
		if code != want || skips != (name == "oracle") || strings.Contains(out.String(), "SKIP-NOTE") {
			t.Errorf("%s: exit %d skips %v, want exit %d", name, code, skips, want)
		}
	}
}

// TestReport: skips of a PASS case are complete; a SPLIT case stops at its first failing variant, so its skips are a
// lower bound (rulings-5 #5).
func TestReport(t *testing.T) {
	for st, want := range map[runner.Status]string{runner.PASS: "complete", runner.SPLIT: "lowerbound"} {
		var out bytes.Buffer
		lower := report(&out, runner.Result{ID: "wide/x", Status: st, Skipped: []string{"s0k1"}})
		if !strings.Contains(out.String(), "SKIP-SENSITIVE wide/x s0k1 "+want+"\n") || lower != (st == runner.SPLIT) {
			t.Errorf("%s: %q lower=%v", st, out.String(), lower)
		}
	}
}
