package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/summary"
)

// TestPerEngine: ABLATE reads the four x53 modes (default is not case), EQUIV-DELTA is full - strict, and GATE reads
// x53 full only (a strict or case score of 314 doesn't pass it).
func TestPerEngine(t *testing.T) {
	passed := map[string]int{"strict": 10, "case": 20, "default": 15, "full": 30}
	runs, err := load(t, func(r Run, s *summary.Summary) {
		s.Passed = passed[r.Mode]
		if r.Engine == "magnet" && r.Mode != "full" {
			s.Passed = 314
		}
		if r.Engine == "gitpod" && r.Mode == "full" {
			s.Passed = 283
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	PerEngine(&out, runs)
	for _, want := range []string{"ABLATE xvt x53 strict=10 case=20 default=15 full=30\n", "EQUIV-DELTA xvt xcur 20\n",
		"GATE gitpod x53 283/314 >=90% PASS\n", "GATE magnet x53 30/314 >=90% FAIL\n", "GATE xvt x53 30/314 >=90% FAIL\n"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("no %q in\n%s", want, out.String())
		}
	}
}
