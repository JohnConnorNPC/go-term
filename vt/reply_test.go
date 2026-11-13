package vt

import (
	"strings"
	"testing"
)

// TestRepliesApartFromInput: with OnReply set, what the terminal says of its
// own accord goes there and what the person typed to OnData; without it,
// both reach OnData as before.
func TestRepliesApartFromInput(t *testing.T) {
	term := newTestTerminal(80, 24)
	var data, reply strings.Builder
	term.OnData = func(s string) { data.WriteString(s) }
	term.OnReply = func(s string) { reply.WriteString(s) }

	term.WriteString("\x1b[c")  // DA1
	term.Input("ls\r", true)    // typed
	term.WriteString("\x1b[6n") // DSR
	if got, want := reply.String(), "\x1b[?1;2c\x1b[1;1R"; got != want {
		t.Errorf("replies %q, want %q", got, want)
	}
	if got := data.String(); got != "ls\r" {
		t.Errorf("input %q, want only what was typed", got)
	}

	plain := newTestTerminal(80, 24)
	var all strings.Builder
	plain.OnData = func(s string) { all.WriteString(s) }
	plain.WriteString("\x1b[c")
	plain.Input("x", true)
	if got := all.String(); got != "\x1b[?1;2cx" {
		t.Errorf("without OnReply: %q", got)
	}
}

// TestXTVersion: CSI > q names the terminal, as configured; other
// parameters are not defined and get nothing.
func TestXTVersion(t *testing.T) {
	term := newTestTerminal(80, 24)
	var out strings.Builder
	term.OnData = func(s string) { out.WriteString(s) }
	term.WriteString("\x1b[>q")
	if got := out.String(); got != "\x1bP>|xterm-go\x1b\\" {
		t.Errorf("default %q", got)
	}
	out.Reset()
	term.Options.XTVersion = "websh(1.2.3)"
	term.WriteString("\x1b[>0q\x1b[>1q")
	if got := out.String(); got != "\x1bP>|websh(1.2.3)\x1b\\" {
		t.Errorf("named %q", got)
	}
}
