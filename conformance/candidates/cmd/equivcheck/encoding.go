package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/norm"
)

// norms are the encoder's own colour-mode constants by encoding.json name (GO S-1: read from norm, never copied).
var norms = map[string]uint32{"CMMask": norm.CMMask, "CMP16": norm.CMP16, "CMP256": norm.CMP256, "CMRGB": norm.CMRGB}

// checkEncoding proves each entry's cited constant line (exactly once in its file) and that the line's last field,
// its hex value, equals the norm constant the entry names; it prints one ENC line per entry and returns the passes and the total.
func checkEncoding(z *zips, path string, out io.Writer) (ok, n int) {
	var f struct {
		Entries []struct {
			ID, Engine, Const string
			Cite              cite
		}
	}
	b, err := os.ReadFile(path)
	if err == nil {
		err = json.Unmarshal(b, &f)
	}
	if err != nil || len(f.Entries) == 0 {
		fmt.Fprintln(out, "ENC FAIL", path, err)
		return 0, 1
	}
	for _, e := range f.Entries {
		n++
		where, got, _ := z.measure(e.Cite)
		v, perr := strconv.ParseUint(e.Cite.Line[strings.LastIndex(e.Cite.Line, " ")+1:], 0, 32)
		want, known := norms[e.Const]
		if got != 1 || perr != nil || !known || uint32(v) != want {
			fmt.Fprintf(out, "ENC FAIL %s %s got=0x%x norm=0x%x lines=%d\n", e.Engine, e.ID, v, want, got)
			continue
		}
		fmt.Fprintf(out, "ENC %s %s @ %s 0x%x == norm.%s\n", e.Engine, e.ID, where, v, e.Const)
		ok++
	}
	return ok, n
}
