package main

import (
	"io/fs"
	"regexp"
	"sort"
	"strings"
)

var getter = regexp.MustCompile(`^func \(\w+ \*?\w+\) (DecPrivateModes|Modes|CoreService)\(\)|^\s+(DecPrivateModes)\s+DecPrivateModes$`)

// modeAPI is `MODEAPI <e> get=<symbols> replay=yes`: the exported mode getters (a method or an exported
// DecPrivateModes field) in the module's non-test Go files, listed as found. replay (re-writing the bytes) works for
// every engine.
func modeAPI(engine string, fsys fs.FS) (string, error) {
	gets := map[string]bool{}
	err := fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(p, ".go") || strings.HasSuffix(p, "_test.go") {
			return err
		}
		b, err := fs.ReadFile(fsys, p)
		for _, l := range strings.Split(string(b), "\n") {
			if m := getter.FindStringSubmatch(l); m != nil {
				gets[m[1]+m[2]] = true
			}
		}
		return err
	})
	return "MODEAPI " + engine + " get=" + names(gets) + " replay=yes", err
}

func names(set map[string]bool) string {
	var s []string
	for k := range set {
		s = append(s, k)
	}
	sort.Strings(s)
	if len(s) == 0 {
		return "none"
	}
	return strings.Join(s, ",")
}
