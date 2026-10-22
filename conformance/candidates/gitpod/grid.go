// Package gitpod adapts github.com/gitpod-io/xterm-go to the conformance runner (#1 D7). It reports only what the
// terminal's public API returns (E3).
package gitpod

import (
	"fmt"

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

// grid fills the buffer fields: the raw cursor (not clamped) and the active buffer's name.
func (e *Engine) grid(st *runner.State) {
	b := e.t.Buffer()
	st.P.Cols, st.P.Rows, st.P.Cursor = e.t.Cols(), e.t.Rows(), [2]int{b.X, b.Y}
	st.P.BaseY, st.P.ViewportY, st.P.Buffer = b.YBase, b.YDisp, "normal"
	if e.t.IsAltBufferActive() {
		st.P.Buffer = "alternate"
	}
}
