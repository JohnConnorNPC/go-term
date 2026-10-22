package magnet

import (
	"testing"

	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/equiv"
)

// read returns modes and opts after writing in with the equiv table in mode.
func read(t *testing.T, mode, in string) (modes, opts map[string]any) {
	t.Helper()
	if err := equiv.Load("../equiv/table.json", mode); err != nil {
		t.Fatal(err)
	}
	st := snap(t, engine(t, 20, 3, "6"), in)
	return st.P.Modes, st.P.Opts
}

// TestModes: bools by identity; the mouse mode through M1-M5 only, magnet's raw spelling in strict.
func TestModes(t *testing.T) {
	for _, c := range []struct {
		mode, in string
		want     any
	}{{"full", "", "none"}, {"strict", "", "NONE"}, {"case", "", "none"}, {"default", "", "NONE"},
		{"full", "\x1b[?1000h", "vt200"}, {"strict", "\x1b[?1000h", "VT200"}, {"case", "\x1b[?1000h", "vt200"}} {
		if m, _ := read(t, c.mode, c.in); m["mouseTrackingMode"] != c.want {
			t.Errorf("%s %q: mouseTrackingMode %v, want %v", c.mode, c.in, m["mouseTrackingMode"], c.want)
		}
	}
	if m, _ := read(t, "strict", "\x1b[4h\x1b[?6h"); m["insertMode"] != true || m["originMode"] != true || m["wraparoundMode"] != true {
		t.Errorf("modes %v, want insertMode, originMode and wraparoundMode true", m)
	}
}
