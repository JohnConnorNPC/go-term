package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/report"
	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/summary"
)

// scores writes the 18 summaries (300/314 each), lic.txt and pin.txt to a temp dir.
func scores(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	files := map[string]string{"lic.txt": "LIC gitpod L1 OK x\nLIC gitpod L2 OK x\nLIC gitpod L3 OK x\nLIC gitpod L4 OK x\nLIC gitpod L11 OK x\n",
		"pin.txt": "PIN gitpod pinned=dae5128cb6b3 head=dae5128cb6b3 SAME\n"}
	for _, r := range report.Runs() {
		b, err := summary.Summary{Total: 314, Passed: 300, Categories: map[string]summary.Cat{"c0": {Total: 314, Passed: 300}}}.JSON()
		if err != nil {
			t.Fatal(err)
		}
		files[r.Engine+"-"+r.Ref+"-"+r.Mode+".json"] = string(b)
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

// TestVerdict: the lines come in a fixed order and end with the verdict; a missing xcur summary is UNMEASURED,
// exit 2, no verdict.
func TestVerdict(t *testing.T) {
	dir := scores(t)
	var out bytes.Buffer
	if code := run([]string{"-scores", dir, "-diff", "../../../golden/diff.json"}, &out); code != 0 {
		t.Fatalf("exit %d:\n%s", code, out.String())
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	order := []string{"SCORE gitpod x53 strict", "CAT gitpod", "ABLATE gitpod", "GATE xvt", "SKIP-NOTE", "DIAG-NOTE", "XCUR-NOTE opts 3", "VERDICT fork"}
	i := 0
	for _, l := range lines {
		if i < len(order) && strings.HasPrefix(l, order[i]) {
			i++
		}
	}
	if i != len(order) || !strings.HasPrefix(lines[len(lines)-1], "VERDICT fork github.com/gitpod-io/xterm-go@dae5128cb6b3") {
		t.Errorf("order reached %d/%d:\n%s", i, len(order), out.String())
	}
	if err := os.Remove(filepath.Join(dir, "gitpod-xcur-full.json")); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	if code := run([]string{"-scores", dir, "-diff", "../../../golden/diff.json"}, &out); code != 2 ||
		!strings.HasPrefix(out.String(), "UNMEASURED xcur") || strings.Contains(out.String(), "VERDICT") {
		t.Errorf("missing xcur: exit %d:\n%s", code, out.String())
	}
}
