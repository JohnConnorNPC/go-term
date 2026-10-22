package gate

import "testing"

func TestArithmetic(t *testing.T) {
	for _, c := range []struct {
		p, t int
		pass bool
		pct  string
		flag string
	}{{899, 1000, false, "89.90", ""}, {900, 1000, true, "90.00", ""}, {9, 10, true, "90.00", ""}, {0, 0, false, "UNMEASURED", ""},
		{2, 3, false, "66.66", ""}, {49, 100, false, "49.00", " FLAG<50"}, {50, 100, false, "50.00", ""}, {1, 2, false, "50.00", ""}} {
		if Pass(c.p, c.t) != c.pass || Pct(c.p, c.t) != c.pct || Flag(c.p, c.t) != c.flag {
			t.Errorf("%d/%d: pass %v pct %s flag %q, want %v %s %q", c.p, c.t, Pass(c.p, c.t), Pct(c.p, c.t), Flag(c.p, c.t), c.pass, c.pct, c.flag)
		}
	}
}

// TestRule: one vector per outcome; the gate reads the full score only.
func TestRule(t *testing.T) {
	fullPass := map[string][2]int{"full": {300, 314}, "strict": {0, 314}}
	fullFail := map[string][2]int{"full": {100, 314}, "strict": {314, 314}, "case": {314, 314}}
	for _, c := range []struct {
		scores map[string][2]int
		lic    bool
		pin    string
		want   string
	}{{fullPass, true, "SAME", "VERDICT fork github.com/gitpod-io/xterm-go@dae5128cb6b3"}, {fullPass, false, "SAME", "VERDICT survey-decides"},
		{fullPass, true, "MOVED", "VERDICT rescore-required"}, {fullFail, true, "SAME", "VERDICT survey-decides"},
		{map[string][2]int{}, true, "SAME", "VERDICT survey-decides"}} {
		if got := Verdict(Gate(c.scores), c.lic, c.pin, "dae5128cb6b3"); got != c.want {
			t.Errorf("%v lic=%v pin=%s: %s, want %s", c.scores, c.lic, c.pin, got, c.want)
		}
	}
}
