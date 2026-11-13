package screen

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
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
