// Package report turns the score summaries and the fact files into the E7 evidence lines and the E8 verdict.
package report

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/summary"
)

// Engines are the candidates in scoring order (rulings-2 #4).
var Engines = []string{"gitpod", "magnet", "xvt"}

// Run names one summary: engine, ref, equiv mode.
type Run struct{ Engine, Ref, Mode string }

// Runs are the 18 summaries the verdict reads (E6): strict and full on both refs, case and default on x53.
func Runs() []Run {
	var rs []Run
	for _, e := range Engines {
		for _, ref := range []string{"x53", "xcur"} {
			rs = append(rs, Run{e, ref, "strict"}, Run{e, ref, "full"})
		}
		rs = append(rs, Run{e, "x53", "case"}, Run{e, "x53", "default"})
	}
	return rs
}

// Load reads every run's summary <dir>/<engine>-<ref>-<mode>.json; a missing, unreadable or empty one is UNMEASURED,
// named by ref and mode first.
func Load(dir string) (map[Run]summary.Summary, error) {
	out := map[Run]summary.Summary{}
	for _, r := range Runs() {
		b, err := os.ReadFile(filepath.Join(dir, r.Engine+"-"+r.Ref+"-"+r.Mode+".json"))
		var s summary.Summary
		if err == nil {
			err = json.Unmarshal(b, &s)
		}
		if err != nil || s.Total == 0 {
			return nil, fmt.Errorf("UNMEASURED %s %s %s: %v", r.Ref, r.Mode, r.Engine, err)
		}
		out[r] = s
	}
	return out, nil
}
