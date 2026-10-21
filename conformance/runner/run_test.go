package runner

import (
	"strings"
	"testing"
)

// recorder is a stub that logs every write, to check the probe twins' byte sequences and their order.
type recorder struct {
	stub
	log *[]string
}

func (r recorder) Write(p []byte) error { *r.log = append(*r.log, string(p)); return nil }

func TestRunProbes(t *testing.T) {
	pr, g := fixture(t)
	var log []string
	// The recorder has no grid, so T4 errors after writing its probes; only the write log matters here.
	_, _ = probe(func(Config) (Engine, error) { return recorder{log: &log}, nil }, g, pr)
	all := strings.Join(log, "|")
	if n := strings.Count(all, pr.Prefix); n != 4 {
		t.Errorf("prefix written %d times, want 4 twins", n)
	}
	for _, p := range pr.T1 {
		if !strings.Contains(all, p[1]) {
			t.Errorf("T1 probe %s never written", p[0])
		}
	}
	if !strings.Contains(all, pr.T3) || !strings.Contains(all, pr.T4.Pen) {
		t.Error("T3/T4 probes missing")
	}
}

func fixture(t *testing.T) (*Probes, *Golden) {
	t.Helper()
	pr, err := LoadProbes("../probes.json")
	if err != nil {
		t.Fatal(err)
	}
	g, err := Load("../golden/x53/c0/smoke-text.json")
	if err != nil {
		t.Fatal(err)
	}
	return pr, g
}

// TestSkippedReported (S2): the variants headless itself is sensitive to are named in the result, not silently passed.
func TestSkippedReported(t *testing.T) {
	pr, _ := fixture(t)
	g, err := Load("../golden/x53/wide/zwsp.json")
	if err != nil || len(g.Sensitive) == 0 {
		t.Fatalf("zwsp golden: %v, sensitive %v", err, g.Sensitive)
	}
	r := RunCase(g, Oracle(pr), pr, Options{Splits: true})
	if r.Status != PASS || strings.Join(r.Skipped, ",") != strings.Join(g.Sensitive, ",") {
		t.Fatalf("status %s skipped %v, want PASS with %v", r.Status, r.Skipped, g.Sensitive)
	}
}
