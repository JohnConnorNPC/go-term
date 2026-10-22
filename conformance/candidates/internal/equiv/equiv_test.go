package equiv

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRead: the shipped table parses (19 rows, E10); an unknown key, an unknown class or an accessor outside the
// closed set is a load error.
func TestRead(t *testing.T) {
	rows, err := Read("../../equiv/table.json")
	if err != nil || len(rows) != 19 {
		t.Fatalf("table.json: %d rows, %v; want 19", len(rows), err)
	}
	for name, body := range map[string]string{
		"unknown key":      `{"rows":[{"id":"X","engine":"e","field":"f","key":"","golden":"none","class":"case","cite":[],"note":1}]}`,
		"unknown class":    `{"rows":[{"id":"X","engine":"e","field":"f","key":"","golden":"none","class":"rule","cite":[]}]}`,
		"unknown accessor": `{"rows":[{"id":"X","engine":"e","field":"f","key":"","golden":"@Options.Foo","class":"source","cite":[]}]}`,
	} {
		p := filepath.Join(t.TempDir(), "t.json")
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := Read(p); err == nil || !strings.Contains(err.Error(), "t.json") {
			t.Errorf("%s: err %v, want a load error", name, err)
		}
	}
}

// TestModes: strict applies no row; case and default keep only their own class; full keeps every row (the gate).
func TestModes(t *testing.T) {
	for _, c := range []struct {
		mode, key string
		want      any
		ok        bool
	}{{"full", "", "none", true}, {"strict", "", nil, false}, {"case", "NONE", "none", true}, {"case", "", nil, false},
		{"default", "", "none", true}, {"default", "NONE", nil, false}} {
		if err := Load("../../equiv/table.json", c.mode); err != nil {
			t.Fatal(err)
		}
		if got, ok := Map("gitpod", "modes.mouseTrackingMode", c.key); got != c.want || ok != c.ok {
			t.Errorf("%s %q: %v %v, want %v %v", c.mode, c.key, got, ok, c.want, c.ok)
		}
	}
	if err := Load("../../equiv/table.json", "loose"); err == nil {
		t.Error("mode loose loaded")
	}
}
