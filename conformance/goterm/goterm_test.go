package goterm

import (
	"testing"

	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

func snap(t *testing.T, in string) runner.State {
	t.Helper()
	e, err := New(runner.Config{Cols: 20, Rows: 3, Scrollback: 100, Unicode: "6"})
	if err != nil {
		t.Fatal(err)
	}
	if err := e.Write([]byte(in)); err != nil {
		t.Fatal(err)
	}
	st, err := e.Snapshot()
	if err != nil {
		t.Fatal(err)
	}
	return st
}

func TestUnicodeOption(t *testing.T) {
	if _, err := New(runner.Config{Cols: 20, Rows: 3, Unicode: "11"}); err == nil {
		t.Error(`unicode "11" accepted; vt has only v6`)
	}
}

// TestModes: the mouse mode lowercased, bools by identity.
func TestModes(t *testing.T) {
	for in, want := range map[string]string{"": "none", "\x1b[?9h": "x10", "\x1b[?1000h": "vt200", "\x1b[?1002h": "drag",
		"\x1b[?1003h": "any", "\x1b[?1000h\x1b[?1000l": "none"} {
		if got := snap(t, in).P.Modes["mouseTrackingMode"]; got != want {
			t.Errorf("%q: mouseTrackingMode %v, want %v", in, got, want)
		}
	}
	m := snap(t, "\x1b[4h\x1b[?6h\x1b[?2004h").P.Modes
	if m["insertMode"] != true || m["originMode"] != true || m["bracketedPasteMode"] != true || m["applicationCursorKeysMode"] != false {
		t.Errorf("modes %v", m)
	}
}

// TestOpts: 6.0.0 semantics. DECSCUSR leaves the options alone (the golden for "\x1b[6 q" is block/false); ?12h
// writes options.cursorBlink.
func TestOpts(t *testing.T) {
	for in, want := range map[string][2]any{"": {"block", false}, "\x1b[6 q": {"block", false}, "\x1b[5 q": {"block", false},
		"\x1b[?12h": {"block", true}} {
		o := snap(t, in).P.Opts
		if o["cursorStyle"] != want[0] || o["cursorBlink"] != want[1] {
			t.Errorf("%q: opts %v, want style %v blink %v", in, o, want[0], want[1])
		}
	}
}

// TestCells: colours by mode, flags, the raw cursor in pending wrap, the alternate buffer.
func TestCells(t *testing.T) {
	st := snap(t, "\x1b[1;31ma\x1b[0;38;5;200mb\x1b[48;2;1;2;3mc\x1b[0md")
	want := []runner.Color{{Mode: "p16", Value: 1}, {Mode: "p256", Value: 200}, {Mode: "p256", Value: 200}, {Mode: "d"}}
	for x, w := range want {
		if got := st.Normal[0].Cells[x].A.FG; got != w {
			t.Errorf("cell %d fg %v, want %v", x, got, w)
		}
	}
	if bg := st.Normal[0].Cells[2].A.BG; bg != (runner.Color{Mode: "rgb", Value: 0x010203}) {
		t.Errorf("cell 2 bg %v", bg)
	}
	if !st.Normal[0].Cells[0].A.Bold || st.Normal[0].Cells[1].A.Bold {
		t.Error("bold not on cell 0 only")
	}
	if c := snap(t, "01234567890123456789").P.Cursor; c != [2]int{20, 0} {
		t.Errorf("pending-wrap cursor %v, want [20 0]", c)
	}
	if st := snap(t, "\x1b[?1049hx"); st.P.Buffer != "alternate" || st.Alt == nil {
		t.Errorf("buffer %q, alt rows %d", st.P.Buffer, len(st.Alt))
	}
}

// TestEvents: replies through OnData, titles, bells, line feeds.
func TestEvents(t *testing.T) {
	st := snap(t, "\x1b]2;one\x07\x1b]0;two\x07\x07\n\x1b[c")
	if st.P.Title != "two" || len(st.P.Titles) != 2 || st.P.Bells != 1 || st.P.LineFeeds != 1 {
		t.Errorf("title %q titles %v bells %d lfs %d", st.P.Title, st.P.Titles, st.P.Bells, st.P.LineFeeds)
	}
	if len(st.P.Replies) != 1 || st.P.Replies[0] != "\x1b[?1;2c" {
		t.Errorf("replies %q, want DA1", st.P.Replies)
	}
}
