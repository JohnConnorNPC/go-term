package runner

import "testing"

// keyed is the oracle with one modes key added to (val != nil) or removed from (val == nil) every snapshot.
type keyed struct {
	*oracle
	key string
	val any
}

func (k keyed) Snapshot() (State, error) {
	st, err := k.oracle.Snapshot()
	m := map[string]any{}
	for kk, v := range st.P.Modes {
		m[kk] = v
	}
	if k.val != nil {
		m[k.key] = k.val
	} else {
		delete(m, k.key)
	}
	st.P.Modes = m
	return st, err
}

// TestGoldenKeys (E3): maps are compared over the golden's keys. An extra engine key passes; a golden key the engine
// lacks is unobserved.
func TestGoldenKeys(t *testing.T) {
	pr, g := fixture(t)
	for _, c := range []struct {
		key    string
		val    any
		status Status
		diff   string
	}{{"synchronizedOutputMode", false, PASS, ""}, {"originMode", nil, FAIL, "unobserved step[0].P.modes.originMode"}} {
		f := func(Config) (Engine, error) { return keyed{&oracle{pr: pr, step: -1}, c.key, c.val}, nil }
		if r := RunCase(g, f, pr, Options{}); r.Status != c.status || r.Diff != c.diff {
			t.Errorf("modes.%s %v: %s %q, want %s %q", c.key, c.val, r.Status, r.Diff, c.status, c.diff)
		}
	}
}
