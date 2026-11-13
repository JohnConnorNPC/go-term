package screen

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/JohnConnorNPC/go-term/vt"
)

// cell0 writes "\x1b[<sgr>mX" into a fresh 4x1 terminal and returns the Fg and
// Bg colour names of cell 0.
func cell0(t *testing.T, sgr string) (fg, bg string) {
	t.Helper()
	opts := vt.NewOptions()
	opts.Cols, opts.Rows = 4, 1
	term := vt.NewTerminal(opts)
	term.WriteString("\x1b[" + sgr + "mX")
	b := term.Buffer()
	line := b.Lines.Get(b.YBase)
	return Colour(line.GetFg(0)), Colour(line.GetBg(0))
}

// TestColourTable pins every colour of Studio's R11 list, written through vt.
func TestColourTable(t *testing.T) {
	names := []string{"black", "red", "green", "brown", "blue", "magenta", "cyan", "white"}
	type tc struct{ sgr, fg, bg string }
	var cases []tc
	for i, n := range names {
		cases = append(cases,
			tc{fmt.Sprint(30 + i), n, Default},
			tc{fmt.Sprint(40 + i), Default, n},
			tc{fmt.Sprint(90 + i), "bright" + n, Default},
			tc{fmt.Sprint(100 + i), Default, "bright" + n})
	}
	for n, hex := range map[int]string{0: "000000", 1: "cd0000", 8: "7f7f7f", 15: "ffffff", 16: "000000", 100: "878700", 244: "808080", 255: "eeeeee"} {
		cases = append(cases, tc{fmt.Sprintf("38;5;%d", n), hex, Default}, tc{fmt.Sprintf("48;5;%d", n), Default, hex})
	}
	cases = append(cases,
		tc{"2", Default, Default}, // dim is ignored
		tc{"39", Default, Default}, tc{"49", Default, Default},
		tc{"31;39", Default, Default},
		tc{"38;2;10;20;30", "0a141e", Default}, tc{"48;2;10;20;30", Default, "0a141e"},
		tc{"38:2::255:0:128", "ff0080", Default},
		tc{"1;31", "red", Default}, // bold does not brighten
	)
	if len(cases) != 56 {
		t.Fatalf("%d cases, want 56", len(cases))
	}
	for _, c := range cases {
		fg, bg := cell0(t, c.sgr)
		if fg != c.fg || bg != c.bg {
			t.Errorf("SGR %s: fg %q bg %q, want %q %q", c.sgr, fg, bg, c.fg, c.bg)
		}
	}
	if fg, bg := cell0(t, "105"); fg != Default || bg != "brightmagenta" {
		t.Errorf("SGR 105 (pyte's bfightmagenta): %q %q", fg, bg)
	}
}

// TestPyte256 compares all 256 entries with pyte 0.8.2's FG_BG_256 dump.
func TestPyte256(t *testing.T) {
	f, err := os.Open("testdata/pyte256.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	n := 0
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if strings.HasPrefix(sc.Text(), "#") {
			continue
		}
		var i int
		var hex string
		if _, err := fmt.Sscanf(sc.Text(), "%d %s", &i, &hex); err != nil {
			t.Fatalf("line %q: %v", sc.Text(), err)
		}
		if i != n {
			t.Fatalf("line %q: index %d, want %d", sc.Text(), i, n)
		}
		if got := Colour(vt.AttrCMP256 | uint32(i)); got != hex {
			t.Errorf("38;5;%d = %q, pyte %q", i, got, hex)
		}
		n++
	}
	if n != 256 {
		t.Fatalf("read %d entries, want 256 (unmeasured)", n)
	}
	t.Logf("VECTORS pyte256 %d/256", n)
}
