// Command conformance scores a terminal engine against the xterm.js headless goldens.
//
//	go run ./conformance/cmd/conformance -engine stub -ref x53
//	go run ./conformance/cmd/conformance -count
//
// Exit status: 0 every case passes (or the counts meet the minimums), 1 a case fails, 2 something is unmeasured.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/JohnConnorNPC/go-term/conformance/runner"
)

func main() {
	root := flag.String("root", "conformance", "the conformance directory (cases/, golden/, probes.json)")
	engine := flag.String("engine", "stub", "engine to score: "+strings.Join(append(runner.Engines(), "oracle"), ", "))
	ref := flag.String("ref", "x53", "golden set: x53 (xterm-headless@5.3.0) or xcur (@xterm/headless latest)")
	splits := flag.Bool("splits", true, "also score the chunk-split variants")
	count := flag.Bool("count", false, "only check the per-category case counts")
	flag.Parse()
	if *count {
		lines, ok := runner.CheckCounts(*root, []string{"x53", "xcur"})
		fmt.Println(strings.Join(lines, "\n"))
		exit(ok, strings.Contains(strings.Join(lines, "\n"), "UNMEASURED"))
	}
	pr, err := runner.LoadProbes(filepath.Join(*root, "probes.json"))
	if err != nil {
		fmt.Println("UNMEASURED", err)
		os.Exit(2)
	}
	runner.Register("oracle", runner.Oracle(pr))
	f, ok := runner.Lookup(*engine)
	gs, err := runner.LoadAll(*root, *ref)
	if !ok || err != nil || len(gs) == 0 {
		fmt.Printf("UNMEASURED engine=%s found=%v goldens=%d err=%v\n", *engine, ok, len(gs), err)
		os.Exit(2)
	}
	pass, total := map[string]int{}, map[string]int{}
	for _, g := range gs {
		r := runner.RunCase(g, f, pr, runner.Options{Splits: *splits})
		total[r.Category]++
		if r.Status == runner.PASS {
			pass[r.Category]++
			fmt.Println("PASS", r.ID)
		} else {
			fmt.Println(r.Status, r.ID, r.Diff)
		}
	}
	var cats []string
	for c := range runner.RequiredCategories {
		cats = append(cats, c)
	}
	sort.Strings(cats)
	all, sum, unmeasured := 0, 0, false
	for _, c := range cats {
		status := ""
		if total[c] == 0 {
			status, unmeasured = " UNMEASURED", true
		}
		fmt.Printf("CAT %s %d/%d%s\n", c, pass[c], total[c], status)
		all, sum = all+pass[c], sum+total[c]
	}
	verdict := "PASS"
	if all != len(gs) {
		verdict = "FAIL"
	}
	fmt.Printf("TOTAL %d/%d %s\n", all, len(gs), verdict)
	exit(all == len(gs) && sum == len(gs), unmeasured)
}

func exit(ok, unmeasured bool) {
	switch {
	case unmeasured:
		os.Exit(2)
	case !ok:
		os.Exit(1)
	}
	os.Exit(0)
}
