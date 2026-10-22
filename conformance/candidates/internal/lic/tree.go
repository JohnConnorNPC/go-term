package lic

import (
	"io/fs"
	"path"
	"strings"
)

// Tree walks every file of the module. L7 lists the licence files other than the root LICENSE (INFO).
func Tree(fsys fs.FS) []Line {
	var others []string
	err := fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		base := strings.ToUpper(path.Base(p))
		licFile := strings.HasPrefix(base, "LICENSE") || strings.HasPrefix(base, "COPYING") || strings.HasPrefix(base, "NOTICE")
		if licFile && p != "LICENSE" {
			others = append(others, p)
		}
		return nil
	})
	if err != nil {
		return []Line{{"L7", "FAIL", "unreadable: " + err.Error()}}
	}
	return []Line{{"L7", "INFO", "files " + list(others)}}
}

func list(s []string) string {
	if len(s) == 0 {
		return "none"
	}
	return strings.Join(s, ",")
}
