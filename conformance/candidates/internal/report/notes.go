package report

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

// XcurNote is `XCUR-NOTE opts <n> <ids…>`: the cases whose xcur golden differs from x53 in an opts.* path
// (golden/diff.json; round-atoms-2 ruling 5), sorted.
func XcurNote(path string) (string, error) {
	var d struct {
		Cases []struct {
			ID string
			P  []string
		}
	}
	b, err := os.ReadFile(path)
	if err == nil {
		err = json.Unmarshal(b, &d)
	}
	if err != nil || len(d.Cases) == 0 {
		return "", fmt.Errorf("UNMEASURED diff %s: %v", path, err)
	}
	var ids []string
	for _, c := range d.Cases {
		for _, p := range c.P {
			if strings.Contains(p, ".P.opts.") {
				ids = append(ids, c.ID)
				break
			}
		}
	}
	sort.Strings(ids)
	return strings.TrimSpace(fmt.Sprintf("XCUR-NOTE opts %d %s", len(ids), strings.Join(ids, " "))), nil
}
