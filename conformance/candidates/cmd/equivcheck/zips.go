package main

import (
	"archive/zip"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// zips reads each cited module's pinned zip (the version is its require line in go.mod, block or one-line form)
// once, as file → lines; mods must be non-nil.
type zips struct {
	proxy, gomod string
	mods         map[string]map[string][]string
}

func (z *zips) module(mod string) (map[string][]string, error) {
	if files, ok := z.mods[mod]; ok {
		return files, nil
	}
	gm, err := os.ReadFile(z.gomod)
	m := regexp.MustCompile(`(?m)^\s*(?:require\s+)?` + regexp.QuoteMeta(mod) + `\s+(\S+)`).FindSubmatch(gm)
	if err != nil || m == nil || strings.ToLower(mod) != mod {
		return nil, fmt.Errorf("no lower-case require line for %s in %s (%v)", mod, z.gomod, err)
	}
	r, err := zip.OpenReader(filepath.Join(z.proxy, mod, "@v", string(m[1])+".zip"))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	files := map[string][]string{}
	for _, f := range r.File {
		b, err := fs.ReadFile(r, f.Name)
		if err != nil {
			return nil, err
		}
		files[strings.TrimPrefix(f.Name, mod+"@"+string(m[1])+"/")] = strings.Split(string(b), "\n")
	}
	z.mods[mod] = files
	return files, nil
}
