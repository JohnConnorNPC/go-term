package screen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JohnConnorNPC/go-term/conformance/runner"
	"github.com/JohnConnorNPC/go-term/vt"
)

// xcurText is the part of a golden's last step that runner.Parse does not
// keep: xterm's text per line (translateToString) and the cursor-hidden flag.
type xcurText struct {
	Steps []struct {
		I    map[string]any `json:"I"`
		Text struct {
			Normal []string `json:"normal"`
			Alt    []string `json:"alt"`
		} `json:"text"`
	} `json:"steps"`
}

func trimAll(rows []string) string {
	out := make([]string, len(rows))
	for i, r := range rows {
		out[i] = strings.TrimRight(r, " ")
	}
	return strings.Join(out, "\n")
}

// checkXcur replays golden g into vt as the gate's adapter does and compares
// Read with the golden's last step: every normal-buffer line (history then
// grid, through a scrollback read) or the alt screen, trimmed; the cursor,
// clamped as Read clamps it; and CursorHidden.
func checkXcur(g *runner.Golden, root string) error {
	raw, err := os.ReadFile(filepath.Join(root, "golden", "xcur", g.ID+".json"))
	if err != nil {
		return err
	}
	var x xcurText
	if err := json.Unmarshal(raw, &x); err != nil || len(x.Steps) == 0 || len(g.Steps) != len(x.Steps) {
		return fmt.Errorf("golden text: %v (%d steps)", err, len(x.Steps))
	}
	o := vt.NewOptions()
	o.Cols, o.Rows, o.Scrollback = g.Init.Cols, g.Init.Rows, g.Init.Scrollback
	term := vt.NewTerminal(o)
	for _, op := range g.Ops {
		switch op.Op {
		case "w":
			term.Write(op.Bytes)
		case "resize":
			term.Resize(op.Cols, op.Rows)
		default:
			term.Reset()
		}
	}
	want, last := g.Steps[len(g.Steps)-1].P, x.Steps[len(x.Steps)-1]
	f := Read(term, Request{View: ViewScrollback, Lines: 1 << 30})
	wantText := last.Text.Normal
	if want.Buffer == "alternate" {
		wantText = last.Text.Alt
	}
	if got, w := trimAll(f.Lines), trimAll(wantText); got != w {
		return fmt.Errorf("lines\n%q\nwant\n%q", got, w)
	}
	wc := Cursor{Row: clamp(want.Cursor[1], 0, f.Rows-1), Col: clamp(want.Cursor[0], 0, f.Cols-1)}
	if *f.Cursor != wc || f.AltScreen != (want.Buffer == "alternate") {
		return fmt.Errorf("cursor %+v alt %v, golden %+v %s", *f.Cursor, f.AltScreen, wc, want.Buffer)
	}
	if hidden, ok := last.I["cursorHidden"].(bool); !ok || hidden != f.CursorHidden {
		return fmt.Errorf("CursorHidden %v, golden %v", f.CursorHidden, last.I["cursorHidden"])
	}
	return nil
}

// TestReadXcur: Read agrees with every @xterm/headless@6.0.0 golden (xcur),
// wide and combining cases (Unicode v6) included. The denominator is the case
// files, as in TestGateXcur.
func TestReadXcur(t *testing.T) {
	const root = "../conformance"
	counts, err := runner.Counts(root)
	if err != nil {
		t.Fatal(err)
	}
	cases := 0
	for _, n := range counts {
		cases += n
	}
	gs, err := runner.LoadAll(root, "xcur")
	if err != nil || len(gs) != cases || cases == 0 {
		t.Fatalf("xcur goldens %d, cases %d, err %v", len(gs), cases, err)
	}
	pass, wide, wideOK := 0, 0, 0
	for _, g := range gs {
		isWide := g.Category() == "wide"
		if isWide {
			wide++
		}
		if err := checkXcur(g, root); err != nil {
			t.Errorf("%s: %v", g.ID, err)
			continue
		}
		pass++
		if isWide {
			wideOK++
		}
	}
	if wide != 30 {
		t.Errorf("%d wide cases, want 30", wide)
	}
	t.Logf("VECTORS xcur read %d/%d (wide %d/%d)", pass, cases, wideOK, wide)
}

// TestStableHashWide: the cursor on either cell of a wide glyph, or on a cell
// holding a combining sequence, masks that whole glyph and nothing else.
func TestStableHashWide(t *testing.T) {
	term := newTerm(6, 1)
	term.WriteString("a中éz")
	for _, tc := range []struct {
		col  int
		want string
	}{
		{0, "\x00中éz "}, {1, "a\x00éz "}, {2, "a\x00éz "},
		{3, "a中\x00z "}, {4, "a中é\x00 "}, {5, "a中éz\x00"},
	} {
		term.WriteString(fmt.Sprintf("\x1b[1;%dH", tc.col+1))
		f := Read(term, Request{})
		if f.Text() != "a中éz " {
			t.Fatalf("text %q", f.Text())
		}
		if f.StableHash != ContentHash(tc.want) {
			t.Errorf("cursor col %d: stable hash is not ContentHash(%q)", tc.col, tc.want)
		}
	}
}
