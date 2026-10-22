// Package gate is the E7 arithmetic and the E8 verdict rule, in integers only (no floats).
package gate

import "fmt"

// Pct is p/t in percent with two decimals, truncated (p*10000/t); UNMEASURED when t is 0.
func Pct(p, t int) string {
	if t <= 0 {
		return "UNMEASURED"
	}
	v := p * 10000 / t
	return fmt.Sprintf("%d.%02d", v/100, v%100)
}

// Pass is the 90 % bar, p*10 >= t*9; never for t == 0.
func Pass(p, t int) bool { return t > 0 && p*10 >= t*9 }

// Flag is " FLAG<50" for a category under 50 % (p*2 < t), else "".
func Flag(p, t int) string {
	if t > 0 && p*2 < t {
		return " FLAG<50"
	}
	return ""
}

// Gate is the E8 gate over one engine's x53 scores by equiv mode: the full table decides (rulings-6); strict, case
// and default are reported only.
func Gate(byMode map[string][2]int) bool {
	s, ok := byMode["full"]
	return ok && Pass(s[0], s[1])
}

// Verdict is the E8 rule: a pin that moved needs a rescore first (rulings-2 #2); else fork gitpod at sha12 iff its gate
// passes and every required licence line is OK; else the survey decides.
func Verdict(gatePass, licOK bool, pin, sha12 string) string {
	switch {
	case pin != "SAME":
		return "VERDICT rescore-required"
	case gatePass && licOK:
		return "VERDICT fork github.com/gitpod-io/xterm-go@" + sha12
	}
	return "VERDICT survey-decides"
}
