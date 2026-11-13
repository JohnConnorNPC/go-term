package vt

import (
	"strings"
	"testing"
)

// TestApc: an APC string reaches the APC handler whole — kitty graphics'
// transport — and the text around it prints; OSC is unaffected, and
// without a handler APC is dropped as before.
func TestApc(t *testing.T) {
	term := newTestTerminal(80, 24)
	var got []string
	term.InputHandler().SetApcHandler(func(d string) bool { got = append(got, d); return true })
	var title string
	term.OnTitleChange = func(s string) { title = s }
	term.WriteString("a\x1b_Gf=100,a=T;iVBORw0KGgo=\x1b\\b\x1b]2;t\x1b\\\x1b_Gq=2;\x1b\\")
	term.WriteString("\x1b_Ga=d,d=" + "A" + "\x1b\\") // in pieces
	if strings.Join(got, "|") != "Gf=100,a=T;iVBORw0KGgo=|Gq=2;|Ga=d,d=A" {
		t.Errorf("apc %q", got)
	}
	if title != "t" || line(term, 0) != "ab" {
		t.Errorf("title %q, line %q", title, line(term, 0))
	}

	plain := newTestTerminal(80, 24)
	plain.WriteString("x\x1b_Gsecret\x1b\\y")
	if line(plain, 0) != "xy" {
		t.Errorf("without a handler: %q", line(plain, 0))
	}
}
