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

// TestTree: a GPL licence file, a GPL Go source (tests too) or a GPL SPDX id fails L7; a GPL quote in a README is INFO
// only; the other licence files and the SPDX ids are listed.
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
		"test file":   {"b_test.go", "package a\n// " + gpl, []string{"L7 FAIL copyleft b_test.go"}},
		"wrapped":     {"b.go", "// under the GNU General Public\n// License, version 2", []string{"L7 FAIL copyleft b.go"}},
		"notes.txt":   {"docs/notes.txt", "unlike the " + gpl, []string{"L7 OK copyleft none", "L7 INFO copyleft-quoted docs/notes.txt"}},
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

// TestFamilies: every copyleft family fails L7 by its name in any case, in a vendored licence file or a source
// header, and by its SPDX id in any case; the same name in a README is only quoted.
func TestFamilies(t *testing.T) {
	names := map[string]string{"GPL": "GNU General Public License", "LGPL": "GNU LESSER GENERAL PUBLIC LICENSE",
		"AGPL": "gnu affero general public license", "MPL": "Mozilla Public License", "EPL": "ECLIPSE PUBLIC LICENSE",
		"EUPL": "European Union Public Licence", "CDDL": "Common Development and Distribution License",
		"SSPL": "server side public license"}
	for id, name := range names {
		for file, want := range map[string]string{"third_party/x/LICENSE": "L7 FAIL copyleft third_party/x/LICENSE",
			"x.c": "L7 FAIL copyleft x.c", "README": "L7 INFO copyleft-quoted README"} {
			if got := lines(fstest.MapFS{file: {Data: []byte("/*\n * Licensed under the " + name + ".\n */")}}); !strings.Contains(got+"\n", want+"\n") {
				t.Errorf("%s in %s: no %q in\n%s", id, file, want, got)
			}
		}
		spdxID := strings.ToLower(id) + "-2.0"
		if got := lines(fstest.MapFS{"y.go": {Data: []byte("// SPDX-License-Identifier: " + spdxID)}}); !strings.Contains(got, "L7 FAIL copyleft y.go") {
			t.Errorf("SPDX %s: %s", spdxID, got)
		}
	}
}
