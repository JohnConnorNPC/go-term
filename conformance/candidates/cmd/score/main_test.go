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
