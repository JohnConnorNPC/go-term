package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var goDirective = regexp.MustCompile(`(?m)^go (\S+)`)

// goLine is `GOLINE <e> go <v> vs 1.23 ABOVE|OK`: the module's go directive against Studio's floor 1.23 (rulings-4),
// compared numerically field by field (1.3 is below 1.23).
func goLine(engine string, gomod []byte) string {
	m := goDirective.FindSubmatch(gomod)
	if m == nil {
		return "GOLINE " + engine + " UNREADABLE vs 1.23 FAIL"
	}
	up, err := above(string(m[1]), []int{1, 23})
	status := "OK"
	if err != nil {
		status = "FAIL " + err.Error()
	} else if up {
		status = "ABOVE"
	}
	return fmt.Sprintf("GOLINE %s go %s vs 1.23 %s", engine, m[1], status)
}

// above reports whether version v (dotted integers) is greater than floor; a field that isn't an integer is an error.
func above(v string, floor []int) (bool, error) {
	for i, f := range strings.Split(v, ".") {
		n, err := strconv.Atoi(f)
		if err != nil {
			return false, fmt.Errorf("go line %q: %w", v, err)
		}
		if i >= len(floor) || n != floor[i] {
			return i >= len(floor) || n > floor[i], nil
		}
	}
	return false, nil
}
