package runner

import (
	"sort"
	"strconv"
	"strings"
	"testing"
)

// leafKind normalises a leaf path: array indices become *, except tuple positions, which stay distinct
// (cursor.0/1, colour mode/value fg.0/fg.1/bg.0/bg.1, RLE run fields c.*.0-3); free-named keys become *.
func leafKind(path string) string {
	seg := strings.Split(path, ".")
	for i, s := range seg {
		if _, err := strconv.Atoi(s); err == nil && !(i > 0 && (seg[i-1] == "cursor" || seg[i-1] == "fg" || seg[i-1] == "bg") || i > 1 && seg[i-2] == "c") {
			seg[i] = "*"
		}
	}
	if m := freeKeys.FindStringSubmatch(strings.Join(seg, ".")); m != nil {
		return m[1] + ".*"
	}
	return strings.Join(seg, ".")
}

// walk visits every non-null scalar leaf in key order (deterministic: Go map order is random).
func walk(v any, path string, f func(path string)) {
	switch t := v.(type) {
	case map[string]any:
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			walk(t[k], strings.TrimPrefix(path+"."+k, "."), f)
		}
	case []any:
		for i, x := range t {
			walk(x, path+"."+strconv.Itoa(i), f)
		}
	case nil:
	default:
		f(path)
	}
}

func TestLeafKind(t *testing.T) {
	for in, want := range map[string]string{"steps.3.P.cursor.1": "steps.*.P.cursor.1", "steps.0.attrs.2.P.fg.0": "steps.*.attrs.*.P.fg.0",
		"steps.1.cells.normal.4.c.2.3": "steps.*.cells.normal.*.c.*.3", "steps.0.P.modes.insertMode": "steps.*.P.modes.*", "B.t2.7": "B.t2.*"} {
		if got := leafKind(in); got != want {
			t.Errorf("leafKind(%s) = %s, want %s", in, got, want)
		}
	}
}
