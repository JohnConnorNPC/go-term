package report

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/summary"
)

const okLic = "LIC gitpod L1 OK x\nLIC gitpod L2 OK x\nLIC gitpod L3 OK x\nLIC gitpod L4 OK x\nLIC gitpod L11 OK x\nLIC gitpod L6 INFO 0/3\n"

// TestRule: fork needs the full gate, the licence lines and SAME; a licence FAIL, a missing L11 or a moved pin changes
// the verdict, and a missing PIN line is UNMEASURED.
func TestRule(t *testing.T) {
	pass := func(r Run, s *summary.Summary) { s.Passed = 290 }
	for _, c := range []struct {
		lic, pin, want string
	}{{okLic, okPin, "VERDICT fork github.com/gitpod-io/xterm-go@dae5128cb6b3"},
		{strings.Replace(okLic, "L2 OK", "L2 FAIL", 1), okPin, "VERDICT survey-decides"},
		{strings.Replace(okLic, "LIC gitpod L11 OK x\n", "", 1), okPin, "VERDICT survey-decides"},
		{okLic, strings.Replace(okPin, "head=dae5128cb6b3 SAME", "head=000000000000 MOVED", 1), "VERDICT rescore-required"},
		{okLic, "PIN magnet pinned=aa0a32cd702b head=aa0a32cd702b SAME\n", "UNMEASURED"}} {
		dir := fixture(t, pass)
		if os.WriteFile(filepath.Join(dir, "lic.txt"), []byte(c.lic), 0o644) != nil ||
			os.WriteFile(filepath.Join(dir, "pin.txt"), []byte(c.pin), 0o644) != nil {
			t.Fatal("fixture write failed")
		}
		runs, err := Load(dir, dir)
		got, err2 := Rule(runs, filepath.Join(dir, "lic.txt"), filepath.Join(dir, "pin.txt"))
		if err != nil || (got != c.want && (err2 == nil || !strings.HasPrefix(err2.Error(), c.want))) {
			t.Errorf("%q: %q %v %v, want %s", c.pin, got, err, err2, c.want)
		}
	}
}
