package vt

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// dumpAll is every buffer line's text, each cell's fg/bg, and the cursor.
func dumpAll(term *Terminal) string {
	b := term.Buffer()
	var s strings.Builder
	for y := 0; y < b.Lines.Length(); y++ {
		l := b.Lines.Get(y)
		s.WriteString(l.TranslateToString(false, 0, -1))
		for x := 0; x < l.Length; x++ {
			fmt.Fprintf(&s, " %x/%x", l.data[x*cellSize+cellFG], l.data[x*cellSize+cellBG])
		}
		fmt.Fprintf(&s, " w%v\n", l.IsWrapped)
	}
	fmt.Fprintf(&s, "cursor %d,%d", b.X, b.Y)
	return s.String()
}

// TestHugeRepeatCounts [#7 r1 M1]: a count of 2^31-1 (the CSI parameter cap)
// for CHT, CBT, IL, DL, SU and SD returns at once with the state of a
// region-sized count, which already saturates (FuzzInstance hung on CBT).
func TestHugeRepeatCounts(t *testing.T) {
	setup := "\x1b[44m" + strings.Repeat("line\r\n", 30) + "abc\x1b[3;20r\x1b[10;30H"
	for _, c := range []struct{ final, small string }{
		{"I", "20"}, {"Z", "20"}, {"L", "11"}, {"M", "11"}, {"S", "18"}, {"T", "18"},
	} {
		small, huge := newTestTerminal(80, 24), newTestTerminal(80, 24)
		small.WriteString(setup + "\x1b[" + c.small + c.final + "x")
		start := time.Now()
		huge.WriteString(setup + "\x1b[2147483647" + c.final + "x")
		if d := time.Since(start); d > time.Second {
			t.Errorf("CSI %s with 2^31-1 took %v", c.final, d)
		}
		if dumpAll(small) != dumpAll(huge) {
			t.Errorf("CSI %s: 2^31-1 differs from %s", c.final, c.small)
		}
	}
}

// TestHugeRep [#7 r1 M1]: REP with 2^31-1 gives what the same clusters written
// out give, for a count past a full buffer that is congruent modulo a line
// (REP keeps the count's remainder past a full buffer instead of allocating
// 2^31-1 clusters).
func TestHugeRep(t *testing.T) {
	for _, c := range []struct {
		cluster       string
		cols, perLine int // clusters a line holds
	}{{"a", 80, 80}, {"日", 79, 39}, {"e\u0301", 80, 80}} {
		perLine := c.perLine
		lines := NewOptions().Scrollback + 24 + 1
		n := perLine*(lines+3) + (2147483647-perLine*(lines+3))%perLine
		written, rep := newTestTerminal(c.cols, 24), newTestTerminal(c.cols, 24)
		written.WriteString("x\r\n" + strings.Repeat(c.cluster, 1+n))
		start := time.Now()
		rep.WriteString("x\r\n" + c.cluster + "\x1b[2147483647b")
		if d := time.Since(start); d > time.Second {
			t.Errorf("%q: REP 2^31-1 took %v", c.cluster, d)
		}
		if dumpAll(written) != dumpAll(rep) {
			t.Errorf("%q at %d cols: REP 2^31-1 differs from %d written", c.cluster, c.cols, 1+n)
		}
	}
}
