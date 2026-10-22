package report

import (
	"fmt"
	"os"
	"strings"

	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/gate"
	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/summary"
)

// required are the licence lines that must be OK for a fork (atom 25; L5-L10 are INFO).
var required = []string{"L1", "L2", "L3", "L4"}

// Rule is the E8 verdict line for gitpod: the gate on x53 full, the licence lines of lic (every required id present
// and OK, no FAIL anywhere) and the pin line of pin (`PIN gitpod pinned=<sha12> head=<sha12> SAME|MOVED`). No gitpod
// LIC or PIN line is UNMEASURED.
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
	state, sha := "", ""
	for _, l := range strings.Split(string(pin), "\n") {
		if f := strings.Fields(l); len(f) == 5 && f[0] == "PIN" && f[1] == "gitpod" {
			state, sha = f[4], strings.TrimPrefix(f[2], "pinned=")
		}
	}
	if err != nil || err2 != nil || len(status) == 0 || state == "" {
		return "", fmt.Errorf("UNMEASURED facts: lic %v (%d lines) pin %v %q", err, len(status), err2, state)
	}
	gitpod := map[string][2]int{}
	for _, m := range []string{"strict", "case", "default", "full"} {
		s := runs[Run{"gitpod", "x53", m}]
		gitpod[m] = [2]int{s.Passed, s.Total}
	}
	return gate.Verdict(gate.Gate(gitpod), licOK, state, sha), nil
}
