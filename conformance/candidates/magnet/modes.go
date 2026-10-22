package magnet

import (
	"strings"

	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/equiv"
	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

// modes fills P.modes under the golden's key names (bools by identity) and P.opts. The mouse mode and both cursor
// options change spelling only through an equiv row; with no row (strict) they are magnet's raw values, a nil
// override being JSON null (rulings-6), so a strict mismatch shows both values.
func (e *Engine) modes(st *runner.State) {
	cs := e.t.CoreService()
	d, m, proto := cs.DecPrivateModes, cs.Modes, e.t.MouseService().ActiveProtocol()
	var style, blink any
	if d.CursorStyle != nil {
		style = *d.CursorStyle
	}
	if d.CursorBlink != nil {
		blink = *d.CursorBlink
	}
	g, ok := equiv.Map("magnet", "modes.mouseTrackingMode", proto)
	mouse := e.resolve(st, "modes.mouseTrackingMode", g, ok, proto)
	g, ok = equiv.Map("magnet", "opts.cursorStyle", override(d.CursorStyle != nil))
	style = e.resolve(st, "opts.cursorStyle", g, ok, style)
	g, ok = equiv.Map("magnet", "opts.cursorBlink", override(d.CursorBlink != nil))
	blink = e.resolve(st, "opts.cursorBlink", g, ok, blink)
	st.P.Modes = map[string]any{"applicationCursorKeysMode": d.ApplicationCursorKeys, "applicationKeypadMode": d.ApplicationKeypad,
		"bracketedPasteMode": d.BracketedPasteMode, "insertMode": m.InsertMode, "mouseTrackingMode": mouse, "originMode": d.Origin,
		"reverseWraparoundMode": d.ReverseWraparound, "sendFocusMode": d.SendFocus,
		"synchronizedOutputMode": d.SynchronizedOutput, "wraparoundMode": d.Wraparound}
	st.P.Opts = map[string]any{"cursorStyle": style, "cursorBlink": blink}
}

// override is the table key of a DECSCUSR override pointer.
func override(set bool) string {
	if set {
		return "override=set"
	}
	return "override=nil"
}

// resolve is the value reported at path: the row's golden value, the override itself for "@override", the
// adapter's own *Options (the one magnet reads, M6/M8) for "@Options.*", raw when no row applies. An accessor outside
// those leaves path unobserved.
func (e *Engine) resolve(st *runner.State, path string, golden any, ok bool, raw any) any {
	s, _ := golden.(string)
	switch {
	case !ok || s == "@override":
		return raw
	case s == "@Options.CursorStyle":
		return e.t.Options.CursorStyle
	case s == "@Options.CursorBlink":
		return e.t.Options.CursorBlink
	case strings.HasPrefix(s, "@"):
		st.Unobserved = append(st.Unobserved, path)
	}
	return golden
}
