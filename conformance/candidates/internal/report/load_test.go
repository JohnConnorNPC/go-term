package report

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/summary"
)

// fixture writes the 18 summaries to a temp dir: each 30/314 in one category, then edit changes any of them.
func fixture(t *testing.T, edit func(Run, *summary.Summary)) string {
	t.Helper()
	dir := t.TempDir()
	for _, r := range Runs() {
		s := summary.Summary{Engine: r.Engine, Ref: r.Ref, Equiv: r.Mode, Total: 314, Passed: 30,
			Categories: map[string]summary.Cat{"c0": {Total: 314, Passed: 30}}}
		if edit != nil {
			edit(r, &s)
		}
		b, err := s.JSON()
		if err == nil {
			err = os.WriteFile(filepath.Join(dir, r.Engine+"-"+r.Ref+"-"+r.Mode+".json"), b, 0o644)
		}
		if err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

// TestLoad: the 18 runs load; a missing xcur or default summary is UNMEASURED, named by ref and mode.
func TestLoad(t *testing.T) {
	dir := fixture(t, nil)
	if runs, err := Load(dir); err != nil || len(runs) != 18 || len(Runs()) != 18 {
		t.Fatalf("%d runs, %v", len(runs), err)
	}
	for name, want := range map[string]string{"magnet-xcur-full.json": "UNMEASURED xcur full", "xvt-x53-default.json": "UNMEASURED x53 default"} {
		dir := fixture(t, nil)
		if err := os.Remove(filepath.Join(dir, name)); err != nil {
			t.Fatal(err)
		}
		if _, err := Load(dir); err == nil || !strings.HasPrefix(err.Error(), want) {
			t.Errorf("without %s: %v, want %s", name, err, want)
		}
	}
}
