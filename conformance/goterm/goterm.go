// Package goterm adapts go-term's own engine (package vt) to the conformance runner. It reports only what vt's public
// API returns, read natively: no equivalence table. The adapter was ported from the #2 candidate adapter
// (conformance/candidates/magnet), whose table rows stay there as the #2 record.
package goterm

import (
	"fmt"
	"strings"

	"github.com/JohnConnorNPC/go-term/conformance/runner"
	"github.com/JohnConnorNPC/go-term/vt"
)

// Engine is one vt terminal.
type Engine struct {
	t  *vt.Terminal
	ev events
}

// New is the runner.Factory. vt has no unicode-version option, so only the goldens' default "6" is accepted.
func New(c runner.Config) (runner.Engine, error) {
	if c.Unicode != "6" {
		return nil, fmt.Errorf("vt has no unicode %q option", c.Unicode)
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
	e.modes(&st)
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

// attr is gen/lib/cells.mjs cellAttr: the colours by mode, the nine flags as vt's getters return them.
func attr(c *vt.CellData) runner.Attr {
	return runner.Attr{FG: color(c.GetFgColorMode(), c.GetFgColor()), BG: color(c.GetBgColorMode(), c.GetBgColor()),
		Bold: c.IsBold(), Italic: c.IsItalic(), Dim: c.IsDim(), Underline: c.IsUnderline(), Blink: c.IsBlink(),
		Inverse: c.IsInverse(), Invisible: c.IsInvisible(), Strike: c.IsStrikethrough(), Overline: c.IsOverline()}
}

// color is cells.mjs color(mode, v): ["d"], or [p16|p256|rgb, v], or ["mode<n>", v] for any other mode value.
func color(mode uint32, v int) runner.Color {
	switch mode {
	case vt.AttrCMDefault:
		return runner.Color{Mode: "d"}
	case vt.AttrCMP16:
		return runner.Color{Mode: "p16", Value: v}
	case vt.AttrCMP256:
		return runner.Color{Mode: "p256", Value: v}
	case vt.AttrCMRGB:
		return runner.Color{Mode: "rgb", Value: v}
	}
	return runner.Color{Mode: fmt.Sprintf("mode%d", mode), Value: v}
}

// modes fills P.modes under the golden's key names and P.opts. The mouse mode is vt's protocol name lowercased
// (NONE, X10, VT200, DRAG, ANY are xterm.js's names in upper case). The cursor options are vt's Options, as
// @xterm/headless 6.0.0 reports terminal.options: DECSCUSR sets only the DEC-mode override and leaves the options
// alone (6.0.0 answers DECRQSS " q" from the options too).
func (e *Engine) modes(st *runner.State) {
	cs := e.t.CoreService()
	d, m := cs.DecPrivateModes, cs.Modes
	st.P.Modes = map[string]any{"applicationCursorKeysMode": d.ApplicationCursorKeys, "applicationKeypadMode": d.ApplicationKeypad,
		"bracketedPasteMode": d.BracketedPasteMode, "insertMode": m.InsertMode,
		"mouseTrackingMode": strings.ToLower(e.t.MouseService().ActiveProtocol()), "originMode": d.Origin,
		"reverseWraparoundMode": d.ReverseWraparound, "sendFocusMode": d.SendFocus,
		"synchronizedOutputMode": d.SynchronizedOutput, "wraparoundMode": d.Wraparound}
	st.P.Opts = map[string]any{"cursorStyle": e.t.Options.CursorStyle, "cursorBlink": e.t.Options.CursorBlink}
}

// events are the terminal's callbacks, kept from construction on as gen/lib/session.mjs keeps them (a reset op
// clears nothing): one reply per OnData call, in order. OnReply stays unset, so vt sends its query replies through
// OnData too, as xterm.js does through onData.
type events struct {
	replies, titles  []string
	bells, lineFeeds int
}

func (ev *events) listen(t *vt.Terminal) {
	t.OnData = func(s string) { ev.replies = append(ev.replies, s) }
	t.OnTitleChange = func(s string) { ev.titles = append(ev.titles, s) }
	t.OnBell = func() { ev.bells++ }
	t.OnLineFeed = func() { ev.lineFeeds++ }
}

// fill copies the counters into st; the title is the last one set, "" before any (session.mjs titles.at(-1) ?? "").
func (ev *events) fill(st *runner.State) {
	st.P.Replies, st.P.Titles = append([]string{}, ev.replies...), append([]string{}, ev.titles...)
	st.P.Bells, st.P.LineFeeds = ev.bells, ev.lineFeeds
	if len(ev.titles) > 0 {
		st.P.Title = ev.titles[len(ev.titles)-1]
	}
}
