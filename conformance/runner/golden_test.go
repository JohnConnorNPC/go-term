package runner

import (
	"os"
	"strings"
	"testing"
)

func TestInSha(t *testing.T) {
	ops := []Op{{Op: "w", Bytes: []byte("a\x1b[1m")}, {Op: "w", Bytes: []byte{0xff}}, {Op: "resize", Cols: 10, Rows: 4}, {Op: "reset"}}
	if got, want := InSha(ops), "4843d07dd5fbb7716b36708c30802b462bbb212d5317d870383b714d8c4fbc1f"; got != want {
		t.Fatalf("InSha = %s, want %s", got, want)
	}
}

func TestExpand(t *testing.T) {
	attrs := []Attr{{}, {Bold: true}}
	var cells []Cell
	for _, run := range [][4]any{{1.0, 0.0, 1.0, "a"}, {1.0, 1.0, 2.0, "中"}, {1.0, 1.0, 0.0, ""}, {1.0, 0.0, 1.0, "é"}, {2.0, 0.0, 1.0, "xy"}, {3.0, 0.0, 1.0, ""}} {
		c, err := Expand(run, attrs)
		if err != nil {
			t.Fatal(err)
		}
		cells = append(cells, c...)
	}
	want := []Cell{{"a", 1, Attr{}}, {"中", 2, Attr{Bold: true}}, {"", 0, Attr{Bold: true}}, {"é", 1, Attr{}},
		{"x", 1, Attr{}}, {"y", 1, Attr{}}, {"", 1, Attr{}}, {"", 1, Attr{}}, {"", 1, Attr{}}}
	if len(cells) != len(want) {
		t.Fatalf("got %d cells, want %d: %v", len(cells), len(want), cells)
	}
	for i := range want {
		if cells[i] != want[i] {
			t.Errorf("cell %d = %v, want %v", i, cells[i], want[i])
		}
	}
	if _, err := Expand([4]any{1.0, 5.0, 1.0, "a"}, attrs); err == nil {
		t.Error("attr index out of range accepted")
	}
}

func TestStale(t *testing.T) {
	data, err := os.ReadFile("../golden/x53/c0/smoke-text.json")
	if err != nil {
		t.Fatal(err)
	}
	g, err := Parse(data)
	if err != nil || g.Stale() {
		t.Fatalf("pristine golden: err=%v stale=%v", err, g != nil && g.Stale())
	}
	for name, edit := range map[string][2]string{"in": {`"in":"68656c6c6f"`, `"in":"68656c6c70"`}, "inSha": {`"inSha": "`, `"inSha": "0`}} {
		bad := strings.Replace(string(data), edit[0], edit[1], 1)
		if bad == string(data) {
			t.Fatalf("%s: edit anchor not found", name)
		}
		if g, err := Parse([]byte(bad)); err != nil || !g.Stale() {
			t.Errorf("%s edited: err=%v, want STALE", name, err)
		}
	}
}

func TestProvenance(t *testing.T) {
	g := &Golden{ID: "c0/x", Ref: "x53", Pins: "p"}
	if r := provenance(g, "x53", "c", "c", "p"); r != "" {
		t.Fatalf("matching provenance reported %q", r)
	}
	for name, r := range map[string]string{"ref": provenance(g, "xcur", "c", "c", "p"), "case": provenance(g, "x53", "c2", "c", "p"),
		"missing": provenance(g, "x53", "", "", "p"), "pins": provenance(g, "x53", "c", "c", "p2")} {
		if r == "" {
			t.Errorf("%s mismatch not reported", name)
		}
	}
}
