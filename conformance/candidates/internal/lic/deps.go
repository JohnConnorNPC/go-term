package lic

import (
	"io/fs"
	"sort"
)

var mitIDs = map[string]bool{"L1": true, "L2": true, "L3": true, "L4": true}

// Deps is L9: each module the candidate's non-test packages import (mods, by path) must pass L1-L4 with its own root
// LICENSE: OK, else FAIL when gated (gitpod, the fork candidate) and INFO for the others. The candidate's test-only
// requires are listed as INFO (e.g. github.com/google/go-cmp, BSD-3-Clause).
func Deps(gated bool, mods map[string]fs.FS, testOnly []string) []Line {
	var paths []string
	for p := range mods {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	var out []Line
	for _, p := range paths {
		var bad []string
		for _, l := range Root(mods[p]) {
			if mitIDs[l.ID] && l.Status != "OK" {
				bad = append(bad, l.ID)
			}
		}
		status := "OK"
		if len(bad) > 0 && gated {
			status = "FAIL"
		} else if len(bad) > 0 {
			status = "INFO"
		}
		out = append(out, Line{"L9", status, p + " not-mit=" + list(bad)})
	}
	if len(paths) == 0 {
		out = append(out, Line{"L9", "OK", "no non-test dependencies"})
	}
	return append(out, Line{"L9", "INFO", "test-only " + list(testOnly)})
}
