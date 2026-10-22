// Command equivcheck proves every citation of equiv/table.json, and the encoding constants of equiv/encoding.json,
// against the pinned module zips in the file proxy, and prints the E7 EQUIV/ENC evidence lines.
//
//	go run ./cmd/equivcheck   (under cand.sh, where the proxy is /cache/goproxy)
//
// Exit 1 on any failure.
package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if !check("/cache/goproxy", "go.mod", "equiv/table.json", "equiv/encoding.json", os.Stdout) {
		os.Exit(1)
	}
}

// check runs both proofs and prints the total line, EQUIV OK n/n or EQUIV FAIL k/n (rows + encoding entries).
func check(proxy, gomod, table, encoding string, out io.Writer) bool {
	z := &zips{proxy: proxy, gomod: gomod, mods: map[string]map[string][]string{}}
	ok, n := checkRows(z, table, out)
	encOK, encN := checkEncoding(z, encoding, out)
	ok, n = ok+encOK, n+encN
	verdict := "FAIL"
	if ok == n && n > 0 {
		verdict = "OK"
	}
	fmt.Fprintf(out, "EQUIV %s %d/%d\n", verdict, ok, n)
	return verdict == "OK"
}
