package screen

import "github.com/JohnConnorNPC/go-term/vt"

// MaxAttrRuns caps AttrRuns as Python's _MAX_ATTR_RUNS does.
const MaxAttrRuns = 2000

// AttrRun is a maximal horizontal run of cells with one style that differs
// from the default style: Studio's attr_runs entry. Row and Col are 0-based
// grid positions and Len counts cells (a wide glyph's 2 cells count 2).
type AttrRun struct {
	Row, Col, Len int
	// Fg and Bg are Colour names.
	Fg, Bg string
	// Reverse is SGR 7, Bold SGR 1, Underscore any underline (SGR 4, 21 and
	// the 4:x styles: pyte only knows SGR 4).
	Reverse, Bold, Underscore bool
}

type style struct {
	fg, bg                    string
	reverse, bold, underscore bool
}

var defaultStyle = style{fg: Default, bg: Default}

// AttrRuns lists the styled runs of t's grid (the live screen of the active
// buffer, never history) as Python's pyte_snap.attr_runs does: row by row,
// cells whose style (fg, bg, reverse, bold, underscore) is the default are
// skipped, and a run extends while the next cell has the same style; the
// list stops after max runs. Other attributes (dim, italic, blink, ...) are
// not part of the style, so dim text alone gives no run.
func AttrRuns(t *vt.Terminal, max int) []AttrRun {
	b := t.Buffer()
	cols, rows := t.Cols(), t.Rows()
	var runs []AttrRun
	cell := vt.NewCellData()
	for y := 0; y < rows && len(runs) < max; y++ {
		line := b.Lines.Get(b.YBase + y)
		if line == nil {
			continue
		}
		n := min(cols, line.Length)
		styleAt := func(x int) style {
			line.LoadCell(x, cell)
			return style{Colour(cell.Fg), Colour(cell.Bg), cell.IsInverse(), cell.IsBold(), cell.IsUnderline()}
		}
		for x := 0; x < n; {
			s := styleAt(x)
			if s == defaultStyle {
				x++
				continue
			}
			start := x
			for x++; x < n && styleAt(x) == s; x++ {
			}
			runs = append(runs, AttrRun{Row: y, Col: start, Len: x - start, Fg: s.fg, Bg: s.bg,
				Reverse: s.reverse, Bold: s.bold, Underscore: s.underscore})
			if len(runs) >= max {
				break
			}
		}
	}
	return runs
}
