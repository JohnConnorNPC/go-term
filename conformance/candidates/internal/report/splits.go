package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/summary"
)

// Splits prints one full run's SKIP line (complete skips; skips of non-PASS cases are a lower bound, rulings-5 #5),
// its SPLITS line (whole-pass = cases whose whole-write run passed, i.e. PASS plus every case failing in a split
// variant), its DIAG line (rulings-8: diagnostic, never gates) and its STATUS line, the cases per class (round 1 S3).
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
	whole, split, first, class := full.Passed, 0, 0, map[string]int{}
	for _, c := range full.Cases {
		class[c.Result]++
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
	fmt.Fprintf(out, "STATUS %s %s full PASS=%d FAIL=%d SPLIT=%d ERROR=%d of %d\n", engine, ref, class["PASS"], class["FAIL"],
		class["SPLIT"], class["ERROR"], full.Total)
}
