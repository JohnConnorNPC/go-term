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

// TestRepPastLineEnd [#9]: REP reads the cell at X minus the preceding
// char's width. Under mode 2027 a zero-width char that joins nothing leaves a
// non-zero join state of width 0, so that cell is X itself: cols after one
// written in the last column, cols+1 after one dropped in the pending-wrap
// column. xterm.js's getString reads past its typed array there, gets ” and
// prints nothing; vt indexed past the line (panic).
func TestRepPastLineEnd(t *testing.T) {
	for _, c := range []struct {
		name, pre string
		x         int
	}{{"at cols", strings.Repeat("0", 11), 12}, {"at cols+1", strings.Repeat("0", 12), 13}} {
		term := newMode2027Terminal(12, 3)
		term.WriteString("\x1b[?2027h" + c.pre + "\x1f́\x1b[b")
		if b := term.Buffer(); b.X != c.x || b.Y != 0 || line(term, 1) != "" {
			t.Errorf("%s: cursor (%d,%d) row 1 %q, want (%d,0) and REP printing nothing", c.name, b.X, b.Y, line(term, 1), c.x)
		}
	}
	// x < 0, any provider: DECAWM off at 2 cols, a second wide char is
	// dropped at X = 1 with the join state of width 2 (x = -1).
	term := newTestTerminal(2, 3)
	term.WriteString("\x1b[?7l中中\x1b[b")
	if b := term.Buffer(); b.X != 1 || b.Y != 0 || line(term, 0) != "中" {
		t.Errorf("2 cols ?7l: cursor (%d,%d) row 0 %q, want (1,0) %q", b.X, b.Y, line(term, 0), "中")
	}
}
