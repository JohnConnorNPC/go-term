package runner

import "testing"

func TestCounts(t *testing.T) {
	lines, ok := CheckCounts("..", []string{"x53", "xcur"})
	for _, l := range lines {
		t.Log(l)
	}
	if !ok {
		t.Fatal("case counts below the required minimums")
	}
	if len(lines) != len(RequiredCategories)+1 {
		t.Fatalf("%d lines, want one per required category (%d) + total", len(lines), len(RequiredCategories))
	}
}
