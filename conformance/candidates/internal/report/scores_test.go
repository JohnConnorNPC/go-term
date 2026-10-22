package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/summary"
)

// TestScores: one SCORE line per run (two-decimal truncated percent), the 6 strict ones labelled, and its CAT lines,
// a category under 50 % flagged.
func TestScores(t *testing.T) {
	runs, err := load(t, func(r Run, s *summary.Summary) {
		if r == (Run{"gitpod", "x53", "full"}) {
			s.Passed, s.Categories = 301, map[string]summary.Cat{"c0": {Total: 100, Passed: 49}, "esc": {Total: 214, Passed: 214}}
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	Scores(&out, runs)
	for _, want := range []string{"SCORE gitpod x53 full 301/314 95.85\n", "CAT gitpod x53 full c0 49/100 49.00 FLAG<50\n",
		"CAT gitpod x53 full esc 214/214 100.00\n", "SCORE xvt x53 default 30/314 9.55\n",
		"SCORE magnet xcur strict 30/314 9.55 strict=raw-encoding (mapping cost, not a quality score)\n"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("no %q in\n%s", want, out.String())
		}
	}
	if n, m := strings.Count(out.String(), "SCORE "), strings.Count(out.String(), "(mapping cost, not a quality score)"); n != 18 || m != 6 {
		t.Errorf("%d SCORE lines, %d labelled, want 18 and the 6 strict ones", n, m)
	}
}
