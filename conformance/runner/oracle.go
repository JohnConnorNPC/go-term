package runner

import (
	"bytes"
	"strings"
)

// Oracle is the reference engine: it replays the golden it is scored on (caseAware), so it scores 100 % by
// construction (it alone gets the golden, via the unexported caseAware). It anchors the stub/mutation tests and the CLI's `-engine oracle`. During probe twins it recognises
// the probe sequences from probes.json and answers with the golden's raw B results.
func Oracle(pr *Probes) Factory {
	return func(Config) (Engine, error) { return &oracle{pr: pr, step: -1}, nil }
}

type oracle struct {
	pr     *Probes
	g      *Golden
	step   int
	pend   []byte
	st     State
	t2     int
	broken bool
}

func (o *oracle) begin(g *Golden) { o.g = g }

func (o *oracle) next(kind string) (Op, bool) {
	if o.g == nil || o.step+1 >= len(o.g.Ops) || o.g.Ops[o.step+1].Op != kind {
		return Op{}, false
	}
	return o.g.Ops[o.step+1], true
}

func (o *oracle) advance() { o.step, o.pend, o.st = o.step+1, nil, o.g.Steps[o.step+1] }

func (o *oracle) Write(p []byte) error {
	if o.g != nil && o.step+1 < len(o.g.Ops) {
		op, ok := o.next("w")
		o.pend = append(o.pend, p...)
		if !ok || !bytes.HasPrefix(op.Bytes, o.pend) {
			o.broken = true
		} else if len(o.pend) == len(op.Bytes) {
			o.advance()
		}
		return nil
	}
	o.probe(string(p))
	return nil
}

func (o *oracle) Resize(int, int) error { return o.other("resize") }
func (o *oracle) Reset() error          { return o.other("reset") }

func (o *oracle) other(kind string) error {
	if _, ok := o.next(kind); ok {
		o.advance()
	} else {
		o.broken = true
	}
	return nil
}

func (o *oracle) addReply(r string) {
	if r != "" {
		o.st.P.Replies = append(append([]string{}, o.st.P.Replies...), r)
	}
}

// probe applies one probe write to the replayed final state, from the golden's B results.
func (o *oracle) probe(s string) {
	b, pr := o.g.B, o.pr
	for _, p := range pr.T1 {
		if s == p[1] {
			o.addReply(b.T1[p[0]])
		}
	}
	switch {
	case s == pr.T2.Step && o.t2 < len(b.T2):
		o.st.P.Cursor[0], o.t2 = b.T2[o.t2], o.t2+1
	case s == pr.T3:
		o.addReply(b.T3)
	case s == pr.T4.Pen:
		o.setRow(func(r *Row) { r.Cells[0].A = b.T4.Pen })
	case len(s) > 1 && s != pr.T2.Start && strings.HasPrefix(pr.T4.Glyphs, s):
		o.setRow(func(r *Row) {
			for x, ch := range b.T4.Glyphs {
				r.Cells[x].Ch = ch
			}
		})
	}
}

// setRow edits a copy of the cursor row, so the golden itself is never changed.
func (o *oracle) setRow(edit func(*Row)) {
	rows := &o.st.Normal
	if o.st.P.Buffer == "alternate" {
		rows = &o.st.Alt
	}
	y := o.st.P.BaseY + o.st.P.Cursor[1]
	if y < 0 || y >= len(*rows) || len((*rows)[y].Cells) == 0 {
		o.broken = true
		return
	}
	cp := append([]Row{}, *rows...)
	cp[y].Cells = append([]Cell{}, cp[y].Cells...)
	edit(&cp[y])
	*rows = cp
}

func (o *oracle) Snapshot() (State, error) {
	if o.broken || o.step < 0 {
		return State{}, nil
	}
	return o.st, nil
}
