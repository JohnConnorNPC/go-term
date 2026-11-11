package sniff

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"unicode/utf8"
	"unsafe"
)

// streams returns every input the other tests use, each as the whole byte stream: the Python vector cases' chunks
// joined, and every conformance golden's write steps joined.
func streams(t testing.TB) [][]byte {
	var out [][]byte
	for _, fam := range []string{"altscreen", "titles", "synthetic"} {
		b, err := os.ReadFile("testdata/python/" + fam + ".json")
		if err != nil {
			t.Fatal(err)
		}
		var f vfile
		if err := json.Unmarshal(b, &f); err != nil {
			t.Fatal(err)
		}
		for _, c := range f.Cases {
			if c.Kind != "stream" {
				continue
			}
			var all []byte
			for _, s := range c.Replay.Chunks {
				p, err := base64.StdEncoding.DecodeString(s)
				if err != nil {
					t.Fatal(err)
				}
				all = append(all, p...)
			}
			out = append(out, all)
		}
	}
	files, _ := filepath.Glob("../conformance/golden/x53/*/*.json")
	for _, fn := range files {
		b, err := os.ReadFile(fn)
		if err != nil {
			t.Fatal(err)
		}
		var g struct{ Steps []goldenStep }
		if err := json.Unmarshal(b, &g); err != nil {
			t.Fatal(err)
		}
		var all []byte
		for _, st := range g.Steps {
			p, _ := hex.DecodeString(st.Op.In)
			all = append(all, p...)
		}
		out = append(out, all)
	}
	if len(out) < 202+len(files) || len(files) == 0 {
		t.Fatalf("%d streams, %d goldens (unmeasured)", len(out), len(files))
	}
	return out
}

// whole feeds p in one call and returns the sniffer, whose whole state (modes, title, partial sequence, decoder) the
// split runs must equal.
func whole(p []byte) *Sniffer {
	s := New(Options{})
	s.Feed(p)
	return s
}

// TestSplitEveryOffset: every stream fed in two chunks split at every byte offset, and fed byte by byte, ends in the
// same state as fed whole (the struct compared, so a later feed cannot tell them apart either). A stream over 8 KiB
// (4 of the Python property-test streams, 30-70 KB) is split at its first and last 256 offsets and 512 spread
// between; every offset of those would take 13.8 G fed bytes.
func TestSplitEveryOffset(t *testing.T) {
	n, sampled := 0, 0
	for _, p := range streams(t) {
		w := whole(p)
		offsets := splitOffsets(len(p))
		if len(offsets) < len(p)+1 {
			sampled++
		}
		for _, k := range offsets {
			s := New(Options{})
			s.Feed(p[:k])
			s.Feed(p[k:])
			if *s != *w {
				t.Fatalf("stream %q split at %d: %+v %q; whole %+v %q", truncate(p), k, s.Modes(), s.title, w.Modes(), w.title)
			}
			n++
		}
		s := New(Options{})
		for i := range p {
			s.Feed(p[i : i+1])
		}
		if *s != *w {
			t.Fatalf("stream %q byte by byte differs", truncate(p))
		}
	}
	t.Logf("%d splits, %d streams sampled", n, sampled)
	if sampled > 4 {
		t.Errorf("%d streams sampled, want at most 4", sampled)
	}
}

func splitOffsets(n int) []int {
	var out []int
	for k := 0; k <= n; k++ {
		if n <= 8<<10 || k < 256 || k > n-256 || k%(n/512) == 0 {
			out = append(out, k)
		}
	}
	return out
}

func truncate(p []byte) []byte {
	if len(p) > 80 {
		return p[:80]
	}
	return p
}

// feedMiB feeds the prefix, then n MiB of fill in 64 KiB chunks from one buffer, and returns the bytes the heap
// grew by (TotalAlloc, so freed garbage counts too).
func feedMiB(s *Sniffer, prefix string, fill byte, n int) uint64 {
	chunk := bytes.Repeat([]byte{fill}, 64<<10)
	var a, b runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&a)
	s.Feed([]byte(prefix))
	for i := 0; i < n*16; i++ {
		s.Feed(chunk)
	}
	runtime.ReadMemStats(&b)
	return b.TotalAlloc - a.TotalAlloc
}

// TestBounds: 64 MiB inside an unterminated OSC title, a 64 MiB title, an unterminated CSI, DCS and APC, and plain
// text: Feed allocates nothing per call (the prefix may allocate once), and the state stays the same size.
func TestBounds(t *testing.T) {
	if size := unsafe.Sizeof(Sniffer{}); size > 8<<10 {
		t.Errorf("Sniffer is %d bytes", size)
	}
	for _, tc := range []struct {
		name, prefix string
		fill         byte
		end          string
		title        string
		have         bool
	}{
		{"unterminated OSC 0", "\x1b]0;", 't', "", "", false},
		{"64 MiB title", "\x1b]0;", 't', "\x07", "", false},
		{"64 MiB title after a title", "\x1b]2;old\x07\x1b]0;", 't', "\x07", "old", true},
		{"unterminated OSC 99", "\x1b]99;", 'x', "", "", false},
		{"unterminated CSI", "\x1b[?", '1', "h", "", false},
		{"unterminated DCS", "\x1bP1$q", 'x', "", "", false},
		{"unterminated APC", "\x1b_", 'x', "", "", false},
		{"text", "", 'x', "", "", false},
		{"utf8 text", "", 0xC3, "", "", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := New(Options{})
			grew := feedMiB(s, tc.prefix, tc.fill, 64)
			s.Feed([]byte(tc.end))
			if grew > 4<<10 {
				t.Errorf("the heap grew by %d bytes over 64 MiB", grew)
			}
			if got, ok := s.Title(); got != tc.title || ok != tc.have {
				t.Errorf("Title() = %.40q %v, want %q %v", got, ok, tc.title, tc.have)
			}
			if m := s.Modes(); m != (Modes{}) {
				t.Errorf("Modes() = %+v", m)
			}
		})
	}
	s := New(Options{})
	chunk := bytes.Repeat([]byte("\x1b]0;tt"), 1<<10)
	if n := testing.AllocsPerRun(100, func() { s.Feed(chunk) }); n != 0 {
		t.Errorf("Feed of 1024 unchanged titles: %v allocs per call", n)
	}
}

// FuzzFeed: any bytes, split anywhere, give the whole-feed state; no panic; the reported modes are in their groups
// and the title is valid UTF-8 within TitleMax; the postamble of the modes re-asserts them on a fresh sniffer.
func FuzzFeed(f *testing.F) {
	for _, p := range [][]byte{
		[]byte("\x1b[?1049;1;2004;1003;1006h\x1b]0;title\x07"),
		[]byte("\x1b]0;junk\x1bx\x1bc\x1b[!p"),
		[]byte("\xc2\x9b?1h\xc2\x9d2;t\xc2\x9c\x1bP$q\x1b\\"),
		[]byte("\x1b]0;h\xc3\xa9\xe2\x82\xac\xff\x07\x1b[?1:2;1016h"),
	} {
		f.Add(p, []byte{3, 1, 7})
	}
	f.Fuzz(func(t *testing.T, p, cuts []byte) {
		w := whole(p)
		s := New(Options{})
		rest := p
		for _, c := range cuts {
			k := int(c) % (len(rest) + 1)
			s.Feed(rest[:k])
			rest = rest[k:]
		}
		s.Feed(rest)
		if *s != *w {
			t.Fatalf("split %v: %+v %q; whole %+v %q", cuts, s.Modes(), s.title, w.Modes(), w.title)
		}
		m := w.Modes()
		inGroup := func(v int, g []int) bool {
			for _, x := range g {
				if v == x {
					return true
				}
			}
			return v == 0
		}
		if !inGroup(m.MouseTrack, MouseTrackModes[:]) || !inGroup(m.MouseEnc, MouseEncModes[:]) {
			t.Fatalf("modes out of group: %+v", m)
		}
		if title, _ := w.Title(); len(title) > DefaultTitleMax || !utf8.ValidString(title) {
			t.Fatalf("title %d bytes, valid %v", len(title), utf8.ValidString(title))
		}
		v := New(Options{})
		v.Feed(Postamble(m))
		m.Alt = false
		if got := v.Modes(); got != m {
			t.Fatalf("Postamble(%+v) re-asserts %+v", m, got)
		}
	})
}
