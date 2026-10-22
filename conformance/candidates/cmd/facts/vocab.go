package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// vocabFields are the golden fields whose spellings the equivalence table maps to (plan-pass-2 S-a).
var vocabFields = []string{"mouseTrackingMode", "cursorStyle"}

// goldenVocab is one `GOLDEN-VOCAB <field> <values> (<source>)` INFO line per field: the union type of that field in
// the reference's typings (xterm-headless@5.3.0's xterm-headless.d.ts, the gate reference).
func goldenVocab(typings string) ([]string, error) {
	b, err := os.ReadFile(typings)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, f := range vocabFields {
		m := regexp.MustCompile(`(?m)^\s*(?:readonly )?` + f + `\??: ('[a-z0-9]+'(?: \| '[a-z0-9]+')*);`).FindSubmatch(b)
		if m == nil {
			return nil, fmt.Errorf("%s: no union type for %s", typings, f)
		}
		values := strings.ReplaceAll(strings.ReplaceAll(string(m[1]), "'", ""), " | ", "|")
		out = append(out, fmt.Sprintf("GOLDEN-VOCAB %s %s (%s)", f, values, typings))
	}
	return out, nil
}
