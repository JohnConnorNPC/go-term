package report

import (
	"fmt"
	"io"

	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/gate"
	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/summary"
)

// PerEngine prints, for each engine, ABLATE (x53, the four modes), EQUIV-DELTA per ref (full - strict), the
// per-ref split lines (Splits), then GATE: x53 full, the only line that decides (rulings-6).
func PerEngine(out io.Writer, runs map[Run]summary.Summary) {
	for _, e := range Engines {
		p := func(ref, mode string) int { return runs[Run{e, ref, mode}].Passed }
		fmt.Fprintf(out, "ABLATE %s x53 strict=%d case=%d default=%d full=%d\n", e, p("x53", "strict"), p("x53", "case"),
			p("x53", "default"), p("x53", "full"))
		for _, ref := range []string{"x53", "xcur"} {
			fmt.Fprintf(out, "EQUIV-DELTA %s %s %d\n", e, ref, p(ref, "full")-p(ref, "strict"))
			Splits(out, e, ref, runs[Run{e, ref, "full"}])
		}
		x := runs[Run{e, "x53", "full"}]
		verdict := "FAIL"
		if gate.Pass(x.Passed, x.Total) {
			verdict = "PASS"
		}
		fmt.Fprintf(out, "GATE %s x53 %d/%d >=90%% %s\n", e, x.Passed, x.Total, verdict)
	}
}
