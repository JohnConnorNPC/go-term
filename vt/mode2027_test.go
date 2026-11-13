package vt

import (
	"strings"
	"testing"
)

// newMode2027Terminal is newTestTerminal with Options.GraphemeMode2027 on,
// for the upstream tests of mode 2027 itself.
func newMode2027Terminal(cols, rows int) *Terminal {
	opts := NewOptions()
	opts.Cols = cols
	opts.Rows = rows
	opts.GraphemeMode2027 = true
	return NewTerminal(opts)
}

// TestMode2027OffByDefault: @xterm/headless 6.0.0 has no mode 2027, so by
// default ?2027h changes no width (Unicode v6 stays) and DECRQM ?2027 answers
// 0, not recognised.
func TestMode2027OffByDefault(t *testing.T) {
	const family = "\U0001F468‍\U0001F469‍\U0001F467" // man ZWJ woman ZWJ girl
	cells := func(term *Terminal) (string, int) {
		b := term.Buffer()
		row := b.Lines.Get(b.YBase)
		var out []string
		for x := 0; x < row.Length; x++ {
			if row.HasContent(x) {
				out = append(out, row.GetString(x))
			}
		}
		return strings.Join(out, "|"), b.X
	}
	for _, text := range []string{family + "x", "❤️x", "\U0001F1E9\U0001F1EAx"} {
		plain := newTestTerminal(20, 3)
		plain.WriteString(text)
		set := newTestTerminal(20, 3)
		set.WriteString("\x1b[?2027h" + text)
		pc, px := cells(plain)
		sc, sx := cells(set)
		if pc != sc || px != sx {
			t.Errorf("%q: with ?2027h cells %q x %d, without %q x %d", text, sc, sx, pc, px)
		}
	}
	term := newTestTerminal(10, 3)
	var out strings.Builder
	term.OnData = func(s string) { out.WriteString(s) }
	for _, seq := range []string{"", "\x1b[?2027h", "\x1b[?2027l"} {
		term.WriteString(seq)
		out.Reset()
		term.WriteString("\x1b[?2027$p")
		if got := out.String(); got != "\x1b[?2027;0$y" {
			t.Errorf("DECRQM after %q = %q, want not recognised", seq, got)
		}
	}
}
