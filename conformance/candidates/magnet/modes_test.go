package magnet

import (
	"maps"
	"strings"
	"testing"

	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/equiv"
	"github.com/JohnConnorNPC/go-term/conformance/runner"
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

// TestOpts: the cursor options through M6-M9; strict reports the raw override (nil when unset); DECSCUSR's override
// is reported in every mode; ?12h writes magnet's Options, read through M8, and a later DECSCUSR overrides it (M9).
func TestOpts(t *testing.T) {
	for _, c := range []struct {
		mode, in     string
		style, blink any
	}{{"full", "", "block", false}, {"strict", "", nil, nil}, {"full", "\x1b[4 q", "underline", false},
		{"strict", "\x1b[4 q", "underline", false}, {"case", "\x1b[4 q", "underline", false}, {"full", "\x1b[?12h", "block", true},
		{"full", "\x1b[?12h\x1b[2 q", "block", false}} {
		if _, o := read(t, c.mode, c.in); o["cursorStyle"] != c.style || o["cursorBlink"] != c.blink {
			t.Errorf("%s %q: opts %v, want cursorStyle %v cursorBlink %v", c.mode, c.in, o, c.style, c.blink)
		}
	}
}

// TestStrictNeverUnobserved (rulings-6): strict has no table, so no path is unobserved, and a golden that differs from
// the raw value is a FAIL that shows both raw values.
func TestStrictNeverUnobserved(t *testing.T) {
	if err := equiv.Load("../equiv/table.json", "strict"); err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		in, path string
		golden   any
		diff     string
	}{{"", "modes.mouseTrackingMode", "none", `step[0].P.modes.mouseTrackingMode: want "none" got "NONE"`},
		{"\x1b[?1000h", "modes.mouseTrackingMode", "vt200", `step[0].P.modes.mouseTrackingMode: want "vt200" got "VT200"`},
		{"", "opts.cursorStyle", "block", `step[0].P.opts.cursorStyle: want "block" got null`},
		{"", "opts.cursorBlink", false, `step[0].P.opts.cursorBlink: want false got null`}} {
		st := snap(t, engine(t, 20, 3, "6"), c.in)
		want := st
		want.P.Modes, want.P.Opts = maps.Clone(st.P.Modes), maps.Clone(st.P.Opts)
		name, key, _ := strings.Cut(c.path, ".")
		m := want.P.Modes
		if name == "opts" {
			m = want.P.Opts
		}
		m[key] = c.golden
		if d := runner.Compare("step[0]", want, st); len(st.Unobserved) != 0 || d != c.diff {
			t.Errorf("strict %q %s: unobserved %v, diff %q, want %q", c.in, c.path, st.Unobserved, d, c.diff)
		}
	}
}
