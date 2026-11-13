package screen

import (
	"fmt"

	"github.com/JohnConnorNPC/go-term/sniff"
)

// SnapshotVersion is the Snapshot.Version Restore accepts. It is go-term's
// own format version: no xterm.js serialize addon format is read or written.
const SnapshotVersion = 1

// Preamble starts every replay, as Studio's raw snapshot (raw.py PREAMBLE):
// reset the attributes, clear the screen, home the cursor. It is not a hard
// reset (ESC c) and not ESC [3J, so it keeps modes and scrollback.
const Preamble = "\x1b[0m\x1b[2J\x1b[H"

// Snapshot is what Restore rebuilds an instance from: Studio's ring of raw
// pty output plus the mode and title state the agent's sniffer reported for
// the whole stream, which covers a mode set that scrolled out of the ring.
// It is a Go value, not a byte encoding; Studio owns its wire form.
type Snapshot struct {
	// Version must be SnapshotVersion.
	Version int
	// Cols and Rows are the terminal size, within New's limits.
	Cols, Rows int
	// Modes is the sniffer's report (sniff.Sniffer.Modes).
	Modes sniff.Modes
	// Title is the sniffer's title, used when TitleSet (sniff.Sniffer.Title's
	// ok). The sniffer cuts a title to its TitleMax (1024 bytes by default),
	// so a longer title restores cut.
	Title    string
	TitleSet bool
	// Ring is the replayed bytes as Studio holds them: no trim of a head cut
	// mid-sequence and no strip of queries; that is Studio's.
	Ring []byte
}

// ReplayBytes is what Restore writes: Preamble, ring, Cancel, then
// sniff.Postamble(m), which re-asserts bracketed paste, DECCKM and the mouse
// modes after the ring.
func ReplayBytes(m sniff.Modes, ring []byte) []byte {
	pa := sniff.Postamble(m)
	b := make([]byte, 0, len(Preamble)+len(ring)+len(Cancel)+len(pa))
	b = append(b, Preamble...)
	b = append(b, ring...)
	b = append(b, Cancel...)
	return append(b, pa...)
}

// Cancel (CAN) separates the ring from the postamble. A ring cut inside a
// sequence leaves it open; xterm.js abandons an open CSI, OSC, DCS or APC on
// CAN without completing it, where the postamble's ESC would end an OSC
// successfully and commit a fragment title. After a whole ring CAN does
// nothing: no handler takes it.
const Cancel = "\x18"

// altOn and altOff switch to and from the alt buffer as Studio's programs do.
const altOn, altOff = "\x1b[?1049h", "\x1b[?1049l"

// Restore returns a new instance of s.Cols x s.Rows (sized by c) fed
// ReplayBytes(s.Modes, s.Ring), whose Read reports s.Modes and, when
// s.TitleSet, s.Title. The postamble does not carry the alt screen, so
// Restore aligns it after the replay: when s.Modes.Alt but the replay ended
// on the normal buffer (the ring holds no switch to the alt buffer: it was
// all drawn there), the replay is fed again to a new instance after
// CSI ? 1049 h; when the replay ended on the alt buffer but s.Modes.Alt is
// false, CSI ? 1049 l follows it. A replay that already agrees is never
// changed: CSI ? 1049 h on the alt buffer does not clear it, so a blind
// prefix would leave the ring's pre-switch output on the alt grid. A ring
// that contradicts s.Modes.Alt even so (it leaves the alt buffer itself) gets
// CSI ? 1049 h last, an empty alt grid. An error from New or a failed
// instance (ErrFailed) gives no instance.
func Restore(s Snapshot, c Config) (*Instance, error) {
	if s.Version != SnapshotVersion {
		return nil, fmt.Errorf("screen: snapshot version %d, want %d", s.Version, SnapshotVersion)
	}
	replay := ReplayBytes(s.Modes, s.Ring)
	i, err := restoreFeed(s, c, "", replay)
	if err != nil {
		return nil, err
	}
	alt, err := i.alt()
	switch {
	case err != nil:
		return nil, err
	case s.Modes.Alt && !alt:
		if i, err = restoreFeed(s, c, altOn, replay); err != nil {
			return nil, err
		}
		alt, err = i.alt()
		if err == nil && !alt {
			// The ring itself leaves the alt buffer (CSI ? 1049 l, ESC c),
			// which the reported state contradicts: enter it last.
			_, err = i.Write([]byte(altOn))
		}
	case !s.Modes.Alt && alt:
		_, err = i.Write([]byte(altOff))
	}
	if err == nil && s.TitleSet {
		err = i.do(func() { i.t.InputHandler().SetTitle(s.Title) })
	}
	if err != nil {
		return nil, err
	}
	return i, nil
}

// restoreFeed is a new instance fed prefix then replay.
func restoreFeed(s Snapshot, c Config, prefix string, replay []byte) (*Instance, error) {
	i, err := New(s.Cols, s.Rows, c)
	if err != nil {
		return nil, err
	}
	if prefix != "" {
		if _, err := i.Write([]byte(prefix)); err != nil {
			return nil, err
		}
	}
	if _, err := i.Write(replay); err != nil {
		return nil, err
	}
	return i, nil
}

// alt says whether the alt buffer is active.
func (i *Instance) alt() (on bool, err error) {
	err = i.do(func() { on = i.t.Buffers().Active() == i.t.Buffers().Alt() })
	return on, err
}
