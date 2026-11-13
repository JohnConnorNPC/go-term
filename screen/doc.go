// Package screen is go-term's typed read API for Browserland Studio: the grid
// text, cursor and modes of a vt.Terminal, attribute runs in pyte's colour
// vocabulary, the blake2b-128 content and stable hashes, row deltas against
// recent frames, and the scrollback caps. Studio builds its JSON from these
// values.
//
// The functions read a *vt.Terminal and take no lock: the caller serialises
// them with Write (FrameCache alone locks itself).
package screen
