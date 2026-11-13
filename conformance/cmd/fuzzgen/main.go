// Command fuzzgen writes random conformance cases for the differential fuzz against @xterm/headless 6.0.0.
//
//	go run ./conformance/cmd/fuzzgen -seed S -n N -root DIR
//
// It writes DIR/cases/fuzz/<S>-<i>.json (the case format of conformance/README "Layout") and copies
// conformance/probes.json to DIR/probes.json and conformance/gen/pins.json to DIR/gen/pins.json, so DIR is a root that
// gen/fuzz.mjs fills with xcur goldens and cmd/conformance scores. The same seed always writes the same files.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	seed := flag.Uint64("seed", 1, "batch seed (the cases of one seed are always the same)")
	n := flag.Int("n", 100, "number of cases")
	root := flag.String("root", "", "output root (cases/fuzz, probes.json, gen/pins.json)")
	src := flag.String("src", "conformance", "the conformance directory to copy probes.json and gen/pins.json from")
	flag.Parse()
	if *root == "" || *n < 1 {
		fmt.Fprintln(os.Stderr, "fuzzgen: -root and -n >= 1 are required")
		os.Exit(2)
	}
	bytes, err := write(*root, *src, *seed, *n)
	if err != nil {
		fmt.Fprintln(os.Stderr, "fuzzgen:", err)
		os.Exit(1)
	}
	fmt.Printf("fuzzgen ok seed=%d cases=%d bytes=%d\n", *seed, *n, bytes)
}

// write generates n cases of seed under root and copies the shared files; it returns the total write bytes.
func write(root, src string, seed uint64, n int) (int, error) {
	dir := filepath.Join(root, "cases", "fuzz")
	if err := os.MkdirAll(filepath.Join(root, "gen"), 0o755); err != nil {
		return 0, err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return 0, err
	}
	for from, to := range map[string]string{"probes.json": "probes.json", "gen/pins.json": "gen/pins.json"} {
		b, err := os.ReadFile(filepath.Join(src, from))
		if err != nil {
			return 0, err
		}
		if err := os.WriteFile(filepath.Join(root, to), b, 0o644); err != nil {
			return 0, err
		}
	}
	total := 0
	for i := range n {
		c := Generate(seed, i)
		total += c.writeBytes()
		b, err := json.Marshal(c)
		if err != nil {
			return 0, err
		}
		if err := os.WriteFile(filepath.Join(dir, fmt.Sprintf("%d-%d.json", seed, i)), append(b, '\n'), 0o644); err != nil {
			return 0, err
		}
	}
	return total, nil
}
