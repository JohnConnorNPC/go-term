package lic

import (
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
