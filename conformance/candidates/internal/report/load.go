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

// Load reads every run's summary <dir>/<engine>-<ref>-<mode>.json and checks that it is that run's (M1): its engine,
// ref and equiv mode are the file's, its total is the number of ref's goldens under root/golden (counted), 0 <= passed
// <= total, and its pin is set and the same in all of the engine's summaries. A missing, unreadable, empty or foreign
// summary is UNMEASURED, named by ref and mode first.
func Load(dir, root string) (map[Run]summary.Summary, error) {
	out, pins := map[Run]summary.Summary{}, map[string]string{}
	for _, r := range Runs() {
		b, err := os.ReadFile(filepath.Join(dir, r.Engine+"-"+r.Ref+"-"+r.Mode+".json"))
		var s summary.Summary
		if err == nil {
			err = json.Unmarshal(b, &s)
		}
		n, gerr := goldens(root, r.Ref)
		if pins[r.Engine] == "" {
			pins[r.Engine] = s.Pin
		}
		switch {
		case err != nil || gerr != nil || s.Total == 0:
			return nil, fmt.Errorf("UNMEASURED %s %s %s: %v %v", r.Ref, r.Mode, r.Engine, err, gerr)
		case s.Engine != r.Engine || s.Ref != r.Ref || s.Equiv != r.Mode:
			return nil, fmt.Errorf("UNMEASURED %s %s %s: the summary is %s %s %s", r.Ref, r.Mode, r.Engine, s.Ref, s.Equiv, s.Engine)
		case s.Total != n || s.Passed < 0 || s.Passed > s.Total:
			return nil, fmt.Errorf("UNMEASURED %s %s %s: passed %d of %d, %d goldens", r.Ref, r.Mode, r.Engine, s.Passed, s.Total, n)
		case s.Pin == "" || s.Pin == "-" || s.Pin != pins[r.Engine]:
			return nil, fmt.Errorf("UNMEASURED %s %s %s: pin %q, the engine's first %q", r.Ref, r.Mode, r.Engine, s.Pin, pins[r.Engine])
		}
		out[r] = s
	}
	return out, nil
}

// goldens counts ref's golden files, root/golden/<ref>/<category>/*.json; none is an error.
func goldens(root, ref string) (int, error) {
	gs, err := filepath.Glob(filepath.Join(root, "golden", ref, "*", "*.json"))
	if err == nil && len(gs) == 0 {
		err = fmt.Errorf("no goldens under %s", filepath.Join(root, "golden", ref))
	}
	return len(gs), err
}
