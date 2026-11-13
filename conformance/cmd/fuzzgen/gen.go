package main

import (
	"encoding/hex"
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
)

// Case is one conformance case file (conformance/README "Layout"). Steps holds only "x", "resize" and "reset" steps.
type Case struct {
	ID         string           `json:"id"`
	Cols       int              `json:"cols"`
	Rows       int              `json:"rows"`
	Scrollback int              `json:"scrollback"`
	Tags       []string         `json:"tags"`
	Steps      []map[string]any `json:"steps"`
}

// Bounds of a generated case.
const (
	MinCols, MaxCols = 2, 60
	MinRows, MaxRows = 1, 20
	Scrollback       = 50
	MaxSteps         = 6
	MaxTokens        = 8    // per write
	MaxString        = 4096 // any OSC/DCS/APC payload, far under both payload caps
)

func (c Case) writeBytes() int {
	n := 0
	for _, s := range c.Steps {
		if x, ok := s["x"].(string); ok {
			n += len(x) / 2
		}
	}
	return n
}

// Generate is case i of seed: deterministic, from its own PCG stream.
func Generate(seed uint64, i int) Case {
	r := rand.New(rand.NewPCG(seed, uint64(i)))
	g := gen{r: r}
	c := Case{ID: fmt.Sprintf("fuzz/%d-%d", seed, i), Cols: g.in(MinCols, MaxCols), Rows: g.in(MinRows, MaxRows),
		Scrollback: Scrollback, Tags: []string{"noop-ok"}}
	for range g.in(1, MaxSteps) {
		switch p := r.IntN(100); {
		case p < 4:
			c.Steps = append(c.Steps, map[string]any{"resize": []int{g.in(MinCols, MaxCols), g.in(MinRows, MaxRows)}})
		case p < 5:
			c.Steps = append(c.Steps, map[string]any{"reset": true})
		default:
			var b strings.Builder
			for range g.in(1, MaxTokens) {
				b.WriteString(g.token())
			}
			c.Steps = append(c.Steps, map[string]any{"x": hex.EncodeToString([]byte(b.String()))})
		}
	}
	return c
}

type gen struct{ r *rand.Rand }

func (g gen) in(lo, hi int) int       { return lo + g.r.IntN(hi-lo+1) }
func (g gen) pick(s ...string) string { return s[g.r.IntN(len(s))] }
func (g gen) chance(percent int) bool { return g.r.IntN(100) < percent }
func (g gen) num(hi int) string       { return strconv.Itoa(g.r.IntN(hi + 1)) }
func (g gen) ascii(lo, hi int) string {
	b := make([]byte, g.in(lo, hi))
	for i := range b {
		b[i] = byte(g.in(0x20, 0x7e))
	}
	return string(b)
}

// st is a string terminator: ST (ESC \) or BEL (BEL ends OSC and DCS in both engines, never APC since #10).
func (g gen) st(belOK bool) string {
	if belOK && g.chance(50) {
		return "\x07"
	}
	return "\x1b\\"
}

func (g gen) token() string {
	switch p := g.r.IntN(100); {
	case p < 22:
		return g.ascii(1, 12)
	case p < 32:
		return g.utf8()
	case p < 42:
		return g.c0()
	case p < 62:
		return g.csi()
	case p < 70:
		return g.sgr()
	case p < 78:
		return g.decset()
	case p < 84:
		return g.osc()
	case p < 87:
		return "\x1bP$q" + g.pick("m", "r", "\"p", "\"q", " q", "s", "t") + g.st(false)
	case p < 89:
		return "\x1b_" + g.ascii(0, 16) + g.st(false)
	default:
		return g.esc()
	}
}

// utf8: 2/3/4-byte, wide, combining, ZWJ, and invalid or truncated sequences.
func (g gen) utf8() string {
	return g.pick("é", "ü", "ß", "Ω", "中", "語", "한", "ｱ", "😀", "👍", "🇦🇺", "é", "ạ̈",
		"👩‍💻", "‍", "​", "️", "❤️", "𠀋",
		"\x80", "\xc3", "\xe4\xb8", "\xf0\x9f\x98", "\xff", "\xc0\xaf", "\xed\xa0\x80")
}

func (g gen) c0() string {
	if g.chance(30) {
		return g.pick("\x18", "\x1a") // CAN, SUB: abort a sequence
	}
	return g.pick("\x00", "\x05", "\x07", "\x08", "\x09", "\x0a", "\x0b", "\x0c", "\x0d", "\x0e", "\x0f", "\x7f")
}

func (g gen) params(n, hi int) string {
	ps := make([]string, n)
	for i := range ps {
		if !g.chance(10) { // an empty param now and then
			ps[i] = g.num(hi)
		}
	}
	return strings.Join(ps, ";")
}

func (g gen) csi() string {
	finals := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "P", "S", "T", "X", "Z", "@",
		"`", "a", "b", "d", "e", "f", "g", "r", "s", "u", "n", "c", "q", "t", "x", "j", "k"}
	priv := ""
	if g.chance(10) {
		priv = g.pick("?", ">", "=")
	}
	inter := ""
	if g.chance(8) {
		inter = g.pick(" ", "$", "\"", "!", "'")
	}
	hi := 30
	if g.chance(5) {
		hi = 100000
	}
	return "\x1b[" + priv + g.params(g.in(0, 3), hi) + inter + g.pick(finals...)
}

func (g gen) sgr() string {
	var p string
	switch g.r.IntN(7) {
	case 0:
		p = g.pick("38", "48", "58") + ";5;" + g.num(255)
	case 1:
		p = g.pick("38", "48", "58") + ";2;" + g.num(255) + ";" + g.num(255) + ";" + g.num(255)
	case 2:
		p = g.pick("38", "48", "58") + ":2::" + g.num(255) + ":" + g.num(255) + ":" + g.num(255)
	case 3:
		p = g.pick("38", "48") + ":5:" + g.num(255)
	case 4:
		p = "4:" + g.num(5)
	default:
		p = g.params(g.in(0, 4), 109)
	}
	return "\x1b[" + p + "m"
}

func (g gen) decset() string {
	m := g.pick("1", "6", "7", "25", "47", "1047", "1048", "1049", "2004", "9", "1000", "1001", "1002", "1003",
		"1006", "1016", "2026", "12", "45", "66", "1004")
	return "\x1b[?" + m + g.pick("h", "l")
}

func (g gen) osc() string {
	switch g.r.IntN(4) {
	case 0:
		return "\x1b]" + g.pick("0", "1", "2") + ";" + g.ascii(0, 20) + g.st(true)
	case 1:
		return "\x1b]8;" + g.pick("", "id=a") + ";http://x/" + g.ascii(0, 8) + g.st(true) + g.ascii(1, 4) + "\x1b]8;;" + g.st(true)
	case 2:
		return "\x1b]52;c;" + g.pick("", "?", "aGVsbG8=") + g.st(true)
	default:
		return "\x1b]" + g.num(200) + ";" + g.ascii(0, 10) + g.st(true)
	}
}

func (g gen) esc() string {
	return g.pick("\x1b7", "\x1b8", "\x1bc", "\x1bD", "\x1bM", "\x1bE", "\x1bH", "\x1b=", "\x1b>", "\x1b(0", "\x1b(B",
		"\x1b)0", "\x1b)B", "\x1b*0", "\x1b+B", "\x1bn", "\x1bo", "\x1b#8", "\x1b%G", "\x1b")
}
