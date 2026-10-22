package report

import (
	"fmt"
	"os"
	"strings"

	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/gate"
	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/summary"
)

// required are the licence lines that must be OK for a fork (atom 25, L11 round 1; L5-L10 are INFO).
var required = []string{"L1", "L2", "L3", "L4", "L11"}

// Rule is the E8 verdict line for gitpod: the gate on x53 full, the licence lines of lic (every required id present
// and OK, no FAIL anywhere) and the pin line of pin (`PIN gitpod pinned=<sha12> head=<sha12> SAME|MOVED`). Every
// summary's pin must end in its engine's pinned sha12 (M3): one that doesn't was scored at another pin, so the verdict
// is rescore-required, as for a pin that MOVED. No gitpod LIC line, or no PIN line for an engine, is UNMEASURED.
func Rule(runs map[Run]summary.Summary, licPath, pinPath string) (string, error) {
	lic, err := os.ReadFile(licPath)
	pin, err2 := os.ReadFile(pinPath)
	status, licOK := map[string]string{}, true
	for _, l := range strings.Split(string(lic), "\n") {
		if f := strings.Fields(l); len(f) >= 4 && f[0] == "LIC" && f[1] == "gitpod" {
			status[f[2]] = f[3]
			licOK = licOK && f[3] != "FAIL"
		}
	}
	for _, id := range required {
		licOK = licOK && status[id] == "OK"
	}
	states, pinned := map[string]string{}, map[string]string{}
	for _, l := range strings.Split(string(pin), "\n") {
		if f := strings.Fields(l); len(f) == 5 && f[0] == "PIN" {
			states[f[1]], pinned[f[1]] = f[4], strings.TrimPrefix(f[2], "pinned=")
		}
	}
	var missing []string
	for _, e := range Engines {
		if states[e] == "" {
			missing = append(missing, e)
		}
	}
	if err != nil || err2 != nil || len(status) == 0 || len(missing) > 0 {
		return "", fmt.Errorf("UNMEASURED facts: lic %v (%d lines) pin %v, no PIN line for %v", err, len(status), err2, missing)
	}
	state := states["gitpod"]
	for r, s := range runs {
		if !strings.HasSuffix(s.Pin, "-"+pinned[r.Engine]) {
			state = "STALE " + r.Engine // scored at another pin: gate.Verdict gives rescore-required
		}
	}
	gitpod := map[string][2]int{}
	for _, m := range []string{"strict", "case", "default", "full"} {
		s := runs[Run{"gitpod", "x53", m}]
		gitpod[m] = [2]int{s.Passed, s.Total}
	}
	return gate.Verdict(gate.Gate(gitpod), licOK, state, pinned["gitpod"]), nil
}
