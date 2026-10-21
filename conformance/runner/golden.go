package runner

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Op is one case step as the golden records it; Bytes is the decoded `in` (what headless's term.write received).
type Op struct {
	Op    string `json:"op"`
	In    string `json:"in"`
	Cols  int    `json:"cols"`
	Rows  int    `json:"rows"`
	Bytes []byte `json:"-"`
}

// Init is the terminal configuration a case starts from.
type Init struct {
	Cols       int    `json:"cols"`
	Rows       int    `json:"rows"`
	Scrollback int    `json:"scrollback"`
	Unicode    string `json:"unicode"`
}

// B holds the raw probe-twin results (D6).
type B struct {
	T1 map[string]string `json:"t1"`
	T2 []int             `json:"t2"`
	T3 string            `json:"t3"`
	T4 struct {
		Glyphs []string `json:"glyphs"`
		Pen    Attr     `json:"pen"`
	} `json:"t4"`
}

// Golden is one parsed golden file.
type Golden struct {
	ID, Ref, InSha, Pins string
	Init                 Init
	Ops                  []Op
	Steps                []State
	B                    B
	Sensitive            []string
	stale                string // provenance mismatch found by LoadAll (S3)
}

// Category is the case's directory (c0, esc, ...).
func (g *Golden) Category() string { return strings.SplitN(g.ID, "/", 2)[0] }

type rawRow struct {
	C [][4]any `json:"c"`
	W bool     `json:"w"`
}

type rawGolden struct {
	B     B      `json:"B"`
	ID    string `json:"id"`
	Init  Init   `json:"init"`
	InSha string `json:"inSha"`
	Pins  string `json:"pins"`
	Ref   string `json:"ref"`
	Split struct {
		Sensitive []string `json:"sensitive"`
		Variants  int      `json:"variants"`
	} `json:"split"`
	Steps []struct {
		Op    Op             `json:"op"`
		P     P              `json:"P"`
		I     map[string]any `json:"I"`
		Attrs []struct {
			P Attr `json:"P"`
			I any  `json:"I"`
		} `json:"attrs"`
		Cells struct {
			Normal []rawRow `json:"normal"`
			Alt    []rawRow `json:"alt"`
		} `json:"cells"`
		Text any `json:"text"`
	} `json:"steps"`
}

// Load reads and parses one golden file.
func Load(path string) (*Golden, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(b)
}

// Parse decodes a golden strictly (unknown keys are errors) and expands its RLE rows.
func Parse(data []byte) (*Golden, error) {
	var r rawGolden
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&r); err != nil {
		return nil, err
	}
	g := &Golden{ID: r.ID, Ref: r.Ref, InSha: r.InSha, Pins: r.Pins, Init: r.Init, B: r.B, Sensitive: r.Split.Sensitive}
	for i, st := range r.Steps {
		op := st.Op
		b, err := hex.DecodeString(op.In)
		if err != nil {
			return nil, fmt.Errorf("%s step %d in: %w", r.ID, i, err)
		}
		op.Bytes = b
		attrs := make([]Attr, len(st.Attrs))
		for k, a := range st.Attrs {
			attrs[k] = a.P
		}
		s := State{P: st.P, I: st.I}
		if s.Normal, err = expandRows(st.Cells.Normal, attrs); err == nil && st.Cells.Alt != nil {
			s.Alt, err = expandRows(st.Cells.Alt, attrs)
		}
		if err != nil {
			return nil, fmt.Errorf("%s step %d: %w", r.ID, i, err)
		}
		g.Ops, g.Steps = append(g.Ops, op), append(g.Steps, s)
	}
	return g, nil
}

func expandRows(rows []rawRow, attrs []Attr) ([]Row, error) {
	out := make([]Row, len(rows))
	for y, r := range rows {
		out[y].Wrapped = r.W
		for _, run := range r.C {
			cells, err := Expand(run, attrs)
			if err != nil {
				return nil, fmt.Errorf("row %d: %w", y, err)
			}
			out[y].Cells = append(out[y].Cells, cells...)
		}
	}
	return out, nil
}

// Expand turns one RLE run [n, attr, width, chars] into cells (D9): chars "" = n empty cells, otherwise one code point
// per cell, or a single cell holding several code points when n is 1.
func Expand(run [4]any, attrs []Attr) ([]Cell, error) {
	n, ok1 := run[0].(float64)
	a, ok2 := run[1].(float64)
	w, ok3 := run[2].(float64)
	t, ok4 := run[3].(string)
	if !ok1 || !ok2 || !ok3 || !ok4 || int(a) < 0 || int(a) >= len(attrs) {
		return nil, fmt.Errorf("bad run %v", run)
	}
	cell := Cell{W: int(w), A: attrs[int(a)]}
	out := make([]Cell, int(n))
	switch cps := []rune(t); {
	case t == "":
	case n == 1:
		cell.Ch = t
	case len(cps) == int(n):
		for i := range out {
			out[i] = cell
			out[i].Ch = string(cps[i])
		}
		return out, nil
	default:
		return nil, fmt.Errorf("run %v: %d code points for %v cells", run, len(cps), n)
	}
	for i := range out {
		out[i] = cell
	}
	return out, nil
}

// InSha is the D3 digest of the ops: W<len>:<bytes> | R<cols>,<rows>; | X;
func InSha(ops []Op) string {
	h := sha256.New()
	for _, o := range ops {
		switch o.Op {
		case "w":
			fmt.Fprintf(h, "W%d:", len(o.Bytes))
			h.Write(o.Bytes)
		case "resize":
			fmt.Fprintf(h, "R%d,%d;", o.Cols, o.Rows)
		default:
			h.Write([]byte("X;"))
		}
	}
	return hex.EncodeToString(h.Sum(nil))
}

// Stale reports a golden whose ops no longer hash to its stored inSha (a hand-edited `in`, or an edited digest), or
// whose provenance does not match the tree (S3).
func (g *Golden) Stale() bool { return g.stale != "" || InSha(g.Ops) != g.InSha }

// provenance (S3) returns why a golden does not belong to the current tree: its ref directory, the current case file
// (sha256 vs the manifest's record), and the current pins.json; "" when all match.
func provenance(g *Golden, dirRef, caseSha, manifestSha, pinsSha string) string {
	switch {
	case g.Ref != dirRef:
		return "ref " + g.Ref + " in directory " + dirRef
	case caseSha == "" || caseSha != manifestSha:
		return "case file " + g.ID + ".json changed or missing since gen (manifest " + manifestSha + ")"
	case g.Pins != pinsSha:
		return "pins " + g.Pins + " != sha256(gen/pins.json) " + pinsSha
	}
	return ""
}
