package sniff

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"regexp"
	"testing"
)

// TestPostamble compares Postamble with the Python reference's _mode_reassert for every combination of the modes it
// writes (testdata/postamble.json, from testdata/postamble_oracle.py).
func TestPostamble(t *testing.T) {
	b, err := os.ReadFile("testdata/postamble.json")
	if err != nil {
		t.Fatal(err)
	}
	var f struct {
		Header map[string]any
		Rows   []struct {
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
	sha := regexp.MustCompile(`^[0-9a-f]{64}$`)
	for _, k := range []string{"agent.py sha256", "altscreen.py sha256"} {
		if v, _ := f.Header[k].(string); !sha.MatchString(v) {
			t.Errorf("header %q = %q, want a sha256", k, v)
		}
	}
	// 2 (2004) x 2 (DECCKM) x 5 (tracking, 0 included) x 3 (encoding, 0 included)
	if want := 2 * 2 * (len(MouseTrackModes) + 1) * (len(MouseEncModes) + 1); len(f.Rows) != want || f.Header["rows"] != float64(want) {
		t.Fatalf("%d rows (header %v), want %d", len(f.Rows), f.Header["rows"], want)
	}
	seen := map[Modes]bool{}
	for _, r := range f.Rows {
		m := Modes{BracketedPaste: r.BracketedPaste, DECCKM: r.AppCursor, MouseTrack: r.MouseTracking, MouseEnc: r.MouseEncoding}
		seen[m] = true
		if got := hex.EncodeToString(Postamble(m)); got != r.Hex {
			t.Errorf("Postamble(%+v)\n got %s\nwant %s", m, got, r.Hex)
		}
	}
	if len(seen) != len(f.Rows) {
		t.Errorf("%d distinct mode sets in %d rows", len(seen), len(f.Rows))
	}
}

// TestPostambleOutOfGroup: a mouse value outside its group sets nothing and resets the whole group.
func TestPostambleOutOfGroup(t *testing.T) {
	for _, m := range []Modes{{MouseTrack: 5}, {MouseTrack: 1006}, {MouseEnc: 1005}, {MouseEnc: -1}} {
		if got, want := string(Postamble(m)), string(Postamble(Modes{})); got != want {
			t.Errorf("Postamble(%+v) = %q, want %q", m, got, want)
		}
	}
	if got := string(Postamble(Modes{Alt: true})); got != string(Postamble(Modes{})) {
		t.Errorf("Alt changed the postamble: %q", got)
	}
}

// TestPostambleHeals: whatever state a viewer is in, the postamble of m leaves it in m (its alternate screen aside).
func TestPostambleHeals(t *testing.T) {
	var all []Modes
	for _, bp := range []bool{false, true} {
		for _, ck := range []bool{false, true} {
			for _, tr := range append([]int{0}, MouseTrackModes[:]...) {
				for _, en := range append([]int{0}, MouseEncModes[:]...) {
					all = append(all, Modes{BracketedPaste: bp, DECCKM: ck, MouseTrack: tr, MouseEnc: en})
				}
			}
		}
	}
	for _, from := range all {
		for _, to := range all {
			s := New(Options{})
			s.Feed([]byte("\x1b[?1049h"))
			s.Feed(Postamble(from))
			if got := s.Modes(); got != (Modes{Alt: true, BracketedPaste: from.BracketedPaste, DECCKM: from.DECCKM, MouseTrack: from.MouseTrack, MouseEnc: from.MouseEnc}) {
				t.Fatalf("Postamble(%+v) gave %+v", from, got)
			}
			s.Feed(Postamble(to))
			want := to
			want.Alt = true
			if got := s.Modes(); got != want {
				t.Fatalf("from %+v, Postamble(%+v) gave %+v", from, to, got)
			}
		}
	}
}
