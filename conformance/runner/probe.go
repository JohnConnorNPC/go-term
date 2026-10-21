package runner

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Probes are the probe byte sequences shared with the generator (conformance/probes.json, D6).
type Probes struct {
	Doc    string      `json:"_doc"`
	Prefix string      `json:"prefix"`
	T1     [][2]string `json:"t1"`
	T2     struct {
		Start string `json:"start"`
		Step  string `json:"step"`
		Max   int    `json:"max"`
	} `json:"t2"`
	T3 string `json:"t3"`
	T4 struct {
		Pen    string `json:"pen"`
		Glyphs string `json:"glyphs"`
	} `json:"t4"`
}

// LoadProbes reads probes.json.
func LoadProbes(path string) (*Probes, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var p Probes
	if err := json.Unmarshal(b, &p); err != nil || len(p.T1) == 0 {
		return nil, fmt.Errorf("probes %s: %v (t1=%d)", path, err, len(p.T1))
	}
	return &p, nil
}

// twin replays the case on a fresh engine, then writes the CAN prefix.
func twin(f Factory, g *Golden, pr *Probes) (Engine, error) {
	e, err := start(f, g)
	if err != nil {
		return nil, err
	}
	for _, op := range g.Ops {
		if err := apply(e, op, nil); err != nil {
			return nil, err
		}
	}
	return e, e.Write([]byte(pr.Prefix))
}

// reply writes s and returns the replies the engine emitted during that write.
func reply(e Engine, s string) (string, error) {
	before, err := e.Snapshot()
	if err != nil {
		return "", err
	}
	if err := e.Write([]byte(s)); err != nil {
		return "", err
	}
	after, err := e.Snapshot()
	if err != nil || len(after.P.Replies) < len(before.P.Replies) {
		return "", fmt.Errorf("replies shrank or %v", err)
	}
	return strings.Join(after.P.Replies[len(before.P.Replies):], ""), nil
}

func activeRow(st State) (Row, error) {
	rows, y := st.Normal, st.P.BaseY+st.P.Cursor[1]
	if st.P.Buffer == "alternate" {
		rows = st.Alt
	}
	if y < 0 || y >= len(rows) {
		return Row{}, fmt.Errorf("cursor row %d outside %d rows", y, len(rows))
	}
	return rows[y], nil
}

// probe runs the four twins T1-T4 exactly as gen/lib/probe*.mjs does.
func probe(f Factory, g *Golden, pr *Probes) (B, error) {
	out := B{T1: map[string]string{}, T2: []int{}}
	e, err := twin(f, g, pr)
	for _, p := range pr.T1 {
		if err == nil {
			out.T1[p[0]], err = reply(e, p[1])
		}
	}
	if err == nil {
		e, err = twin(f, g, pr)
	}
	var st State
	if err == nil {
		if err = e.Write([]byte(pr.T2.Start)); err == nil {
			st, err = e.Snapshot()
		}
	}
	// T2: cursorX after each HT, starting from CR (the start is not in the list), at most Max HTs.
	for i := 0; err == nil && i < min(st.P.Cols, pr.T2.Max); i++ {
		var s State
		if err = e.Write([]byte(pr.T2.Step)); err == nil {
			if s, err = e.Snapshot(); err == nil {
				out.T2 = append(out.T2, s.P.Cursor[0])
			}
		}
	}
	if err == nil {
		e, err = twin(f, g, pr)
	}
	if err == nil {
		out.T3, err = reply(e, pr.T3)
	}
	if err == nil {
		out.T4.Pen, out.T4.Glyphs, err = penGlyphs(f, g, pr)
	}
	return out, err
}

func penGlyphs(f Factory, g *Golden, pr *Probes) (Attr, []string, error) {
	e, err := twin(f, g, pr)
	if err != nil {
		return Attr{}, nil, err
	}
	if err := e.Write([]byte(pr.T4.Pen)); err != nil {
		return Attr{}, nil, err
	}
	st, err := e.Snapshot()
	row, err2 := activeRow(st)
	if err != nil || err2 != nil || len(row.Cells) == 0 {
		return Attr{}, nil, fmt.Errorf("pen: %v %v", err, err2)
	}
	pen := row.Cells[0].A
	gl := pr.T4.Glyphs[:min(len(pr.T4.Glyphs), 1+st.P.Cols)]
	if err := e.Write([]byte(gl)); err != nil {
		return pen, nil, err
	}
	if st, err = e.Snapshot(); err != nil {
		return pen, nil, err
	}
	if row, err = activeRow(st); err != nil || len(row.Cells) < len(gl)-1 {
		return pen, nil, fmt.Errorf("glyphs: %v", err)
	}
	glyphs := make([]string, len(gl)-1)
	for x := range glyphs {
		glyphs[x] = row.Cells[x].Ch
	}
	return pen, glyphs, nil
}
