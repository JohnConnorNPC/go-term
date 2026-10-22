package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

// blind reports nothing it can see: every snapshot lists lineFeeds as unobserved.
type blind struct{}

func (blind) Write([]byte) error    { return nil }
func (blind) Resize(int, int) error { return nil }
func (blind) Reset() error          { return nil }
func (blind) Snapshot() (runner.State, error) {
	return runner.State{Unobserved: []string{"lineFeeds"}}, nil
}

// tree is a temp conformance root: links to the real cases, probes and gen/, and a copy of the goldens, which edit
// changes.
func tree(t *testing.T, edit func(root string)) string {
	t.Helper()
	root, src := t.TempDir(), "../../.."
	err := os.CopyFS(filepath.Join(root, "golden"), os.DirFS(filepath.Join(src, "golden")))
	for _, p := range []string{"cases", "probes.json", "gen"} {
		abs, aerr := filepath.Abs(filepath.Join(src, p))
		if err == nil {
			err = aerr
		}
		if err == nil {
			err = os.Symlink(abs, filepath.Join(root, p))
		}
	}
	if err != nil {
		t.Fatal(err)
	}
	edit(root)
	return root
}

func score(t *testing.T, root, engine string) (code int, out string, wrote bool) {
	var b bytes.Buffer
	path := filepath.Join(t.TempDir(), "s.json")
	code = run([]string{"-engine", engine, "-root", root, "-table", "../../equiv/table.json", "-out", path}, &b)
	_, err := os.Stat(path)
	return code, b.String(), err == nil
}

// TestUnobservedLine (ruling A): an unobserved path is a FAIL line, never UNMEASURED; the summary is written.
func TestUnobservedLine(t *testing.T) {
	runner.Register("blind", func(runner.Config) (runner.Engine, error) { return blind{}, nil })
	code, out, wrote := score(t, "../../..", "blind")
	if code != 1 || !wrote || strings.Contains(out, "UNMEASURED") || !strings.Contains(out, "FAIL c0/smoke-text unobserved step[0].P.lineFeeds") {
		t.Errorf("blind: exit %d wrote %v:\n%.400s", code, wrote, out)
	}
}

// TestUnmeasured (ruling A): a missing golden or a STALE one is UNMEASURED (exit 2, no file), never "unobserved".
func TestUnmeasured(t *testing.T) {
	for name, edit := range map[string]func(string){
		"missing golden": func(r string) {
			if err := os.Remove(filepath.Join(r, "golden/x53/c0/smoke-text.json")); err != nil {
				t.Fatal(err)
			}
		},
		"stale golden": func(r string) { stale(t, filepath.Join(r, "golden/x53/c0/smoke-text.json")) },
	} {
		code, out, wrote := score(t, tree(t, edit), "stub")
		if code != 2 || wrote || !strings.Contains(out, "UNMEASURED") || strings.Contains(out, "unobserved") {
			t.Errorf("%s: exit %d wrote %v:\n%.400s", name, code, wrote, out)
		}
	}
}

// stale changes one hex digit of the golden's first `in`, so its ops no longer hash to inSha.
func stale(t *testing.T, p string) {
	b, err := os.ReadFile(p)
	i := bytes.Index(b, []byte(`"in":"`))
	if err != nil || i < 0 {
		t.Fatalf("%s: %v", p, err)
	}
	b[i+6] ^= 1
	if err := os.WriteFile(p, b, 0o644); err != nil {
		t.Fatal(err)
	}
}
