package report

import (
	"fmt"
	"io"
	"sort"

	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/equiv"
	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/gate"
	"github.com/JohnConnorNPC/go-term/conformance/candidates/internal/summary"
)

// Scores prints SCORE and every CAT line of each run, in Runs order. Strict is a raw-encoding agreement (the mapping
// cost), not a quality score (go-v4 item 4): each strict SCORE line ends in equiv.StrictLabel; only GATE decides.
func Scores(out io.Writer, runs map[Run]summary.Summary) {
	for _, r := range Runs() {
		s := runs[r]
		fmt.Fprintf(out, "SCORE %s %s %s %d/%d %s%s\n", r.Engine, r.Ref, r.Mode, s.Passed, s.Total, gate.Pct(s.Passed, s.Total),
			equiv.Label(r.Mode))
		var cats []string
		for c := range s.Categories {
			cats = append(cats, c)
		}
		sort.Strings(cats)
		for _, c := range cats {
			v := s.Categories[c]
			fmt.Fprintf(out, "CAT %s %s %s %s %d/%d %s%s\n", r.Engine, r.Ref, r.Mode, c, v.Passed, v.Total,
				gate.Pct(v.Passed, v.Total), gate.Flag(v.Passed, v.Total))
		}
	}
}
