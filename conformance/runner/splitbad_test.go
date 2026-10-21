package runner

import "testing"

// splitbad is the oracle, except that a one-byte case write (not a probe byte) corrupts its later snapshots:
// it passes whole writes and must be caught as SPLIT.
type splitbad struct {
	*oracle
	one bool
}

func (s *splitbad) Write(p []byte) error {
	if len(p) == 1 && p[0] != 0x18 && p[0] != '\t' && p[0] != '\r' {
		s.one = true
	}
	return s.oracle.Write(p)
}

func (s *splitbad) Snapshot() (State, error) {
	st, err := s.oracle.Snapshot()
	if s.one {
		st.P.Cursor[0]++
	}
	return st, err
}

func TestRunSplit(t *testing.T) {
	pr, g := fixture(t)
	f := func(Config) (Engine, error) { return &splitbad{oracle: &oracle{pr: pr, step: -1}}, nil }
	if r := RunCase(g, f, pr, Options{Splits: false}); r.Status != PASS {
		t.Fatalf("whole writes: %s %s, want PASS", r.Status, r.Diff)
	}
	if r := RunCase(g, f, pr, Options{Splits: true}); r.Status != SPLIT {
		t.Fatalf("split variants: %s %s, want SPLIT", r.Status, r.Diff)
	}
}
