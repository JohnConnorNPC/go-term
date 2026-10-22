package count

import (
	"strings"
	"testing"

	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

// TestCheck: one result per case file passes; one dropped or one extra names its category.
func TestCheck(t *testing.T) {
	want, err := runner.Counts("../../..")
	if err != nil {
		t.Fatal(err)
	}
	var rs []runner.Result
	for c, n := range want {
		for i := 0; i < n; i++ {
			rs = append(rs, runner.Result{Category: c})
		}
	}
	if err := Check("../../..", rs); err != nil {
		t.Fatalf("all cases: %v", err)
	}
	extra := append(append([]runner.Result{}, rs...), runner.Result{Category: "c0"})
	for name, c := range map[string][]runner.Result{"dropped": dropOne(rs, "esc"), "extra": extra} {
		if err := Check("../../..", c); err == nil || !strings.Contains(err.Error(), "UNMEASURED count") {
			t.Errorf("%s: %v, want an UNMEASURED count error", name, err)
		}
	}
	if err := Check("../../..", dropOne(rs, "esc")); err == nil || !strings.Contains(err.Error(), " esc ") {
		t.Errorf("dropped esc: %v, want it named", err)
	}
}

func dropOne(rs []runner.Result, cat string) []runner.Result {
	for i, r := range rs {
		if r.Category == cat {
			return append(append([]runner.Result{}, rs[:i]...), rs[i+1:]...)
		}
	}
	return rs
}
