package xvt

import (
	"strings"
	"testing"

	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

func snap(t *testing.T, in string) runner.State {
	t.Helper()
	e, err := New(runner.Config{Cols: 10, Rows: 3, Scrollback: 100, Unicode: "6"})
	if err != nil {
		t.Fatal(err)
	}
	if err := e.Write([]byte(in)); err != nil {
		t.Fatal(err)
	}
	st, err := e.Snapshot()
	if err != nil {
		t.Fatal(err)
	}
	return st
}

// TestGrid: size, cursor and the active screen; the paths x/vt has no API for are listed, not made up.
func TestGrid(t *testing.T) {
	st := snap(t, "ab")
	if st.P.Cols != 10 || st.P.Rows != 3 || st.P.Cursor != [2]int{2, 0} {
		t.Errorf("cols %d rows %d cursor %v, want 10x3 at [2 0]", st.P.Cols, st.P.Rows, st.P.Cursor)
	}
	for _, u := range []string{"attrs", "viewportY", "w"} {
		if !strings.Contains(strings.Join(st.Unobserved, " "), u) {
			t.Errorf("Unobserved %v lacks %s", st.Unobserved, u)
		}
	}
	if st := snap(t, "x\x1b[?1049h"); st.P.Buffer != "alternate" {
		t.Errorf("alt: buffer %q, want alternate", st.P.Buffer)
	}
}
