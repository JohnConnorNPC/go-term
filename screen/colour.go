package screen

import (
	"fmt"

	"github.com/JohnConnorNPC/go-term/vt"
)

// Default is the colour name of an unset colour (SGR 39/49 or never set).
const Default = "default"

// paletteNames are the 8 names of SGR 30-37 / 40-47, as pyte spells them
// (brown, not yellow). SGR 90-97 / 100-107 add the prefix "bright"; pyte's bg
// 105 typo "bfightmagenta" is not reproduced (Studio's named exception).
var paletteNames = [8]string{"black", "red", "green", "brown", "blue", "magenta", "cyan", "white"}

// palette256 is pyte's 256-colour table (graphics.FG_BG_256) as 6-hex strings:
// the 16 xterm-like entries, the 6x6x6 cube over 00 5f 87 af d7 ff, then the
// grey ramp 08+10i.
var palette256 = func() (t [256]string) {
	first16 := [16][3]int{
		{0x00, 0x00, 0x00}, {0xcd, 0x00, 0x00}, {0x00, 0xcd, 0x00}, {0xcd, 0xcd, 0x00},
		{0x00, 0x00, 0xee}, {0xcd, 0x00, 0xcd}, {0x00, 0xcd, 0xcd}, {0xe5, 0xe5, 0xe5},
		{0x7f, 0x7f, 0x7f}, {0xff, 0x00, 0x00}, {0x00, 0xff, 0x00}, {0xff, 0xff, 0x00},
		{0x5c, 0x5c, 0xff}, {0xff, 0x00, 0xff}, {0x00, 0xff, 0xff}, {0xff, 0xff, 0xff},
	}
	for i, c := range first16 {
		t[i] = fmt.Sprintf("%02x%02x%02x", c[0], c[1], c[2])
	}
	steps := [6]int{0x00, 0x5f, 0x87, 0xaf, 0xd7, 0xff}
	for i := 0; i < 216; i++ {
		t[16+i] = fmt.Sprintf("%02x%02x%02x", steps[(i/36)%6], steps[(i/6)%6], steps[i%6])
	}
	for i := 0; i < 24; i++ {
		v := 8 + 10*i
		t[232+i] = fmt.Sprintf("%02x%02x%02x", v, v, v)
	}
	return t
}()

// Colour names a vt colour word (the Fg or Bg of a cell; only its colour mode
// and value bits are read) in Studio's vocabulary, which is pyte's:
//   - "default" when no colour is set;
//   - a 16-colour SGR (30-37, 40-47, 90-97, 100-107) gives black red green
//     brown blue magenta cyan white, or the same with the prefix "bright";
//   - 38;5;n / 48;5;n gives pyte's 256 table entry as 6-hex lowercase, for every
//     n, 0-15 included (n=1 is "cd0000", not "red");
//   - 38;2;r;g;b / 48;2;r;g;b gives "%02x%02x%02x".
//
// Dim (SGR 2) lives in other bits and is ignored, as pyte has no dim.
func Colour(word uint32) string {
	switch word & vt.AttrCMMask {
	case vt.AttrCMP16:
		n := word & 0x0f
		if n >= 8 {
			return "bright" + paletteNames[n&7]
		}
		return paletteNames[n]
	case vt.AttrCMP256:
		return palette256[word&0xff]
	case vt.AttrCMRGB:
		return fmt.Sprintf("%06x", word&vt.AttrRGBMask)
	}
	return Default
}
