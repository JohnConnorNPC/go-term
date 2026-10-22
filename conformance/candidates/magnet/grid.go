// Package magnet adapts github.com/0magnet/xterm-go/vt to the conformance runner (#1 D7). It reports only what the
// terminal's public API returns (E3): attributes go through norm.
package magnet

import (
	"fmt"

	vt "github.com/0magnet/xterm-go/vt"
	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/norm"
	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

// Engine is one magnet terminal.
type Engine struct {
	t  *vt.Terminal
	ev events
}

// New is the runner.Factory. magnet has no unicode-version option, so only the goldens' default "6" is accepted.
func New(c runner.Config) (runner.Engine, error) {
	if c.Unicode != "6" {
		return nil, fmt.Errorf("magnet has no unicode %q option", c.Unicode)
	}
	o := vt.NewOptions()
	o.Cols, o.Rows, o.Scrollback = c.Cols, c.Rows, c.Scrollback
	e := &Engine{t: vt.NewTerminal(o)}
	e.ev.listen(e.t)
	return e, nil
}

func (e *Engine) Write(p []byte) error {
	e.t.Write(p)
	return nil
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
	e.ev.fill(&st)
	return st, nil
}

// grid fills the buffer fields: the raw cursor (not clamped), the normal buffer always, the alternate one while active.
func (e *Engine) grid(st *runner.State) {
	bs := e.t.Buffers()
	b := bs.Active()
	st.P.Cols, st.P.Rows, st.P.Cursor = e.t.Cols(), e.t.Rows(), [2]int{b.X, b.Y}
	st.P.BaseY, st.P.ViewportY, st.P.Buffer = b.YBase, b.YDisp, "normal"
	st.Normal = rows(bs.Normal())
	st.P.NormalLength = len(st.Normal)
	if b == bs.Alt() {
		st.P.Buffer, st.Alt = "alternate", rows(bs.Alt())
	}
}

func rows(b *vt.Buffer) []runner.Row {
	out := make([]runner.Row, b.Lines.Length())
	cell := vt.NewCellData()
	for y := range out {
		l := b.Lines.Get(y)
		out[y].Wrapped = l.IsWrapped
		for x := 0; x < l.Length; x++ {
			l.LoadCell(x, cell)
			out[y].Cells = append(out[y].Cells, runner.Cell{Ch: cell.GetChars(), W: cell.GetWidth(), A: attr(cell)})
		}
	}
	return out
}

// attr copies the cell's IBufferCell-style getters into norm.Raw; magnet's flag getters already return bools.
func attr(c *vt.CellData) runner.Attr {
	return norm.Attr(norm.Raw{
		FgMode: c.GetFgColorMode(), BgMode: c.GetBgColorMode(), Fg: c.GetFgColor(), Bg: c.GetBgColor(),
		Bold: c.IsBold(), Italic: c.IsItalic(), Dim: c.IsDim(), Underline: c.IsUnderline(),
		Blink: c.IsBlink(), Inverse: c.IsInverse(), Invisible: c.IsInvisible(),
		Strike: c.IsStrikethrough(), Overline: c.IsOverline(),
	})
}
