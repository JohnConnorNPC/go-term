package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/lic"
	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/modcache"
)

var require = regexp.MustCompile(`(?m)^\s*(?:require\s+)?(\S+\.\S+/\S+) v\S+`)

// check is the module's root, tree and dependency lines; test-only = the go.mod requires its packages don't import.
func check(gated bool, mod string) ([]lic.Line, error) {
	dir, err := modcache.Dir(mod)
	imports, err2 := modcache.Imports(mod + "/...")
	gomod, err3 := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil || err2 != nil || err3 != nil {
		return nil, fmt.Errorf("%s: %v %v %v", mod, err, err2, err3)
	}
	mods, used := map[string]fs.FS{}, map[string]bool{}
	for _, m := range imports {
		used[m] = true
		if m == mod {
			continue
		}
		d, err := modcache.Dir(m)
		if err != nil {
			return nil, err
		}
		mods[m] = os.DirFS(d)
	}
	var testOnly []string
	for _, m := range require.FindAllSubmatch(gomod, -1) {
		if !used[string(m[1])] {
			testOnly = append(testOnly, string(m[1]))
		}
	}
	fsys := os.DirFS(dir)
	return append(append(lic.Root(fsys), lic.Tree(fsys)...), lic.Deps(gated, mods, testOnly)...), nil
}
