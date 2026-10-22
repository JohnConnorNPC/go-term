package lic

import (
	"strings"
	"testing"
	"testing/fstest"
)

const gpl = "GNU GENERAL PUBLIC LICENSE Version 3"

func lines(fsys fstest.MapFS) string {
	var out []string
	for _, l := range Tree(fsys) {
		out = append(out, l.ID+" "+l.Status+" "+l.Value)
	}
	return strings.Join(out, "\n")
}

// TestTree: the licence files other than the root LICENSE are listed.
func TestTree(t *testing.T) {
	base := func() fstest.MapFS {
		return fstest.MapFS{"LICENSE": {Data: []byte(mitText)}, "a.go": {Data: []byte("package a\n// SPDX-License-Identifier: MIT\n")}}
	}
	for name, c := range map[string]struct {
		file, body string
		want       []string
	}{
		"clean":       {"b.go", "package a\n", []string{"L7 INFO files none"}},
		"sub/COPYING": {"sub/COPYING", gpl, []string{"L7 INFO files sub/COPYING"}},
	} {
		fsys := base()
		fsys[c.file] = &fstest.MapFile{Data: []byte(c.body)}
		got := lines(fsys)
		for _, w := range c.want {
			if !strings.Contains(got+"\n", w+"\n") {
				t.Errorf("%s: no %q in\n%s", name, w, got)
			}
		}
	}
}
