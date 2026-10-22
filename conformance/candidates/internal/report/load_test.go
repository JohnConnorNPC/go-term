package report

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/summary"
)

// pins are the fixture's summary pins, each ending in its engine's sha12 of okPin.
var pins = map[string]string{"gitpod": "github.com/gitpod-io/xterm-go@v0.0.0-20260907130418-dae5128cb6b3",
	"magnet": "github.com/0magnet/xterm-go@v0.0.1-0.20261013223808-aa0a32cd702b",
	"xvt":    "github.com/charmbracelet/x/vt@v0.1.1-0.20261018003744-cb897cb4ed72"}

const okPin = "PIN gitpod pinned=dae5128cb6b3 head=dae5128cb6b3 SAME\nPIN magnet pinned=aa0a32cd702b head=aa0a32cd702b SAME\n" +
	"PIN xvt pinned=cb897cb4ed72 head=67fa755d101a SAME\n"

// fixture writes 314 empty goldens per ref under <dir>/golden and the 18 summaries, each its run's (engine, ref,
// mode, pin) at 30/314 in one category; edit changes any of them.
func fixture(t *testing.T, edit func(Run, *summary.Summary)) string {
	t.Helper()
	dir := t.TempDir()
	for _, ref := range []string{"x53", "xcur"} {
		d := filepath.Join(dir, "golden", ref, "c0")
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
		for i := 0; i < 314; i++ {
			if err := os.WriteFile(filepath.Join(d, fmt.Sprintf("%03d.json", i)), nil, 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}
	for _, r := range Runs() {
		s := summary.Summary{Engine: r.Engine, Ref: r.Ref, Equiv: r.Mode, Pin: pins[r.Engine], Total: 314, Passed: 30,
			Categories: map[string]summary.Cat{"c0": {Total: 314, Passed: 30}}}
		if edit != nil {
			edit(r, &s)
		}
		b, err := s.JSON()
		if err == nil {
			err = os.WriteFile(filepath.Join(dir, r.Engine+"-"+r.Ref+"-"+r.Mode+".json"), b, 0o644)
		}
		if err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

// load is Load over a fixture, its goldens in the same dir.
func load(t *testing.T, edit func(Run, *summary.Summary)) (map[Run]summary.Summary, error) {
	t.Helper()
	dir := fixture(t, edit)
	return Load(dir, dir)
}

// at applies f to run's summary only.
func at(run Run, f func(*summary.Summary)) func(Run, *summary.Summary) {
	return func(r Run, s *summary.Summary) {
		if r == run {
			f(s)
		}
	}
}

// TestLoad: the 18 runs load; an empty summary (0 cases) or a missing xcur or default one is UNMEASURED, named by
// ref and mode.
func TestLoad(t *testing.T) {
	if runs, err := load(t, nil); err != nil || len(runs) != 18 || len(Runs()) != 18 {
		t.Fatalf("%d runs, %v", len(runs), err)
	}
	empty := at(Run{"magnet", "x53", "case"}, func(s *summary.Summary) { s.Total, s.Passed = 0, 0 })
	if _, err := load(t, empty); err == nil || !strings.HasPrefix(err.Error(), "UNMEASURED x53 case magnet") {
		t.Errorf("an empty summary: %v, want UNMEASURED x53 case magnet", err)
	}
	for name, want := range map[string]string{"magnet-xcur-full.json": "UNMEASURED xcur full", "xvt-x53-default.json": "UNMEASURED x53 default"} {
		dir := fixture(t, nil)
		if err := os.Remove(filepath.Join(dir, name)); err != nil {
			t.Fatal(err)
		}
		if _, err := Load(dir, dir); err == nil || !strings.HasPrefix(err.Error(), want) {
			t.Errorf("without %s: %v, want %s", name, err, want)
		}
	}
}

// TestIdentity (M1): a summary that is not its run's is UNMEASURED: another ref, mode or engine (the oracle too), a
// total other than the goldens counted, passed outside 0..total, an unset pin, or pins that differ within an engine.
func TestIdentity(t *testing.T) {
	gx := Run{"gitpod", "x53", "full"}
	allXvt := func(pin string) func(Run, *summary.Summary) {
		return func(r Run, s *summary.Summary) {
			if r.Engine == "xvt" {
				s.Pin = pin
			}
		}
	}
	for name, c := range map[string]struct {
		edit func(Run, *summary.Summary)
		want string
	}{
		"case in the full slot":    {at(gx, func(s *summary.Summary) { s.Equiv = "case" }), "UNMEASURED x53 full gitpod"},
		"magnet in gitpod's slot":  {at(gx, func(s *summary.Summary) { s.Engine = "magnet" }), "UNMEASURED x53 full gitpod"},
		"the oracle":               {at(Run{"magnet", "x53", "strict"}, func(s *summary.Summary) { s.Engine, s.Pin = "oracle", "-" }), "UNMEASURED x53 strict magnet"},
		"total 313 of 314 goldens": {at(gx, func(s *summary.Summary) { s.Total = 313 }), "UNMEASURED x53 full gitpod"},
		"passed > total":           {at(gx, func(s *summary.Summary) { s.Passed = 315 }), "UNMEASURED x53 full gitpod"},
		"passed < 0":               {at(gx, func(s *summary.Summary) { s.Passed = -1 }), "UNMEASURED x53 full gitpod"},
		"pin - in every xvt run":   {allXvt("-"), "UNMEASURED x53 strict xvt"},
		"no pin in any xvt run":    {allXvt(""), "UNMEASURED x53 strict xvt"},
		"mixed pins":               {at(Run{"xvt", "xcur", "full"}, func(s *summary.Summary) { s.Pin = pins["gitpod"] }), "UNMEASURED xcur full xvt"},
	} {
		if _, err := load(t, c.edit); err == nil || !strings.HasPrefix(err.Error(), c.want) {
			t.Errorf("%s: %v, want %s", name, err, c.want)
		}
	}
	dir := fixture(t, nil)
	b, err := os.ReadFile(filepath.Join(dir, "gitpod-xcur-full.json"))
	if err == nil {
		err = os.WriteFile(filepath.Join(dir, "gitpod-x53-full.json"), b, 0o644)
	}
	if _, lerr := Load(dir, dir); err != nil || lerr == nil || !strings.HasPrefix(lerr.Error(), "UNMEASURED x53 full gitpod") {
		t.Errorf("the xcur file copied into the x53 slot: %v %v", err, lerr)
	}
	for rm, want := range map[string]string{"golden/xcur/c0/000.json": "UNMEASURED xcur strict gitpod", "golden": "UNMEASURED x53 strict gitpod"} {
		dir := fixture(t, nil)
		if err := os.RemoveAll(filepath.Join(dir, rm)); err != nil {
			t.Fatal(err)
		}
		if _, err := Load(dir, dir); err == nil || !strings.HasPrefix(err.Error(), want) {
			t.Errorf("without %s: %v, want %s", rm, err, want)
		}
	}
}
