package runner

import "testing"

// unobs is the oracle with State.Unobserved set to paths.
type unobs struct {
	*oracle
	paths []string
}

func (u *unobs) Snapshot() (State, error) {
	st, err := u.oracle.Snapshot()
	st.Unobserved = u.paths
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
