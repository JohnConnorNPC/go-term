package vt

import (
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
