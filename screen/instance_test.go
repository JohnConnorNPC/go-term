package screen

import (
	"testing"

	"github.com/JohnConnorNPC/go-term/vt"
)

// TestInstanceNew: New takes cols 1..2000 and rows 1..1000 (one column
// reads as vt.MinimumCols) and refuses anything else; the payload cap is
// Config's, DefaultPayloadLimit by default.
func TestInstanceNew(t *testing.T) {
	for _, s := range [][2]int{{0, 24}, {80, 0}, {-1, 24}, {2001, 24}, {80, 1001}} {
		if i, err := New(s[0], s[1], Config{}); err == nil || i != nil {
			t.Errorf("New(%d, %d) accepted", s[0], s[1])
		}
	}
	for _, s := range [][2]int{{1, 1}, {2000, 1000}, {80, 24}} {
		i, err := New(s[0], s[1], Config{})
		if err != nil {
			t.Fatalf("New(%d, %d): %v", s[0], s[1], err)
		}
		if want := max(s[0], vt.MinimumCols); i.t.Cols() != want || i.t.Rows() != s[1] {
			t.Errorf("New(%d, %d) is %dx%d", s[0], s[1], i.t.Cols(), i.t.Rows())
		}
	}
	if _, err := New(80, 24, Config{Unicode: "11"}); err == nil {
		t.Errorf("Unicode 11 accepted")
	}
	i, _ := New(80, 24, Config{})
	j, _ := New(80, 24, Config{PayloadLimit: 5})
	if i.t.Options.PayloadLimit != DefaultPayloadLimit || j.t.Options.PayloadLimit != 5 {
		t.Errorf("payload limits %d, %d", i.t.Options.PayloadLimit, j.t.Options.PayloadLimit)
	}
}

// TestInstanceWriteRead: Write takes the bytes, Read reads them, Resize
// refuses a size outside the limits and leaves the size, and re-caps the
// history for the new width as Resize does; Delta is the instance's cache.
func TestInstanceWriteRead(t *testing.T) {
	i, _ := New(100, 5, Config{})
	if n, err := i.Write([]byte("hello\x1b[1mbold")); n != 13 || err != nil {
		t.Fatalf("Write = %d, %v", n, err)
	}
	f, err := i.Read(Request{Attrs: true})
	if err != nil || f.Trimmed()[0] != "hellobold" || len(f.AttrRuns) != 1 {
		t.Fatalf("Read = %q, %d runs, %v", f.Trimmed()[0], len(f.AttrRuns), err)
	}
	if _, delta := i.Delta(&f, "", false); delta {
		t.Errorf("delta with no since")
	}
	if changed, delta := i.Delta(&f, f.ContentHash, false); !delta || len(changed) != 0 {
		t.Errorf("delta against itself: %v %v", changed, delta)
	}
	for _, s := range [][2]int{{0, 5}, {2001, 5}, {100, 1001}} {
		if err := i.Resize(s[0], s[1]); err == nil || i.t.Cols() != 100 || i.t.Rows() != 5 {
			t.Errorf("Resize(%d, %d): %v, now %dx%d", s[0], s[1], err, i.t.Cols(), i.t.Rows())
		}
	}
	for k := 0; k < 1100; k++ {
		i.Write([]byte("x\r\n"))
	}
	if err := i.Resize(101, 5); err != nil {
		t.Fatal(err)
	}
	if f, _ := i.Read(Request{View: ViewScrollback, Lines: 5000}); f.HistoryLines != 990 || f.Cols != 101 {
		t.Errorf("after resize to 101 cols: %d history lines, %d cols", f.HistoryLines, f.Cols)
	}
}
