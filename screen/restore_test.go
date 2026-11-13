package screen

import (
	"errors"
	"strings"
	"testing"

	"github.com/JohnConnorNPC/go-term/sniff"
)

func restore(t *testing.T, s Snapshot) *Instance {
	t.Helper()
	if s.Version == 0 {
		s.Version = SnapshotVersion
	}
	if s.Cols == 0 {
		s.Cols, s.Rows = 10, 3
	}
	i, err := Restore(s, Config{})
	if err != nil {
		t.Fatalf("Restore: %v", err)
	}
	return i
}

func readTrimmed(t *testing.T, i *Instance) (Frame, string) {
	t.Helper()
	f, err := i.Read(Request{})
	if err != nil {
		t.Fatal(err)
	}
	return f, strings.Join(f.Trimmed(), "|")
}

// TestRestoreBasic: versions and sizes are checked, the injected title wins
// over the ring's, and the alt screen is aligned with the reported state in
// both directions without touching a replay that already agrees.
func TestRestoreBasic(t *testing.T) {
	for _, v := range []int{0, 2} {
		if i, err := Restore(Snapshot{Version: v, Cols: 10, Rows: 3}, Config{}); err == nil || i != nil {
			t.Errorf("version %d: %v, %v; want an error and no instance", v, i, err)
		}
	}
	if i, err := Restore(Snapshot{Version: SnapshotVersion}, Config{}); err == nil || i != nil {
		t.Errorf("0x0: %v, %v; want an error", i, err)
	}

	// Title: the ring's title stays unless one is injected.
	f, _ := readTrimmed(t, restore(t, Snapshot{Ring: []byte("\x1b]2;ring\x07")}))
	if f.Title != "ring" {
		t.Errorf("title %q, want the ring's", f.Title)
	}
	f, _ = readTrimmed(t, restore(t, Snapshot{Ring: []byte("\x1b]2;ring\x07"), Title: "agent", TitleSet: true}))
	if f.Title != "agent" {
		t.Errorf("title %q, want the injected one", f.Title)
	}

	// Alt reported, the ring has no switch (1049h scrolled out): the ring is
	// drawn on the alt grid, the normal buffer stays empty.
	i := restore(t, Snapshot{Modes: sniff.Modes{Alt: true}, Ring: []byte("top\r\nmid")})
	f, got := readTrimmed(t, i)
	if !f.AltScreen || got != "top|mid|" || *f.Cursor != (Cursor{Row: 1, Col: 3}) {
		t.Errorf("alt seed: alt %v grid %q cursor %+v", f.AltScreen, got, *f.Cursor)
	}
	if _, err := i.Write([]byte(altOff)); err != nil {
		t.Fatal(err)
	}
	if _, got := readTrimmed(t, i); got != "||" {
		t.Errorf("alt seed: normal buffer %q, want empty", got)
	}

	// Alt reported and the ring enters it after some output: the replay
	// already agrees and is kept as is (a ?1049h prefix would leave "pre" on
	// the alt grid).
	ring := []byte("pre\x1b[?1049hpost")
	ref, _ := New(10, 3, Config{})
	if _, err := ref.Write(ring); err != nil {
		t.Fatal(err)
	}
	_, want := readTrimmed(t, ref)
	f, got = readTrimmed(t, restore(t, Snapshot{Modes: sniff.Modes{Alt: true}, Ring: ring}))
	if !f.AltScreen || got != want || want != "   post||" {
		t.Errorf("alt kept: alt %v grid %q, plain feed %q", f.AltScreen, got, want)
	}

	// Alt not reported but the ring ends on the alt buffer: 1049l follows,
	// the normal grid shows.
	f, got = readTrimmed(t, restore(t, Snapshot{Ring: ring}))
	if f.AltScreen || got != "pre||" {
		t.Errorf("alt off: alt %v grid %q", f.AltScreen, got)
	}

	// Alt reported but the ring leaves the alt buffer itself: alt is entered
	// last.
	f, got = readTrimmed(t, restore(t, Snapshot{Modes: sniff.Modes{Alt: true}, Ring: []byte("a\x1b[?1049hb\x1b[?1049l")}))
	if !f.AltScreen || got != "||" {
		t.Errorf("contradicting ring: alt %v grid %q", f.AltScreen, got)
	}

	// A failed instance gives no instance: New's errors pass through too.
	if _, err := Restore(Snapshot{Version: SnapshotVersion, Cols: MaxCols + 1, Rows: 3}, Config{}); err == nil || errors.Is(err, ErrFailed) {
		t.Errorf("too wide: %v", err)
	}
}
