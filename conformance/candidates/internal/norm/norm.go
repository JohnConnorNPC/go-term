// Package norm is the schema encoding (E3a): it writes one cell's attributes, as an xterm.js port's IBufferCell-style
// getters return them, in the golden's form, exactly as gen/lib/cells.mjs cellAttr does. It converts nothing else.
package norm

import (
	"fmt"

	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

// The xterm.js colour-mode bits (common/buffer/Constants.ts Attributes.CM_*). equivcheck proves each candidate's own
// constants equal these (equiv/encoding.json), so the adapters pass their GetFgColorMode() values through unchanged.
const (
	CMMask uint32 = 0x3000000
	CMP16  uint32 = 0x1000000
	CMP256 uint32 = 0x2000000
	CMRGB  uint32 = 0x3000000
)

// Raw is one cell as the getters return it: colour mode and colour (GetFgColorMode/GetFgColor, Bg likewise) and the
// nine flags (IsBold … IsOverline) in cells.mjs FLAGS order.
type Raw struct {
	FgMode, BgMode                                                            uint32
	Fg, Bg                                                                    int
	Bold, Italic, Dim, Underline, Blink, Inverse, Invisible, Strike, Overline bool
}

// On is cells.mjs's flag test `if (c[m]())`: a getter that returns its flag bit (xterm.js, gitpod) is set when non-zero.
func On(bit uint32) bool { return bit != 0 }

// color is cells.mjs color(mode, v): ["d"], or [p16|p256|rgb, v], or ["mode<n>", v] for any other mode value.
func color(mode uint32, v int) runner.Color {
	names := map[uint32]string{CMP16: "p16", CMP256: "p256", CMRGB: "rgb"}
	switch name, ok := names[mode]; {
	case mode == 0:
		return runner.Color{Mode: "d"}
	case ok:
		return runner.Color{Mode: name, Value: v}
	default:
		return runner.Color{Mode: fmt.Sprintf("mode%d", mode), Value: v}
	}
}

// Attr is cellAttr's P part.
func Attr(r Raw) runner.Attr {
	return runner.Attr{FG: color(r.FgMode, r.Fg), BG: color(r.BgMode, r.Bg), Bold: r.Bold, Italic: r.Italic, Dim: r.Dim,
		Underline: r.Underline, Blink: r.Blink, Inverse: r.Inverse, Invisible: r.Invisible, Strike: r.Strike, Overline: r.Overline}
}
