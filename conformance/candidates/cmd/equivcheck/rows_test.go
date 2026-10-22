package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// srcs is the fixture module: six assignments in c.go (round-atoms-2 ruling 1: a sixth one must fail a count of 5).
var srcs = map[string]string{
	"a.go": "package m\n\tAttrCMP16      uint32 = 0x1000000\n\tAttrCMP256     uint32 = 0x2000001\n",
	"c.go": "package m\n" + strings.Repeat("\th.MouseTrackingMode = \"X\"\n", 6),
}

// table writes a one-file table: rows OK (a line cite) and SIX (a count cite of 5 over c.go).
func table(t *testing.T) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "table.json")
	body := `{"rows":[
 {"id":"OK","engine":"e","field":"f","key":"","golden":"none","class":"default","cite":[{"module":"example.com/m","file":"a.go","line":"AttrCMP16      uint32 = 0x1000000"}]},
 {"id":"SIX","engine":"e","field":"f","key":"X","golden":"x","class":"case","cite":[{"module":"example.com/m","files":"c.go","pattern":"MouseTrackingMode *= ","count":5}]}]}`
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRows(t *testing.T) {
	var out bytes.Buffer
	ok, n := checkRows(zipsFor(fixture(t, srcs)), table(t), &out)
	for _, want := range []string{`EQUIV e OK "" => none @ a.go:2`, `EQUIV FAIL SIX`, `got=6 want=5`, `EQUIV e SIX "X" => x @ count(MouseTrackingMode *= )=6`} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("output lacks %q:\n%s", want, out.String())
		}
	}
	if ok != 1 || n != 2 {
		t.Errorf("ok=%d n=%d, want 1/2", ok, n)
	}
	if ok, n := checkRows(zipsFor(fixture(t, srcs)), "nope.json", &out); ok != 0 || n != 1 {
		t.Errorf("unreadable table: ok=%d n=%d, want 0/1", ok, n)
	}
}
