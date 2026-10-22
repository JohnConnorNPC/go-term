package magnet

import (
	"strings"

	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/equiv"
	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

// modes fills P.modes under the golden's key names (bools by identity). The mouse mode changes spelling only through
// an equiv row; with no row (strict) it is magnet's raw value (rulings-6), so a strict mismatch shows both values.
func (e *Engine) modes(st *runner.State) {
	cs := e.t.CoreService()
	d, m, proto := cs.DecPrivateModes, cs.Modes, e.t.MouseService().ActiveProtocol()
	g, ok := equiv.Map("magnet", "modes.mouseTrackingMode", proto)
	mouse := e.resolve(st, "modes.mouseTrackingMode", g, ok, proto)
	st.P.Modes = map[string]any{"applicationCursorKeysMode": d.ApplicationCursorKeys, "applicationKeypadMode": d.ApplicationKeypad,
		"bracketedPasteMode": d.BracketedPasteMode, "insertMode": m.InsertMode, "mouseTrackingMode": mouse, "originMode": d.Origin,
		"reverseWraparoundMode": d.ReverseWraparound, "sendFocusMode": d.SendFocus,
		"synchronizedOutputMode": d.SynchronizedOutput, "wraparoundMode": d.Wraparound}
}

// resolve is the value reported at path: the row's golden value, raw for "@override" or when no row applies. An
// accessor magnet can't read here (any other @…) leaves path unobserved.
func (e *Engine) resolve(st *runner.State, path string, golden any, ok bool, raw any) any {
	s, _ := golden.(string)
	switch {
	case !ok || s == "@override":
		return raw
	case strings.HasPrefix(s, "@"):
		st.Unobserved = append(st.Unobserved, path)
	}
	return golden
}
