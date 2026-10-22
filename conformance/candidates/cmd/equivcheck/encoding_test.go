package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestEncoding (GO S-1): a cited constant must equal norm's own constant; 0x2000001 is not norm.CMP256, and an unknown
// constant name fails too.
func TestEncoding(t *testing.T) {
	p := filepath.Join(t.TempDir(), "encoding.json")
	body := `{"entries":[
 {"id":"e-CMP16","engine":"e","const":"CMP16","cite":{"module":"example.com/m","file":"a.go","line":"AttrCMP16      uint32 = 0x1000000"}},
 {"id":"e-CMP256","engine":"e","const":"CMP256","cite":{"module":"example.com/m","file":"a.go","line":"AttrCMP256     uint32 = 0x2000001"}},
 {"id":"e-X","engine":"e","const":"CMX","cite":{"module":"example.com/m","file":"a.go","line":"AttrCMP16      uint32 = 0x1000000"}}]}`
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	ok, n := checkEncoding(zipsFor(fixture(t, srcs)), p, &out)
	for _, want := range []string{"ENC e e-CMP16 @ a.go:2 0x1000000 == norm.CMP16", "ENC FAIL e e-CMP256 got=0x2000001 norm=0x2000000", "ENC FAIL e e-X"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("output lacks %q:\n%s", want, out.String())
		}
	}
	if ok != 1 || n != 3 {
		t.Errorf("ok=%d n=%d, want 1/3", ok, n)
	}
}
