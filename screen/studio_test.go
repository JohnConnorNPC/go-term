package screen

import (
	"bytes"
	"encoding/hex"
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
