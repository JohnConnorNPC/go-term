package runner

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// wideRoot is a root holding only the committed wide category (its cases, xcur goldens, the manifest, pins and
// probes, byte copies): it has headless-sensitive variants and misses every other required category.
func wideRoot(t *testing.T) string {
	t.Helper()
	dst := t.TempDir()
	cp := func(rel string) {
		b, err := os.ReadFile(filepath.Join("..", rel))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Dir(filepath.Join(dst, rel)), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dst, rel), b, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for _, f := range []string{"probes.json", "gen/pins.json", "golden/manifest.json"} {
		cp(f)
	}
	ms, _ := filepath.Glob(filepath.Join("..", "golden", "xcur", "wide", "*.json"))
	if len(ms) == 0 {
		t.Fatal("no xcur wide goldens (unmeasured)")
	}
	for _, m := range ms {
		name := filepath.Base(m)
		cp(filepath.Join("golden", "xcur", "wide", name))
		cp(filepath.Join("cases", "wide", name))
	}
	return dst
}

// TestCLIFuzz: on a wide-only root the default run skips the headless-sensitive variants and exits 2 (the other
// required categories are UNMEASURED); -fuzz scores them (hsens equals the default's skipped=, no SKIP or CAT line)
// and exits 0 for the oracle, which is split-safe; a failing engine exits 1.
func TestCLIFuzz(t *testing.T) {
	root := wideRoot(t)
	run := func(args ...string) (string, int) {
		var out bytes.Buffer
		code := CLI(append([]string{"-root", root}, args...), &out, nil)
		return out.String(), code
	}
	def, code := run("-engine", "oracle")
	m := regexp.MustCompile(`(?m)^TOTAL (\d+)/(\d+) PASS skipped=(\d+)$`).FindStringSubmatch(def)
	if code != 2 || m == nil || m[1] != m[2] || m[3] == "0" || !strings.Contains(def, " UNMEASURED") {
		t.Fatalf("default oracle: exit %d, last lines %q, want exit 2 with UNMEASURED categories", code, tail(def))
	}
	fz, code := run("-engine", "oracle", "-fuzz")
	want := "TOTAL " + m[1] + "/" + m[2] + " PASS hsens=" + m[3]
	if code != 0 || !strings.HasSuffix(strings.TrimSpace(fz), want) || strings.Contains(fz, "\nSKIP ") || strings.Contains(fz, "\nCAT ") {
		t.Fatalf("-fuzz oracle: exit %d, last lines %q, want %q and no SKIP/CAT lines", code, tail(fz), want)
	}
	st, code := run("-engine", "stub", "-fuzz", "-splits=false")
	if code != 1 || !regexp.MustCompile(`TOTAL 0/\d+ FAIL hsens=0\n$`).MatchString(st) {
		t.Fatalf("-fuzz stub: exit %d, last lines %q, want exit 1 and TOTAL 0/n FAIL", code, tail(st))
	}
	t.Logf("CLI -fuzz wide: %s (default: %s, exit 2)", want, m[0])
}

func tail(s string) string {
	l := strings.Split(strings.TrimSpace(s), "\n")
	return strings.Join(l[max(0, len(l)-3):], " | ")
}
