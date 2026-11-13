package vt

import (
	"strings"
	"testing"
)

func decodeSixel(t *testing.T, data string) *SixelImage {
	t.Helper()
	p := DefaultSixelPalette()
	d := NewSixelDecoder(&p, false)
	d.Write(utf32(data))
	return d.Image()
}

func utf32(s string) []uint32 {
	out := make([]uint32, 0, len(s))
	for _, r := range s {
		out = append(out, uint32(r)) // #nosec G115 -- test data
	}
	return out
}

// px is the RGBA of a pixel.
func px(img *SixelImage, x, y int) [4]byte {
	o := (y*img.Width + x) * 4
	return [4]byte(img.Pix[o : o+4])
}

func TestSixelBasic(t *testing.T) {
	// Register 1 red; "~" sets all six pixels, twice.
	img := decodeSixel(t, "#1;2;100;0;0#1~~")
	if img == nil || img.Width != 2 || img.Height != 6 {
		t.Fatalf("image %+v", img)
	}
	for y := range 6 {
		if got := px(img, 1, y); got != [4]byte{255, 0, 0, 255} {
			t.Errorf("pixel 1,%d = %v", y, got)
		}
	}
}

// TestSixelBitOrder: the low bit of a sixel is its top pixel.
func TestSixelBitOrder(t *testing.T) {
	// "@" is 0x40: bit 0 only. "_" is 0x60: bit 5 only.
	img := decodeSixel(t, "#2;2;0;100;0@_")
	if px(img, 0, 0)[3] != 255 || px(img, 0, 1)[3] != 0 {
		t.Errorf("@ drew %v %v", px(img, 0, 0), px(img, 0, 1))
	}
	if px(img, 1, 5)[3] != 255 || px(img, 1, 0)[3] != 0 {
		t.Errorf("_ drew %v %v", px(img, 1, 5), px(img, 1, 0))
	}
}

func TestSixelRepeatCarriageReturnNewline(t *testing.T) {
	// Three green across, back to the start, one blue over the first; then
	// the next band, two pixels in.
	img := decodeSixel(t, "#1;2;0;100;0!3~$#2;2;0;0;100@-#1??~")
	if img.Width != 3 || img.Height != 12 {
		t.Fatalf("size %dx%d", img.Width, img.Height)
	}
	if got := px(img, 0, 0); got != [4]byte{0, 0, 255, 255} {
		t.Errorf("overdrawn pixel %v", got)
	}
	if got := px(img, 0, 1); got != [4]byte{0, 255, 0, 255} {
		t.Errorf("pixel under it %v", got)
	}
	if got := px(img, 2, 6); got != [4]byte{0, 255, 0, 255} {
		t.Errorf("second band %v", got)
	}
	if got := px(img, 0, 6); got[3] != 0 {
		t.Errorf("unset pixel %v", got)
	}
}

func TestSixelRasterAttributesTruncate(t *testing.T) {
	img := decodeSixel(t, "\"1;1;3;4#1!10~-~")
	if img.Width != 3 || img.Height != 4 {
		t.Fatalf("size %dx%d, want the raster's 3x4", img.Width, img.Height)
	}
	// And they give the size even where nothing is drawn.
	img = decodeSixel(t, "\"1;1;20;10#1~")
	if img.Width != 20 || img.Height != 10 {
		t.Fatalf("size %dx%d, want 20x10", img.Width, img.Height)
	}
}

func TestSixelHLS(t *testing.T) {
	for _, c := range []struct {
		h    int
		want [4]byte
	}{
		{0, [4]byte{0, 0, 255, 255}},   // DEC's hue: blue at 0
		{120, [4]byte{255, 0, 0, 255}}, // red
		{240, [4]byte{0, 255, 0, 255}}, // green
	} {
		img := decodeSixel(t, "#5;1;"+itoa(c.h)+";50;100~")
		if got := px(img, 0, 0); got != c.want {
			t.Errorf("hue %d: %v, want %v", c.h, got, c.want)
		}
	}
}

func TestSixelLimits(t *testing.T) {
	img := decodeSixel(t, "#1!5000~")
	if img.Width != SixelMaxWidth {
		t.Errorf("width %d, want it cut at %d", img.Width, SixelMaxWidth)
	}
	img = decodeSixel(t, "#1"+strings.Repeat("~-", 1000))
	if img.Height != SixelMaxHeight {
		t.Errorf("height %d, want it cut at %d", img.Height, SixelMaxHeight)
	}
	if img := decodeSixel(t, ""); img != nil {
		t.Errorf("an empty picture is %+v", img)
	}
}

// TestSixelChunking: the data may arrive split anywhere, even inside a
// number.
func TestSixelChunking(t *testing.T) {
	data := "\"1;1;7;12#12;2;10;20;30#12!7~-#3;1;120;50;100!3N$!5w"
	whole := decodeSixel(t, data)
	runes := utf32(data)
	for cut := range runes {
		p := DefaultSixelPalette()
		d := NewSixelDecoder(&p, false)
		d.Write(runes[:cut])
		d.Write(runes[cut:])
		img := d.Image()
		if string(img.Pix) != string(whole.Pix) {
			t.Fatalf("split at %d decodes differently", cut)
		}
	}
}

func sixelTerminal(t *testing.T) (*Terminal, *[]*SixelImage, *strings.Builder) {
	t.Helper()
	term := newTestTerminal(80, 24)
	var imgs []*SixelImage
	term.OnSixel = func(img *SixelImage) { imgs = append(imgs, img) }
	var reply strings.Builder
	term.OnReply = func(s string) { reply.WriteString(s) }
	return term, &imgs, &reply
}

func TestSixelThroughTheTerminal(t *testing.T) {
	term, imgs, _ := sixelTerminal(t)
	term.WriteString("\x1b[41m\x1bP0;1;0q#1;2;100;100;100A\x1b\\")
	term.WriteString("\x1bPq#1~\x1b\\")
	if len(*imgs) != 2 {
		t.Fatalf("%d pictures", len(*imgs))
	}
	a, b := (*imgs)[0], (*imgs)[1]
	if !a.Transparent || b.Transparent {
		t.Errorf("P2: transparent %v %v", a.Transparent, b.Transparent)
	}
	if a.Attr.GetBgColor() != 1 {
		t.Errorf("background attribute %d, want red (1)", a.Attr.GetBgColor())
	}
	// The register defined in the first picture is still defined in the
	// second; a reset puts them back.
	if got := px(b, 0, 0); got != [4]byte{255, 255, 255, 255} {
		t.Errorf("register 1 in the next picture: %v", got)
	}
	term.WriteString("\x1bc\x1bPq#1~\x1b\\")
	if got := px((*imgs)[2], 0, 0); got == [4]byte{255, 255, 255, 255} {
		t.Errorf("register 1 survived a reset")
	}
	// Fill paints only what was not drawn.
	a.Fill(0x102030)
	if px(a, 0, 0) != [4]byte{0x10, 0x20, 0x30, 255} || px(a, 0, 1) != [4]byte{255, 255, 255, 255} {
		t.Errorf("fill %v %v", px(a, 0, 0), px(a, 0, 1))
	}
	// The DCS $ q handler is still DECRQSS, not sixel.
	n := len(*imgs)
	term.WriteString("\x1bP$qm\x1b\\")
	if len(*imgs) != n {
		t.Errorf("DECRQSS decoded as sixel")
	}
}

func TestSixelAdvertised(t *testing.T) {
	term, imgs, reply := sixelTerminal(t)
	term.WriteString("\x1b[c")
	if got := reply.String(); got != "\x1b[?62;4;9;22c" {
		t.Errorf("DA1 with sixel %q", got)
	}

	// Nobody to draw: no 4, and nothing decoded.
	plain := newTestTerminal(80, 24)
	var r strings.Builder
	plain.OnReply = func(s string) { r.WriteString(s) }
	plain.WriteString("\x1b[c\x1b[?1;1S")
	if got := r.String(); got != "\x1b[?1;2c" {
		t.Errorf("DA1 without a sixel sink %q", got)
	}

	// Switched off: the same, though someone is listening.
	term.Options.Sixel = false
	reply.Reset()
	term.WriteString("\x1b[c\x1bPq#1~\x1b\\")
	if got := reply.String(); got != "\x1b[?1;2c" || len(*imgs) != 0 {
		t.Errorf("DA1 with sixel off %q, %d pictures", got, len(*imgs))
	}
}

func TestXTSMGRAPHICS(t *testing.T) {
	term, _, reply := sixelTerminal(t)
	term.OnSixelGeometry = func() (int, int) { return 800, 5000 }
	for _, c := range []struct{ in, want string }{
		{"\x1b[?1;1S", "\x1b[?1;0;256S"},
		{"\x1b[?1;4S", "\x1b[?1;0;256S"},
		{"\x1b[?1;3;16S", "\x1b[?1;3S"},
		{"\x1b[?2;1S", "\x1b[?2;0;800;4096S"},
		{"\x1b[?2;4S", "\x1b[?2;0;4096;4096S"},
		{"\x1b[?3;1S", "\x1b[?3;1S"},
	} {
		reply.Reset()
		term.WriteString(c.in)
		if got := reply.String(); got != c.want {
			t.Errorf("%q: %q, want %q", c.in, got, c.want)
		}
	}
}
