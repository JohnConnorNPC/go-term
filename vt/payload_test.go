package vt

import (
	"strings"
	"testing"
)

func limitedTerminal(limit int) *Terminal {
	o := NewOptions()
	o.PayloadLimit = limit
	return NewTerminal(o)
}

// TestPayloadLimit: Options.PayloadLimit drops an OSC, DCS or APC string
// longer than the limit whole, across writes, and the next one parses; the
// default is xterm.js's 10,000,000.
func TestPayloadLimit(t *testing.T) {
	if p := NewTerminal(nil).inputHandler.parser; p.payloadLimit != 10000000 || p.oscParser.limit != 10000000 || p.dcsParser.limit != 10000000 {
		t.Errorf("default limit %d", p.payloadLimit)
	}

	term := limitedTerminal(10)
	var title string
	term.OnTitleChange = func(s string) { title = s }
	term.WriteString("\x1b]2;01234")
	term.WriteString("56789\x07")
	if title != "0123456789" {
		t.Errorf("10 units: title %q", title)
	}
	term.WriteString("\x1b]2;012345")
	term.WriteString("6789X\x07\x1b]2;ok\x07")
	if title != "ok" {
		t.Errorf("11 units then ok: title %q", title)
	}

	term = limitedTerminal(1)
	var replies []string
	term.OnData = func(s string) { replies = append(replies, s) }
	term.WriteString("\x1bP$qm\x1b\\\x1bP$q\"p\x1b\\\x1bP$qm\x1b\\")
	if got := strings.Join(replies, "|"); got != "\x1bP1$r0m\x1b\\|\x1bP1$r0m\x1b\\" {
		t.Errorf("DECRQSS replies %q", got)
	}

	term = limitedTerminal(3)
	var apc []string
	term.InputHandler().SetApcHandler(func(d string) bool { apc = append(apc, d); return true })
	term.WriteString("\x1b_Gab\x1b\\\x1b_Ga")
	term.WriteString("bc\x1b\\\x1b_Gx\x1b\\")
	if got := strings.Join(apc, "|"); got != "Gab|Gx" {
		t.Errorf("apc %q", got)
	}
}
