// Package modcache asks the go command (run inside cand.sh) where a pinned module was extracted and verified, and
// which modules a package pattern imports outside its tests.
package modcache

import (
	"os/exec"
	"strings"
)

// Dir is the module cache directory of mod at the version go.mod requires.
func Dir(mod string) (string, error) { return golist("-m", "-f", "{{.Dir}}", mod) }

// Imports are the distinct modules the non-test packages matching pattern import (std excluded), in go list order.
func Imports(pattern string) ([]string, error) {
	out, err := golist("-deps", "-f", "{{with .Module}}{{.Path}}{{end}}", pattern)
	seen, mods := map[string]bool{}, []string{}
	for _, m := range strings.Fields(out) {
		if !seen[m] {
			seen[m] = true
			mods = append(mods, m)
		}
	}
	return mods, err
}

func golist(args ...string) (string, error) {
	out, err := exec.Command("go", append([]string{"list"}, args...)...).Output()
	return strings.TrimSpace(string(out)), err
}
