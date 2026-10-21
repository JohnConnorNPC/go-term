package runner

import (
	"strings"
	"testing"
)

func TestCompare(t *testing.T) {
	base := func() State {
		return State{P: P{Cols: 2, Rows: 1, Buffer: "normal", Modes: map[string]any{"insertMode": false}},
			Normal: []Row{{Cells: []Cell{{Ch: "a", W: 1}, {Ch: "", W: 1}}}}}
	}
	if d := Compare("step[0]", base(), base()); d != "" {
		t.Fatalf("equal states differ: %s", d)
	}
	cases := map[string]func(*State){
		`step[0].normal[0][0].ch: want "a" got "b"`:       func(s *State) { s.Normal[0].Cells[0].Ch = "b" },
		`step[0].P.cursor.x: want 0 got 1`:                func(s *State) { s.P.Cursor[0] = 1 },
		`step[0].P.modes.insertMode: want false got true`: func(s *State) { s.P.Modes["insertMode"] = true },
		`step[0].normal[0][1].attr`:                       func(s *State) { s.Normal[0].Cells[1].A.Bold = true },
		`step[0].alt: want false got true`:                func(s *State) { s.Alt = []Row{} },
	}
	for want, edit := range cases {
		got := base()
		edit(&got)
		if d := Compare("step[0]", base(), got); !strings.HasPrefix(d, want) {
			t.Errorf("diff %q, want prefix %q", d, want)
		}
	}
}
