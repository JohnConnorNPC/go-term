package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCheck: the total line counts rows and encoding entries; a failing encoding entry fails the whole check.
func TestCheck(t *testing.T) {
	proxy, gomod := fixture(t, srcs)
	dir := t.TempDir()
	tbl, enc := filepath.Join(dir, "t.json"), filepath.Join(dir, "e.json")
	row := `{"rows":[{"id":"OK","engine":"e","field":"f","key":"","golden":"none","class":"default","cite":[{"module":"example.com/m","file":"a.go","line":"AttrCMP16      uint32 = 0x1000000"}]}]}`
	for name, entry := range map[string]string{"OK 2/2": `"CMP16"`, "FAIL 1/2": `"CMP256"`} {
		e := `{"entries":[{"id":"x","engine":"e","const":` + entry + `,"cite":{"module":"example.com/m","file":"a.go","line":"AttrCMP16      uint32 = 0x1000000"}}]}`
		if os.WriteFile(tbl, []byte(row), 0o644) != nil || os.WriteFile(enc, []byte(e), 0o644) != nil {
			t.Fatal("fixture write failed")
		}
		var out bytes.Buffer
		if got := check(proxy, gomod, tbl, enc, &out); got != (name == "OK 2/2") || !strings.Contains(out.String(), "EQUIV "+name) {
			t.Errorf("%s: check %v\n%s", name, got, out.String())
		}
	}
}
