package vt

import (
	"strings"
	"testing"
)

// TestOSC66Ignored: kitty's text sizing protocol (OSC 66) is not supported, as
// in xterm.js. Its payload must not reach the screen and the cursor must not
// move, which is also what makes the protocol's CPR-based detection report
// "unsupported" so that applications fall back to plain text.
func TestOSC66Ignored(t *testing.T) {
	for _, seq := range []string{
		"\x1b]66;s=2;hello\x07",
		"\x1b]66;w=2; \x07",
		"\x1b]66;s=3:w=4:n=1:d=2:v=2:h=1;sized\x1b\\",
		"\x1b]66;w=3;മലയാളം\x07",
		"\x1b]66;s=2;" + strings.Repeat("x", 4096) + "\x07",
	} {
		term := newTestTerminal(20, 3)
		term.WriteString("a" + seq + "b")
		if got := line(term, 0); got != "ab" {
			t.Errorf("%q: row0 = %q", seq[:min(len(seq), 24)], got)
		}
		if b := term.Buffer(); b.X != 2 || b.Y != 0 {
			t.Errorf("%q: cursor = %d,%d", seq[:min(len(seq), 24)], b.X, b.Y)
		}
	}
}
