package runner

import "testing"

// unobs is the oracle with State.Unobserved set to paths; with afterCAN it lists them only once the probe prefix is
// written, so the whole-write steps pass and only the probe twins see the list.
type unobs struct {
	*oracle
	paths         []string
	afterCAN, can bool
}

func (u *unobs) Write(p []byte) error {
	u.can = u.can || string(p) == u.pr.Prefix
	return u.oracle.Write(p)
}

func (u *unobs) Snapshot() (State, error) {
	st, err := u.oracle.Snapshot()
	if !u.afterCAN || u.can {
		st.Unobserved = u.paths
	}
	return st, err
}

// golden is a one-cell state with one mode and an empty opts map.
func golden() State {
	return State{P: P{Cols: 2, Rows: 1, Modes: map[string]any{"insertMode": false}, Opts: map[string]any{}},
		Normal: []Row{{Cells: []Cell{{Ch: "a", W: 1}}}}}
}

// diffs checks Compare(golden, golden + Unobserved [entry]) against want per entry ("" = PASS).
func diffs(t *testing.T, want map[string]string) {
	t.Helper()
	for entry, d := range want {
		got := golden()
		got.Unobserved = []string{entry}
		if got := Compare("step[0]", golden(), got); got != d {
			t.Errorf("Unobserved [%s]: diff %q, want %q", entry, got, d)
		}
	}
}

// TestUnobserved: a listed P path the golden has fails whatever the engine's value, in Compare's diff form; a listed
// path the golden lacks (an absent key, an empty map) passes, and so does an empty list.
func TestUnobserved(t *testing.T) {
	diffs(t, map[string]string{
		"lineFeeds":        "unobserved step[0].P.lineFeeds",
		"cursor.x":         "unobserved step[0].P.cursor.x",
		"modes":            "unobserved step[0].P.modes",
		"modes.insertMode": "unobserved step[0].P.modes.insertMode",
		"foo":              "unobserved step[0].P.foo (unknown path)",
		"modes.noSuchKey":  "",
		"opts":             "",
	})
	if d := Compare("step[0]", golden(), golden()); d != "" {
		t.Errorf("no Unobserved list: diff %q, want PASS", d)
	}
	pr, g := fixture(t)
	f := func(Config) (Engine, error) {
		return &unobs{oracle: &oracle{pr: pr, step: -1}, paths: []string{"lineFeeds"}}, nil
	}
	if r := RunCase(g, f, pr, Options{}); r.Status != FAIL || r.Diff != "unobserved step[0].P.lineFeeds" {
		t.Errorf("RunCase with lineFeeds unobserved: %s %q, want FAIL", r.Status, r.Diff)
	}
}

// TestUnobservedCells: wrapped and attrs name the golden's first row and first cell (w is no name: in Compare's diffs
// ".w" is a cell's width); a golden with no rows passes.
func TestUnobservedCells(t *testing.T) {
	diffs(t, map[string]string{
		"wrapped": "unobserved step[0].normal[0].wrapped",
		"attrs":   "unobserved step[0].normal[0][0].attr",
		"w":       "unobserved step[0].P.w (unknown path)",
	})
	empty := State{P: golden().P}
	if d := Compare("step[0]", empty, State{P: empty.P, Unobserved: []string{"wrapped", "attrs"}}); d != "" {
		t.Errorf("wrapped/attrs with no golden rows: diff %q, want PASS", d)
	}
}

// TestUnobservedProbe: a probe twin that reads a listed path (listed only after the CAN prefix, so every whole-write
// step passes) is a FAIL in the B diff form, even when the golden reply is "", and never a "B: " error.
func TestUnobservedProbe(t *testing.T) {
	pr, g := fixture(t)
	cp := *g
	cp.B.T1 = map[string]string{}
	for k, v := range g.B.T1 {
		cp.B.T1[k] = v
	}
	cp.B.T1[pr.T1[0][0]] = ""
	for entry, want := range map[string]string{
		"replies":  "unobserved B.t1." + pr.T1[0][0] + ": replies",
		"cursor.x": "unobserved B.t2: cursor.x",
		"attrs":    "unobserved B.t4.pen: attrs",
	} {
		f := func(Config) (Engine, error) {
			return &unobs{oracle: &oracle{pr: pr, step: -1}, paths: []string{entry}, afterCAN: true}, nil
		}
		if r := RunCase(&cp, f, pr, Options{}); r.Status != FAIL || r.Diff != want {
			t.Errorf("%s after CAN: %s %q, want FAIL %q", entry, r.Status, r.Diff, want)
		}
	}
}
