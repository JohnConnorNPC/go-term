package sniff

import "strconv"

// MouseTrackModes and MouseEncModes are the members of the two mouse groups, in the order Postamble writes them.
var (
	MouseTrackModes = [...]int{9, 1000, 1002, 1003}
	MouseEncModes   = [...]int{1006, 1016}
)

// Postamble returns the bytes that re-assert m on a viewer, to append after a snapshot's grid: CSI ? 2004 h|l,
// CSI ? 1 h|l, then for each mouse group (tracking, then encoding) a reset of every member that is not active followed
// by the set of the active one. The order matters: xterm.js maps a reset of any member to the group's default, so a
// reset written after the set would undo it. A value outside its group counts as 0 (nothing set). The alternate
// screen is not re-asserted: the snapshot's grid replay carries it. The result is idempotent for a viewer already in m.
func Postamble(m Modes) []byte {
	b := make([]byte, 0, 64)
	b = mode(b, 2004, m.BracketedPaste)
	b = mode(b, 1, m.DECCKM)
	b = group(b, MouseTrackModes[:], m.MouseTrack)
	return group(b, MouseEncModes[:], m.MouseEnc)
}

func mode(b []byte, n int, on bool) []byte {
	b = append(b, "\x1b[?"...)
	b = strconv.AppendInt(b, int64(n), 10)
	if on {
		return append(b, 'h')
	}
	return append(b, 'l')
}

func group(b []byte, members []int, active int) []byte {
	in := false
	for _, n := range members {
		if n == active {
			in = true
		} else {
			b = mode(b, n, false)
		}
	}
	if in {
		b = mode(b, active, true)
	}
	return b
}
