// Package magnet adapts github.com/0magnet/xterm-go/vt to the conformance runner (#1 D7). It reports only what the
// terminal's public API returns (E3).
package magnet

import (
	vt "github.com/0magnet/xterm-go/vt"
	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

// Engine is one magnet terminal.
type Engine struct {
	t *vt.Terminal
}

// New is the runner.Factory.
func New(c runner.Config) (runner.Engine, error) {
	o := vt.NewOptions()
	o.Cols, o.Rows, o.Scrollback = c.Cols, c.Rows, c.Scrollback
	return &Engine{t: vt.NewTerminal(o)}, nil
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
	return st, nil
}

// grid fills the buffer fields: the raw cursor (not clamped) and the active buffer's name.
func (e *Engine) grid(st *runner.State) {
	bs := e.t.Buffers()
	b := bs.Active()
	st.P.Cols, st.P.Rows, st.P.Cursor = e.t.Cols(), e.t.Rows(), [2]int{b.X, b.Y}
	st.P.BaseY, st.P.ViewportY, st.P.Buffer = b.YBase, b.YDisp, "normal"
	if b == bs.Alt() {
		st.P.Buffer = "alternate"
	}
}
