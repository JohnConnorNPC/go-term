package sniff

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
)

// The vector files are byte-for-byte copies of Browserland Studio's vectors/python/ (A2a: what the Python reference
// sniffers return, recorded from its own tests) and vectors/headless/ (A2b: the same inputs replayed through
// @xterm/headless 6.0.0, the release the page runs), copied at studio 766fcb6. A regeneration there needs a new copy here.
var vectorPins = map[string]string{
	"python/altscreen.json":   "c5221a92bec87e0a2791380b9d836c019888977aa17833b5510d264cc9c397ed",
	"python/titles.json":      "58c0384d54a2e57622a0e96f557d8908e83b8abb22ad29e0e90ce8c3c2cce791",
	"python/synthetic.json":   "d3ef7b87b05939878688cb0ce80dc6ef49eb757f44b7c71c7ca254a028f86803",
	"headless/altscreen.json": "6a60401833fddd30145a54dce28e7173b5c70447459d4b2805d77ee6c61d0612",
	"headless/titles.json":    "63cb5a8595e91d3d20d941a34a9e6d908fdcf0e068ea2e1c26a1bfbd8e061879",
	"headless/synthetic.json": "d5407c330659f18390c2d328c8592287ded4bc39acc8ca3253fbda9ee5634046",
}

// pythonDivergences are the Python vectors this package does not follow, by the owner's ruling that xterm.js wins
// where the two differ: a bare ESC inside an OSC title payload. The Python sniffer aborts the title there; xterm.js
// ends the OSC successfully and commits it ("junk", "another title"), and so does this package. Each listed case
// must still differ from Python (a listed case that agrees is a failure), and every other case must agree.
var pythonDivergences = map[string]string{
	"tests/test_titles.py::test_st_terminated_survives_any_split[21]#0": "a bare ESC inside the payload ends the OSC and commits the title (xterm.js); Python aborts",
	"tests/test_titles.py::test_esc_abort_inside_payload#0":             "a bare ESC inside the payload ends the OSC and commits the title (xterm.js); Python aborts",
}

func readVectors(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestVectorPins(t *testing.T) {
	names := make([]string, 0, len(vectorPins))
	for n := range vectorPins {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		sum := sha256.Sum256(readVectors(t, n))
		if got := hex.EncodeToString(sum[:]); got != vectorPins[n] {
			t.Errorf("%s: sha256 %s, pinned %s", n, got, vectorPins[n])
		}
	}
}

// vcase is one stream case of a vector file: a Python case's ops, or a headless case's per-chunk states.
type vcase struct {
	ID    string            `json:"id"`
	Entry string            `json:"entry"`
	Kind  string            `json:"kind"`
	Ops   []json.RawMessage `json:"ops"`

	Replay struct {
		Chunks []string `json:"chunks_b64"`
	} `json:"replay"`

	Chunks []hchunk `json:"chunks"`
	Final  *hchunk  `json:"final"`
}

type hchunk struct {
	Buffer string `json:"buffer"`
	Modes  struct {
		DECCKM   *bool   `json:"applicationCursorKeysMode"`
		BP       *bool   `json:"bracketedPasteMode"`
		Tracking *string `json:"mouseTrackingMode"`
		Encoding *string `json:"mouse_encoding"`
	} `json:"modes"`
	Titles []string `json:"titles"`
	Title  *string  `json:"title"`
}

type vfile struct {
	Header struct {
		Family string          `json:"family"`
		Source json.RawMessage `json:"source"` // a path in python/, {file, sha256} in headless/
	} `json:"header"`
	Cases []vcase `json:"cases"`
}

func loadVectors(t *testing.T, name, family string) vfile {
	t.Helper()
	var f vfile
	if err := json.Unmarshal(readVectors(t, name), &f); err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	if f.Header.Family != family {
		t.Fatalf("%s: family %q, want %q", name, f.Header.Family, family)
	}
	if len(f.Cases) == 0 {
		t.Fatalf("%s: 0 cases (unmeasured)", name)
	}
	return f
}

// streamCases returns the sniffer cases of a family: every case of altscreen and titles, and the cap cases of
// synthetic (the entries are the two Python sniffer classes).
func streamCases(t *testing.T, f vfile) []vcase {
	var out []vcase
	for _, c := range f.Cases {
		switch c.Entry {
		case "altscreen.DecModeSniffer", "altscreen.AltScreenSniffer", "titles.OscTitleSniffer":
			if c.Kind != "stream" {
				t.Fatalf("%s: kind %q", c.ID, c.Kind)
			}
			out = append(out, c)
		default:
			if f.Header.Family != "synthetic" {
				t.Fatalf("%s: entry %q in family %s", c.ID, c.Entry, f.Header.Family)
			}
		}
	}
	return out
}

// untag decodes one value of the A2a tagged encoding as far as the sniffer cases use it: {"$d": [[k, v], ...]} (a
// dict, as ordered pairs), {"$b": base64} (bytes) and plain JSON scalars. Anything else is an error, never skipped.
func untag(raw json.RawMessage) (any, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return nil, err
	}
	return untagValue(v)
}

type pair struct {
	k string
	v any
}

func untagValue(v any) (any, error) {
	switch x := v.(type) {
	case nil, bool, string:
		return x, nil
	case json.Number:
		return x.Int64()
	case map[string]any:
		if len(x) != 1 {
			return nil, fmt.Errorf("an untagged object with %d keys", len(x))
		}
		if b, ok := x["$b"].(string); ok {
			return base64.StdEncoding.Strict().DecodeString(b)
		}
		ps, ok := x["$d"].([]any)
		if !ok {
			return nil, fmt.Errorf("unknown tag in %v", x)
		}
		d := []pair{}
		for _, p := range ps {
			kv, ok := p.([]any)
			if !ok || len(kv) != 2 {
				return nil, fmt.Errorf("$d pair %v", p)
			}
			k, ok := kv[0].(string)
			if !ok {
				return nil, fmt.Errorf("$d key %v", kv[0])
			}
			val, err := untagValue(kv[1])
			if err != nil {
				return nil, err
			}
			d = append(d, pair{k, val})
		}
		return d, nil
	}
	return nil, fmt.Errorf("cannot untag %T", v)
}

// pyOp is one decoded feed op: the chunk, the call's result and the state after.
type pyOp struct {
	chunk  []byte
	result any
	state  []pair
}

func pyOps(t *testing.T, c vcase) []pyOp {
	t.Helper()
	var out []pyOp
	for i, raw := range c.Ops {
		var op []json.RawMessage
		if err := json.Unmarshal(raw, &op); err != nil || len(op) != 4 {
			t.Fatalf("%s op %d: not [method, args, result, state]", c.ID, i)
		}
		var method string
		if err := json.Unmarshal(op[0], &method); err != nil || method != "feed" {
			t.Fatalf("%s op %d: method %s", c.ID, i, op[0])
		}
		args, err1 := untag(op[1])
		res, err2 := untag(op[2])
		st, err3 := untag(op[3])
		a, ok1 := args.([]pair)
		s, ok2 := st.([]pair)
		if err1 != nil || err2 != nil || err3 != nil || !ok1 || !ok2 || len(a) != 1 {
			t.Fatalf("%s op %d: %v %v %v", c.ID, i, err1, err2, err3)
		}
		chunk, ok := a[0].v.([]byte)
		if !ok {
			t.Fatalf("%s op %d: the chunk is %T", c.ID, i, a[0].v)
		}
		out = append(out, pyOp{chunk, res, s})
	}
	if len(out) == 0 {
		t.Fatalf("%s: no ops", c.ID)
	}
	return out
}

func pyModes(t *testing.T, id string, st []pair) Modes {
	t.Helper()
	var m Modes
	seen := 0
	for _, p := range st {
		var ok bool
		switch p.k {
		case "alt_screen":
			m.Alt, ok = p.v.(bool)
		case "app_cursor":
			m.DECCKM, ok = p.v.(bool)
		case "bracketed_paste":
			m.BracketedPaste, ok = p.v.(bool)
		case "mouse_tracking":
			var n int64
			n, ok = p.v.(int64)
			m.MouseTrack = int(n)
		case "mouse_encoding":
			var n int64
			n, ok = p.v.(int64)
			m.MouseEnc = int(n)
		}
		if !ok {
			t.Fatalf("%s: state key %q = %v", id, p.k, p.v)
		}
		seen++
	}
	if seen != 5 {
		t.Fatalf("%s: %d state keys, want 5", id, seen)
	}
	return m
}

// replayPython feeds a Python case op by op and returns the first difference ("" when every op agrees).
func replayPython(t *testing.T, c vcase) string {
	t.Helper()
	s := New(Options{TitleMax: maxPayload})
	var want string
	var have bool
	for i, op := range pyOps(t, c) {
		s.Feed(op.chunk)
		switch c.Entry {
		case "titles.OscTitleSniffer":
			if len(op.state) != 0 {
				t.Fatalf("%s op %d: titles state %v", c.ID, i, op.state)
			}
			// feed returns the new title when it changed, else None; the current title is the last one returned.
			switch r := op.result.(type) {
			case string:
				want, have = r, true
			case nil:
			default:
				t.Fatalf("%s op %d: result %T", c.ID, i, r)
			}
			if got, ok := s.Title(); got != want || ok != have {
				return fmt.Sprintf("op %d: Title() = %q, %v; Python %q, %v", i, got, ok, want, have)
			}
		default:
			wm := pyModes(t, c.ID, op.state)
			if c.Entry == "altscreen.AltScreenSniffer" {
				if r, ok := op.result.(bool); !ok || r != wm.Alt {
					t.Fatalf("%s op %d: AltScreenSniffer.feed returned %v, state alt %v", c.ID, i, op.result, wm.Alt)
				}
			} else if op.result != nil {
				t.Fatalf("%s op %d: feed returned %v", c.ID, i, op.result)
			}
			if got := s.Modes(); got != wm {
				return fmt.Sprintf("op %d: Modes() = %+v; Python %+v", i, got, wm)
			}
		}
	}
	return ""
}

// TestPythonVectors replays every Python sniffer case (altscreen 118, titles 80, the 4 synthetic cap cases) op by op
// and compares the modes, or the title, after every feed.
func TestPythonVectors(t *testing.T) {
	counts := map[string]int{}
	for _, fam := range []string{"altscreen", "titles", "synthetic"} {
		for _, c := range streamCases(t, loadVectors(t, "python/"+fam+".json", fam)) {
			counts[fam]++
			diff := replayPython(t, c)
			why, listed := pythonDivergences[c.ID]
			switch {
			case listed && diff == "":
				t.Errorf("%s: listed as a divergence (%s) but agrees with Python", c.ID, why)
			case !listed && diff != "":
				t.Errorf("%s: %s", c.ID, diff)
			case listed:
				t.Logf("divergence %s: %s", c.ID, diff)
			}
		}
	}
	if counts["altscreen"] != 118 || counts["titles"] != 80 || counts["synthetic"] != 4 {
		t.Errorf("case counts %v, want altscreen 118, titles 80, synthetic 4", counts)
	}
}

// TestPythonCapCases names the four synthetic cases: the caps this package keeps from the Python reference.
func TestPythonCapCases(t *testing.T) {
	var ids []string
	for _, c := range streamCases(t, loadVectors(t, "python/synthetic.json", "synthetic")) {
		ids = append(ids, c.ID)
	}
	sort.Strings(ids)
	want := "synthetic/mode-carry-cap#0 synthetic/mode-carry-cap#1 synthetic/title-payload-cap#0 synthetic/title-payload-cap#1"
	if got := strings.Join(ids, " "); got != want {
		t.Errorf("synthetic sniffer cases %q, want %q", got, want)
	}
}
