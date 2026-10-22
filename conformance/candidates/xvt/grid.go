// Package xvt adapts github.com/charmbracelet/x/vt to the conformance runner (#1 D7). It reports only what the
// emulator's public API returns (E3). Every path that API can't show is listed in State.Unobserved.
package xvt

import (
	"fmt"

	"github.com/JohnConnorNPC/go-term/conformance/runner"
	uv "github.com/charmbracelet/ultraviolet"
	vt "github.com/charmbracelet/x/vt"
)

// gridUnobserved: x/vt has no wrap flag (w), no viewport (viewportY), and no table rows for its colours (attrs, E10).
var gridUnobserved = []string{"attrs", "viewportY", "w"}

// Engine is one x/vt emulator.
type Engine struct {
	e  *vt.Emulator
	ev events
}

// New is the runner.Factory. x/vt has no unicode-version option, so only the goldens' default "6" is accepted.
func New(c runner.Config) (runner.Engine, error) {
	if c.Unicode != "6" {
		return nil, fmt.Errorf("x/vt has no unicode %q option", c.Unicode)
	}
	e := &Engine{e: vt.NewEmulator(c.Cols, c.Rows)}
	e.e.SetScrollbackSize(c.Scrollback)
	e.ev.listen(e.e)
	return e, nil
}

func (e *Engine) Write(p []byte) error {
	_, err := e.e.Write(p)
	return err
}

func (e *Engine) Resize(cols, rows int) error {
	e.e.Resize(cols, rows)
	return nil
}

// Reset is RIS (rulings-3 Q6): x/vt has no reset call.
func (e *Engine) Reset() error {
	return e.Write([]byte("\x1bc"))
}

// Snapshot reads the state after the writes so far.
func (e *Engine) Snapshot() (runner.State, error) {
	var st runner.State
	err := e.grid(&st)
	if ferr := e.ev.fill(&st); err == nil {
		err = ferr
	}
	return st, err
}

// grid fills the buffer fields. The normal screen can't be read while the alternate one is active, so then Normal
// and normalLength are not reported (normalLength unobserved) instead of being made up.
func (e *Engine) grid(st *runner.State) error {
	p := e.e.CursorPosition()
	st.P.Cols, st.P.Rows, st.P.Cursor, st.P.BaseY = e.e.Width(), e.e.Height(), [2]int{p.X, p.Y}, e.e.ScrollbackLen()
	st.Unobserved = append(st.Unobserved, gridUnobserved...)
	st.P.Buffer = "normal"
	screen, err := rows(e.e.Width(), e.e.Height(), e.e.CellAt)
	if e.e.IsAltScreen() {
		st.P.Buffer, st.Alt = "alternate", screen
		st.Unobserved = append(st.Unobserved, "normalLength")
		return err
	}
	back, err2 := rows(e.e.Width(), e.e.ScrollbackLen(), e.e.ScrollbackCellAt)
	st.Normal = append(back, screen...)
	st.P.NormalLength = len(st.Normal)
	if err == nil {
		err = err2
	}
	return err
}

// rows reads h rows of w cells through at: text and width by identity, attributes left zero (unobserved). A cell
// the API doesn't return is an error, never an empty cell.
func rows(w, h int, at func(x, y int) *uv.Cell) ([]runner.Row, error) {
	out := make([]runner.Row, h)
	for y := range out {
		for x := 0; x < w; x++ {
			cell := at(x, y)
			if cell == nil {
				return nil, fmt.Errorf("x/vt returned no cell at %d,%d", x, y)
			}
			out[y].Cells = append(out[y].Cells, runner.Cell{Ch: cell.Content, W: cell.Width})
		}
	}
	return out, nil
}
