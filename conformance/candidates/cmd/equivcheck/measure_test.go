package main

import "testing"

func zipsFor(proxy, gomod string) *zips {
	return &zips{proxy: proxy, gomod: gomod, mods: map[string]map[string][]string{}}
}

// TestMeasure: a line citation matches trimmed whole lines (a tab-indented source line included) and counts them; a
// count citation counts pattern lines over the non-test .go files or one file; anything unreadable is -1.
func TestMeasure(t *testing.T) {
	z := zipsFor(fixture(t, map[string]string{
		"a.go":      "package m\n\tAttrCMP16      uint32 = 0x1000000\ny = 2\ny = 2\n",
		"b.go":      "package m\ny = 3\n",
		"a_test.go": "package m\ny = 4\n",
		"README.md": "y = 5\n",
	}))
	five := 5
	for _, c := range []struct {
		c         cite
		got, want int
	}{
		{cite{Module: "example.com/m", File: "a.go", Line: "AttrCMP16      uint32 = 0x1000000"}, 1, 1},
		{cite{Module: "example.com/m", File: "a.go", Line: "AttrCMP16 uint32 = 0x1000000"}, 0, 1},
		{cite{Module: "example.com/m", File: "a.go", Line: "y = 2"}, 2, 1},
		{cite{Module: "example.com/m", Files: "*", Pattern: "^y = ", Count: &five}, 3, 5},
		{cite{Module: "example.com/m", Files: "a.go", Pattern: "^y = ", Count: &five}, 2, 5},
		{cite{Module: "example.com/m", Files: "nope.go", Pattern: "^y = ", Count: &five}, -1, 5},
		{cite{Module: "example.com/m", Files: "*", Pattern: "(", Count: &five}, -1, 5},
		{cite{Module: "example.com/gone", File: "a.go", Line: "y = 2"}, -1, 1},
	} {
		if where, got, want := z.measure(c.c); got != c.got || want != c.want {
			t.Errorf("%+v: got=%d want=%d (%s), expected got=%d want=%d", c.c, got, want, where, c.got, c.want)
		}
	}
	if where, _, _ := z.measure(cite{Module: "example.com/m", File: "a.go", Line: "y = 2"}); where != "a.go:4" {
		t.Errorf("where %s, want the last matching line a.go:4", where)
	}
}
