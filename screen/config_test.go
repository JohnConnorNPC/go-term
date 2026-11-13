package screen

import (
	"fmt"
	"strings"
	"testing"

	"github.com/JohnConnorNPC/go-term/vt"
)

// fill writes lines "L0".."L<n-1>", one per row, ending on the last one.
func fill(term *vt.Terminal, n int) {
	var b strings.Builder
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteString("\r\n")
		}
		fmt.Fprintf(&b, "L%d", i)
	}
	term.WriteString(b.String())
}

// history reads the whole history: its length and its oldest line.
func history(term *vt.Terminal) (int, string) {
	f := Read(term, Request{View: ViewScrollback, Lines: 1 << 30})
	if f.HistoryLines == 0 {
		return 0, ""
	}
	return f.HistoryLines, strings.TrimRight(f.Lines[0], " ")
}

func TestHistoryCap(t *testing.T) {
	var c Config
	for cols, want := range map[int]int{1: 1000, 80: 1000, 100: 1000, 101: 990, 1000: 100, 2000: 50, 200000: 1} {
		if got := c.HistoryCap(cols); got != want {
			t.Errorf("HistoryCap(%d) = %d, want %d", cols, got, want)
		}
	}
	if got := (Config{ScrollbackLines: 50}).HistoryCap(80); got != 50 {
		t.Errorf("lines 50: %d", got)
	}
	if got := (Config{ScrollbackCells: 5000}).HistoryCap(100); got != 50 {
		t.Errorf("cells 5000 at 100 cols: %d", got)
	}
}

// TestScrollbackCap: the caps hold in the terminal's own history, at the
// limit and one past it.
func TestScrollbackCap(t *testing.T) {
	const rows = 4
	for _, tc := range []struct {
		name    string
		cols    int
		c       Config
		written int // history lines written
		want    int
	}{
		{"1000 lines at the limit", 80, Config{}, 1000, 1000},
		{"1001 lines, one past", 80, Config{}, 1001, 1000},
		{"cols 100: 100000 cells = 1000 lines", 100, Config{}, 1001, 1000},
		{"cols 101: 990 lines", 101, Config{}, 991, 990},
		{"cols 101 at the limit", 101, Config{}, 990, 990},
		{"lines configurable", 80, Config{ScrollbackLines: 50}, 51, 50},
		{"cells configurable", 100, Config{ScrollbackCells: 5000}, 51, 50},
	} {
		term, err := NewTerminal(tc.cols, rows, tc.c)
		if err != nil {
			t.Fatal(err)
		}
		fill(term, tc.written+rows)
		n, oldest := history(term)
		if n != tc.want || oldest != fmt.Sprintf("L%d", tc.written-tc.want) {
			t.Errorf("%s: history %d oldest %q, want %d from L%d", tc.name, n, oldest, tc.want, tc.written-tc.want)
		}
	}
}

// TestResizeCap: Resize re-applies the cap for the new width.
func TestResizeCap(t *testing.T) {
	term, err := NewTerminal(100, 4, Config{})
	if err != nil {
		t.Fatal(err)
	}
	fill(term, 1000+4)
	if n, _ := history(term); n != 1000 {
		t.Fatalf("at 100 cols: %d", n)
	}
	if err := Resize(term, 101, 4, Config{}); err != nil {
		t.Fatal(err)
	}
	if n, oldest := history(term); n != 990 || oldest != "L10" {
		t.Errorf("at 101 cols: history %d oldest %q, want 990 from L10", n, oldest)
	}
	if term.Options.Scrollback != 990 {
		t.Errorf("scrollback option %d", term.Options.Scrollback)
	}
	fill(term, 50)
	if n, _ := history(term); n != 990 {
		t.Errorf("after more output: %d", n)
	}
	if err := Resize(term, 80, 4, Config{}); err != nil || term.Options.Scrollback != 1000 {
		t.Errorf("back to 80 cols: scrollback %d err %v", term.Options.Scrollback, err)
	}
}

func TestConfigUnicode(t *testing.T) {
	for _, u := range []string{"", "6"} {
		if _, err := NewTerminal(80, 24, Config{Unicode: u}); err != nil {
			t.Errorf("Unicode %q: %v", u, err)
		}
	}
	if _, err := NewTerminal(80, 24, Config{Unicode: "11"}); err == nil {
		t.Error("Unicode 11 accepted")
	}
	term, _ := NewTerminal(80, 24, Config{})
	if err := Resize(term, 80, 24, Config{Unicode: "11"}); err == nil {
		t.Error("Resize with Unicode 11 accepted")
	}
}
