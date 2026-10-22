package lic

import (
	"io/fs"
	"path"
	"regexp"
	"strings"
)

var (
	copyleft = regexp.MustCompile(`GNU GENERAL PUBLIC|GNU LESSER|Mozilla Public`)
)

// Tree walks every file of the module. L7 lists the licence files other than the root LICENSE (INFO) and FAILs on a
// copyleft text in a licence file or a non-test Go source; a copyleft quote anywhere else (a README) is INFO.
func Tree(fsys fs.FS) []Line {
	var others, hard, soft []string
	err := fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, err := fs.ReadFile(fsys, p)
		base := strings.ToUpper(path.Base(p))
		licFile := strings.HasPrefix(base, "LICENSE") || strings.HasPrefix(base, "COPYING") || strings.HasPrefix(base, "NOTICE")
		if licFile && p != "LICENSE" {
			others = append(others, p)
		}
		if src := strings.HasSuffix(p, ".go") && !strings.HasSuffix(p, "_test.go"); copyleft.Match(b) && (licFile || src) {
			hard = append(hard, p)
		} else if copyleft.Match(b) {
			soft = append(soft, p)
		}
		return err
	})
	if err != nil {
		return []Line{{"L7", "FAIL", "unreadable: " + err.Error()}}
	}
	status := "OK"
	if len(hard) > 0 {
		status = "FAIL"
	}
	return []Line{{"L7", "INFO", "files " + list(others)}, {"L7", status, "copyleft " + list(hard)},
		{"L7", "INFO", "copyleft-quoted " + list(soft)}}
}

func list(s []string) string {
	if len(s) == 0 {
		return "none"
	}
	return strings.Join(s, ",")
}
