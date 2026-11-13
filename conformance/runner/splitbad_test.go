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

// TestRunScoreSensitive (#9 -fuzz): with every split variant marked headless-sensitive, the default scoring skips them
// all and passes an engine that breaks one-byte writes; ScoreSensitive scores them against the whole-write golden and
// catches it as SPLIT, naming each variant it scored. The oracle (split-safe) passes both ways.
func TestRunScoreSensitive(t *testing.T) {
	pr, g := fixture(t)
	g.Sensitive = nil
	for _, v := range variants(g) {
		g.Sensitive = append(g.Sensitive, v.id)
	}
	bad := func(Config) (Engine, error) { return &splitbad{oracle: &oracle{pr: pr, step: -1}}, nil }
	if r := RunCase(g, bad, pr, Options{Splits: true}); r.Status != PASS || len(r.Skipped) != len(g.Sensitive) || r.Rescored != nil {
		t.Fatalf("default: %s skipped %d rescored %v, want PASS with all %d skipped", r.Status, len(r.Skipped), r.Rescored, len(g.Sensitive))
	}
	if r := RunCase(g, bad, pr, Options{Splits: true, ScoreSensitive: true}); r.Status != SPLIT || r.Skipped != nil {
		t.Fatalf("ScoreSensitive: %s %s skipped %v, want SPLIT and none skipped", r.Status, r.Diff, r.Skipped)
	}
	r := RunCase(g, Oracle(pr), pr, Options{Splits: true, ScoreSensitive: true})
	if r.Status != PASS || len(r.Rescored) != len(g.Sensitive) || r.Skipped != nil {
		t.Fatalf("oracle ScoreSensitive: %s %s rescored %d of %d", r.Status, r.Diff, len(r.Rescored), len(g.Sensitive))
	}
}
