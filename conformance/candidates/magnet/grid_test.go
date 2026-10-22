package magnet

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

// TestGridCells: cell text, width and attributes through norm; the wrap flag of the continuation line; the normal
// buffer kept while the alternate one is active.
func TestGridCells(t *testing.T) {
	st := snap(t, engine(t, 20, 3, "6"), "a\x1b[31mb")
	a, b := st.Normal[0].Cells[0], st.Normal[0].Cells[1]
	if a.Ch != "a" || a.W != 1 || a.A.FG != (runner.Color{Mode: "d"}) || b.Ch != "b" || b.A.FG != (runner.Color{Mode: "p16", Value: 1}) {
		t.Errorf("cells %+v %+v, want a default, b p16 1", a, b)
	}
	if st := snap(t, engine(t, 20, 3, "6"), "012345678901234567890"); st.Normal[0].Wrapped || !st.Normal[1].Wrapped {
		t.Errorf("wrapped %v %v, want false true", st.Normal[0].Wrapped, st.Normal[1].Wrapped)
	}
	st = snap(t, engine(t, 20, 3, "6"), "x\x1b[?1049h")
	if st.Alt == nil || len(st.Normal) == 0 || st.Normal[0].Cells[0].Ch != "x" {
		t.Errorf("alt %v normal %v, want both, normal keeping x", st.Alt != nil, st.Normal)
	}
}
