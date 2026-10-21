package runner

import (
	"fmt"
	"path/filepath"
	"sort"
)

// RequiredCategories maps every category issue #1 requires to its minimum case count; MinTotal is the overall bar.
// The -count CLI and TestCounts both use these (one list, S9).
var RequiredCategories = map[string]int{
	"c0": 22, "esc": 25, "csi": 60, "sgr": 40, "osc": 20, "dcs": 11, "decmodes": 40, "altscreen": 20, "scrollback": 20, "wide": 30, "resize": 26,
}

const MinTotal = 300

// Counts returns the number of case files per category under root/cases (split variants are not cases).
func Counts(root string) (map[string]int, error) {
	files, err := filepath.Glob(filepath.Join(root, "cases", "*", "*.json"))
	if err != nil {
		return nil, err
	}
	out := map[string]int{}
	for _, f := range files {
		out[filepath.Base(filepath.Dir(f))]++
	}
	return out, nil
}

// CheckCounts compares the case counts, and the golden counts of every ref, with the required minimums.
// A category with no cases, or an unreadable directory, is UNMEASURED (never a pass).
func CheckCounts(root string, refs []string) (lines []string, ok bool) {
	cases, err := Counts(root)
	if err != nil || len(cases) == 0 {
		return []string{fmt.Sprintf("cases UNMEASURED (%v)", err)}, false
	}
	ok = true
	total := 0
	var cats []string
	for c := range RequiredCategories {
		cats = append(cats, c)
	}
	sort.Strings(cats)
	for _, c := range cats {
		n, min := cases[c], RequiredCategories[c]
		total += n
		status := "ok"
		switch {
		case n == 0:
			status, ok = "UNMEASURED", false
		case n < min:
			status, ok = fmt.Sprintf("FAIL %d < %d", n, min), false
		}
		for _, ref := range refs {
			gs, _ := filepath.Glob(filepath.Join(root, "golden", ref, c, "*.json"))
			if len(gs) != n {
				status, ok = fmt.Sprintf("%s; %s goldens %d != cases %d", status, ref, len(gs), n), false
			}
		}
		lines = append(lines, fmt.Sprintf("%s %d %s", c, n, status))
	}
	for c := range cases {
		if _, known := RequiredCategories[c]; !known {
			lines, ok = append(lines, fmt.Sprintf("%s %d FAIL unknown category", c, cases[c])), false
		}
	}
	verdict := "ok"
	if total < MinTotal {
		verdict, ok = fmt.Sprintf("FAIL < %d", MinTotal), false
	}
	return append(lines, fmt.Sprintf("total %d >= %d %s", total, MinTotal, verdict)), ok
}
