package goterm

import (
	"sync/atomic"
	"testing"

	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

// TestGateXcur is go-term's conformance gate (docs/decisions/0003-engine.md): every case against the
// @xterm/headless@6.0.0 goldens (xcur), its whole write, its probe twins and every chunk-split variant, must PASS.
// The denominator is the case files: the goldens loaded must number exactly the cases counted on disk, so a missing
// golden or case directory is a FAIL, never a smaller total. Under -race (raceOn) the split variants run for every
// 8th case only, to keep the race step's time sane; the plain test step always runs them all.
func TestGateXcur(t *testing.T) {
	const root = ".."
	if lines, ok := runner.CheckCounts(root, []string{"xcur"}); !ok {
		t.Fatalf("case counts: %v", lines)
	}
	counts, err := runner.Counts(root)
	if err != nil {
		t.Fatal(err)
	}
	cases := 0
	for _, n := range counts {
		cases += n
	}
	pr, err := runner.LoadProbes(root + "/probes.json")
	if err != nil {
		t.Fatal(err)
	}
	gs, err := runner.LoadAll(root, "xcur")
	if err != nil || len(gs) != cases || cases == 0 {
		t.Fatalf("xcur goldens %d, cases %d, err %v", len(gs), cases, err)
	}
	var pass atomic.Int64
	t.Run("cases", func(t *testing.T) {
		for i, g := range gs {
			splits := !raceOn || i%8 == 0
			t.Run(g.ID, func(t *testing.T) {
				t.Parallel()
				r := runner.RunCase(g, New, pr, runner.Options{Splits: splits})
				if r.Status != runner.PASS {
					t.Errorf("%s %s: %s", r.Status, r.ID, r.Diff)
					return
				}
				pass.Add(1)
			})
		}
	})
	t.Logf("GATE goterm xcur %d/%d", pass.Load(), cases)
	if int(pass.Load()) != cases {
		t.Errorf("GATE goterm xcur %d/%d, want every case", pass.Load(), cases)
	}
}
