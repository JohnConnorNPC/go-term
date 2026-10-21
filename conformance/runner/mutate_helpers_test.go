package runner

import (
	"encoding/json"
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

func at(v any, s string) any {
	if m, ok := v.(map[string]any); ok {
		return m[s]
	}
	i, _ := strconv.Atoi(s)
	return v.([]any)[i]
}

func put(v any, s string, x any) {
	if m, ok := v.(map[string]any); ok {
		m[s] = x
		return
	}
	i, _ := strconv.Atoi(s)
	v.([]any)[i] = x
}

func num(v any) int { n, _ := v.(json.Number).Int64(); return int(n) }

// mutateAt changes the leaf at path so that the golden still parses and means something else. It returns false when this
// occurrence has no such change (a run with chars cannot grow; a run index needs another attr entry with a different P).
func mutateAt(doc map[string]any, path string) bool {
	seg := strings.Split(path, ".")
	var parent any = doc
	for _, s := range seg[:len(seg)-1] {
		parent = at(parent, s)
	}
	last, v := seg[len(seg)-1], at(parent, seg[len(seg)-1])
	run := len(seg) > 2 && seg[len(seg)-3] == "c"
	var nv any
	switch {
	case run && last == "0":
		if at(parent, "3") != "" {
			return false
		}
		nv = json.Number(strconv.Itoa(num(v) + 1))
	case run && last == "1":
		attrs := at(at(at(doc, "steps"), seg[1]), "attrs").([]any)
		for j := range attrs {
			if show(at(attrs[j], "P")) != show(at(attrs[num(v)], "P")) {
				nv = json.Number(strconv.Itoa(j))
				break
			}
		}
	case run && last == "3":
		t, r := v.(string), []rune(v.(string))
		switch {
		case t == "" && num(at(parent, "0")) == 1:
			nv = "Z"
		case t != "" && r[len(r)-1] != 'Z':
			r[len(r)-1] = 'Z'
			nv = string(r)
		case t != "":
			r[len(r)-1] = 'Y'
			nv = string(r)
		}
	default:
		switch t := v.(type) {
		case bool:
			nv = !t
		case string:
			nv = t + "Z"
		case json.Number:
			nv = json.Number(strconv.Itoa(num(t) + 1))
		}
	}
	if nv == nil {
		return false
	}
	put(parent, last, nv)
	return true
}

func TestLeafKind(t *testing.T) {
	for in, want := range map[string]string{"steps.3.P.cursor.1": "steps.*.P.cursor.1", "steps.0.attrs.2.P.fg.0": "steps.*.attrs.*.P.fg.0",
		"steps.1.cells.normal.4.c.2.3": "steps.*.cells.normal.*.c.*.3", "steps.0.P.modes.insertMode": "steps.*.P.modes.*", "B.t2.7": "B.t2.*"} {
		if got := leafKind(in); got != want {
			t.Errorf("leafKind(%s) = %s, want %s", in, got, want)
		}
	}
}
