// Package summary turns the runner's per-case results into the E6 score file: one class per case, the totals per
// category, and the split variants headless itself is sensitive to.
package summary

import (
	"regexp"
	"strings"

	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

// errPrefixes is the one closed list of runner Diff prefixes that mean the engine itself errored, each cited at
// 8a67250 in 2/binding.md: run.go's "step[%d]: engine error: %v", "step[%d]: snapshot error: %v", "engine: "+err and
// "B: "+err. internal/safe's "panic: …" and "timeout" arrive inside the first two.
var errPrefixes = []string{"engine error: ", "snapshot error: ", "engine: ", "B: "}

var stepPrefix = regexp.MustCompile(`^step\[[0-9]+\]: `)

// Classify is a non-STALE result's E6 class, PASS|FAIL|SPLIT|ERROR, and the variant a split diff names (ruling B).
// Only a SPLIT diff starts with "<vid>: " (the runner's v.id+": "); it is cut, then an optional "step[i]: ", and
// the rest is matched against errPrefixes: an engine error inside a split variant is an ERROR, not a split mismatch.
// "unobserved …" and every other diff are FAIL (SPLIT when a variant was cut); a SPLIT diff without its "<vid>: "
// names no variant, so it is a FAIL.
func Classify(r runner.Result) (class, variant string) {
	if r.Status == runner.PASS {
		return "PASS", ""
	}
	d := r.Diff
	if r.Status == runner.SPLIT {
		vid, rest, found := strings.Cut(d, ": ")
		if !found {
			return "FAIL", ""
		}
		variant, d = vid, rest
	}
	d = Rest(d)
	for _, p := range errPrefixes {
		if strings.HasPrefix(d, p) {
			return "ERROR", variant
		}
	}
	if variant != "" {
		return "SPLIT", variant
	}
	return "FAIL", ""
}

// Rest is a diff without its optional leading "step[i]: ", the form errPrefixes (and verdict's DIAG count) match.
func Rest(diff string) string { return stepPrefix.ReplaceAllString(diff, "") }
