// Package gitpod adapts github.com/gitpod-io/xterm-go to the conformance runner (#1 D7). It reports only what the
// terminal's public API returns (E3): attributes go through norm.
package gitpod

import (
	"fmt"

	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/norm"
	"github.com/JohnConnorNPC/go-term/conformance/runner"
	xterm "github.com/gitpod-io/xterm-go"
)

// Engine is one gitpod terminal.
type Engine struct {
	t *xterm.Terminal
}

// New is the runner.Factory. gitpod has no unicode-version option, so only the goldens' default "6" is accepted.
func New(c runner.Config) (runner.Engine, error) {
	if c.Unicode != "6" {
		return nil, fmt.Errorf("gitpod has no unicode %q option", c.Unicode)
	}
	return &Engine{t: xterm.New(xterm.WithCols(c.Cols), xterm.WithRows(c.Rows), xterm.WithScrollback(c.Scrollback))}, nil
}

func (e *Engine) Write(p []byte) error {
	_, err := e.t.Write(p)
	return err
}

func (e *Engine) Resize(cols, rows int) error {
	e.t.Resize(cols, rows)
	return nil
}

func (e *Engine) Reset() error {
	e.t.Reset()
	return nil
}

// Snapshot reads the state after the writes so far.
func (e *Engine) Snapshot() (runner.State, error) {
	var st runner.State
	e.grid(&st)
	return st, nil
}

// grid fills the buffer fields: the raw cursor (not clamped), the normal buffer always, the alternate one while active.
func (e *Engine) grid(st *runner.State) {
	b := e.t.Buffer()
	st.P.Cols, st.P.Rows, st.P.Cursor = e.t.Cols(), e.t.Rows(), [2]int{b.X, b.Y}
	st.P.BaseY, st.P.ViewportY, st.P.Buffer = b.YBase, b.YDisp, "normal"
	st.Normal = rows(e.t.NormalBuffer())
	st.P.NormalLength = len(st.Normal)
	if e.t.IsAltBufferActive() {
		st.P.Buffer, st.Alt = "alternate", rows(e.t.AltBuffer())
	}
}

func rows(b *xterm.Buffer) []runner.Row {
	out := make([]runner.Row, b.Lines.Length())
	cell := xterm.NewCellData()
	for y := range out {
		l := b.Lines.Get(y)
		out[y].Wrapped = l.IsWrapped
		for x := 0; x < l.Len; x++ {
			l.LoadCell(x, cell)
			out[y].Cells = append(out[y].Cells, runner.Cell{Ch: cell.GetChars(), W: cell.GetWidth(), A: attr(cell)})
		}
	}
	return out
}

// attr copies the cell's IBufferCell-style getters into norm.Raw; gitpod's flag getters return the flag bit (norm.On).
func attr(c *xterm.CellData) runner.Attr {
	return norm.Attr(norm.Raw{
		FgMode: c.GetFgColorMode(), BgMode: c.GetBgColorMode(), Fg: c.GetFgColor(), Bg: c.GetBgColor(),
		Bold: norm.On(c.IsBold()), Italic: norm.On(c.IsItalic()), Dim: norm.On(c.IsDim()), Underline: norm.On(c.IsUnderline()),
		Blink: norm.On(c.IsBlink()), Inverse: norm.On(c.IsInverse()), Invisible: norm.On(c.IsInvisible()),
		Strike: norm.On(c.IsStrikethrough()), Overline: norm.On(c.IsOverline()),
	})
}
