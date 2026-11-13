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

// TestPayloadBuffersReleased: a long OSC, DCS or APC string does not keep
// its buffer once it ends; a short one keeps it for the next.
func TestPayloadBuffersReleased(t *testing.T) {
	term := NewTerminal(nil)
	term.InputHandler().SetApcHandler(func(string) bool { return true })
	p := term.inputHandler.parser
	osc := p.oscParser.handlers[2][0].(*OscHandler)
	dcs := p.dcsParser.handlers[p.identifier(FunctionID{Intermediates: "$", Final: "q"}, [2]int{0x40, 0x7e})][0].(*DcsHandlerFunc)
	long := strings.Repeat("x", 100<<10)
	term.WriteString("\x1b]2;" + long + "\x07\x1bP$q" + long + "\x1b\\\x1b_" + long + "\x1b\\")
	if cap(osc.data) != 0 || cap(dcs.data) != 0 || cap(p.apcBuf) != 0 {
		t.Errorf("long: caps osc %d dcs %d apc %d", cap(osc.data), cap(dcs.data), cap(p.apcBuf))
	}
	term.WriteString("\x1b]2;short\x07\x1bP$qm\x1b\\\x1b_short\x1b\\")
	if cap(osc.data) == 0 || cap(dcs.data) == 0 || cap(p.apcBuf) == 0 {
		t.Errorf("short: caps osc %d dcs %d apc %d", cap(osc.data), cap(dcs.data), cap(p.apcBuf))
	}
}
