package summary

import (
	"testing"

	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

// TestBuild: totals by case, a SPLIT and an ERROR counted as not passed, skips flagged complete only on PASS cases,
// and a required category with no results is an error, never 0/0.
func TestBuild(t *testing.T) {
	rs := []runner.Result{
		{ID: "c0/b", Category: "c0", Status: runner.SPLIT, Diff: "s0k1: step[0].P.cursor.x: want 1 got 2", Skipped: []string{"bytes"}},
		{ID: "c0/a", Category: "c0", Status: runner.PASS, Skipped: []string{"s0k2"}},
		{ID: "esc/a", Category: "esc", Status: runner.FAIL, Diff: "step[0]: engine error: panic: boom"},
	}
	s, err := Build(Summary{}, rs, []string{"c0", "esc"})
	if err != nil || s.Total != 3 || s.Passed != 1 || s.Categories["c0"] != (Cat{Passed: 1, Total: 2}) {
		t.Fatalf("totals %d/%d %v %v", s.Passed, s.Total, s.Categories, err)
	}
	if s.Cases[0].ID != "c0/a" || s.Cases[1].Result != "SPLIT" || s.Cases[1].Variant != "s0k1" || s.Cases[2].Result != "ERROR" ||
		s.Cases[2].Detail != "step[0]: engine error: panic: boom" {
		t.Errorf("cases %+v", s.Cases)
	}
	if len(s.Skipped) != 2 || !s.Skipped[0].Complete || s.Skipped[1].Complete {
		t.Errorf("skipped %+v, want c0/a complete, c0/b a lower bound", s.Skipped)
	}
	if _, err := Build(Summary{}, rs, []string{"c0", "esc", "osc"}); err == nil {
		t.Error("category osc without results was accepted")
	}
}
