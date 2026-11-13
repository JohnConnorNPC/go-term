package screen

import (
	"fmt"
	"reflect"
	"testing"
)

// splitStreams cover each kind of sequence a chunk boundary can cut.
var splitStreams = map[string]string{
	"csi":     "\x1b[1;31mred\x1b[0m \x1b[5;10Hxy\x1b[2;4r\x1b[?1049h alt\x1b[?1049l\x1b[38;2;1;2;3mT\x1b[48;5;200mB\x1b[m\x1b[?1h\x1b[?25l\x1b[3A\x1b[2Kz",
	"osc":     "\x1b]0;title one\x07a\x1b]2;two\x1b\\b\x1b]8;;http://x/\x1b\\link\x1b]8;;\x1b\\\x1b]10;?\x07\x1b]4;1;?\x1b\\c",
	"dcs":     "\x1bP$qm\x1b\\\x1bP$q\"q\x1b\\\x1bP$qr\x1b\\\x1bPzz\x1b\\x",
	"apc":     "a\x1b_Gf=1;abc\x1b\\b\x1b_bel\x07still\x1b\\c",
	"utf8":    "é ü ñ ∑ € ½ 𝄞 🎉 \xc3\xa9\xe2\x82\xac\xf0\x9f\x8e\x89 \xff\xc3( end",
	"wide":    "日本語 ｱ漢字\r\n1234567890123456789日x\r\n123456789012345678日y 全角",
	"combine": "e\u0301 a\u0308\u0323 n\u0303x\r\n1234567890123456789e\u0301z",
	"cancel":  "\x1b[1\x18x\x1b]2;t\x1aY\x1bP$qm\x18Z\x1b[31\x1b[32mG",
	// U+2014, U+200B ZWSP, U+200D ZWJ, U+20000 (cjk-ext-b): each has a held
	// continuation byte 0x80 that a write can end on. @xterm/headless 6.0.0
	// loses the character there (its goldens wide/cjk-ext-b, wide/zwj-family,
	// wide/zwsp are split.sensitive); vt does not (#14).
	"quirk":  "a—b\u200bc👨\u200d👩d𠀀e",
	"scroll": "1\r\n2\r\n3\r\n4\r\n5\r\n6\r\n7\r\n8\x1b[2;3r\x1b[3;1H\x1bM\x1bD\x1b[S\x1b[T\x1b[r\x1b[?1049h\x1b[2J",
}

type splitState struct {
	Frame          Frame
	Title, Replies string
}

// feed writes chunks to a fresh 20x6 instance and returns its state.
func feed(t *testing.T, chunks ...string) splitState {
	t.Helper()
	i, err := New(20, 6, Config{})
	if err != nil {
		t.Fatal(err)
	}
	var s splitState
	i.t.OnTitleChange = func(v string) { s.Title = v }
	i.t.OnData = func(v string) { s.Replies += v }
	for _, c := range chunks {
		if _, err := i.Write([]byte(c)); err != nil {
			t.Fatal(err)
		}
	}
	if s.Frame, err = i.Read(Request{View: ViewScrollback, Lines: 1000, Attrs: true}); err != nil {
		t.Fatal(err)
	}
	return s
}

// TestInstanceSplit [#7]: each stream written whole and split at every offset
// into two writes, and one byte per write, gives the same state: the frame
// (history, grid, cursor, modes, attribute runs, hashes), the title and the
// replies, with no exception (the quirk stream included).
func TestInstanceSplit(t *testing.T) {
	n := 0
	for name, s := range splitStreams {
		whole := feed(t, s)
		variants := map[string][]string{}
		bytes := make([]string, len(s))
		for k := 0; k < len(s); k++ { // bytes, not runes: cut inside UTF-8 too
			bytes[k] = s[k : k+1]
			if k > 0 {
				variants[fmt.Sprintf("at %d", k)] = []string{s[:k], s[k:]}
			}
		}
		variants["bytewise"] = bytes
		for v, chunks := range variants {
			n++
			if got := feed(t, chunks...); !reflect.DeepEqual(got, whole) {
				t.Errorf("%s split %s:\n got  %q %q %q\n want %q %q %q", name, v, got.Frame.Lines, got.Title, got.Replies, whole.Frame.Lines, whole.Title, whole.Replies)
			}
		}
	}
	t.Logf("%d streams, %d split variants equal to the whole write", len(splitStreams), n)
}
