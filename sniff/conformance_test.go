package sniff

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// conformanceDiffs lists the golden steps this package is known to differ on, by "<ref>/<case id>"; each must still
// differ. Empty: every step of every case agrees.
var conformanceDiffs = map[string]string{}

type goldenStep struct {
	Op struct {
		Op string `json:"op"`
		In string `json:"in"`
	} `json:"op"`
	P struct {
		Buffer string `json:"buffer"`
		Title  string `json:"title"`
		Modes  struct {
			DECCKM   *bool   `json:"applicationCursorKeysMode"`
			BP       *bool   `json:"bracketedPasteMode"`
			Tracking *string `json:"mouseTrackingMode"`
		} `json:"modes"`
	} `json:"P"`
	I struct {
		MouseEncoding *string `json:"mouseEncoding"`
	} `json:"I"`
}

// TestConformance runs the sniffer over every conformance case's golden (x53 = xterm-headless 5.3.0, xcur =
// @xterm/headless 6.0.0), step by step from the exact bytes term.write received (op.in), and compares after every
// step: P.buffer, P.modes DECCKM / bracketed paste / mouseTrackingMode, I.mouseEncoding, and P.title (the last
// onTitleChange, "" before any). A resize step feeds nothing. Every case file must have a golden in both refs.
func TestConformance(t *testing.T) {
	cases, err := filepath.Glob("../conformance/cases/*/*.json")
	if err != nil || len(cases) == 0 {
		t.Fatalf("no conformance cases (unmeasured): %v", err)
	}
	steps, lit := map[string]int{}, map[string]int{}
	for _, ref := range []string{"x53", "xcur"} {
		for _, cf := range cases {
			rel, _ := filepath.Rel("../conformance/cases", cf)
			id := ref + "/" + strings.TrimSuffix(filepath.ToSlash(rel), ".json")
			diff, n, l := runGolden(t, filepath.Join("../conformance/golden", ref, rel))
			steps[ref] += n
			lit[ref] += l
			why, listed := conformanceDiffs[id]
			switch {
			case listed && diff == "":
				t.Errorf("%s: listed (%s) but agrees", id, why)
			case !listed && diff != "":
				t.Errorf("%s: %s", id, diff)
			}
		}
	}
	t.Logf("%d cases, steps %v, steps with a mode on or a title %v", len(cases), steps, lit)
	if steps["x53"] == 0 || steps["x53"] != steps["xcur"] || lit["x53"] == 0 || lit["xcur"] == 0 {
		t.Errorf("steps %v, lit %v: unmeasured or uneven", steps, lit)
	}
}

// runGolden returns the first difference, the steps compared, and how many of them had a mode on or a title.
func runGolden(t *testing.T, path string) (string, int, int) {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var g struct {
		Steps []goldenStep `json:"steps"`
	}
	if err := json.Unmarshal(b, &g); err != nil {
		t.Fatalf("%s: %v", path, err)
	}
	if len(g.Steps) == 0 {
		t.Fatalf("%s: no steps", path)
	}
	s := New(Options{TitleMax: maxPayload})
	lit := 0
	for i, st := range g.Steps {
		switch st.Op.Op {
		case "w":
			p, err := hex.DecodeString(st.Op.In)
			if err != nil {
				t.Fatalf("%s step %d: %v", path, i, err)
			}
			s.Feed(p)
		case "resize":
		default:
			t.Fatalf("%s step %d: op %q", path, i, st.Op.Op)
		}
		m := st.P.Modes
		if m.DECCKM == nil || m.BP == nil || m.Tracking == nil || st.I.MouseEncoding == nil {
			t.Fatalf("%s step %d: a mode is unmeasured", path, i)
		}
		h := hchunk{Buffer: st.P.Buffer}
		h.Modes.DECCKM, h.Modes.BP, h.Modes.Tracking, h.Modes.Encoding = m.DECCKM, m.BP, m.Tracking, st.I.MouseEncoding
		want := xtModes(t, path, h)
		if got := s.Modes(); got != want {
			return fmt.Sprintf("step %d: Modes() = %+v; golden %+v", i, got, want), i + 1, lit
		}
		if want != (Modes{}) || st.P.Title != "" {
			lit++
		}
		if got, _ := s.Title(); got != st.P.Title {
			return fmt.Sprintf("step %d: Title() = %.60q; golden %.60q", i, got, st.P.Title), i + 1, lit
		}
	}
	return "", len(g.Steps), lit
}
