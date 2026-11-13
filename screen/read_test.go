package screen

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
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
	BaseY   int            `json:"base_y"`
	Rows    []string       `json:"rows"`
	History []string       `json:"history"`
	Attrs   []hlRun        `json:"attrs"`
	Modes   map[string]any `json:"modes"`
	Title   *string        `json:"title"`
}

// hlRun is one of replay.js's attrs runs: a maximal same-style run of
// viewport cells, default cells skipped; fg and bg are xterm's [mode, value]
// ("default" -1, "palette" index, "rgb" packed) and flags the set is*() bits.
type hlRun struct {
	Row   int      `json:"row"`
	Col   int      `json:"col"`
	Len   int      `json:"len"`
	Fg    [2]any   `json:"fg"`
	Bg    [2]any   `json:"bg"`
	Flags []string `json:"flags"`
}

type hlCase struct {
	ID          string  `json:"id"`
	Cols        int     `json:"cols"`
	Rows        int     `json:"rows"`
	InputSHA256 string  `json:"input_sha256"`
	Final       hlState `json:"final"`
	// OutputState is headless over the python case's output bytes, written
	// whole at the case geometry (present when that output is bytes).
	OutputState *hlState `json:"output_state"`
}

// pyBytes decodes a python vector value {"$b": base64}; ok is false for any
// other encoding ($digest, null, a dict).
func pyBytes(raw json.RawMessage) (b []byte, ok bool) {
	var v struct {
		B *string `json:"$b"`
	}
	if json.Unmarshal(raw, &v) != nil || v.B == nil {
		return nil, false
	}
	b, err := base64.StdEncoding.DecodeString(*v.B)
	return b, err == nil
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

// hlAttrs is replay.js attrs() over vt's cells: the viewport row by row, a
// cell with default fg, bg and no flag skipped, a run extended while the next
// cell has the same style. Cells past a line's length read as default.
func hlAttrs(term *vt.Terminal) []hlRun {
	b := term.Buffer()
	colour := func(mode uint32, v int) [2]any {
		switch mode {
		case vt.AttrCMDefault:
			return [2]any{"default", float64(-1)}
		case vt.AttrCMRGB:
			return [2]any{"rgb", float64(v)}
		}
		return [2]any{"palette", float64(v)}
	}
	var runs []hlRun
	cell := vt.NewCellData()
	for y := 0; y < term.Rows(); y++ {
		line := b.Lines.Get(b.YBase + y)
		var cur *hlRun
		for x := 0; x < term.Cols(); x++ {
			if line == nil || x >= line.Length {
				cur = nil
				continue
			}
			line.LoadCell(x, cell)
			var flags []string
			for _, f := range []struct {
				name string
				on   bool
			}{{"bold", cell.IsBold()}, {"dim", cell.IsDim()}, {"italic", cell.IsItalic()},
				{"underline", cell.IsUnderline()}, {"blink", cell.IsBlink()}, {"inverse", cell.IsInverse()},
				{"invisible", cell.IsInvisible()}, {"strikethrough", cell.IsStrikethrough()},
				{"overline", cell.IsOverline()}} {
				if f.on {
					flags = append(flags, f.name)
				}
			}
			r := hlRun{Row: y, Col: x, Len: 1, Fg: colour(cell.GetFgColorMode(), cell.GetFgColor()),
				Bg: colour(cell.GetBgColorMode(), cell.GetBgColor()), Flags: flags}
			if r.Fg[0] == "default" && r.Bg[0] == "default" && len(flags) == 0 {
				cur = nil
				continue
			}
			if cur != nil && fmt.Sprint(cur.Fg, cur.Bg, cur.Flags) == fmt.Sprint(r.Fg, r.Bg, r.Flags) {
				cur.Len++
				continue
			}
			runs = append(runs, r)
			cur = &runs[len(runs)-1]
		}
	}
	return runs
}

// hlModes is replay.js modes(): xterm.js's t.modes under its key names plus
// the mouse encoding, from vt's state (as conformance/goterm reports P.modes).
func hlModes(term *vt.Terminal) map[string]any {
	cs := term.CoreService()
	d := cs.DecPrivateModes
	return map[string]any{"applicationCursorKeysMode": d.ApplicationCursorKeys, "applicationKeypadMode": d.ApplicationKeypad,
		"bracketedPasteMode": d.BracketedPasteMode, "insertMode": cs.Modes.InsertMode,
		"mouseTrackingMode": strings.ToLower(term.MouseService().ActiveProtocol()), "originMode": d.Origin,
		"reverseWraparoundMode": d.ReverseWraparound, "sendFocusMode": d.SendFocus,
		"synchronizedOutputMode": d.SynchronizedOutput, "wraparoundMode": d.Wraparound,
		"mouse_encoding": term.MouseService().ActiveEncoding()}
}

// hlLines is replay.js's translateToString(true) over lines [from, to) of
// buf: xterm trims only never-written cells, so written spaces stay.
func hlLines(buf *vt.Buffer, from, to int) []string {
	out := []string{}
	for y := from; y < to; y++ {
		s := ""
		if l := buf.Lines.Get(y); l != nil {
			s = l.TranslateToString(true, 0, -1)
		}
		out = append(out, s)
	}
	return out
}

// checkHeadless compares term with a headless final state: the grid rows as
// xterm gives them (vt's translateToString) and as Read gives them (Trimmed,
// against the headless rows space-trimmed), the cursor (clamped as Read
// does), the active buffer, on the normal buffer the whole history (vt's
// lines and a scrollback Read), the attrs runs, the modes and the title
// (headless null, no title event, reads as vt's "").
func checkHeadless(term *vt.Terminal, want hlState) error {
	f := Read(term, Request{})
	b := term.Buffer()
	if got := hlLines(b, b.YBase, b.YBase+f.Rows); !reflect.DeepEqual(got, want.Rows) {
		return fmt.Errorf("rows %q, headless %q", got, want.Rows)
	}
	if got := strings.Join(f.Trimmed(), "\n"); got != trimAll(want.Rows) {
		return fmt.Errorf("Read rows %q, headless %q", f.Trimmed(), want.Rows)
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
		if s.HistoryLines != want.BaseY || strings.Join(h, "\n") != trimAll(want.History) {
			return fmt.Errorf("history %d %q, headless base_y %d %q", s.HistoryLines, h, want.BaseY, want.History)
		}
		if got := hlLines(term.Buffers().Normal(), 0, s.HistoryLines); !reflect.DeepEqual(got, append([]string{}, want.History...)) {
			return fmt.Errorf("history lines %q, headless %q", got, want.History)
		}
	}
	if got := hlAttrs(term); !reflect.DeepEqual(normRuns(got), normRuns(want.Attrs)) {
		return fmt.Errorf("attrs %v, headless %v", got, want.Attrs)
	}
	if got := hlModes(term); !reflect.DeepEqual(got, want.Modes) {
		return fmt.Errorf("modes %v, headless %v", got, want.Modes)
	}
	wt := ""
	if want.Title != nil {
		wt = *want.Title
	}
	if f.Title != wt {
		return fmt.Errorf("title %q, headless %q", f.Title, wt)
	}
	return nil
}

// normRuns makes an empty flag list and an absent one compare equal.
func normRuns(rs []hlRun) []hlRun {
	out := make([]hlRun, len(rs))
	for i, r := range rs {
		if len(r.Flags) == 0 {
			r.Flags = nil
		}
		out[i] = r
	}
	return out
}

// knownHeadless are the headless cases go-term's vt is known to miss, with
// the reason. Each must still miss, so a fix shows up here.
var knownHeadless = map[string]string{}

// TestHeadlessVectors replays every A2b case (Studio's vectors/headless: the
// A2a inputs through @xterm/headless 6.0.0) chunk by chunk and compares the
// final state, and feeds every case's Python output bytes whole and compares
// output_state (what the page shows for a pyte re-emit or a raw snapshot).
// The sniffer families' files are sniff's (sniffData).
func TestHeadlessVectors(t *testing.T) {
	var report []string
	known, outOK, outN := 0, 0, 0
	for _, fam := range []struct {
		name, dir string
		want      int
	}{{"textgrid", "", 115}, {"hashes", "", 74}, {"colours", "", 53}, {"pyte", "", 26}, {"raw", "", 32},
		{"synthetic", sniffData, 57}, {"altscreen", sniffData, 118}, {"titles", sniffData, 80}} {
		py := map[string]pyCase{}
		for _, c := range loadPython(t, fam.dir+"python/"+fam.name+".json") {
			py[c.ID] = c
		}
		var f struct {
			Cases []hlCase `json:"cases"`
		}
		if err := json.Unmarshal(readVectors(t, fam.dir+"headless/"+fam.name+".json"), &f); err != nil {
			t.Fatal(err)
		}
		ok := 0
		for _, c := range f.Cases {
			if c.OutputState != nil {
				outN++
				out, isBytes := pyBytes(py[c.ID].Output)
				if !isBytes {
					t.Errorf("headless/%s %s: output_state but the python output is not bytes", fam.name, c.ID)
				} else {
					term := newTerm(c.Cols, c.Rows)
					term.Write(out)
					if err := checkHeadless(term, *c.OutputState); err != nil {
						t.Errorf("headless/%s %s output_state: %v", fam.name, c.ID, err)
					} else {
						outOK++
					}
				}
			}
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
	// divergence.json header "compared_with_output_state".
	if outN != 115 {
		t.Errorf("%d output_state cases, want 115", outN)
	}
	t.Logf("VECTORS headless %s output_state %d/%d (known misses %d, not counted as passes)",
		strings.Join(report, " "), outOK, outN, known)
}
