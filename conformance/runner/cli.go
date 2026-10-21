package runner

import (
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"
)

// CLI is cmd/conformance's body: it registers factories (extra named engines), parses args and scores one engine, writing
// its lines to stdout. It returns the exit status: 0 every case passes (or the counts meet the minimums), 1 a case
// fails, 2 something is unmeasured.
func CLI(args []string, stdout io.Writer, factories map[string]Factory) int {
	for name, f := range factories {
		Register(name, f)
	}
	fs := flag.NewFlagSet("conformance", flag.ContinueOnError)
	root := fs.String("root", "conformance", "the conformance directory (cases/, golden/, probes.json)")
	engine := fs.String("engine", "stub", "engine to score: "+strings.Join(append(Engines(), "oracle"), ", "))
	ref := fs.String("ref", "x53", "golden set: x53 (xterm-headless@5.3.0) or xcur (@xterm/headless latest)")
	splits := fs.Bool("splits", true, "also score the chunk-split variants")
	count := fs.Bool("count", false, "only check the per-category case counts")
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	if *count {
		lines, ok := CheckCounts(*root, []string{"x53", "xcur"})
		fmt.Fprintln(stdout, strings.Join(lines, "\n"))
		return exitCode(ok, strings.Contains(strings.Join(lines, "\n"), "UNMEASURED"))
	}
	pr, err := LoadProbes(filepath.Join(*root, "probes.json"))
	if err != nil {
		fmt.Fprintln(stdout, "UNMEASURED", err)
		return 2
	}
	Register("oracle", Oracle(pr))
	f, ok := Lookup(*engine)
	gs, err := LoadAll(*root, *ref)
	if !ok || err != nil || len(gs) == 0 {
		fmt.Fprintf(stdout, "UNMEASURED engine=%s found=%v goldens=%d err=%v\n", *engine, ok, len(gs), err)
		return 2
	}
	pass, total, skipped := map[string]int{}, map[string]int{}, 0
	for _, g := range gs {
		r := RunCase(g, f, pr, Options{Splits: *splits})
		total[r.Category]++
		if len(r.Skipped) > 0 {
			skipped += len(r.Skipped)
			fmt.Fprintln(stdout, "SKIP", r.ID, strings.Join(r.Skipped, ","), "(headless itself is split-sensitive)")
		}
		if r.Status == PASS {
			pass[r.Category]++
			fmt.Fprintln(stdout, "PASS", r.ID)
		} else {
			fmt.Fprintln(stdout, r.Status, r.ID, r.Diff)
		}
	}
	var cats []string
	for c := range RequiredCategories {
		cats = append(cats, c)
	}
	sort.Strings(cats)
	all, sum, unmeasured := 0, 0, false
	for _, c := range cats {
		status := ""
		if total[c] == 0 {
			status, unmeasured = " UNMEASURED", true
		}
		fmt.Fprintf(stdout, "CAT %s %d/%d%s\n", c, pass[c], total[c], status)
		all, sum = all+pass[c], sum+total[c]
	}
	verdict := "PASS"
	if all != len(gs) {
		verdict = "FAIL"
	}
	fmt.Fprintf(stdout, "TOTAL %d/%d %s skipped=%d\n", all, len(gs), verdict, skipped)
	return exitCode(all == len(gs) && sum == len(gs), unmeasured)
}

func exitCode(ok, unmeasured bool) int {
	switch {
	case unmeasured:
		return 2
	case !ok:
		return 1
	}
	return 0
}
