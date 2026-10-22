package modcache

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

// TestModcache: gitpod's directory holds its LICENSE; x/vt imports ultraviolet; gitpod imports no module (go-cmp is
// test-only).
func TestModcache(t *testing.T) {
	dir, err := Dir("github.com/gitpod-io/xterm-go")
	if _, serr := os.Stat(filepath.Join(dir, "LICENSE")); err != nil || serr != nil {
		t.Fatalf("gitpod dir %q: %v %v", dir, err, serr)
	}
	vt, err := Imports("github.com/charmbracelet/x/vt/...")
	if err != nil || !slices.Contains(vt, "github.com/charmbracelet/ultraviolet") {
		t.Errorf("x/vt imports %v %v", vt, err)
	}
	if g, err := Imports("github.com/gitpod-io/xterm-go/..."); err != nil || slices.Contains(g, "github.com/google/go-cmp") {
		t.Errorf("gitpod imports %v %v", g, err)
	}
}
