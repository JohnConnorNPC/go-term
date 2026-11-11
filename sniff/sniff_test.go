package sniff

import (
	"strings"
	"testing"
)

func feed(s *Sniffer, chunks ...string) {
	for _, c := range chunks {
		s.Feed([]byte(c))
	}
}

func TestModes(t *testing.T) {
	for _, tc := range []struct {
		name string
		in   []string
		want Modes
	}{
		{"power-on", nil, Modes{}},
		{"alt 1049", []string{"\x1b[?1049h"}, Modes{Alt: true}},
		{"alt 47", []string{"\x1b[?47h"}, Modes{Alt: true}},
		{"alt 1047 then 1049l", []string{"\x1b[?1047h", "\x1b[?1049l"}, Modes{}},
		{"decckm", []string{"\x1b[?1h"}, Modes{DECCKM: true}},
		{"decckm off", []string{"\x1b[?1h\x1b[?1l"}, Modes{}},
		{"2004", []string{"\x1b[?2004h"}, Modes{BracketedPaste: true}},
		{"combined", []string{"\x1b[?1049;1;2004h"}, Modes{Alt: true, DECCKM: true, BracketedPaste: true}},
		{"leading zeros", []string{"\x1b[?0001h"}, Modes{DECCKM: true}},
		{"no substring", []string{"\x1b[?470h\x1b[?10h\x1b[?12004h"}, Modes{}},
		{"mouse left to right", []string{"\x1b[?1002;1003h"}, Modes{MouseTrack: 1003}},
		{"mouse group reset by any member", []string{"\x1b[?1003h\x1b[?9l"}, Modes{}},
		{"mouse reset then set", []string{"\x1b[?1003h\x1b[?1000l\x1b[?1002h"}, Modes{MouseTrack: 1002}},
		{"encoding", []string{"\x1b[?1006h"}, Modes{MouseEnc: 1006}},
		{"encoding 1016 then 1006l", []string{"\x1b[?1016h\x1b[?1006l"}, Modes{}},
		{"1005 and 1015 not tracked", []string{"\x1b[?1005h\x1b[?1015h"}, Modes{}},
		{"RIS", []string{"\x1b[?1049;1;2004;1003;1006h", "\x1bc"}, Modes{}},
		{"RIS then set", []string{"\x1b[?1003h\x1bc\x1b[?1002h"}, Modes{MouseTrack: 1002}},
		{"DECSTR", []string{"\x1b[?1049;1;2004;1003;1006h", "\x1b[!p"}, Modes{Alt: true, MouseTrack: 1003, MouseEnc: 1006}},
		{"not private", []string{"\x1b[1049h\x1b[4h"}, Modes{}},
		{"other prefix", []string{"\x1b[>1h"}, Modes{}},
		{"prefix after param ignored", []string{"\x1b[1?1h"}, Modes{}},
		{"intermediate", []string{"\x1b[?1$h"}, Modes{}},
		{"subparam", []string{"\x1b[?1049:5h"}, Modes{Alt: true}},
		{"empty params", []string{"\x1b[?;;1h"}, Modes{DECCKM: true}},
		{"32 params", []string{"\x1b[?" + strings.Repeat("0;", 31) + "1h"}, Modes{DECCKM: true}},
		{"33rd param dropped", []string{"\x1b[?" + strings.Repeat("0;", 32) + "1h"}, Modes{}},
		{"clamped", []string{"\x1b[?99999999999999999999h"}, Modes{}},
		{"C0 inside CSI", []string{"\x1b[?10\n49h"}, Modes{Alt: true}},
		{"CAN aborts", []string{"\x1b[?10\x1849h"}, Modes{}},
		{"ESC restarts", []string{"\x1b[?1\x1b[?2004h"}, Modes{BracketedPaste: true}},
		{"8-bit CSI as UTF-8", []string{"\xc2\x9b?1h"}, Modes{DECCKM: true}},
		{"raw 0x9b dropped", []string{"\x9b?1h"}, Modes{}},
		{"invalid byte inside CSI dropped", []string{"\x1b[?10\x8049h"}, Modes{Alt: true}},
		{"non-ASCII inside CSI aborts", []string{"\x1b[?10é49h"}, Modes{}},
		{"inside OSC", []string{"\x1b]0;\x1b[?1h"}, Modes{DECCKM: true}}, // the bare ESC ends the OSC
		{"inside OSC payload text", []string{"\x1b]2;x[?1h\x07"}, Modes{}},
		{"inside DCS", []string{"\x1bP1$q[?1h\x1b\\"}, Modes{}},
		{"after DCS, ESC c is RIS", []string{"\x1b[?1h\x1bP$qm\x1bc"}, Modes{}},
		{"inside APC", []string{"\x1b_\x9b?1h\x1b\\"}, Modes{}},
		{"seq at cap", []string{"\x1b[?" + strings.Repeat("0", 61) + "1049h"}, Modes{Alt: true}},
		{"seq past cap", []string{"\x1b[?" + strings.Repeat("0", 62) + "1049h"}, Modes{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := New(Options{})
			feed(s, tc.in...)
			if got := s.Modes(); got != tc.want {
				t.Errorf("Modes() = %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestTitle(t *testing.T) {
	type want struct {
		title string
		ok    bool
	}
	for _, tc := range []struct {
		name string
		in   []string
		want want
	}{
		{"none", []string{"plain"}, want{"", false}},
		{"OSC 0 BEL", []string{"\x1b]0;hello\x07"}, want{"hello", true}},
		{"OSC 2 ST", []string{"\x1b]2;hello\x1b\\"}, want{"hello", true}},
		{"OSC 1 ignored", []string{"\x1b]1;icon\x07"}, want{"", false}},
		{"OSC 00", []string{"\x1b]00;zero\x07"}, want{"zero", true}},
		{"OSC 20 ignored", []string{"\x1b]20;x\x07"}, want{"", false}},
		{"no id", []string{"\x1b];x\x07"}, want{"", false}},
		{"bad id aborts", []string{"\x1b]0a;x\x07"}, want{"", false}},
		{"id only", []string{"\x1b]2\x07"}, want{"", true}},
		{"empty", []string{"\x1b]0;\x07"}, want{"", true}},
		{"bare ESC commits", []string{"\x1b]0;junk\x1bx"}, want{"junk", true}},
		{"C1 ST", []string{"\x1b]0;a\xc2\x9c"}, want{"a", true}},
		{"C1 OSC", []string{"\xc2\x9d0;b\x07"}, want{"b", true}},
		{"CAN abandons", []string{"\x1b]0;a\x18\x07"}, want{"", false}},
		{"C1 control abandons", []string{"\x1b]0;a\xc2\x85\x07"}, want{"", false}},
		{"C0 dropped", []string{"\x1b]0;a\tb\x01c\x07"}, want{"abc", true}},
		{"DEL kept", []string{"\x1b]2;a\x7fb\x07"}, want{"a\x7fb", true}}, // xterm.js's OSC_STRING puts 0x7F too
		{"utf8", []string{"\x1b]0;h\xc3\xa9\xe2\x82\xac\x07"}, want{"hé€", true}},
		{"invalid utf8 dropped", []string{"\x1b]0;a\xffb\xc3c\x07"}, want{"abc", true}},
		{"split utf8", []string{"\x1b]0;\xe2\x82", "\xac\x07"}, want{"€", true}},
		{"last wins", []string{"\x1b]0;a\x07\x1b]2;b\x07"}, want{"b", true}},
		{"kept over RIS", []string{"\x1b]0;a\x07\x1bc"}, want{"a", true}},
		{"payload at cap", []string{"\x1b]0;" + strings.Repeat("t", 4096) + "\x07"}, want{strings.Repeat("t", 1024), true}},
		{"payload past cap", []string{"\x1b]0;" + strings.Repeat("t", 4097) + "\x07"}, want{"", false}},
		{"old title kept past cap", []string{"\x1b]0;a\x07\x1b]0;" + strings.Repeat("t", 5000) + "\x07"}, want{"a", true}},
		{"rune cut", []string{"\x1b]0;" + strings.Repeat("x", 1023) + "é\x07"}, want{strings.Repeat("x", 1023), true}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := New(Options{})
			feed(s, tc.in...)
			if title, ok := s.Title(); title != tc.want.title || ok != tc.want.ok {
				t.Errorf("Title() = %q, %v; want %q, %v", title, ok, tc.want.title, tc.want.ok)
			}
		})
	}
}

func TestTitleMax(t *testing.T) {
	s := New(Options{TitleMax: 3})
	feed(s, "\x1b]0;aé\x07")
	if got, _ := s.Title(); got != "aé" {
		t.Errorf("TitleMax 3: %q", got)
	}
	feed(s, "\x1b]0;abé\x07")
	if got, _ := s.Title(); got != "ab" {
		t.Errorf("TitleMax 3, rune cut: %q", got)
	}
	s = New(Options{TitleMax: 4096})
	feed(s, "\x1b]0;"+strings.Repeat("t", 4096)+"\x07")
	if got, _ := s.Title(); len(got) != 4096 {
		t.Errorf("TitleMax 4096: len %d", len(got))
	}
}

func TestReset(t *testing.T) {
	s := New(Options{TitleMax: 2})
	feed(s, "\x1b]0;abc\x07\x1b[?1049;1003h\x1b]0;partial\xe2")
	s.Reset()
	if m := s.Modes(); m != (Modes{}) {
		t.Errorf("Modes after Reset: %+v", m)
	}
	if title, ok := s.Title(); ok || title != "" {
		t.Errorf("Title after Reset: %q %v", title, ok)
	}
	feed(s, "\x82\xac\x07\x1b]0;xyz\x07") // the partial OSC and code point are gone; TitleMax is kept
	if title, ok := s.Title(); title != "xy" || !ok {
		t.Errorf("after Reset: %q %v", title, ok)
	}
}
