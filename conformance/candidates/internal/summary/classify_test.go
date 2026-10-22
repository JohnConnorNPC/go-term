package summary

import (
	"testing"

	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

// TestClassify (ruling B): one vector per error prefix, an error inside a split, a split mismatch, and the diffs
// that stay FAIL (a plain diff, an unknown prefix, an unobserved path).
func TestClassify(t *testing.T) {
	F, S := runner.FAIL, runner.SPLIT
	for _, c := range []struct {
		status         runner.Status
		diff           string
		class, variant string
	}{
		{F, "step[2]: engine error: x", "ERROR", ""},
		{F, "step[0]: snapshot error: x", "ERROR", ""},
		{F, "engine: x", "ERROR", ""},
		{F, "B: x", "ERROR", ""},
		{S, "v1: step[0]: engine error: x", "ERROR", "v1"},
		{S, "v1: step[0].P.cursor.x: want 1 got 2", "SPLIT", "v1"},
		{F, "step[0].P.cursor.x: want 1 got 2", "FAIL", ""},
		{F, "boom: x", "FAIL", ""},
		{F, "unobserved B.t1.?1: replies", "FAIL", ""},
		{S, "no variant prefix", "FAIL", ""},
		{runner.PASS, "", "PASS", ""},
	} {
		if class, v := Classify(runner.Result{Status: c.status, Diff: c.diff}); class != c.class || v != c.variant {
			t.Errorf("%s %q: %s %q, want %s %q", c.status, c.diff, class, v, c.class, c.variant)
		}
	}
}
