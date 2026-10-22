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

// pins are the candidates' pins as go.mod requires them; pin.txt's sha12s are their last 12 hex.
var pins = map[string]string{"gitpod": "github.com/gitpod-io/xterm-go@v0.0.0-20260907130418-dae5128cb6b3",
	"magnet": "github.com/0magnet/xterm-go@v0.0.1-0.20261013223808-aa0a32cd702b",
	"xvt":    "github.com/charmbracelet/x/vt@v0.1.1-0.20261018003744-cb897cb4ed72"}

// scores writes the 18 summaries (each its run's engine, ref, mode and pin; 95 % of the ref's goldens, counted under
// ../../../golden), lic.txt and pin.txt to a temp dir.
func scores(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	files := map[string]string{"lic.txt": "LIC gitpod L1 OK x\nLIC gitpod L2 OK x\nLIC gitpod L3 OK x\nLIC gitpod L4 OK x\nLIC gitpod L11 OK x\n",
		"pin.txt": "PIN gitpod pinned=dae5128cb6b3 head=dae5128cb6b3 SAME\nPIN magnet pinned=aa0a32cd702b head=aa0a32cd702b SAME\n" +
			"PIN xvt pinned=cb897cb4ed72 head=67fa755d101a SAME\n"}
	for _, r := range report.Runs() {
		gs, err := filepath.Glob("../../../golden/" + r.Ref + "/*/*.json")
		n := len(gs)
		b, jerr := summary.Summary{Engine: r.Engine, Ref: r.Ref, Equiv: r.Mode, Pin: pins[r.Engine], Total: n, Passed: n * 19 / 20,
			Categories: map[string]summary.Cat{"c0": {Total: n, Passed: n * 19 / 20}}}.JSON()
		if err != nil || jerr != nil || n == 0 {
			t.Fatal(n, err, jerr)
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
	if code := run([]string{"-scores", dir, "-root", "../../..", "-diff", "../../../golden/diff.json"}, &out); code != 0 {
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
	if code := run([]string{"-scores", dir, "-root", "../../..", "-diff", "../../../golden/diff.json"}, &out); code != 2 ||
		!strings.HasPrefix(out.String(), "UNMEASURED xcur") || strings.Contains(out.String(), "VERDICT") {
		t.Errorf("missing xcur: exit %d:\n%s", code, out.String())
	}
}

// TestForeign (M1): the xcur file copied into the x53 slot is UNMEASURED, exit 2, no verdict.
func TestForeign(t *testing.T) {
	dir := scores(t)
	b, err := os.ReadFile(filepath.Join(dir, "gitpod-xcur-full.json"))
	if err == nil {
		err = os.WriteFile(filepath.Join(dir, "gitpod-x53-full.json"), b, 0o644)
	}
	var out bytes.Buffer
	if code := run([]string{"-scores", dir, "-root", "../../..", "-diff", "../../../golden/diff.json"}, &out); err != nil || code != 2 ||
		!strings.HasPrefix(out.String(), "UNMEASURED x53 full gitpod") || strings.Contains(out.String(), "VERDICT") {
		t.Errorf("xcur in the x53 slot: %v exit %d:\n%s", err, code, out.String())
	}
}
