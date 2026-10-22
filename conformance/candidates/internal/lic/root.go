// Package lic checks a candidate module's licence (atoms 25-27) from its files, read as an fs.FS rooted at the
// module (the verified extraction of its pinned zip). Every finding is one LIC line.
package lic

import (
	"io/fs"
	"strings"
)

// Line is `LIC <engine> <ID> <Status> <Value>`; Status is OK, FAIL or INFO (INFO never gates).
type Line struct{ ID, Status, Value string }

// mit are the three MIT phrases L2-L4, matched against the LICENSE with its whitespace normalised.
var mit = [][2]string{{"L2", "Permission is hereby granted, free of charge"},
	{"L3", "The above copyright notice and this permission notice shall be included"}, {"L4", `THE SOFTWARE IS PROVIDED "AS IS"`}}

// Root checks the module's root LICENSE: L1 present, L2-L4 the MIT phrases.
func Root(fsys fs.FS) []Line {
	b, err := fs.ReadFile(fsys, "LICENSE")
	if err != nil {
		return []Line{{"L1", "FAIL", "no root LICENSE: " + err.Error()}}
	}
	text := strings.Join(strings.Fields(string(b)), " ")
	out := []Line{{"L1", "OK", "LICENSE"}}
	for _, p := range mit {
		status := "FAIL"
		if strings.Contains(text, p[1]) {
			status = "OK"
		}
		out = append(out, Line{p[0], status, p[1]})
	}
	return out
}
