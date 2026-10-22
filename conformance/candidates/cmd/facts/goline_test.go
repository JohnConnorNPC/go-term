package main

import "testing"

func TestGoLine(t *testing.T) {
	for in, want := range map[string]string{"go 1.23\n": "GOLINE e go 1.23 vs 1.23 OK", "module x\n\ngo 1.24.2\n": "GOLINE e go 1.24.2 vs 1.23 ABOVE",
		"go 1.3\n": "GOLINE e go 1.3 vs 1.23 OK", "go 1.25\n": "GOLINE e go 1.25 vs 1.23 ABOVE", "module x\n": "GOLINE e UNREADABLE vs 1.23 FAIL",
		"go 1.2x\n": `GOLINE e go 1.2x vs 1.23 FAIL go line "1.2x": strconv.Atoi: parsing "2x": invalid syntax`} {
		if got := goLine("e", []byte(in)); got != want {
			t.Errorf("%q: %s, want %s", in, got, want)
		}
	}
}
