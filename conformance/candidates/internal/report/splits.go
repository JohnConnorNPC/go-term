package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/summary"
)

// Splits prints one full run's SKIP line (complete skips; skips of non-PASS cases are a lower bound, rulings-5 #5),
// its SPLITS line (whole-pass = cases whose whole-write run passed, i.e. PASS plus every case failing in a split
// variant) and its DIAG line (rulings-8: diagnostic, never gates).
func Splits(out io.Writer, engine, ref string, full summary.Summary) {
	complete, lower := 0, 0
	for _, sk := range full.Skipped {
		if sk.Complete {
			complete++
		} else {
			lower++
		}
	}
	fmt.Fprintf(out, "SKIP %s %s pass=%d nonpass>=%d\n", engine, ref, complete, lower)
	whole, split, first := full.Passed, 0, 0
	for _, c := range full.Cases {
		if c.Variant != "" {
			whole++
		}
		if c.Result == "SPLIT" {
			split++
		}
		if c.Result == "FAIL" && strings.HasPrefix(summary.Rest(c.Detail), "unobserved ") {
			first++
		}
	}
	fmt.Fprintf(out, "SPLITS %s %s whole-pass=%d split-fail=%d\n", engine, ref, whole, split)
	fmt.Fprintf(out, "DIAG %s %s unobserved-first %d/%d\n", engine, ref, first, full.Total)
}
