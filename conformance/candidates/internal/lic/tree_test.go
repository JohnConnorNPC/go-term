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

// TestTree: a GPL licence file, a GPL Go source or a GPL SPDX id fails L7; a GPL quote in a README or a test file is
// INFO only; the other licence files and the SPDX ids are listed.
func TestTree(t *testing.T) {
	base := func() fstest.MapFS {
		return fstest.MapFS{"LICENSE": {Data: []byte(mitText)}, "a.go": {Data: []byte("package a\n// SPDX-License-Identifier: MIT\n")}}
	}
	for name, c := range map[string]struct {
		file, body string
		want       []string
	}{
		"clean":       {"b.go", "package a\n", []string{"L7 INFO files none", "L7 OK copyleft none", "L8 INFO spdx MIT"}},
		"sub/COPYING": {"sub/COPYING", gpl, []string{"L7 INFO files sub/COPYING", "L7 FAIL copyleft sub/COPYING"}},
		"go source":   {"b.go", "package a\n// " + gpl, []string{"L7 FAIL copyleft b.go"}},
		"spdx gpl":    {"c.c", "/* SPDX-License-Identifier: GPL-2.0 */", []string{"L7 FAIL copyleft c.c", "L8 INFO spdx GPL-2.0,MIT"}},
		"readme":      {"README.md", "unlike the " + gpl, []string{"L7 OK copyleft none", "L7 INFO copyleft-quoted README.md"}},
		"test file":   {"b_test.go", "package a\n// " + gpl, []string{"L7 OK copyleft none", "L7 INFO copyleft-quoted b_test.go"}},
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
