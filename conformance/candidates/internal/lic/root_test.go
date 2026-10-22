package lic

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"
	"testing/fstest"
)

const mitText = `MIT License

Copyright (c) 2026 Ona

Permission is hereby granted, free of charge, to any person obtaining a copy of this software.

The above copyright notice and this permission
notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND.
`

func status(lines []Line) map[string]string {
	m := map[string]string{}
	for _, l := range lines {
		m[l.ID] = l.Status + " " + l.Value
	}
	return m
}

// TestRootMIT: an MIT text passes L1-L4 (L3 split over two lines still matches); each missing part fails its line;
// no LICENSE fails L1.
func TestRootMIT(t *testing.T) {
	for name, c := range map[string]struct{ text, id, want string }{
		"mit":           {mitText, "L3", "OK"},
		"no permission": {strings.Replace(mitText, "Permission is hereby granted", "Use is allowed", 1), "L2", "FAIL"},
		"no condition":  {strings.Replace(mitText, "shall be included", "may be dropped", 1), "L3", "FAIL"},
		"no disclaimer": {strings.Replace(mitText, `PROVIDED "AS IS"`, "PROVIDED", 1), "L4", "FAIL"},
	} {
		got := status(Root(fstest.MapFS{"LICENSE": {Data: []byte(c.text)}}))
		if !strings.HasPrefix(got[c.id], c.want) || (name == "mit" && (got["L1"][:2] != "OK" || got["L2"][:2] != "OK" || got["L4"][:2] != "OK")) {
			t.Errorf("%s: %v, want %s %s", name, got, c.id, c.want)
		}
	}
	if got := status(Root(fstest.MapFS{"README": {Data: []byte("x")}})); !strings.HasPrefix(got["L1"], "FAIL") {
		t.Errorf("no LICENSE: %v", got)
	}
}

// TestRootInfo: the copyright lines verbatim, the xterm.js lines carried (0/3 gitpod-shaped, 3/3 xterm.js-shaped)
// and the sha256 of the bytes as they are (not normalised).
func TestRootInfo(t *testing.T) {
	got := status(Root(fstest.MapFS{"LICENSE": {Data: []byte(mitText)}}))
	sum := sha256.Sum256([]byte(mitText))
	if got["L5"] != "INFO Copyright (c) 2026 Ona" || got["L6"] != "INFO 0/3" || got["L10"] != "INFO "+hex.EncodeToString(sum[:])[:12] {
		t.Errorf("gitpod-shaped: %v", got)
	}
	xj := strings.Join(xtermjs, "\n") + "\n\n" + mitText
	if got := status(Root(fstest.MapFS{"LICENSE": {Data: []byte(xj)}})); got["L6"] != "INFO 3/3" {
		t.Errorf("xterm.js-shaped: %v", got)
	}
}

// TestRootRestrictions (L11): MIT with a rider that takes rights away fails, in any case and across a line break;
// plain MIT passes with no restriction named.
func TestRootRestrictions(t *testing.T) {
	for rider, want := range map[string]string{"": "OK restrictions none",
		"\"Commons Clause\" License Condition v1.0": "FAIL restrictions commons clause",
		"For NON-COMMERCIAL use only.":              "FAIL restrictions non-commercial",
		"You may not\nsell the Software.":           "FAIL restrictions may not sell",
		"Noncommercial":                             "FAIL restrictions noncommercial",
		"It may not be sold.":                       "FAIL restrictions may not be sold",
		"You shall not sell it.":                    "FAIL restrictions shall not sell",
		"It may not be used for commercial gain.":   "FAIL restrictions not be used for commercial"} {
		if got := status(Root(fstest.MapFS{"LICENSE": {Data: []byte(mitText + rider)}})); got["L11"] != want || got["L2"][:2] != "OK" {
			t.Errorf("MIT + %q: L11 %q, want %q", rider, got["L11"], want)
		}
	}
}
