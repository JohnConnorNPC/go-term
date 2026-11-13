package screen

import (
	"strings"

	"github.com/JohnConnorNPC/go-term/sniff"
	"github.com/JohnConnorNPC/go-term/vt"
)

// View is what a read covers: the grid alone, or history rows above it.
type View int

const (
	// ViewScreen reads the grid (the live screen of the active buffer).
	ViewScreen View = iota
	// ViewScrollback prepends up to Request.Lines rows of the normal buffer's
	// history. On the alt screen it reads as ViewScreen (a full-screen app's
	// grid is the whole story), and Frame.View says so.
	ViewScrollback
)

// String gives Studio's view names, "screen" and "scrollback".
func (v View) String() string {
	if v == ViewScrollback {
		return "scrollback"
	}
	return "screen"
}

// Request selects what Read returns.
type Request struct {
	View  View
	Lines int  // history rows wanted on ViewScrollback; <= 0 gives none
	Attrs bool // fill Frame.AttrRuns (grid rows only, also on ViewScrollback)
}

// Frame is one typed read of a terminal: Studio's screen_text fields.
type Frame struct {
	// Lines are HistoryLines history rows, oldest first, then Rows grid rows.
	// Each is padded with spaces to the full width, as Python's textgrid and
	// pyte render it: a wide glyph is 1 rune over 2 cells, so its row is 1
	// rune shorter; a cell holding a combining sequence gives all its runes.
	// Use Trimmed for right-trimmed rows (xterm's translateToString(true)).
	Lines []string
	// Cols and Rows are the grid size.
	Cols, Rows int
	// View is the view produced: ViewScrollback reads as ViewScreen on the
	// alt screen.
	View View
	// Cursor is grid-relative (history rows are not counted), clamped to the
	// grid: a cursor in the pending-wrap position (column == Cols) reads as
	// column Cols-1. Read always sets it; nil is for a caller's own degraded
	// frame, and gives StableHash == ContentHash.
	Cursor *Cursor
	// CursorHidden is DECTCEM off (CSI ? 25 l). The cursor is still reported.
	CursorHidden bool
	// AltScreen is true while the alt buffer is active; AppCursor is DECCKM
	// (application cursor keys).
	AltScreen, AppCursor bool
	// BracketedPaste is mode 2004. MouseTrack is the mouse tracking mode (0
	// off, 9, 1000, 1002 or 1003) and MouseEnc the report encoding (0 the
	// default, 1006 or 1016), numbered as sniff.Modes; Modes gives all five.
	BracketedPaste       bool
	MouseTrack, MouseEnc int
	// Title is the window title the last OSC 0 or 2 set; "" when none did.
	Title string
	// HistoryLines is how many of Lines are history rows.
	HistoryLines int
	// AttrRuns are AttrRuns(t, MaxAttrRuns): the styled runs of the grid rows
	// (Row is a grid row, history is never styled); nil unless Request.Attrs.
	AttrRuns []AttrRun
	// ContentHash is ContentHash(Text()). StableHash is the same text with the
	// cursor's glyph replaced by one NUL: the glyph whose cells cover the
	// cursor cell (a wide glyph covers 2 cells, a combining sequence is all
	// its runes), which for 1-cell 1-rune glyphs is exactly
	// StableHash(Text(), Cursor-as-rune-index, HistoryLines).
	ContentHash, StableHash string
}

// Modes is the frame's DEC private-mode state in the sniff package's terms,
// the state a restore re-asserts.
func (f *Frame) Modes() sniff.Modes {
	return sniff.Modes{Alt: f.AltScreen, DECCKM: f.AppCursor, BracketedPaste: f.BracketedPaste,
		MouseTrack: f.MouseTrack, MouseEnc: f.MouseEnc}
}

// vt's mouse protocol and encoding names as DECSET numbers; any other name
// (NONE, DEFAULT) is 0.
var (
	mouseTrackModes = map[string]int{"X10": 9, "VT200": 1000, "DRAG": 1002, "ANY": 1003}
	mouseEncModes   = map[string]int{"SGR": 1006, "SGR_PIXELS": 1016}
)

// Text is Studio's text: Lines joined with "\n".
func (f *Frame) Text() string { return strings.Join(f.Lines, "\n") }

// Trimmed returns Lines with trailing spaces removed from each row.
func (f *Frame) Trimmed() []string {
	out := make([]string, len(f.Lines))
	for i, l := range f.Lines {
		out[i] = strings.TrimRight(l, " ")
	}
	return out
}

// Read reads t: the live screen of the active buffer at its base (not a
// viewport the user scrolled), the cursor, the modes, history rows on
// ViewScrollback, the title, the hashes, and attribute runs when r.Attrs. It does not
// lock t.
func Read(t *vt.Terminal, r Request) Frame {
	bufs := t.Buffers()
	active := bufs.Active()
	cols, rows := t.Cols(), t.Rows()
	core := t.CoreService()
	f := Frame{
		Cols:         cols,
		Rows:         rows,
		View:         ViewScreen,
		CursorHidden: core.IsCursorHidden,
		AltScreen:    active == bufs.Alt(),
		AppCursor:    core.DecPrivateModes.ApplicationCursorKeys,

		BracketedPaste: core.DecPrivateModes.BracketedPasteMode,
		MouseTrack:     mouseTrackModes[t.MouseService().ActiveProtocol()],
		MouseEnc:       mouseEncModes[t.MouseService().ActiveEncoding()],
		Title:          t.InputHandler().WindowTitle(),
	}
	f.Cursor = &Cursor{Row: clamp(active.Y, 0, rows-1), Col: clamp(active.X, 0, cols-1)}
	if r.View == ViewScrollback && !f.AltScreen {
		f.View = ViewScrollback
		normal := bufs.Normal()
		if r.Lines > 0 {
			f.HistoryLines = min(r.Lines, normal.YBase)
		}
		for i := normal.YBase - f.HistoryLines; i < normal.YBase; i++ {
			s, _, _ := rowText(normal.Lines.Get(i), -1)
			f.Lines = append(f.Lines, s)
		}
	}
	maskRowIdx, maskStart, maskEnd := -1, 0, 0
	for y := 0; y < rows; y++ {
		cx := -1
		if y == f.Cursor.Row {
			cx = f.Cursor.Col
		}
		s, start, end := rowText(active.Lines.Get(active.YBase+y), cx)
		if cx >= 0 && start >= 0 {
			maskRowIdx, maskStart, maskEnd = len(f.Lines), start, end
		}
		f.Lines = append(f.Lines, s)
	}
	f.ContentHash = ContentHash(f.Text())
	f.StableHash = f.ContentHash
	if maskRowIdx >= 0 {
		f.StableHash = ContentHash(maskRow(f.Lines, maskRowIdx, maskStart, maskEnd))
	}
	if r.Attrs {
		f.AttrRuns = AttrRuns(t, MaxAttrRuns)
	}
	return f
}

// rowText renders a buffer line as xterm's translateToString(false) does:
// each glyph's string (a blank cell is a space), stepping over the cells a
// wide glyph covers. For cell column cx >= 0 it also returns the byte span
// [start, end) of the glyph covering that cell; start is -1 when no glyph
// covers it.
func rowText(line *vt.BufferLine, cx int) (s string, start, end int) {
	start = -1
	if line == nil {
		return "", start, 0
	}
	var b strings.Builder
	b.Grow(line.Length)
	for x := 0; x < line.Length; {
		g := line.GetString(x)
		if g == "" {
			g = string(rune(vt.WhitespaceCellChar))
		}
		adv := line.GetWidth(x)
		if adv < 1 {
			adv = 1
		}
		if cx >= x && cx < x+adv {
			start, end = b.Len(), b.Len()+len(g)
		}
		b.WriteString(g)
		x += adv
	}
	return b.String(), start, end
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
