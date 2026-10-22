// Package lic checks a candidate module's licence (atoms 25-27) from its files, read as an fs.FS rooted at the
// module (the verified extraction of its pinned zip). Every finding is one LIC line.
package lic

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"strings"
)

// Line is `LIC <engine> <ID> <Status> <Value>`; Status is OK, FAIL or INFO (INFO never gates).
type Line struct{ ID, Status, Value string }

// mit are the three MIT phrases L2-L4, matched against the LICENSE with its whitespace normalised.
var mit = [][2]string{{"L2", "Permission is hereby granted, free of charge"},
	{"L3", "The above copyright notice and this permission notice shall be included"}, {"L4", `THE SOFTWARE IS PROVIDED "AS IS"`}}

// restrictions are riders that take rights away from MIT (L11), matched in the LICENSE lower-cased, whitespace
// normalised.
var restrictions = []string{"commons clause", "non-commercial", "noncommercial", "may not sell", "may not be sold",
	"shall not sell", "not be used for commercial"}

// xtermjs are the three copyright lines of the xterm.js LICENSE (gen/node_modules/xterm/LICENSE), each a NOTICE duty.
var xtermjs = []string{"Copyright (c) 2017-2019, The xterm.js authors (https://github.com/xtermjs/xterm.js)",
	"Copyright (c) 2014-2016, SourceLair Private Company (https://www.sourcelair.com)",
	"Copyright (c) 2012-2013, Christopher Jeffrey (https://github.com/chjj/)"}

// Root checks the module's root LICENSE: L1 present, L2-L4 the MIT phrases, L11 no known restriction added to MIT;
// L5 its copyright lines verbatim, L6 how many of the xterm.js lines it carries, L10 the sha256 of its bytes (12 hex),
// all INFO.
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
	var riders []string
	for _, r := range restrictions {
		if strings.Contains(strings.ToLower(text), r) {
			riders = append(riders, r)
		}
	}
	status := "OK"
	if len(riders) > 0 {
		status = "FAIL"
	}
	out = append(out, Line{"L11", status, "restrictions " + list(riders)})
	var copyrights []string
	carried := 0
	for _, l := range strings.Split(string(b), "\n") {
		if l = strings.TrimSpace(l); strings.HasPrefix(l, "Copyright") {
			copyrights = append(copyrights, l)
		}
		for _, x := range xtermjs {
			if l == x {
				carried++
			}
		}
	}
	sum := sha256.Sum256(b)
	return append(out, Line{"L5", "INFO", strings.Join(copyrights, " | ")}, Line{"L6", "INFO", fmt.Sprintf("%d/3", carried)},
		Line{"L10", "INFO", hex.EncodeToString(sum[:])[:12]})
}
