package main

import (
	"strings"
	"testing"
)

// TestCheck (real modules): gitpod is MIT with Ona's line and none of the xterm.js lines, has no non-test imports and
// go-cmp test-only; magnet carries the three xterm.js lines.
func TestCheck(t *testing.T) {
	text := func(gated bool, mod string) string {
		lines, err := check(gated, mod)
		if err != nil {
			t.Fatal(err)
		}
		var s []string
		for _, l := range lines {
			s = append(s, l.ID+" "+l.Status+" "+l.Value)
		}
		return strings.Join(s, "\n") + "\n"
	}
	g := text(true, "github.com/gitpod-io/xterm-go")
	for _, want := range []string{"L1 OK LICENSE\n", "L4 OK", "L5 INFO Copyright (c) 2026 Ona\n", "L6 INFO 0/3\n", "L10 INFO e625cb3be4cc\n",
		"L7 OK copyleft none\n", "L11 OK restrictions none\n",
		"L9 OK no non-test dependencies\n", "L9 INFO test-only github.com/google/go-cmp\n"} {
		if !strings.Contains(g, want) {
			t.Errorf("gitpod lacks %q:\n%s", want, g)
		}
	}
	if m := text(false, "github.com/0magnet/xterm-go"); !strings.Contains(m, "L6 INFO 3/3\n") || !strings.Contains(m, "L10 INFO b569f629d00f\n") {
		t.Errorf("magnet:\n%s", m)
	}
}
