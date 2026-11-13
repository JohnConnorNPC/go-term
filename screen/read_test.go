package screen

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/JohnConnorNPC/go-term/vt"
)

func newTerm(cols, rows int) *vt.Terminal {
	opts := vt.NewOptions()
	opts.Cols, opts.Rows = cols, rows
	return vt.NewTerminal(opts)
}

func TestReadGrid(t *testing.T) {
	term := newTerm(6, 3)
	term.WriteString("ab\r\n\x1b[?1h\x1b[?25l")
	f := Read(term, Request{})
	if got := f.Text(); got != "ab    \n      \n      " {
		t.Errorf("text %q, want rows padded to 6", got)
	}
	if f.Cols != 6 || f.Rows != 3 || f.View != ViewScreen || f.HistoryLines != 0 {
		t.Errorf("frame %+v", f)
	}
	if *f.Cursor != (Cursor{Row: 1, Col: 0}) || !f.CursorHidden || !f.AppCursor || f.AltScreen {
		t.Errorf("cursor %+v hidden %v app %v alt %v", *f.Cursor, f.CursorHidden, f.AppCursor, f.AltScreen)
	}
	if f.ContentHash != ContentHash(f.Text()) {
		t.Error("content hash is not ContentHash(Text())")
	}
	if f.StableHash != StableHash(f.Text(), f.Cursor, 0) || f.StableHash == f.ContentHash {
		t.Error("stable hash: the cursor cell (a blank, inside the padded row) is not masked")
	}
	if got := strings.Join(f.Trimmed(), "|"); got != "ab||" {
		t.Errorf("Trimmed %q", got)
	}
}

// TestReadPendingWrap: after writing the last column the cursor sits at
// column == cols (pending wrap); it reads as cols-1, as Python clamps it.
func TestReadPendingWrap(t *testing.T) {
	term := newTerm(4, 2)
	term.WriteString("abcd")
	if x := term.Buffer().X; x != 4 {
		t.Fatalf("vt cursor x %d, want 4 (pending wrap)", x)
	}
	f := Read(term, Request{})
	if *f.Cursor != (Cursor{Row: 0, Col: 3}) {
		t.Errorf("cursor %+v, want 0,3", *f.Cursor)
	}
	if f.StableHash != ContentHash("abc\x00\n    ") {
		t.Error("the clamped cursor does not mask the last cell")
	}
}

func TestReadScrollback(t *testing.T) {
	term := newTerm(5, 2)
	term.WriteString("l1\r\nl2\r\nl3\r\nl4\r\nl5")
	f := Read(term, Request{View: ViewScrollback, Lines: 2})
	if f.View != ViewScrollback || f.HistoryLines != 2 {
		t.Fatalf("view %v history %d", f.View, f.HistoryLines)
	}
	if got := strings.Join(f.Trimmed(), "|"); got != "l2|l3|l4|l5" {
		t.Errorf("lines %q", got)
	}
	if *f.Cursor != (Cursor{Row: 1, Col: 2}) {
		t.Errorf("cursor %+v is not grid-relative", *f.Cursor)
	}
	if f.StableHash != StableHash(f.Text(), f.Cursor, f.HistoryLines) {
		t.Error("stable hash not shifted by HistoryLines")
	}
	if g := Read(term, Request{View: ViewScrollback, Lines: 99}); g.HistoryLines != 3 {
		t.Errorf("Lines 99: history %d, want all 3", g.HistoryLines)
	}
	if g := Read(term, Request{View: ViewScrollback}); g.HistoryLines != 0 || g.View != ViewScrollback || len(g.Lines) != 2 {
		t.Errorf("Lines 0: %+v", g)
	}
	term.ScrollLines(-2) // a user-scrolled viewport does not move the read
	if g := Read(term, Request{}); strings.Join(g.Trimmed(), "|") != "l4|l5" {
		t.Errorf("scrolled viewport changed the read: %q", g.Lines)
	}
	term.WriteString("\x1b[?1049h")
	g := Read(term, Request{View: ViewScrollback, Lines: 2})
	if g.View != ViewScreen || g.HistoryLines != 0 || !g.AltScreen || len(g.Lines) != 2 {
		t.Errorf("alt screen: view %v history %d alt %v", g.View, g.HistoryLines, g.AltScreen)
	}
}

// hlState is the part of a headless vector state this package reads.
type hlState struct {
	Buffer string `json:"buffer"`
	Cursor struct {
		Row int `json:"row"`
		Col int `json:"col"`
	} `json:"cursor"`
	BaseY   int      `json:"base_y"`
	Rows    []string `json:"rows"`
	History []string `json:"history"`
	Attrs   []struct {
		Row   int      `json:"row"`
		Col   int      `json:"col"`
		Len   int      `json:"len"`
		Flags []string `json:"flags"`
	} `json:"attrs"`
}

type hlCase struct {
	ID          string  `json:"id"`
	Cols        int     `json:"cols"`
	Rows        int     `json:"rows"`
	InputSHA256 string  `json:"input_sha256"`
	Final       hlState `json:"final"`
}

// replayHeadless feeds a headless case's input (the python case's replay
// chunks, matched by input_sha256) chunk by chunk into a fresh terminal at the
// case's geometry, as Studio's replay.js does.
func replayHeadless(t *testing.T, c hlCase, py map[string]pyCase) *vt.Terminal {
	t.Helper()
	p, ok := py[c.ID]
	if !ok {
		t.Fatalf("%s: no python case", c.ID)
	}
	sum := sha256.Sum256(p.Replay)
	if hex.EncodeToString(sum[:]) != c.InputSHA256 {
		t.Fatalf("%s: replay sha256 %x, headless input_sha256 %s", c.ID, sum, c.InputSHA256)
	}
	var r struct {
		Chunks []string `json:"chunks_b64"`
	}
	mustUnmarshal(t, p.Replay, &r)
	term := newTerm(c.Cols, c.Rows)
	for _, ch := range r.Chunks {
		b, err := base64.StdEncoding.DecodeString(ch)
		if err != nil {
			t.Fatalf("%s: %v", c.ID, err)
		}
		term.Write(b)
	}
	return term
}

// checkHeadless compares a Read of term with a headless final state: grid
// rows (trimmed), cursor (clamped as Read does), the active buffer and, on
// the normal buffer, the whole history through a scrollback read.
func checkHeadless(term *vt.Terminal, want hlState) error {
	f := Read(term, Request{})
	if got := strings.Join(f.Trimmed(), "\n"); got != strings.Join(want.Rows, "\n") {
		return fmt.Errorf("rows %q, headless %q", f.Trimmed(), want.Rows)
	}
	wc := Cursor{Row: clamp(want.Cursor.Row, 0, f.Rows-1), Col: clamp(want.Cursor.Col, 0, f.Cols-1)}
	if *f.Cursor != wc {
		return fmt.Errorf("cursor %+v, headless %+v", *f.Cursor, wc)
	}
	if f.AltScreen != (want.Buffer == "alternate") {
		return fmt.Errorf("alt %v, headless buffer %s", f.AltScreen, want.Buffer)
	}
	if !f.AltScreen {
		s := Read(term, Request{View: ViewScrollback, Lines: 1 << 30})
		h := s.Trimmed()[:s.HistoryLines]
		if s.HistoryLines != want.BaseY || strings.Join(h, "\n") != strings.Join(want.History, "\n") {
			return fmt.Errorf("history %d %q, headless base_y %d %q", s.HistoryLines, h, want.BaseY, want.History)
		}
	}
	return nil
}

// knownHeadless are the headless cases go-term's vt is known to miss, with
// the reason. Each must still miss, so a fix shows up here.
var knownHeadless = map[string]string{
	// vt collects APC (ESC _) as an OSC string, so BEL ends it and "ok"
	// prints; in 6.0.0 APC is the SOS/PM/APC ignore state, which only ST
	// ends, so "ok" is swallowed (an engine fix, not the read API).
	"tests/test_textgrid.py::test_control_string_families_are_skipped#2": "APC ended by BEL",
}

func TestHeadlessVectors(t *testing.T) {
	var report []string
	known := 0
	for _, fam := range []struct {
		name string
		want int
	}{{"textgrid", 115}, {"hashes", 74}} {
		py := map[string]pyCase{}
		for _, c := range loadPython(t, "python/"+fam.name+".json") {
			py[c.ID] = c
		}
		var f struct {
			Cases []hlCase `json:"cases"`
		}
		if err := json.Unmarshal(readVectors(t, "headless/"+fam.name+".json"), &f); err != nil {
			t.Fatal(err)
		}
		ok := 0
		for _, c := range f.Cases {
			err := checkHeadless(replayHeadless(t, c, py), c.Final)
			if why, isKnown := knownHeadless[c.ID]; isKnown {
				if err == nil {
					t.Errorf("headless/%s %s: known miss (%s) now passes: drop it from knownHeadless", fam.name, c.ID, why)
				}
				known++
				continue
			}
			if err != nil {
				t.Errorf("headless/%s %s: %v", fam.name, c.ID, err)
				continue
			}
			ok++
		}
		if len(f.Cases) != fam.want {
			t.Errorf("headless/%s: %d cases, want %d", fam.name, len(f.Cases), fam.want)
		}
		report = append(report, fmt.Sprintf("%s %d/%d", fam.name, ok, len(f.Cases)))
	}
	if known != len(knownHeadless) {
		t.Errorf("%d known misses seen, %d listed", known, len(knownHeadless))
	}
	t.Logf("VECTORS headless %s (known misses %d, not counted as passes)", strings.Join(report, " "), known)
}
