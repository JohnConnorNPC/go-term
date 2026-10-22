package scoring

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/summary"
	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

// results is one PASS per case file (runner.Counts), except the first c0 case, which gets c0.
func results(t *testing.T, c0 runner.Status) []runner.Result {
	n, err := runner.Counts("../../..")
	if err != nil {
		t.Fatal(err)
	}
	var rs []runner.Result
	for c, k := range n {
		for i := 0; i < k; i++ {
			st := runner.PASS
			if c == "c0" && i == 0 {
				st = c0
			}
			rs = append(rs, runner.Result{ID: fmt.Sprintf("%s/%03d", c, i), Category: c, Status: st})
		}
	}
	return rs
}

// TestWrite: an all-pass run exits 0 and writes the summary (TOTAL labelled for strict only); a failing case exits 1;
// a run with one result dropped (internal/count) or without a required category is UNMEASURED: exit 2 and no file.
func TestWrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "s.json")
	var out bytes.Buffer
	all := results(t, runner.PASS)
	if code := Write(&out, summary.Summary{Engine: "e", Equiv: "full"}, "../../..", path, all); code != 0 {
		t.Fatalf("all pass: exit %d\n%s", code, out.String())
	}
	if b, err := os.ReadFile(path); err != nil || !strings.Contains(string(b), `"engine": "e"`) || strings.Contains(out.String(), "mapping cost") {
		t.Fatalf("summary %s %v, output %s", b, err, out.String())
	}
	out.Reset()
	if code := Write(&out, summary.Summary{Equiv: "strict"}, "../../..", path, all); code != 0 ||
		!strings.HasSuffix(out.String(), "/"+fmt.Sprint(len(all))+" strict=raw-encoding (mapping cost, not a quality score)\n") {
		t.Errorf("strict: exit %d, output %s", code, out.String())
	}
	if code := Write(&out, summary.Summary{}, "../../..", path, results(t, runner.FAIL)); code != 1 {
		t.Errorf("one FAIL: exit %d, want 1", code)
	}
	dropped := filepath.Join(dir, "d.json")
	if code := Write(&out, summary.Summary{}, "../../..", dropped, all[1:]); code != 2 || fileExists(dropped) ||
		!strings.Contains(out.String(), "UNMEASURED count") {
		t.Errorf("one result dropped: exit %d file %v, want 2 UNMEASURED count and no file", code, fileExists(dropped))
	}
	var noC0 []runner.Result
	for _, r := range all {
		if r.Category != "c0" {
			noC0 = append(noC0, r)
		}
	}
	missing := filepath.Join(dir, "m.json")
	if code := Write(&out, summary.Summary{}, "../../..", missing, noC0); code != 2 || fileExists(missing) {
		t.Errorf("missing category: exit %d file %v, want 2 and no file", code, fileExists(missing))
	}
}

func fileExists(p string) bool { _, err := os.Stat(p); return err == nil }
