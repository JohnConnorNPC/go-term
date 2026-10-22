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

// TestGrid: cell text and width by identity; the paths x/vt has no API for are listed, not made up.
func TestGrid(t *testing.T) {
	st := snap(t, "ab")
	if c := st.Normal[0].Cells; c[0].Ch != "a" || c[1].Ch != "b" || c[0].W != 1 || st.P.Cols != 10 || st.P.Rows != 3 {
		t.Errorf("cells %+v cols %d rows %d, want a b width 1, 10x3", c[:2], st.P.Cols, st.P.Rows)
	}
	for _, u := range []string{"attrs", "viewportY", "w"} {
		if !strings.Contains(strings.Join(st.Unobserved, " "), u) {
			t.Errorf("Unobserved %v lacks %s", st.Unobserved, u)
		}
	}
	if st := snap(t, "x\x1b[?1049h"); st.P.Buffer != "alternate" || st.Normal != nil || !strings.Contains(strings.Join(st.Unobserved, " "), "normalLength") {
		t.Errorf("alt: buffer %q normal %v unobserved %v, want alternate, no normal rows, normalLength unobserved", st.P.Buffer, st.Normal != nil, st.Unobserved)
	}
}
