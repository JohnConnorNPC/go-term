package screen

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// pyRuns decodes a python attr_runs output: a list of $d dicts.
func pyRuns(t *testing.T, raw json.RawMessage) []AttrRun {
	t.Helper()
	var list []json.RawMessage
	mustUnmarshal(t, raw, &list)
	out := []AttrRun{}
	for _, item := range list {
		d, ok := pyDict(t, item)
		if !ok {
			t.Fatalf("run %s is not a dict", item)
		}
		var r AttrRun
		for k, v := range map[string]any{"row": &r.Row, "col": &r.Col, "len": &r.Len, "fg": &r.Fg, "bg": &r.Bg,
			"reverse": &r.Reverse, "bold": &r.Bold, "underscore": &r.Underscore} {
			mustUnmarshal(t, d[k], v)
		}
		out = append(out, r)
	}
	return out
}

// runPyAttrRuns replays a python pyte.attr_runs case (its replay chunks at its
// geometry) and returns AttrRuns with the case's cap, and the expected runs.
func runPyAttrRuns(t *testing.T, c pyCase) (got, want []AttrRun) {
	t.Helper()
	var r struct {
		Chunks     []string `json:"chunks_b64"`
		Cols, Rows int
	}
	mustUnmarshal(t, c.Replay, &r)
	args, _ := pyDict(t, c.Args)
	var max int
	mustUnmarshal(t, args["cap"], &max)
	term := newTerm(r.Cols, r.Rows)
	for _, ch := range r.Chunks {
		b, err := base64.StdEncoding.DecodeString(ch)
		if err != nil {
			t.Fatal(err)
		}
		term.Write(b)
	}
	got = AttrRuns(term, max)
	if got == nil {
		got = []AttrRun{}
	}
	return got, pyRuns(t, c.Output)
}

// TestPythonColours: Studio's python/colours.json, pyte.attr_runs over one
// SGR each, exact; the file's named exception (bg 105: pyte "bfightmagenta",
// go "brightmagenta") is applied to the expected value.
func TestPythonColours(t *testing.T) {
	var f struct {
		Named struct{ SGR, Python, Go string } `json:"named_exception"`
	}
	mustUnmarshal(t, readVectors(t, "python/colours.json"), &f)
	if f.Named.SGR != "105" || f.Named.Python != "bfightmagenta" || f.Named.Go != "brightmagenta" {
		t.Fatalf("named_exception %+v", f.Named)
	}
	n, ok := 0, 0
	for _, c := range loadPython(t, "python/colours.json") {
		if c.Entry != "pyte.attr_runs" {
			continue
		}
		n++
		got, want := runPyAttrRuns(t, c)
		for i := range want {
			if want[i].Bg == f.Named.Python {
				want[i].Bg = f.Named.Go
			}
		}
		if fmt.Sprint(got) != fmt.Sprint(want) {
			t.Errorf("%s: %+v, want %+v", c.ID, got, want)
			continue
		}
		ok++
	}
	if n != 53 {
		t.Errorf("%d cases, want 53", n)
	}
	t.Logf("VECTORS python colours %d/%d", ok, n)
}

// TestPythonPyteAttrRuns: the pyte.attr_runs cases of python/pyte.json
// (styles, grouping, the cap) are exact too.
func TestPythonPyteAttrRuns(t *testing.T) {
	n, ok := 0, 0
	for _, c := range loadPython(t, "python/pyte.json") {
		if c.Entry != "pyte.attr_runs" {
			continue
		}
		n++
		got, want := runPyAttrRuns(t, c)
		if fmt.Sprint(got) != fmt.Sprint(want) {
			t.Errorf("%s: %+v, want %+v", c.ID, got, want)
			continue
		}
		ok++
	}
	if n != 11 {
		t.Errorf("%d cases, want 11", n)
	}
	t.Logf("VECTORS python pyte attr_runs %d/%d", ok, n)
}

// TestAttrRunsCap: 2000 runs fit; a screen with 2001 possible runs gives the
// first 2000, in row order.
func TestAttrRunsCap(t *testing.T) {
	for _, tc := range []struct{ runs, want int }{{2000, 2000}, {2001, 2000}, {1999, 1999}} {
		// 40 cols: 20 one-cell reverse runs per row, separated by plain cells.
		rows := (tc.runs + 19) / 20
		term := newTerm(40, rows)
		var b strings.Builder
		for i := 0; i < tc.runs; i++ {
			if i > 0 && i%20 == 0 {
				b.WriteString("\r\n")
			}
			b.WriteString("\x1b[7mX\x1b[0m.")
		}
		term.WriteString(b.String())
		got := AttrRuns(term, MaxAttrRuns)
		if len(got) != tc.want {
			t.Errorf("%d runs on screen: %d, want %d", tc.runs, len(got), tc.want)
			continue
		}
		last := got[len(got)-1]
		if wantRow, wantCol := (tc.want-1)/20, 2*((tc.want-1)%20); last.Row != wantRow || last.Col != wantCol || !last.Reverse || last.Len != 1 {
			t.Errorf("%d runs: last %+v, want row %d col %d", tc.runs, last, wantRow, wantCol)
		}
		if f := Read(term, Request{Attrs: true}); len(f.AttrRuns) != tc.want {
			t.Errorf("Read Attrs: %d runs", len(f.AttrRuns))
		}
	}
}

func TestAttrRuns(t *testing.T) {
	for _, tc := range []struct{ name, in, want string }{
		{"dim alone is no run", "\x1b[2mab", "[]"},
		{"italic alone is no run", "\x1b[3mab", "[]"},
		{"dim does not split a run", "\x1b[31ma\x1b[2mb", "[{0 0 2 red default false false false}]"},
		{"wide glyph counts 2 cells", "\x1b[1m中x", "[{0 0 3 default default false true false}]"},
		{"curly underline is underscore", "\x1b[4:3mab", "[{0 0 2 default default false false true}]"},
		{"bold does not brighten", "\x1b[1;31mab", "[{0 0 2 red default false true false}]"},
		{"bg erase is a run", "\x1b[44m\x1b[K", "[{0 0 6 default blue false false false}]"},
		{"two rows", "\x1b[7ma\r\nb", "[{0 0 1 default default true false false} {1 0 1 default default true false false}]"},
	} {
		term := newTerm(6, 2)
		term.WriteString(tc.in)
		if got := fmt.Sprint(AttrRuns(term, MaxAttrRuns)); got != tc.want {
			t.Errorf("%s: %s, want %s", tc.name, got, tc.want)
		}
	}
	term := newTerm(6, 2)
	term.WriteString("\x1b[7ma")
	if f := Read(term, Request{}); f.AttrRuns != nil {
		t.Error("AttrRuns filled without Request.Attrs")
	}
}
