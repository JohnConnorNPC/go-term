// Package count is the denominator check (E5): a run's results must cover exactly the case files, per category.
package count

import (
	"fmt"
	"sort"

	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

// Check returns nil when every category's result count, and so the total, equals its case files under root/cases
// (runner.Counts, #1's counter); else the UNMEASURED line naming the first category that differs.
func Check(root string, rs []runner.Result) error {
	want, err := runner.Counts(root)
	if err != nil || len(want) == 0 {
		return fmt.Errorf("UNMEASURED count cases: %v", err)
	}
	got := map[string]int{}
	for _, r := range rs {
		got[r.Category]++
	}
	var cats []string
	for c := range want {
		cats = append(cats, c)
	}
	for c := range got {
		if _, ok := want[c]; !ok {
			cats = append(cats, c)
		}
	}
	sort.Strings(cats)
	for _, c := range cats {
		if got[c] != want[c] {
			return fmt.Errorf("UNMEASURED count %s %d!=%d", c, got[c], want[c])
		}
	}
	return nil
}
