package vt

import (
	"encoding/hex"
	"strings"
	"testing"
)

// Narrowing a full scrollback to the minimum width used to panic with an
// index of -1.
//
// Fit clamps to MinimumCols/MinimumRows, so 2x1 is a size the terminal is
// expected to accept, and a browser hands it one whenever the container is
// measured mid-layout — an absolutely positioned terminal on an otherwise
// empty page reports almost no size for a frame.
//
// The buffer has to be full for this to bite. Reflowing wraps each original
// line into many, and only once originalLines+countToInsert exceeds the
// scrollback does the loop writing them back run out of room and walk its
// index below zero. A near-empty buffer always has somewhere to put them,
// which is why a short write does not reproduce it.
func fillScrollback(t *testing.T, term *Terminal) {
	t.Helper()
	opts := NewOptions()
	var b strings.Builder
	// Comfortably more lines than the scrollback holds, each long enough to
	// wrap into a great many lines once the terminal is two columns wide.
	for i := 0; i < opts.Scrollback+400; i++ {
		b.WriteString(strings.Repeat("the quick brown fox ", 12))
		b.WriteString("\r\n")
	}
	term.WriteString(b.String())
}

func TestNarrowFullScrollbackToMinimum(t *testing.T) {
	term := NewTerminal(NewOptions())
	fillScrollback(t, term)

	term.Resize(MinimumCols, MinimumRows)
	if got := term.Cols(); got != MinimumCols {
		t.Errorf("Cols = %d, want %d", got, MinimumCols)
	}
	term.Resize(80, 24)
	if got, want := term.Cols(), 80; got != want {
		t.Errorf("Cols after widening = %d, want %d", got, want)
	}
}

// The same walk, arrived at one column at a time, since each step reflows
// against a different amount of already-wrapped content.
func TestNarrowFullScrollbackOneColumnAtATime(t *testing.T) {
	term := NewTerminal(NewOptions())
	fillScrollback(t, term)
	for cols := 80; cols >= MinimumCols; cols-- {
		term.Resize(cols, 24)
	}
	for cols := MinimumCols; cols <= 80; cols++ {
		term.Resize(cols, 24)
	}
}

// A one-row terminal has no room for the wrapped remainder of anything, and
// has to keep taking writes afterwards.
func TestSingleRowNarrowStillAcceptsWrites(t *testing.T) {
	term := NewTerminal(NewOptions())
	fillScrollback(t, term)
	term.Resize(MinimumCols, 1)
	term.WriteString(strings.Repeat("y", 500))
}

// TestCombiningAtPendingWrap: a combining mark after a C0 control joins
// nothing, and in the pending-wrap column it used to be written one cell past
// the line (index out of range; xterm.js's typed array drops the write). gh #7
// FuzzInstance found it at 2 columns. xterm.js still moves the cursor to
// cols+1 there [#9]; a char in the next write wraps from it (the print-start
// wide-char check must not read cell cols), and a second zero-width char,
// joining nothing either, wraps and takes column 0.
func TestCombiningAtPendingWrap(t *testing.T) {
	for _, cols := range []int{2, 80} {
		term := newTestTerminal(cols, 3)
		full := strings.Repeat("0", cols)
		term.WriteString(full + "\x1f\u0301")
		if b := term.Buffer(); b.X != cols+1 || b.Y != 0 {
			t.Errorf("%d cols: cursor (%d,%d), want (%d,0)", cols, b.X, b.Y, cols+1)
		}
		term.WriteString("x") // a separate write: print starts at X = cols+1
		if b := term.Buffer(); line(term, 0) != full || line(term, 1) != "x" || b.X != 1 || b.Y != 1 {
			t.Errorf("%d cols: %q, %q, cursor (%d,%d), want x at (1,1)", cols, line(term, 0), line(term, 1), b.X, b.Y)
		}
		two := newTestTerminal(cols, 3)
		two.WriteString(full + "\x1f\u200d\u200d")
		if b := two.Buffer(); b.X != 1 || b.Y != 1 {
			t.Errorf("%d cols, two zero-width: cursor (%d,%d), want (1,1)", cols, b.X, b.Y)
		}
	}
}

// TestFuzzCase1794598569: steps 0-1 of the differential fuzz case
// fuzz/1794598569-6 (12x11): the second write ends a row in the pending-wrap
// column, then OSC 8 and a ZWJ that joins nothing. Headless 6.0.0's golden
// for that step has cursor [13, 1].
func TestFuzzCase1794598569(t *testing.T) {
	term := newTestTerminal(12, 11)
	for _, h := range []string{"405539323b437d2a4c6c2a",
		"1b5b33383b353b3235356d1b5b34383a323a3a3133343a3136353a34306d20632a26337b6e583c5c1b5d383b3b687474703a2f2f782f71447d395a6a3f075c2f381b5d383b3b07e2808d"} {
		b, err := hex.DecodeString(h)
		if err != nil {
			t.Fatal(err)
		}
		term.Write(b)
	}
	if b := term.Buffer(); b.X != 13 || b.Y != 1 {
		t.Errorf("cursor (%d,%d), headless (13,1)", b.X, b.Y)
	}
	if got := line(term, 1); got != "c*&3{nX<\\\\/8" {
		t.Errorf("row 1 %q", got)
	}
}
