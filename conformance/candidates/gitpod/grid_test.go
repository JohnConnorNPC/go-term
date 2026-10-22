package gitpod

import (
	"testing"

	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

func engine(t *testing.T, cols, rows int, unicode string) runner.Engine {
	t.Helper()
	e, err := New(runner.Config{Cols: cols, Rows: rows, Scrollback: 100, Unicode: unicode})
	if err != nil {
		t.Fatal(err)
	}
	return e
}

func snap(t *testing.T, e runner.Engine, in string) runner.State {
	t.Helper()
	if err := e.Write([]byte(in)); err != nil {
		t.Fatal(err)
	}
	st, err := e.Snapshot()
	if err != nil {
		t.Fatal(err)
	}
	return st
}

// TestGrid: the raw cursor (10 chars on 10 cols leaves x at 10, as xterm.js's cursorX does), the active buffer, and
// the unicode-version guard.
func TestGrid(t *testing.T) {
	if st := snap(t, engine(t, 10, 3, "6"), "0123456789"); st.P.Cursor != [2]int{10, 0} || st.P.Cols != 10 || st.P.Rows != 3 {
		t.Errorf("cursor %v cols %d rows %d, want [10 0] 10 3", st.P.Cursor, st.P.Cols, st.P.Rows)
	}
	if st := snap(t, engine(t, 10, 3, "6"), "\x1b[?1049h"); st.P.Buffer != "alternate" {
		t.Errorf("buffer %q after ?1049h, want alternate", st.P.Buffer)
	}
	if _, err := New(runner.Config{Cols: 10, Rows: 3, Unicode: "11"}); err == nil {
		t.Error("unicode 11 accepted")
	}
}
