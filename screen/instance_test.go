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
