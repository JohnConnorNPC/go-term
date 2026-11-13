package screen

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
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

// TestReplayBytes: for every mode combination of the Python reference's
// postamble (sniff/testdata/postamble.json), ReplayBytes is Preamble, the
// ring, CAN, then exactly that row's bytes.
func TestReplayBytes(t *testing.T) {
	b, err := os.ReadFile("../sniff/testdata/postamble.json")
	if err != nil {
		t.Fatal(err)
	}
	var f struct {
		Rows []struct {
			BracketedPaste bool   `json:"bracketed_paste"`
			AppCursor      bool   `json:"app_cursor"`
			MouseTracking  int    `json:"mouse_tracking"`
			MouseEncoding  int    `json:"mouse_encoding"`
			Hex            string `json:"hex"`
		}
	}
	if err := json.Unmarshal(b, &f); err != nil {
		t.Fatal(err)
	}
	if len(f.Rows) != 60 {
		t.Fatalf("%d rows, want 60 (2 x 2 x 5 x 3)", len(f.Rows))
	}
	ring := []byte("ring\x1b[31m")
	for _, r := range f.Rows {
		pa, err := hex.DecodeString(r.Hex)
		if err != nil || len(pa) == 0 {
			t.Fatalf("row %+v: hex %v", r, err)
		}
		m := sniff.Modes{BracketedPaste: r.BracketedPaste, DECCKM: r.AppCursor, MouseTrack: r.MouseTracking, MouseEnc: r.MouseEncoding}
		want := Preamble + string(ring) + "\x18" + string(pa)
		for _, alt := range []bool{false, true} { // the alt screen is not part of the postamble
			m.Alt = alt
			if got := string(ReplayBytes(m, ring)); got != want {
				t.Errorf("modes %+v: %q, want %q", m, got, want)
			}
		}
	}
}

// restoreStreams are the split streams plus streams that set every mode the
// postamble carries, on either buffer.
func restoreStreams() map[string]string {
	s := map[string]string{
		"modes":    "\x1b[?2004h\x1b[?1h\x1b[?1002h\x1b[?1006hx\x1b]2;t\x07\x1b[1;32mgreen\x1b[0m",
		"altmodes": "pre\x1b[?1049h\x1b[?1000h\x1b[?1016h\x1b]0;alt title\x1b\\\x1b[7malt\x1b[m text",
		"mouseoff": "\x1b[?9h\x1b[?1006h\x1b[?1003l\x1b[?2004h\x1b[?2004l\x1b[?1h\x1b[!pdone",
	}
	for k, v := range splitStreams {
		s[k] = v
	}
	return s
}

func sniffed(s string) (sniff.Modes, string, bool) {
	sn := sniff.New(sniff.Options{})
	sn.Feed([]byte(s))
	title, ok := sn.Title()
	return sn.Modes(), title, ok
}

func readAll(t *testing.T, i *Instance) Frame {
	t.Helper()
	f, err := i.Read(Request{View: ViewScrollback, Lines: 1 << 20, Attrs: true})
	if err != nil {
		t.Fatal(err)
	}
	return f
}

// TestRestoreRoundTrip: a stream S fed to A, and B restored from S with the
// modes and title the sniffer reports for S, read the same (lines with
// history, cursor, attribute runs, hashes, modes, title). A ring that lost
// the mode sets (scrolled out) still reports the injected modes. An empty
// ring and every head and tail cut of each stream restore without error,
// report the injected modes and take further output.
func TestRestoreRoundTrip(t *testing.T) {
	streams := restoreStreams()
	cuts := 0
	for name, s := range streams {
		m, title, ok := sniffed(s)
		a, _ := New(20, 6, Config{})
		if _, err := a.Write([]byte(s)); err != nil {
			t.Fatal(err)
		}
		snap := Snapshot{Version: SnapshotVersion, Cols: 20, Rows: 6, Modes: m, Title: title, TitleSet: ok, Ring: []byte(s)}
		b, err := Restore(snap, Config{})
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		fa, fb := readAll(t, a), readAll(t, b)
		if !reflect.DeepEqual(fa, fb) {
			t.Errorf("%s: restored frame differs\n fed      %+v\n restored %+v", name, fa, fb)
		}
		if fb.Modes() != m {
			t.Errorf("%s: modes %+v, sniffed %+v", name, fb.Modes(), m)
		}
		for k := 0; k <= len(s); k++ {
			for _, ring := range []string{s[k:], s[:k]} {
				snap.Ring = []byte(ring)
				checkUsable(t, fmt.Sprintf("%s cut %d %q", name, k, ring), snap)
				noTitle := snap
				noTitle.Title, noTitle.TitleSet = "", false
				checkUsable(t, fmt.Sprintf("%s cut %d %q, no title", name, k, ring), noTitle)
				cuts++
			}
		}
	}

	// The mode sets are before the ring (scrolled out): only the injected
	// state carries them.
	gone := "\x1b[?2004h\x1b[?1h\x1b[?1003h\x1b[?1006h\x1b]2;gone\x07\x1b[?1049h"
	for _, tail := range []string{"visible\r\nmore", "\x1b[?1049lnormal again"} {
		m, title, ok := sniffed(gone + tail)
		a, _ := New(20, 6, Config{})
		if _, err := a.Write([]byte(gone + tail)); err != nil {
			t.Fatal(err)
		}
		b, err := Restore(Snapshot{Version: SnapshotVersion, Cols: 20, Rows: 6, Modes: m, Title: title, TitleSet: ok, Ring: []byte(tail)}, Config{})
		if err != nil {
			t.Fatal(err)
		}
		fa, fb := readAll(t, a), readAll(t, b)
		if fb.Modes() != m || fb.Modes() != fa.Modes() || fb.Title != "gone" || fb.Text() != fa.Text() {
			t.Errorf("scrolled out %q: restored modes %+v title %q text %q; fed %+v %q %q",
				tail, fb.Modes(), fb.Title, fb.Text(), fa.Modes(), fa.Title, fa.Text())
		}
	}

	// An empty ring, for every mode combination.
	for _, alt := range []bool{false, true} {
		for _, track := range []int{0, 9, 1000, 1002, 1003} {
			for _, enc := range []int{0, 1006, 1016} {
				m := sniff.Modes{Alt: alt, DECCKM: track == 1000, BracketedPaste: enc == 1006, MouseTrack: track, MouseEnc: enc}
				checkUsable(t, fmt.Sprintf("empty %+v", m), Snapshot{Version: SnapshotVersion, Cols: 20, Rows: 6, Modes: m, Title: "x", TitleSet: true})
				cuts++
			}
		}
	}
	t.Logf("restore: %d streams, %d cuts and empty rings", len(streams), cuts)
}

// checkUsable restores snap, which must give its modes and its title, or
// without one the title of the ring fed alone (a sequence the ring was cut in
// commits no fragment), and then show further output.
func checkUsable(t *testing.T, name string, snap Snapshot) {
	t.Helper()
	i, err := Restore(snap, Config{})
	if err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	want := snap.Title
	if !snap.TitleSet {
		plain, _ := New(snap.Cols, snap.Rows, Config{})
		if _, err := plain.Write(snap.Ring); err != nil {
			t.Fatal(err)
		}
		want = readAll(t, plain).Title
	}
	f := readAll(t, i)
	if f.Modes() != snap.Modes || f.Title != want {
		t.Errorf("%s: modes %+v title %q, want %+v %q", name, f.Modes(), f.Title, snap.Modes, want)
	}
	if _, err := i.Write([]byte("\x1b[2J\x1b[Hok")); err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	if f, err = i.Read(Request{}); err != nil || f.Trimmed()[0] != "ok" {
		t.Errorf("%s: after restore, row 0 %q (%v)", name, f.Trimmed()[0], err)
	}
}
