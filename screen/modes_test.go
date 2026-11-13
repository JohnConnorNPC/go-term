package screen

import (
	"testing"

	"github.com/JohnConnorNPC/go-term/sniff"
)

// TestReadModes: Read reports bracketed paste, the mouse tracking mode and
// encoding (numbered as sniff.Modes), DECCKM, the alt screen and the title,
// step by step, and after every step Frame.Modes equals what the sniffer
// reports for the same stream.
func TestReadModes(t *testing.T) {
	steps := []struct {
		in    string
		want  sniff.Modes
		title string
	}{
		{"", sniff.Modes{}, ""},
		{"\x1b[?2004h", sniff.Modes{BracketedPaste: true}, ""},
		{"\x1b[?1h", sniff.Modes{BracketedPaste: true, DECCKM: true}, ""},
		{"\x1b[?9h", sniff.Modes{BracketedPaste: true, DECCKM: true, MouseTrack: 9}, ""},
		{"\x1b[?1000h", sniff.Modes{BracketedPaste: true, DECCKM: true, MouseTrack: 1000}, ""},
		{"\x1b[?1002h", sniff.Modes{BracketedPaste: true, DECCKM: true, MouseTrack: 1002}, ""},
		{"\x1b[?1003h\x1b[?1006h", sniff.Modes{BracketedPaste: true, DECCKM: true, MouseTrack: 1003, MouseEnc: 1006}, ""},
		{"\x1b[?1016h", sniff.Modes{BracketedPaste: true, DECCKM: true, MouseTrack: 1003, MouseEnc: 1016}, ""},
		// A reset of any member of a group clears the group.
		{"\x1b[?9l\x1b[?1006l", sniff.Modes{BracketedPaste: true, DECCKM: true}, ""},
		{"\x1b[?1049h\x1b]2;one\x07", sniff.Modes{Alt: true, BracketedPaste: true, DECCKM: true}, "one"},
		{"\x1b[?2004l\x1b[?1l\x1b]0;two\x1b\\", sniff.Modes{Alt: true}, "two"},
		{"\x1b[?1049l\x1b[?1000h\x1b[?1016h", sniff.Modes{MouseTrack: 1000, MouseEnc: 1016}, "two"},
		{"\x1bc", sniff.Modes{}, "two"}, // RIS resets the modes, not the title
	}
	term := newTerm(10, 3)
	sn := sniff.New(sniff.Options{})
	for i, s := range steps {
		term.WriteString(s.in)
		sn.Feed([]byte(s.in))
		f := Read(term, Request{})
		if got := f.Modes(); got != s.want {
			t.Errorf("step %d %q: modes %+v, want %+v", i, s.in, got, s.want)
		}
		if got := f.Modes(); got != sn.Modes() {
			t.Errorf("step %d %q: modes %+v, sniffer %+v", i, s.in, got, sn.Modes())
		}
		if f.Title != s.title {
			t.Errorf("step %d %q: title %q, want %q", i, s.in, f.Title, s.title)
		}
		if f.AltScreen != f.Modes().Alt || f.AppCursor != f.Modes().DECCKM {
			t.Errorf("step %d: Modes disagrees with AltScreen/AppCursor", i)
		}
	}
}
