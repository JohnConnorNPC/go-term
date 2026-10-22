package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/equiv"
)

// checkRows proves every citation of every table row and prints one EQUIV line per row (its file:line and count
// results), preceded by an EQUIV FAIL line per failing citation; it returns the passing rows and the total. An
// unreadable table counts as one failing row.
func checkRows(z *zips, table string, out io.Writer) (ok, n int) {
	rows, err := equiv.Read(table)
	if err != nil {
		fmt.Fprintln(out, "EQUIV FAIL", err)
		return 0, 1
	}
	for _, r := range rows {
		var where []string
		good := true
		for _, raw := range r.Cite {
			var c cite
			err := json.Unmarshal(raw, &c)
			w, got, want := z.measure(c)
			if err != nil || got != want {
				fmt.Fprintf(out, "EQUIV FAIL %s %s got=%d want=%d %v\n", r.ID, raw, got, want, err)
				good = false
			}
			where = append(where, w)
		}
		fmt.Fprintf(out, "EQUIV %s %s %q => %v @ %s\n", r.Engine, r.ID, r.Key, r.Golden, strings.Join(where, ","))
		n++
		if good {
			ok++
		}
	}
	return ok, n
}
