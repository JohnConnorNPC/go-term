package vt

import (
	"fmt"
	"strings"
	"testing"
)

func newTestTerminal(cols, rows int) *Terminal {
	opts := NewOptions()
	opts.Cols = cols
	opts.Rows = rows
	return NewTerminal(opts)
}

// line returns the trimmed text of viewport row y.
func line(term *Terminal, y int) string {
	b := term.Buffer()
	return b.Lines.Get(b.YBase+y).TranslateToString(true, 0, -1)
}

func TestTermPrintAndWrap(t *testing.T) {
	term := newTestTerminal(10, 5)
	term.WriteString("abcdefghijklmno")
	if got := line(term, 0); got != "abcdefghij" {
		t.Fatalf("row0 = %q", got)
	}
	if got := line(term, 1); got != "klmno" {
		t.Fatalf("row1 = %q", got)
	}
	if !term.Buffer().Lines.Get(1).IsWrapped {
		t.Fatal("row1 should be wrapped")
	}
	if term.Buffer().X != 5 || term.Buffer().Y != 1 {
		t.Fatalf("cursor = %d,%d", term.Buffer().X, term.Buffer().Y)
	}
}

func TestTermCrLf(t *testing.T) {
	term := newTestTerminal(20, 5)
	term.WriteString("hello\r\nworld")
	if line(term, 0) != "hello" || line(term, 1) != "world" {
		t.Fatalf("rows = %q, %q", line(term, 0), line(term, 1))
	}
	// bare LF keeps the column
	term = newTestTerminal(20, 5)
	term.WriteString("abc\ndef")
	if got := line(term, 1); got != "   def" {
		t.Fatalf("bare LF row1 = %q", got)
	}
	// convertEol treats LF as CRLF
	opts := NewOptions()
	opts.Cols = 20
	opts.Rows = 5
	opts.ConvertEol = true
	term = NewTerminal(opts)
	term.WriteString("abc\ndef")
	if got := line(term, 1); got != "def" {
		t.Fatalf("convertEol row1 = %q", got)
	}
}

func TestTermCursorMovement(t *testing.T) {
	term := newTestTerminal(80, 24)
	term.WriteString("\x1b[5;10H")
	if term.Buffer().Y != 4 || term.Buffer().X != 9 {
		t.Fatalf("CUP: %d,%d", term.Buffer().X, term.Buffer().Y)
	}
	term.WriteString("\x1b[2A") // up 2
	if term.Buffer().Y != 2 {
		t.Fatalf("CUU: y=%d", term.Buffer().Y)
	}
	term.WriteString("\x1b[3B") // down 3
	if term.Buffer().Y != 5 {
		t.Fatalf("CUD: y=%d", term.Buffer().Y)
	}
	term.WriteString("\x1b[7C") // forward 7
	if term.Buffer().X != 16 {
		t.Fatalf("CUF: x=%d", term.Buffer().X)
	}
	term.WriteString("\x1b[100D") // backward, clamps at 0
	if term.Buffer().X != 0 {
		t.Fatalf("CUB clamp: x=%d", term.Buffer().X)
	}
	term.WriteString("\x1b[999;999H") // clamps to viewport
	if term.Buffer().Y != 23 || term.Buffer().X != 79 {
		t.Fatalf("CUP clamp: %d,%d", term.Buffer().X, term.Buffer().Y)
	}
}

func TestTermEraseInLine(t *testing.T) {
	term := newTestTerminal(20, 5)
	term.WriteString("hello world welt")
	term.WriteString("\x1b[1;7H") // on 'w' of world
	term.WriteString("\x1b[K")    // erase to right
	// the typed space at col 6 is real content and survives trimming
	if got := line(term, 0); got != "hello " {
		t.Fatalf("EL0 = %q", got)
	}
	term = newTestTerminal(20, 5)
	term.WriteString("hello world welt")
	term.WriteString("\x1b[1;7H")
	term.WriteString("\x1b[1K") // erase to left (inclusive)
	if got := line(term, 0); got != "       orld welt" {
		t.Fatalf("EL1 = %q", got)
	}
	term.WriteString("\x1b[2K")
	if got := line(term, 0); got != "" {
		t.Fatalf("EL2 = %q", got)
	}
}

func TestTermEraseInDisplay(t *testing.T) {
	term := newTestTerminal(10, 4)
	term.WriteString("aaa\r\nbbb\r\nccc\r\nddd")
	term.WriteString("\x1b[2;2H")
	term.WriteString("\x1b[J") // erase below
	if line(term, 0) != "aaa" || line(term, 1) != "b" || line(term, 2) != "" || line(term, 3) != "" {
		t.Fatalf("ED0: %q %q %q %q", line(term, 0), line(term, 1), line(term, 2), line(term, 3))
	}

	term = newTestTerminal(10, 4)
	term.WriteString("aaa\r\nbbb\r\nccc\r\nddd")
	term.WriteString("\x1b[2;2H")
	term.WriteString("\x1b[1J") // erase above
	if line(term, 0) != "" || line(term, 1) != "  b" || line(term, 2) != "ccc" {
		t.Fatalf("ED1: %q %q %q", line(term, 0), line(term, 1), line(term, 2))
	}

	term.WriteString("\x1b[2J")
	for y := 0; y < 4; y++ {
		if line(term, y) != "" {
			t.Fatalf("ED2 row%d = %q", y, line(term, y))
		}
	}
}

func TestTermScrollbackAndED3(t *testing.T) {
	term := newTestTerminal(10, 3)
	for i := 0; i < 10; i++ {
		term.WriteString(fmt.Sprintf("line%d\r\n", i))
	}
	b := term.Buffer()
	if b.YBase != 8 {
		t.Fatalf("ybase = %d", b.YBase)
	}
	if b.Lines.Length() != 11 {
		t.Fatalf("lines = %d", b.Lines.Length())
	}
	// erase scrollback
	term.WriteString("\x1b[3J")
	if b.YBase != 0 || b.Lines.Length() != 3 {
		t.Fatalf("after ED3: ybase=%d lines=%d", b.YBase, b.Lines.Length())
	}
	if line(term, 0) != "line8" {
		t.Fatalf("viewport top = %q", line(term, 0))
	}
}

func TestTermSGRColors(t *testing.T) {
	term := newTestTerminal(20, 5)
	term.WriteString("\x1b[31mrot\x1b[0m")
	l := term.Buffer().Lines.Get(0)
	fg := l.GetFg(0)
	if fg&AttrCMMask != AttrCMP16 || fg&AttrPColorMask != 1 {
		t.Fatalf("fg = %08x", fg)
	}
	// after SGR0 attrs are reset
	attr := term.InputHandler().GetAttrData()
	if attr.Fg != 0 || attr.Bg != 0 {
		t.Fatalf("attrs not reset: %08x %08x", attr.Fg, attr.Bg)
	}

	// 256 colors, semicolon form ("rot" ended at x=3)
	term.WriteString("\x1b[38;5;123mx")
	fg = l.GetFg(3)
	if fg&AttrCMMask != AttrCMP256 || fg&AttrPColorMask != 123 {
		t.Fatalf("p256 fg = %08x", fg)
	}

	// RGB, subparam form
	term.WriteString("\x1b[0m\x1b[38:2::10:20:30my")
	fg = l.GetFg(4)
	if fg&AttrCMMask != AttrCMRGB || fg&AttrRGBMask != 10<<16|20<<8|30 {
		t.Fatalf("rgb fg = %08x", fg)
	}

	// RGB, semicolon form
	term.WriteString("\x1b[0m\x1b[48;2;40;50;60mz")
	bg := l.GetBg(5)
	if bg&AttrCMMask != AttrCMRGB || bg&AttrRGBMask != 40<<16|50<<8|60 {
		t.Fatalf("rgb bg = %08x", bg)
	}

	// bold + underline flags
	term.WriteString("\x1b[0m\x1b[1;4mw")
	fg = l.GetFg(6)
	if fg&FgBold == 0 || fg&FgUnderline == 0 {
		t.Fatalf("bold/underline fg = %08x", fg)
	}
}

func TestTermSGRUnderlineStyle(t *testing.T) {
	term := newTestTerminal(20, 5)
	term.WriteString("\x1b[4:3m") // curly
	attr := term.InputHandler().GetAttrData()
	if attr.Extended.UnderlineStyle() != UnderlineCurly {
		t.Fatalf("style = %d", attr.Extended.UnderlineStyle())
	}
	term.WriteString("\x1b[24m")
	if attr := term.InputHandler().GetAttrData(); attr.IsUnderline() {
		t.Fatal("underline should be off")
	}
	// 4: with empty subparam defaults to single
	term.WriteString("\x1b[4:m")
	if got := term.InputHandler().GetAttrData().Extended.UnderlineStyle(); got != UnderlineSingle {
		t.Fatalf("default style = %d", got)
	}
}

func TestTermInsertDeleteLines(t *testing.T) {
	term := newTestTerminal(10, 4)
	term.WriteString("aaa\r\nbbb\r\nccc\r\nddd")
	term.WriteString("\x1b[2;1H\x1b[1L") // insert line at row 2
	if line(term, 0) != "aaa" || line(term, 1) != "" || line(term, 2) != "bbb" || line(term, 3) != "ccc" {
		t.Fatalf("IL: %q %q %q %q", line(term, 0), line(term, 1), line(term, 2), line(term, 3))
	}
	term.WriteString("\x1b[2;1H\x1b[1M") // delete it again
	if line(term, 0) != "aaa" || line(term, 1) != "bbb" || line(term, 2) != "ccc" || line(term, 3) != "" {
		t.Fatalf("DL: %q %q %q %q", line(term, 0), line(term, 1), line(term, 2), line(term, 3))
	}
}

func TestTermInsertDeleteChars(t *testing.T) {
	term := newTestTerminal(10, 2)
	term.WriteString("abcdef")
	term.WriteString("\x1b[1;3H\x1b[2@") // insert 2 blanks at 'c'
	if got := line(term, 0); got != "ab  cdef" {
		t.Fatalf("ICH = %q", got)
	}
	term.WriteString("\x1b[2P") // delete them
	if got := line(term, 0); got != "abcdef" {
		t.Fatalf("DCH = %q", got)
	}
	term.WriteString("\x1b[3X") // erase 3 chars
	if got := line(term, 0); got != "ab   f" {
		t.Fatalf("ECH = %q", got)
	}
}

func TestTermScrollRegion(t *testing.T) {
	term := newTestTerminal(10, 5)
	term.WriteString("aaa\r\nbbb\r\nccc\r\nddd\r\neee")
	term.WriteString("\x1b[2;4r") // margins rows 2-4
	// cursor is homed by DECSTBM
	if term.Buffer().Y != 0 || term.Buffer().X != 0 {
		t.Fatalf("cursor after DECSTBM: %d,%d", term.Buffer().X, term.Buffer().Y)
	}
	// move to bottom margin and LF -> scrolls only the region
	term.WriteString("\x1b[4;1H\n")
	if line(term, 0) != "aaa" || line(term, 1) != "ccc" || line(term, 2) != "ddd" || line(term, 3) != "" || line(term, 4) != "eee" {
		t.Fatalf("region scroll: %q %q %q %q %q", line(term, 0), line(term, 1), line(term, 2), line(term, 3), line(term, 4))
	}
}

func TestTermAltBuffer(t *testing.T) {
	term := newTestTerminal(10, 3)
	term.WriteString("normal")
	term.WriteString("\x1b[?1049h") // alt buffer
	if term.Buffers().Active() != term.Buffers().Alt() {
		t.Fatal("alt buffer not active")
	}
	term.WriteString("\x1b[Halt!")
	if got := line(term, 0); got != "alt!" {
		t.Fatalf("alt row0 = %q", got)
	}
	term.WriteString("\x1b[?1049l") // back
	if term.Buffers().Active() != term.Buffers().Normal() {
		t.Fatal("normal buffer not active")
	}
	if got := line(term, 0); got != "normal" {
		t.Fatalf("normal row0 = %q", got)
	}
}

func TestTermDeviceReports(t *testing.T) {
	term := newTestTerminal(80, 24)
	var out strings.Builder
	term.OnData = func(data string) { out.WriteString(data) }

	term.WriteString("\x1b[5;10H\x1b[6n") // DSR CPR
	if out.String() != "\x1b[5;10R" {
		t.Fatalf("CPR = %q", out.String())
	}
	out.Reset()
	term.WriteString("\x1b[c") // DA1
	if out.String() != "\x1b[?1;2c" {
		t.Fatalf("DA1 = %q", out.String())
	}
	out.Reset()
	term.WriteString("\x1b[>c") // DA2
	if out.String() != "\x1b[>0;276;0c" {
		t.Fatalf("DA2 = %q", out.String())
	}
	out.Reset()
	term.WriteString("\x1b[?1$p") // DECRQM app cursor keys -> reset
	if out.String() != "\x1b[?1;2$y" {
		t.Fatalf("DECRQM = %q", out.String())
	}
	out.Reset()
	term.WriteString("\x1bP$qr\x1b\\") // DECRQSS margins
	if out.String() != "\x1bP1$r1;24r\x1b\\" {
		t.Fatalf("DECRQSS = %q", out.String())
	}
}

func TestTermCharsetLineDrawing(t *testing.T) {
	term := newTestTerminal(10, 2)
	term.WriteString("\x1b(0lqk\x1b(B")
	if got := line(term, 0); got != "┌─┐" {
		t.Fatalf("line drawing = %q", got)
	}
	// SO/SI with G1
	term = newTestTerminal(10, 2)
	term.WriteString("\x1b)0a\x0eq\x0fq")
	if got := line(term, 0); got != "a─q" {
		t.Fatalf("SO/SI = %q", got)
	}
}

func TestTermTitle(t *testing.T) {
	term := newTestTerminal(10, 2)
	var title string
	term.OnTitleChange = func(t string) { title = t }
	term.WriteString("\x1b]0;my title\x07")
	if title != "my title" {
		t.Fatalf("title = %q", title)
	}
	if term.InputHandler().WindowTitle() != "my title" {
		t.Fatal("window title not stored")
	}
}

func TestTermTabStops(t *testing.T) {
	term := newTestTerminal(40, 2)
	term.WriteString("\ta")
	if term.Buffer().X != 9 {
		t.Fatalf("tab: x=%d", term.Buffer().X)
	}
	// set custom stop at 3, clear all defaults first
	term.WriteString("\r\x1b[3g")      // clear all tab stops
	term.WriteString("\x1b[1;4H\x1bH") // HTS at col 3
	term.WriteString("\r\t")
	if term.Buffer().X != 3 {
		t.Fatalf("custom tab: x=%d", term.Buffer().X)
	}
	// next tab goes to end of line
	term.WriteString("\t")
	if term.Buffer().X != 39 {
		t.Fatalf("tab to end: x=%d", term.Buffer().X)
	}
}

func TestTermRepeatPrecedingCharacter(t *testing.T) {
	term := newTestTerminal(20, 2)
	term.WriteString("a\x1b[4b")
	if got := line(term, 0); got != "aaaaa" {
		t.Fatalf("REP = %q", got)
	}
	// REP after non-printable is a NOOP
	term = newTestTerminal(20, 2)
	term.WriteString("a\r\n\x1b[4b")
	if got := line(term, 1); got != "" {
		t.Fatalf("REP after LF = %q", got)
	}
}

func TestTermScreenAlignmentPattern(t *testing.T) {
	term := newTestTerminal(4, 3)
	term.WriteString("\x1b#8")
	for y := 0; y < 3; y++ {
		if got := line(term, y); got != "EEEE" {
			t.Fatalf("DECALN row%d = %q", y, got)
		}
	}
}

func TestTermWideChars(t *testing.T) {
	term := newTestTerminal(10, 3)
	term.WriteString("汉字ab")
	l := term.Buffer().Lines.Get(0)
	if got := l.TranslateToString(true, 0, -1); got != "汉字ab" {
		t.Fatalf("wide = %q", got)
	}
	if l.GetWidth(0) != 2 || l.GetWidth(1) != 0 || l.GetWidth(4) != 1 {
		t.Fatalf("widths: %d %d %d", l.GetWidth(0), l.GetWidth(1), l.GetWidth(4))
	}
	// wide char at odd last column wraps early leaving an empty cell
	term = newTestTerminal(5, 3)
	term.WriteString("abcd汉")
	if got := line(term, 0); got != "abcd" {
		t.Fatalf("row0 = %q", got)
	}
	if got := line(term, 1); got != "汉" {
		t.Fatalf("row1 = %q", got)
	}
	if !term.Buffer().Lines.Get(1).IsWrapped {
		t.Fatal("row1 should be wrapped")
	}
}

func TestTermCombiningChars(t *testing.T) {
	term := newTestTerminal(10, 2)
	term.WriteString("éx") // e + combining acute
	l := term.Buffer().Lines.Get(0)
	if got := l.TranslateToString(true, 0, -1); got != "éx" {
		t.Fatalf("combining = %q", got)
	}
	if term.Buffer().X != 2 {
		t.Fatalf("cursor x = %d", term.Buffer().X)
	}
}

func TestTermReverseIndex(t *testing.T) {
	term := newTestTerminal(10, 3)
	term.WriteString("aaa\r\nbbb\r\nccc")
	term.WriteString("\x1b[1;1H\x1bM") // RI at top -> scroll down
	if line(term, 0) != "" || line(term, 1) != "aaa" || line(term, 2) != "bbb" {
		t.Fatalf("RI: %q %q %q", line(term, 0), line(term, 1), line(term, 2))
	}
}

func TestTermSaveRestoreCursor(t *testing.T) {
	term := newTestTerminal(20, 5)
	term.WriteString("\x1b[3;5H\x1b[31m\x1b7") // position + red + save
	term.WriteString("\x1b[H\x1b[0m")          // home + reset
	term.WriteString("\x1b8")                  // restore
	if term.Buffer().Y != 2 || term.Buffer().X != 4 {
		t.Fatalf("restored cursor: %d,%d", term.Buffer().X, term.Buffer().Y)
	}
	attr := term.InputHandler().GetAttrData()
	if attr.Fg&AttrCMMask != AttrCMP16 || attr.Fg&AttrPColorMask != 1 {
		t.Fatalf("restored fg = %08x", attr.Fg)
	}
}

func TestTermInsertMode(t *testing.T) {
	term := newTestTerminal(10, 2)
	term.WriteString("abcdef")
	term.WriteString("\x1b[1;3H\x1b[4hXY\x1b[4l") // IRM on, insert XY, off
	if got := line(term, 0); got != "abXYcdef" {
		t.Fatalf("IRM = %q", got)
	}
}

func TestTermHyperlink(t *testing.T) {
	term := newTestTerminal(30, 2)
	term.WriteString("\x1b]8;;http://example.com\x1b\\link\x1b]8;;\x1b\\")
	l := term.Buffer().Lines.Get(0)
	cell := l.LoadCell(0, &CellData{AttributeData: AttributeData{Extended: NewExtendedAttrs()}})
	if cell.Extended == nil || cell.Extended.URLID == 0 {
		t.Fatal("no url id on link cell")
	}
	link, ok := term.OscLinkService().GetLinkData(cell.Extended.URLID)
	if !ok || link.URI != "http://example.com" {
		t.Fatalf("link = %+v ok=%v", link, ok)
	}
	// after the closing OSC 8 the attr has no url
	if term.InputHandler().GetAttrData().Extended.URLID != 0 {
		t.Fatal("url id should be cleared")
	}
}

func TestTermMouseReports(t *testing.T) {
	term := newTestTerminal(80, 24)
	var out strings.Builder
	term.OnData = func(data string) { out.WriteString(data) }
	var bin strings.Builder
	term.OnBinary = func(data string) { bin.WriteString(data) }

	// no protocol active -> no report
	if term.MouseService().TriggerMouseEvent(&MouseEvent{Col: 5, Row: 5, Button: MouseButtonLeft, Action: MouseActionDown}) {
		t.Fatal("report sent without protocol")
	}

	// VT200 + DEFAULT encoding -> binary report
	term.WriteString("\x1b[?1000h")
	if !term.MouseService().TriggerMouseEvent(&MouseEvent{Col: 5, Row: 5, Button: MouseButtonLeft, Action: MouseActionDown}) {
		t.Fatal("report not sent")
	}
	if bin.String() != "\x1b[M &&" { // 32+0, 32+6, 32+6
		t.Fatalf("default report = %q", bin.String())
	}

	// SGR encoding
	term.WriteString("\x1b[?1006h")
	out.Reset()
	term.MouseService().TriggerMouseEvent(&MouseEvent{Col: 5, Row: 5, Button: MouseButtonLeft, Action: MouseActionDown})
	if out.String() != "\x1b[<0;6;6M" {
		t.Fatalf("sgr report = %q", out.String())
	}
	out.Reset()
	term.MouseService().TriggerMouseEvent(&MouseEvent{Col: 5, Row: 5, Button: MouseButtonLeft, Action: MouseActionUp})
	if out.String() != "\x1b[<0;6;6m" {
		t.Fatalf("sgr release = %q", out.String())
	}
}

func TestTermDecPrivateModes(t *testing.T) {
	term := newTestTerminal(20, 5)
	term.WriteString("\x1b[?1h")
	if !term.CoreService().DecPrivateModes.ApplicationCursorKeys {
		t.Fatal("DECCKM not set")
	}
	term.WriteString("\x1b[?1l")
	if term.CoreService().DecPrivateModes.ApplicationCursorKeys {
		t.Fatal("DECCKM not reset")
	}
	term.WriteString("\x1b[?2004h")
	if !term.CoreService().DecPrivateModes.BracketedPasteMode {
		t.Fatal("bracketed paste not set")
	}
	term.WriteString("\x1b[?25l")
	if !term.CoreService().IsCursorHidden {
		t.Fatal("cursor not hidden")
	}
	term.WriteString("\x1b[?7l")
	if term.CoreService().DecPrivateModes.Wraparound {
		t.Fatal("DECAWM not reset")
	}
	// without wraparound, output past the margin stays on the line
	term.WriteString("\x1b[1;18Hxxxxxxx")
	if term.Buffer().Y != 0 {
		t.Fatalf("no-wrap y = %d", term.Buffer().Y)
	}
}

func TestTermOscColor(t *testing.T) {
	term := newTestTerminal(20, 5)
	var events []ColorEvent
	term.OnColor = func(e []ColorEvent) { events = append(events, e...) }

	term.WriteString("\x1b]4;17;#aabbcc\x07")
	if len(events) != 1 || events[0].Type != ColorRequestSet || events[0].Index != 17 ||
		events[0].Color != [3]int{0xaa, 0xbb, 0xcc} {
		t.Fatalf("osc4 = %+v", events)
	}
	events = nil
	term.WriteString("\x1b]10;?\x07")
	if len(events) != 1 || events[0].Type != ColorRequestReport || events[0].Index != SpecialColorForeground {
		t.Fatalf("osc10 = %+v", events)
	}
	events = nil
	term.WriteString("\x1b]104\x07")
	if len(events) != 1 || events[0].Type != ColorRequestRestore || events[0].Index != -1 {
		t.Fatalf("osc104 = %+v", events)
	}
}

func TestParseColorFormats(t *testing.T) {
	cases := []struct {
		in   string
		want [3]int
		ok   bool
	}{
		{"rgb:12/34/56", [3]int{0x12, 0x34, 0x56}, true},
		{"rgb:1/2/3", [3]int{0x11, 0x22, 0x33}, true},
		{"rgb:ffff/0000/8080", [3]int{255, 0, 128}, true},
		{"#123", [3]int{0x10, 0x20, 0x30}, true},
		{"#aabbcc", [3]int{0xaa, 0xbb, 0xcc}, true},
		{"#AABBCC", [3]int{0xaa, 0xbb, 0xcc}, true},
		{"#aaabbbccc", [3]int{0xaa, 0xbb, 0xcc}, true},
		{"nonsense", [3]int{}, false},
		{"", [3]int{}, false},
	}
	for _, c := range cases {
		got, ok := ParseColor(c.in)
		if ok != c.ok || got != c.want {
			t.Fatalf("ParseColor(%q) = %v,%v want %v,%v", c.in, got, ok, c.want, c.ok)
		}
	}
	if s := ToRgbString([3]int{255, 0, 128}, 16); s != "rgb:ffff/0000/8080" {
		t.Fatalf("ToRgbString = %q", s)
	}
}

func TestWcwidth(t *testing.T) {
	cases := []struct {
		cp   uint32
		want int
	}{
		{'a', 1},
		{0x00, 0},
		{0x0301, 0},  // combining acute
		{0x4E16, 2},  // 世
		{0x1F600, 1}, // emoji (V6 tables: width 1)
		{0x20000, 2}, // CJK ext B
		{0xE0100, 0}, // variation selector
		{0x1160, 0},  // hangul jungseong filler (combining range)
		{0x303F, 1},  // explicitly narrow inside wide block
	}
	for _, c := range cases {
		if got := Wcwidth(c.cp); got != c.want {
			t.Fatalf("Wcwidth(%#x) = %d want %d", c.cp, got, c.want)
		}
	}
	if got := GetStringCellWidth("ab世"); got != 4 {
		t.Fatalf("GetStringCellWidth = %d", got)
	}
	if got := GetStringCellWidth("é"); got != 1 {
		t.Fatalf("combining width = %d", got)
	}
}

func TestTermFullReset(t *testing.T) {
	term := newTestTerminal(10, 3)
	term.WriteString("abc\x1b[31m\x1b[?1h")
	term.WriteString("\x1bc") // RIS
	if got := line(term, 0); got != "" {
		t.Fatalf("after RIS row0 = %q", got)
	}
	if term.CoreService().DecPrivateModes.ApplicationCursorKeys {
		t.Fatal("modes not reset")
	}
	attr := term.InputHandler().GetAttrData()
	if attr.Fg != 0 {
		t.Fatal("attrs not reset")
	}
}

func TestTermCursorStyle(t *testing.T) {
	term := newTestTerminal(10, 3)
	term.WriteString("\x1b[4 q") // steady underline
	dm := &term.CoreService().DecPrivateModes
	if dm.CursorStyle == nil || *dm.CursorStyle != "underline" {
		t.Fatalf("style = %v", dm.CursorStyle)
	}
	if dm.CursorBlink == nil || *dm.CursorBlink {
		t.Fatalf("blink = %v", dm.CursorBlink)
	}
	term.WriteString("\x1b[0 q")
	if dm.CursorStyle != nil || dm.CursorBlink != nil {
		t.Fatal("style not reset")
	}
}
