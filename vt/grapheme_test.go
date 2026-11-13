package vt

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"testing"
)

// TestGraphemeBreakConformance runs the Unicode test file for UAX #29 grapheme
// cluster boundaries (testdata/GraphemeBreakTest.txt, the version the tables
// are generated from) against the break rules.
func TestGraphemeBreakConformance(t *testing.T) {
	f, err := os.Open("testdata/GraphemeBreakTest.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close() //nolint:errcheck
	sc := bufio.NewScanner(f)
	n := 0
	for sc.Scan() {
		line, _, _ := strings.Cut(sc.Text(), "#")
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		// fields alternate: ÷|× cp ÷|× cp ... ÷
		var want, got []bool // break before each code point after the first
		var st uint32
		for i := 1; i < len(fields); i += 2 {
			v, err := strconv.ParseUint(fields[i], 16, 32)
			if err != nil {
				t.Fatalf("%q: %v", line, err)
			}
			p := propsOf(uint32(v))
			brk := true
			if i > 1 {
				brk = graphemeBreak(st, p)
				want = append(want, fields[i-1] == "÷")
				got = append(got, brk)
			}
			st = nextGraphemeState(st, p, !brk)
		}
		for i := range want {
			if want[i] != got[i] {
				t.Errorf("%s: boundary %d: break=%v, want %v", strings.TrimSpace(line), i+1, got[i], want[i])
				break
			}
		}
		n++
	}
	if n < 700 {
		t.Fatalf("only %d test lines read", n)
	}
}

func TestGraphemeWidths(t *testing.T) {
	for _, tc := range []struct {
		cp   uint32
		want int
	}{
		{'a', 1}, {0x4E00, 2}, {0x1F600, 2}, {0x0301, 0}, {0x200D, 0},
		{0x1F3FB, 2}, {0x1F1E6, 2}, {0xFE0F, 0}, {0x2764, 1}, {0x1160, 0},
	} {
		if got := GraphemeWidth(tc.cp); got != tc.want {
			t.Errorf("GraphemeWidth(%U) = %d, want %d", tc.cp, got, tc.want)
		}
	}
}

// TestMode2027 drives Print through the terminal with the mode off and on.
func TestMode2027(t *testing.T) {
	const family = "\U0001F468\u200D\U0001F469\u200D\U0001F467" // man ZWJ woman ZWJ girl
	for _, tc := range []struct {
		name, text string
		on         bool
		cells      []string // the cells holding content, in order
		x          int      // cursor column afterwards
	}{
		{"ascii", "ab", true, []string{"a", "b"}, 2},
		{"combining", "e\u0301x", true, []string{"e\u0301", "x"}, 2},
		{"zwj family", family + "x", true, []string{family, "x"}, 3},
		{"zwj family, mode off", family, false, []string{"\U0001F468\u200D", "\U0001F469\u200D", "\U0001F467"}, 3}, // UnicodeV6 has emoji one cell wide
		{"vs16 widens", "\u2764\uFE0Fx", true, []string{"\u2764\uFE0F", "x"}, 3},
		{"vs16, mode off", "\u2764\uFE0Fx", false, []string{"\u2764\uFE0F", "x"}, 2},
		{"vs16 after a letter stays narrow", "a\uFE0Fx", true, []string{"a\uFE0F", "x"}, 2},
		{"flag", "\U0001F1E9\U0001F1EA\U0001F1EB\U0001F1F7", true, []string{"\U0001F1E9\U0001F1EA", "\U0001F1EB\U0001F1F7"}, 4},
		{"skin tone", "\U0001F44B\U0001F3FD", true, []string{"\U0001F44B\U0001F3FD"}, 2},
		{"keycap", "1\uFE0F\u20E3", true, []string{"1\uFE0F\u20E3"}, 2},
		{"devanagari conjunct", "\u0915\u094D\u0937", true, []string{"\u0915\u094D\u0937"}, 1},
		{"hangul jamo", "\u1100\u1161\u11A8", true, []string{"\u1100\u1161\u11A8"}, 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			term := newTestTerminal(20, 3)
			if tc.on {
				term.WriteString("\x1b[?2027h")
			}
			term.WriteString(tc.text)
			b := term.Buffer()
			row := b.Lines.Get(b.YBase)
			var cells []string
			for x := 0; x < 20; x++ {
				if row.HasContent(x) {
					cells = append(cells, row.GetString(x))
				}
			}
			if strings.Join(cells, "|") != strings.Join(tc.cells, "|") {
				t.Errorf("cells = %q, want %q", cells, tc.cells)
			}
			if b.X != tc.x {
				t.Errorf("cursor x = %d, want %d", b.X, tc.x)
			}
		})
	}
}

// TestMode2027WrapOnVS16: VS16 widening a cluster in the last column moves it
// to the next line, as the spec allows with autowrap on.
func TestMode2027WrapOnVS16(t *testing.T) {
	term := newTestTerminal(4, 3)
	term.WriteString("\x1b[?2027habc\u2764\uFE0F")
	if got := line(term, 0); got != "abc" {
		t.Errorf("row0 = %q", got)
	}
	if got := line(term, 1); got != "\u2764\uFE0F" {
		t.Errorf("row1 = %q", got)
	}
	if b := term.Buffer(); b.X != 2 || b.Y != 1 {
		t.Errorf("cursor = %d,%d", b.X, b.Y)
	}
}

func TestMode2027Report(t *testing.T) {
	term := newTestTerminal(10, 3)
	var out strings.Builder
	term.OnData = func(s string) { out.WriteString(s) }
	ask := func() string {
		out.Reset()
		term.WriteString("\x1b[?2027$p")
		return out.String()
	}
	if got := ask(); got != "\x1b[?2027;2$y" {
		t.Errorf("initially %q", got)
	}
	term.WriteString("\x1b[?2027h")
	if got := ask(); got != "\x1b[?2027;1$y" {
		t.Errorf("after set %q", got)
	}
	term.WriteString("\x1b[!p") // DECSTR resets it
	if got := ask(); got != "\x1b[?2027;2$y" {
		t.Errorf("after DECSTR %q", got)
	}
	term.WriteString("\x1b[?2027h\x1bc") // so does RIS
	if got := ask(); got != "\x1b[?2027;2$y" {
		t.Errorf("after RIS %q", got)
	}
}
