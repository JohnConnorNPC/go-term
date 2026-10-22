package report

import (
	"os"
	"path/filepath"
	"testing"
)

// TestXcurNote: only opts.* paths count; on the real diff.json the three cases of round-atoms-2 ruling 5.
func TestXcurNote(t *testing.T) {
	p := filepath.Join(t.TempDir(), "diff.json")
	body := `{"id":"diff","cases":[{"id":"a/x","P":["steps.0.P.opts.cursorStyle"],"B":[],"I":[]},{"id":"b/y","P":["steps.0.P.modes.insertMode"],"B":[],"I":[]}]}`
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	for path, want := range map[string]string{p: "XCUR-NOTE opts 1 a/x",
		"../../../golden/diff.json": "XCUR-NOTE opts 3 dcs/decrqss-sp-q decmodes/decscusr-bar decmodes/decscusr-reset"} {
		if got, err := XcurNote(path); err != nil || got != want {
			t.Errorf("%s: %q %v, want %q", path, got, err, want)
		}
	}
	if _, err := XcurNote("nope.json"); err == nil {
		t.Error("a missing diff.json is not UNMEASURED")
	}
}
