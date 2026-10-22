package xvt

import (
	"slices"
	"sort"
	"strings"
	"testing"
	"time"

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
		if !slices.Contains(st.Unobserved, u) {
			t.Errorf("Unobserved %v lacks %s", st.Unobserved, u)
		}
	}
	if st := snap(t, "x\x1b[?1049h"); st.P.Buffer != "alternate" || st.Normal != nil || !slices.Contains(st.Unobserved, "normalLength") {
		t.Errorf("alt: buffer %q normal %v unobserved %v, want alternate, no normal rows, normalLength unobserved", st.P.Buffer, st.Normal != nil, st.Unobserved)
	}
}

// TestEvents: a DA1 write returns (the reply pipe is drained), titles and bells by identity, and the exact list of
// paths x/vt can't show (GO S-3: modes too).
func TestEvents(t *testing.T) {
	st := snap(t, "\x1b[c\x1b]2;a\a\a")
	if st.P.Title != "a" || st.P.Bells != 1 {
		t.Errorf("title %q bells %d, want a 1", st.P.Title, st.P.Bells)
	}
	u := append([]string{}, st.Unobserved...)
	sort.Strings(u)
	if strings.Join(u, " ") != "attrs lineFeeds modes opts replies viewportY w" {
		t.Errorf("Unobserved %v", u)
	}
}

// TestDrainEnd: once the reply pipe ends (here the emulator is closed), every snapshot is an error, not a state.
func TestDrainEnd(t *testing.T) {
	e, err := New(runner.Config{Cols: 10, Rows: 3, Scrollback: 100, Unicode: "6"})
	if err != nil {
		t.Fatal(err)
	}
	if err := e.(*Engine).e.Close(); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 100; i++ {
		if _, err := e.Snapshot(); err != nil {
			if !strings.Contains(err.Error(), "reply drain ended") {
				t.Fatalf("snapshot error %v, want the drain's end", err)
			}
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("snapshots after the drain ended still succeed")
}
