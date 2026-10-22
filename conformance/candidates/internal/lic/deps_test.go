package lic

import (
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"
)

// TestDeps: an MIT dependency is OK; one without the MIT grant is FAIL for the gated candidate and INFO for another;
// test-only requires are only listed.
func TestDeps(t *testing.T) {
	mods := map[string]fs.FS{"example.com/mit": fstest.MapFS{"LICENSE": {Data: []byte(mitText)}},
		"example.com/bsd": fstest.MapFS{"LICENSE": {Data: []byte("Redistribution and use in source and binary forms")}}}
	for gated, want := range map[bool]string{true: "L9 FAIL example.com/bsd not-mit=L2,L3,L4", false: "L9 INFO example.com/bsd not-mit=L2,L3,L4"} {
		var got []string
		for _, l := range Deps(gated, mods, []string{"github.com/google/go-cmp"}) {
			got = append(got, l.ID+" "+l.Status+" "+l.Value)
		}
		all := strings.Join(got, "\n")
		if !strings.Contains(all, want) || !strings.Contains(all, "L9 OK example.com/mit not-mit=none") ||
			!strings.Contains(all, "L9 INFO test-only github.com/google/go-cmp") {
			t.Errorf("gated=%v:\n%s", gated, all)
		}
	}
	if l := Deps(true, nil, nil); l[0].Status != "OK" || l[1].Value != "test-only none" {
		t.Errorf("no deps: %v", l)
	}
}
