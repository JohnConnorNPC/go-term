package lic

import (
	"io/fs"
	"path"
	"regexp"
	"sort"
	"strings"
)

// copyleftNames are the copyleft licence names, matched case-insensitively on text whose whitespace and comment marks
// are collapsed to one space: GPL, LGPL and AGPL ("GNU [Lesser|Affero] General Public License"), MPL, EPL, EUPL,
// CDDL, SSPL.
var copyleftNames = []string{`general public licen[cs]e`, `mozilla public licen[cs]e`, `eclipse public licen[cs]e`,
	`european union public licen[cs]e`, `common development and distribution licen[cs]e`, `server side public licen[cs]e`}

var (
	copyleft     = regexp.MustCompile(`(?i)` + strings.Join(copyleftNames, "|"))
	marks        = regexp.MustCompile(`[\s/#*]+`)
	spdx         = regexp.MustCompile(`SPDX-License-Identifier:\s*([A-Za-z0-9.+-]+)`)
	copyleftSPDX = regexp.MustCompile(`(?i)^(GPL|LGPL|AGPL|MPL|EPL|EUPL|CDDL|SSPL)-`)
	docs         = regexp.MustCompile(`(?i)(^README|^CHANGELOG|^CHANGES|\.md|\.markdown|\.rst|\.txt|\.adoc|\.html)$`)
)

// Tree walks every file of the module. L7 lists the licence files other than the root LICENSE (INFO) and FAILs on a
// copyleft name in any file but documentation (a licence file, a source of any kind, tests included, a file with an
// SPDX header), or on a copyleft SPDX id anywhere; a copyleft name quoted in documentation (README*, *.md, *.txt, …)
// is INFO. L8 lists the distinct SPDX ids (INFO).
func Tree(fsys fs.FS) []Line {
	var others, hard, soft []string
	ids := map[string]bool{}
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
		header := false
		for _, m := range spdx.FindAllSubmatch(b, -1) {
			ids[string(m[1])], header = true, true
			if copyleftSPDX.Match(m[1]) {
				hard = append(hard, p)
			}
		}
		named := copyleft.Match(marks.ReplaceAll(b, []byte(" ")))
		if doc := docs.MatchString(path.Base(p)) && !licFile && !header; named && !doc {
			hard = append(hard, p)
		} else if named {
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
	var idList []string
	for id := range ids {
		idList = append(idList, id)
	}
	sort.Strings(idList)
	return []Line{{"L7", "INFO", "files " + list(others)}, {"L7", status, "copyleft " + list(hard)},
		{"L7", "INFO", "copyleft-quoted " + list(soft)}, {"L8", "INFO", "spdx " + list(idList)}}
}

func list(s []string) string {
	if len(s) == 0 {
		return "none"
	}
	return strings.Join(s, ",")
}
