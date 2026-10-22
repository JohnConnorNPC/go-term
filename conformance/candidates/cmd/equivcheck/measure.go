package main

import (
	"fmt"
	"regexp"
	"strings"
)

// cite is one citation (equiv.Row.Cite): a line {module, file, line}, or a count {module, files, pattern, count}.
type cite struct {
	Module, File, Line, Files, Pattern string
	Count                              *int
}

// measure returns where c points and its got/want counts. A line must match exactly one whole line of its file (both
// sides trimmed); a count is the number of lines matching its pattern over the module's non-test .go files, or over
// its one file. An unreadable zip, file or pattern gives got=-1: unmeasured, never a 0 that could equal want.
func (z *zips) measure(c cite) (where string, got, want int) {
	files, err := z.module(c.Module)
	if err != nil {
		return "unreadable: " + err.Error(), -1, 1
	}
	if c.Count == nil {
		got, where = 0, c.File+":?"
		for i, l := range files[c.File] {
			if strings.TrimSpace(l) == strings.TrimSpace(c.Line) {
				got++
				where = fmt.Sprintf("%s:%d", c.File, i+1)
			}
		}
		return where, got, 1
	}
	re, err := regexp.Compile(c.Pattern)
	if err != nil || (c.Files != "*" && files[c.Files] == nil) {
		return fmt.Sprintf("unreadable: %s %v", c.Files, err), -1, *c.Count
	}
	for name, lines := range files {
		if (c.Files == "*" && strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go")) || name == c.Files {
			for _, l := range lines {
				if re.MatchString(l) {
					got++
				}
			}
		}
	}
	return fmt.Sprintf("count(%s)=%d", c.Pattern, got), got, *c.Count
}
