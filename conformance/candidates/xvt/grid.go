// Package xvt adapts github.com/charmbracelet/x/vt to the conformance runner (#1 D7). It reports only what the
// emulator's public API returns (E3). The paths that API can't show are listed in State.Unobserved.
package xvt

import (
	"github.com/JohnConnorNPC/go-term/conformance/runner"
	vt "github.com/charmbracelet/x/vt"
)

// gridUnobserved: x/vt has no wrap flag (w), no viewport (viewportY), and no table rows for its colours (attrs, E10).
var gridUnobserved = []string{"attrs", "viewportY", "w"}

// Engine is one x/vt emulator.
type Engine struct {
	e *vt.Emulator
}

// New is the runner.Factory.
func New(c runner.Config) (runner.Engine, error) {
	e := &Engine{e: vt.NewEmulator(c.Cols, c.Rows)}
	e.e.SetScrollbackSize(c.Scrollback)
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
	return st, err
}

// grid fills the scalar buffer fields: cursor, size, scrollback length (baseY) and the active screen's name.
func (e *Engine) grid(st *runner.State) error {
	p := e.e.CursorPosition()
	st.P.Cols, st.P.Rows, st.P.Cursor, st.P.BaseY = e.e.Width(), e.e.Height(), [2]int{p.X, p.Y}, e.e.ScrollbackLen()
	st.Unobserved = append(st.Unobserved, gridUnobserved...)
	st.P.Buffer = "normal"
	if e.e.IsAltScreen() {
		st.P.Buffer = "alternate"
	}
	return nil
}
