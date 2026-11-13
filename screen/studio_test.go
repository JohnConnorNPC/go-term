package screen

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/JohnConnorNPC/go-term/sniff"
)

// TestRawReplay checks ReplayBytes against Studio's Python snapshot, whose
// bytes are raw.render(ring, evicted) + _mode_reassert(sniffer)
// (old/webterm/agent/agent.py, "payload += _mode_reassert"). For each
// raw.render case whose output is Preamble + ring (Studio trimmed or
// stripped nothing), ReplayBytes(m, ring) must be that output, one CAN, then
// the Python postamble of m, for every mode set of sniff's postamble.json
// (Python's _mode_reassert bytes). CAN is go-term's one-byte difference
// (README Restore). Studio's trim and strip are Studio's (#8 rulings Q1):
// those cases are counted, not compared.
func TestRawReplay(t *testing.T) {
	var pa struct {
		Rows []struct {
			BracketedPaste bool   `json:"bracketed_paste"`
			AppCursor      bool   `json:"app_cursor"`
			MouseTracking  int    `json:"mouse_tracking"`
			MouseEncoding  int    `json:"mouse_encoding"`
			Hex            string `json:"hex"`
		} `json:"rows"`
	}
	mustUnmarshal(t, readVectors(t, sniffData+"postamble.json"), &pa)
	if len(pa.Rows) == 0 {
		t.Fatal("postamble.json: no rows (unmeasured)")
	}
	exact, trimmed, digest, total := 0, 0, 0, 0
	for _, c := range loadPython(t, "python/raw.json") {
		if c.Entry != "raw.render" {
			continue
		}
		total++
		args, ok := pyDict(t, c.Args)
		ring, okRing := pyBytes(args["ring_bytes"])
		out, okOut := pyBytes(c.Output)
		switch {
		case !ok || !okRing:
			t.Fatalf("%s: no ring_bytes", c.ID)
		case !okOut:
			digest++ // $digest: the output is stored as a sha256 only
			continue
		case !bytes.Equal(out, append([]byte(Preamble), ring...)):
			trimmed++
			continue
		}
		exact++
		for _, r := range pa.Rows {
			post, err := hex.DecodeString(r.Hex)
			if err != nil {
				t.Fatal(err)
			}
			m := sniff.Modes{BracketedPaste: r.BracketedPaste, DECCKM: r.AppCursor, MouseTrack: r.MouseTracking,
				MouseEnc: r.MouseEncoding}
			want := append(append(append([]byte{}, out...), Cancel...), post...)
			if got := ReplayBytes(m, ring); !bytes.Equal(got, want) {
				t.Errorf("%s modes %+v: ReplayBytes %q, python snapshot + CAN %q", c.ID, m, got, want)
			}
		}
	}
	// python3 over vectors/python/raw.json: 31 raw.render, 6 equal to
	// PREAMBLE + ring_bytes, 24 trimmed or stripped, 1 $digest.
	if exact != 6 || trimmed != 24 || digest != 1 || total != 31 {
		t.Errorf("raw.render: exact %d trimmed %d digest %d of %d, want 6 24 1 of 31", exact, trimmed, digest, total)
	}
	t.Logf("REPLAY raw.render %d/%d = python snapshot + CAN (%d mode sets); %d trimmed or stripped by Studio, %d digest",
		exact, total, len(pa.Rows), trimmed, digest)
}

// pyScreen is the part of a Python screen read this package compares: the
// text, the cursor, history_lines and (hashes only) alt_screen.
type pyScreen struct {
	text         string
	cursor       *Cursor
	historyLines int
	alt          *bool
}

// decodePyScreen decodes a textgrid.render output (a str: text only) or a
// render_screen / _render_screen_text output ({"$d": ...}); ok is false for a
// $digest.
func decodePyScreen(t *testing.T, raw json.RawMessage) (p pyScreen, ok bool) {
	t.Helper()
	if json.Unmarshal(raw, &p.text) == nil {
		return p, true
	}
	d, isDict := pyDict(t, raw)
	if !isDict {
		return p, false
	}
	mustUnmarshal(t, d["text"], &p.text)
	if c, has := pyDict(t, d["cursor"]); has {
		p.cursor = &Cursor{}
		mustUnmarshal(t, c["row"], &p.cursor.Row)
		mustUnmarshal(t, c["col"], &p.cursor.Col)
	}
	if h, has := d["history_lines"]; has {
		mustUnmarshal(t, h, &p.historyLines)
	}
	if a, has := d["alt_screen"]; has {
		p.alt = new(bool)
		mustUnmarshal(t, a, p.alt)
	}
	return p, true
}

// pythonDiff returns the fields where go-term's Read differs from a Python
// screen value, in divergence.json's projection: trimmed rows, cursor,
// history_lines, alt_screen.
func pythonDiff(f Frame, p pyScreen) []string {
	var d []string
	py := strings.Split(p.text, "\n")
	if trimAll(py) != strings.Join(f.Trimmed(), "\n") {
		d = append(d, "text")
	}
	if p.cursor != nil && *p.cursor != *f.Cursor {
		d = append(d, "cursor")
	}
	if p.cursor != nil && p.historyLines != f.HistoryLines {
		d = append(d, "history_lines")
	}
	if p.alt != nil && *p.alt != f.AltScreen {
		d = append(d, "alt_screen")
	}
	return d
}

// divergence is Studio's headless/divergence.json: where headless and the
// Python engine that recorded a vector differ (pinned in vectorPins).
type divergence struct {
	Header struct {
		Uncompared []struct {
			ID     string `json:"id"`
			Reason string `json:"reason"`
		} `json:"uncompared"`
		WithOutputState int `json:"compared_with_output_state"`
	} `json:"header"`
	Differ []struct {
		ID      string                     `json:"id"`
		Entry   string                     `json:"entry"`
		Side    string                     `json:"side"`
		Against string                     `json:"against"`
		Fields  map[string]json.RawMessage `json:"fields"`
	} `json:"differ"`
}

func loadDivergence(t *testing.T) divergence {
	t.Helper()
	var d divergence
	mustUnmarshal(t, readVectors(t, "headless/divergence.json"), &d)
	if len(d.Differ) == 0 {
		t.Fatal("divergence.json: no differ rows (unmeasured)")
	}
	return d
}

// pythonScreens runs every textgrid.render, textgrid.render_screen and
// hashes._render_screen_text case through vt (the python replay chunks, as
// headless was fed: no trim of an evicted head, no keyframe) and Read at the
// call's view, and returns the differing fields by "family:id", the counts of
// compared cases by entry and the $digest ids.
func pythonScreens(t *testing.T) (diff map[string][]string, compared map[string]int, digests []string) {
	t.Helper()
	diff, compared = map[string][]string{}, map[string]int{}
	for _, fam := range []string{"textgrid", "hashes"} {
		for _, c := range loadPython(t, "python/"+fam+".json") {
			switch c.Entry {
			case "textgrid.render", "textgrid.render_screen", "hashes._render_screen_text":
			default:
				continue
			}
			p, ok := decodePyScreen(t, c.Output)
			if !ok {
				digests = append(digests, fam+":"+c.ID)
				continue
			}
			var r struct {
				Chunks     []string `json:"chunks_b64"`
				Cols, Rows int
			}
			mustUnmarshal(t, c.Replay, &r)
			req := Request{}
			if a, isDict := pyDict(t, c.Args); isDict {
				var view string
				if v, has := a["view"]; has {
					mustUnmarshal(t, v, &view)
				}
				if view == "scrollback" {
					req.View = ViewScrollback
					mustUnmarshal(t, a["lines"], &req.Lines)
				}
			}
			term := newTerm(r.Cols, r.Rows)
			for _, ch := range r.Chunks {
				b, err := base64.StdEncoding.DecodeString(ch)
				if err != nil {
					t.Fatalf("%s: %v", c.ID, err)
				}
				term.Write(b)
			}
			compared[c.Entry]++
			if d := pythonDiff(Read(term, req), p); len(d) > 0 {
				diff[fam+":"+c.ID] = d
			}
		}
	}
	return diff, compared, digests
}

// TestPythonTextgrid compares go-term with Python's own screen values. Where
// they differ, headless differs from Python too: the differing ids and
// fields are exactly divergence.json's rows against python for these
// entries (go-term equals headless there, TestHeadlessVectors).
func TestPythonTextgrid(t *testing.T) {
	diff, compared, digests := pythonScreens(t)
	listed := map[string][]string{}
	for _, d := range loadDivergence(t).Differ {
		switch d.Entry {
		case "textgrid.render", "textgrid.render_screen", "hashes._render_screen_text":
		default:
			continue
		}
		if d.Against != "python" {
			continue
		}
		var fs []string
		for f := range d.Fields {
			fs = append(fs, f)
		}
		sort.Strings(fs)
		listed[d.ID] = fs
	}
	for id, fs := range diff {
		sort.Strings(fs)
		if l, ok := listed[id]; !ok || strings.Join(l, ",") != strings.Join(fs, ",") {
			t.Errorf("%s: go-term differs from python in %v, divergence.json lists %v", id, fs, l)
		}
	}
	for id, l := range listed {
		if _, ok := diff[id]; !ok {
			t.Errorf("%s: divergence.json lists %v, go-term equals python", id, l)
		}
	}
	n := 0
	for _, v := range compared {
		n += v
	}
	t.Logf("PYTHON screens render %d render_screen %d _render_screen_text %d: %d match, %d differ (= divergence.json %d); %d digest",
		compared["textgrid.render"], compared["textgrid.render_screen"], compared["hashes._render_screen_text"],
		n-len(diff), len(diff), len(listed), len(digests))
}

// sniffDiff replays a Python sniffer case (altscreen.DecModeSniffer,
// AltScreenSniffer, titles.OscTitleSniffer) through package sniff op by op,
// as sniff's own vector test does, and returns the differing field ("modes"
// or "title") or "" when every op agrees. TitleMax is 4096, the Python
// payload cap (sniff's maxPayload), so no Python title is cut.
func sniffDiff(t *testing.T, c pyCase) string {
	t.Helper()
	s := sniff.New(sniff.Options{TitleMax: 4096})
	want, have := "", false
	if len(c.Ops) == 0 {
		t.Fatalf("%s: no ops", c.ID)
	}
	for i, op := range c.Ops {
		args, ok := pyDict(t, op[1])
		if !ok || len(args) != 1 {
			t.Fatalf("%s op %d: args %s", c.ID, i, op[1])
		}
		for _, v := range args {
			b, isBytes := pyBytes(v)
			if !isBytes {
				t.Fatalf("%s op %d: chunk %s", c.ID, i, v)
			}
			s.Feed(b)
		}
		if c.Entry == "titles.OscTitleSniffer" {
			var r *string
			mustUnmarshal(t, op[2], &r)
			if r != nil {
				want, have = *r, true
			}
			if got, ok := s.Title(); got != want || ok != have {
				return "title"
			}
			continue
		}
		st, ok := pyDict(t, op[3])
		if !ok || len(st) != 5 {
			t.Fatalf("%s op %d: state %s", c.ID, i, op[3])
		}
		var m sniff.Modes
		mustUnmarshal(t, st["alt_screen"], &m.Alt)
		mustUnmarshal(t, st["app_cursor"], &m.DECCKM)
		mustUnmarshal(t, st["bracketed_paste"], &m.BracketedPaste)
		mustUnmarshal(t, st["mouse_tracking"], &m.MouseTrack)
		mustUnmarshal(t, st["mouse_encoding"], &m.MouseEnc)
		if s.Modes() != m {
			return "modes"
		}
	}
	return ""
}

// studioEntries are the A2a producers with no go-term counterpart: Studio's
// byte producers (raw snapshot, pyte re-emit, textgrid's trim), checked
// through output_state where their output is bytes (#9 rulings-0 R2).
var studioEntries = map[string]bool{"raw.render": true, "raw._trim": true, "pyte.render": true,
	"pyte.emit_screen": true, "pyte._color_params": true, "textgrid._trim_for_screen": true}

// TestStudioLedger gives every A2a case (Studio's vectors/python, 746) exactly
// one disposition, each re-checked here:
//   - GO: go-term reproduces Python's value (attr_runs, the hashes, the
//     screen reads of TestPythonTextgrid, the sniffers through package sniff);
//   - LISTED: go-term differs from Python, divergence.json lists the case
//     against python with those fields, and go-term equals headless there;
//   - STUDIO: a Studio producer (studioEntries); its bytes, when bytes, are fed
//     to vt and equal headless's output_state;
//   - UNCOMPARED: divergence.json header "uncompared", and pyte.available.
//
// Every A2b case passes outright (TestHeadlessVectors). go-term's own
// difference from Studio's snapshot, the CAN byte, is TestRawReplay's.
func TestStudioLedger(t *testing.T) {
	div := loadDivergence(t)
	uncompared := map[string]bool{}
	for _, u := range div.Header.Uncompared {
		uncompared[u.ID] = true
	}
	listed := map[string]string{}
	for _, d := range div.Differ {
		if d.Against != "python" {
			continue
		}
		var fs []string
		for f := range d.Fields {
			fs = append(fs, f)
		}
		sort.Strings(fs)
		listed[d.ID] = strings.Join(fs, ",")
	}
	screens, _, _ := pythonScreens(t)
	disp := map[string]int{}
	seen := map[string]bool{}
	reached := map[string]bool{} // python rows met in the GO or LISTED branch
	listedGO := 0
	outState, colourSame, colourLost := 0, 0, 0
	for _, fam := range []struct{ name, dir string }{{"textgrid", ""}, {"colours", ""}, {"hashes", ""}, {"pyte", ""},
		{"raw", ""}, {"altscreen", sniffData}, {"synthetic", sniffData}, {"titles", sniffData}} {
		py := map[string]pyCase{}
		cases := loadPython(t, fam.dir+"python/"+fam.name+".json")
		for _, c := range cases {
			py[c.ID] = c
		}
		var hf struct {
			Cases []hlCase `json:"cases"`
		}
		mustUnmarshal(t, readVectors(t, fam.dir+"headless/"+fam.name+".json"), &hf)
		hl := map[string]hlCase{}
		for _, c := range hf.Cases {
			hl[c.ID] = c
		}
		for _, c := range cases {
			key := fam.name + ":" + c.ID
			seen[key] = true
			var diff string // the fields go-term differs from Python in, "" for GO
			switch {
			case uncompared[key] || c.Entry == "pyte.available":
				disp["UNCOMPARED"]++
				continue
			case studioEntries[c.Entry]:
				out, isBytes := pyBytes(c.Output)
				h, ok := hl[c.ID]
				switch {
				case !isBytes && c.Entry == "pyte._color_params":
					// Studio's re-emit table: go-term reads the params back.
					switch same := colourParamsBack(t, c); {
					case same:
						colourSame++
					case colourLostNames[colourParamsName(t, c)]:
						colourLost++
					default:
						t.Errorf("%s: go-term reads %s's params back as another colour", key, c.Entry)
					}
				case !isBytes:
					t.Errorf("%s: %s output is not bytes and not uncompared", key, c.Entry)
				case !ok || h.OutputState == nil:
					t.Errorf("%s: no headless output_state", key)
				default:
					term := newTerm(h.Cols, h.Rows)
					term.Write(out)
					if err := checkHeadless(term, *h.OutputState); err != nil {
						t.Errorf("%s output_state: %v", key, err)
					}
					outState++
				}
				disp["STUDIO"]++
				continue
			case c.Entry == "pyte.attr_runs":
				got, want := runPyAttrRuns(t, c)
				for i := range want {
					if fam.name == "colours" && want[i].Bg == "bfightmagenta" { // colours.json's named exception
						want[i].Bg = "brightmagenta"
					}
				}
				if fmt.Sprint(got) != fmt.Sprint(want) {
					diff = "attrs"
				}
			case c.Entry == "hashes._content_hash" || c.Entry == "hashes._stable_hash":
				var want string
				mustUnmarshal(t, c.Output, &want)
				if pyHash(t, c) != want {
					diff = "value"
				}
			case c.Entry == "textgrid.render" || c.Entry == "textgrid.render_screen" ||
				c.Entry == "hashes._render_screen_text":
				diff = strings.Join(sortedCopy(screens[key]), ",")
			case c.Entry == "altscreen.DecModeSniffer" || c.Entry == "altscreen.AltScreenSniffer" ||
				c.Entry == "titles.OscTitleSniffer":
				diff = sniffDiff(t, c)
			default:
				t.Errorf("%s: entry %q has no disposition", key, c.Entry)
				continue
			}
			if diff == "" {
				// headless differs from Python here and go-term agrees with
				// Python: only where go-term follows Python by design.
				if listed[key] != "" {
					switch {
					case fam.name == "colours": // Colour maps xterm's palette to pyte's names (SA1 Q1)
					case strings.HasPrefix(c.ID, "synthetic/") && strings.Contains(c.ID, "-cap#"): // sniff keeps Python's caps
					default:
						t.Errorf("%s: divergence.json lists %s, but go-term equals Python", key, listed[key])
					}
					listedGO++
					reached[key] = true
				}
				disp["GO"]++
				continue
			}
			h, ok := hl[c.ID]
			switch {
			case listed[key] == "":
				t.Errorf("%s: go-term differs from Python in %s, not in divergence.json", key, diff)
			case listed[key] != diff:
				t.Errorf("%s: go-term differs from Python in %s, divergence.json lists %s", key, diff, listed[key])
			case !ok:
				t.Errorf("%s: listed, but no headless case", key)
			default:
				if err := checkHeadless(replayHeadless(t, h, py), h.Final); err != nil {
					t.Errorf("%s: listed, but go-term differs from headless: %v", key, err)
				}
				if c.Entry == "titles.OscTitleSniffer" {
					s := sniff.New(sniff.Options{TitleMax: 4096})
					var r struct {
						Chunks []string `json:"chunks_b64"`
					}
					mustUnmarshal(t, c.Replay, &r)
					for _, ch := range r.Chunks {
						b, _ := base64.StdEncoding.DecodeString(ch)
						s.Feed(b)
					}
					if got, ok := s.Title(); h.Final.Title == nil || !ok || got != *h.Final.Title {
						t.Errorf("%s: listed, but sniff's title %q differs from headless", key, got)
					}
				}
			}
			if listed[key] != "" {
				reached[key] = true
			}
			disp["LISTED"]++
		}
	}
	pyRows := 0
	for _, d := range div.Differ {
		if !seen[d.ID] {
			t.Errorf("divergence.json %s: no such A2a case", d.ID)
		}
		if d.Against == "python" {
			pyRows++
			if !reached[d.ID] {
				t.Errorf("divergence.json %s (against python): never compared (not GO or LISTED)", d.ID)
			}
		}
	}
	n := disp["GO"] + disp["LISTED"] + disp["STUDIO"] + disp["UNCOMPARED"]
	if n != 746 || len(seen) != 746 {
		t.Errorf("%d dispositions over %d cases, want 746", n, len(seen))
	}
	// Each disposition is pinned, not only the sum: a case moving between
	// two of them keeps 746. The vector files are pinned by sha256.
	want := map[string]int{"GO": 546, "LISTED": 27, "STUDIO": 167, "UNCOMPARED": 6}
	for k, w := range want {
		if disp[k] != w {
			t.Errorf("ledger %s %d, want %d", k, disp[k], w)
		}
	}
	if listedGO != 11 || disp["LISTED"]+listedGO != pyRows {
		t.Errorf("listed GO %d (want 11) + LISTED %d != %d divergence.json rows against python", listedGO, disp["LISTED"], pyRows)
	}
	if outState != div.Header.WithOutputState || colourSame+colourLost+outState != disp["STUDIO"] {
		t.Errorf("STUDIO %d: output_state %d (divergence.json %d) + colour params %d+%d", disp["STUDIO"], outState,
			div.Header.WithOutputState, colourSame, colourLost)
	}
	if colourLost != 2*len(colourLostNames) {
		t.Errorf("colour params lost %d, want %d (fg and bg of %d names)", colourLost, 2*len(colourLostNames), len(colourLostNames))
	}
	t.Logf("LEDGER a2a GO %d (%d listed, go-term follows Python) LISTED %d STUDIO %d (output_state %d, _color_params read back %d same + %d default for brown/brightbrown/bfightmagenta) UNCOMPARED %d = %d; python rows %d reached; a2b 555 (TestHeadlessVectors); list testdata/headless/divergence.json %.16s; divergences can-separator (TestRawReplay)",
		disp["GO"], listedGO, disp["LISTED"], disp["STUDIO"], outState, colourSame, colourLost, disp["UNCOMPARED"], n, pyRows, vectorPins["headless/divergence.json"])
}

// colourLostNames are the pyte names Studio's _color_params has no params
// for: it emits 39/49, so a re-emitted screen shows them in the default
// colour. pyte names SGR 33 "brown", so this loses yellow text in Studio's
// replay (a Studio matter; go-term emits no screens).
var colourLostNames = map[string]bool{"brown": true, "brightbrown": true, "bfightmagenta": true}

// colourParamsName is the colour name a pyte._color_params case is about.
func colourParamsName(t *testing.T, c pyCase) string {
	t.Helper()
	var a struct {
		D [][2]json.RawMessage `json:"$d"`
	}
	mustUnmarshal(t, c.Args, &a)
	if len(a.D) != 2 {
		t.Fatalf("%s: args %s", c.ID, c.Args)
	}
	var name string
	mustUnmarshal(t, a.D[0][1], &name)
	return name
}

// colourParamsBack feeds Python's SGR params for a colour (fg or bg) before
// one glyph and reports whether go-term reads that cell as the same colour.
func colourParamsBack(t *testing.T, c pyCase) bool {
	t.Helper()
	var a struct {
		D [][2]json.RawMessage `json:"$d"`
	}
	mustUnmarshal(t, c.Args, &a)
	var isBg bool
	mustUnmarshal(t, a.D[1][1], &isBg)
	var params []int
	mustUnmarshal(t, c.Output, &params)
	if len(params) == 0 {
		t.Fatalf("%s: no params", c.ID)
	}
	ps := make([]string, len(params))
	for i, p := range params {
		ps[i] = fmt.Sprint(p)
	}
	term := newTerm(4, 1)
	term.WriteString("\x1b[" + strings.Join(ps, ";") + "mX")
	got := Default
	if runs := AttrRuns(term, MaxAttrRuns); len(runs) == 1 {
		got = runs[0].Fg
		if isBg {
			got = runs[0].Bg
		}
	} else if len(runs) > 1 {
		t.Fatalf("%s: %d runs for one glyph", c.ID, len(runs))
	}
	return got == colourParamsName(t, c)
}

func sortedCopy(ss []string) []string {
	out := append([]string(nil), ss...)
	sort.Strings(out)
	return out
}
