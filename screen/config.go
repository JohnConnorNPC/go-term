package screen

import (
	"fmt"

	"github.com/JohnConnorNPC/go-term/vt"
)

const (
	// DefaultScrollbackLines is xterm.js's default scrollback and Python's
	// _MAX_HISTORY_LINES.
	DefaultScrollbackLines = 1000
	// DefaultScrollbackCells is Python's _MAX_HISTORY_CELLS: history holds at
	// most this many cells, so a wide grid keeps fewer lines.
	DefaultScrollbackCells = 100000
)

// Config is how Studio sizes a terminal's history and widths.
type Config struct {
	// ScrollbackLines and ScrollbackCells bound the history: HistoryCap.
	// Zero or less takes the default.
	ScrollbackLines, ScrollbackCells int
	// Unicode is the width table: "6" (or ""), what the page uses (no
	// unicode11 addon). Version 11 is not supported yet.
	Unicode string
}

func (c Config) lines() int {
	if c.ScrollbackLines <= 0 {
		return DefaultScrollbackLines
	}
	return c.ScrollbackLines
}

func (c Config) cells() int {
	if c.ScrollbackCells <= 0 {
		return DefaultScrollbackCells
	}
	return c.ScrollbackCells
}

// HistoryCap is how many history lines a grid cols wide keeps, as Python's
// textgrid hist_cap: min(lines, max(1, cells/cols)). At the defaults that is
// 1000 lines up to 100 cols, 990 at 101 cols, 100 at 1000 cols.
func (c Config) HistoryCap(cols int) int {
	return min(c.lines(), max(1, c.cells()/max(1, cols)))
}

func (c Config) check() error {
	if c.Unicode != "" && c.Unicode != "6" {
		return fmt.Errorf("screen: Unicode %q: only version 6 widths are supported (v11 is a later issue)", c.Unicode)
	}
	return nil
}

// NewTerminal returns a vt terminal of cols x rows whose scrollback is
// c.HistoryCap(cols). The terminal's Options are vt's defaults otherwise.
func NewTerminal(cols, rows int, c Config) (*vt.Terminal, error) {
	if err := c.check(); err != nil {
		return nil, err
	}
	o := vt.NewOptions()
	o.Cols, o.Rows = max(cols, vt.MinimumCols), max(rows, vt.MinimumRows)
	o.Scrollback = c.HistoryCap(o.Cols)
	return vt.NewTerminal(o), nil
}

// Resize resizes t to cols x rows and re-applies c.HistoryCap for the new
// width: a wider grid keeps fewer history lines, and the oldest go.
func Resize(t *vt.Terminal, cols, rows int, c Config) error {
	if err := c.check(); err != nil {
		return err
	}
	t.Options.Scrollback = c.HistoryCap(max(cols, vt.MinimumCols))
	t.Resize(cols, rows)
	return nil
}
