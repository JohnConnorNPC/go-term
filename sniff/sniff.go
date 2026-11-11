// Package sniff watches a terminal's output stream for the state a snapshot cannot carry: the DEC private modes a
// viewer must re-assert (alternate screen, application cursor keys, bracketed paste, mouse tracking and its encoding)
// and the window title (OSC 0 and 2). Postamble turns reported modes into the bytes appended to a snapshot.
//
// The sniffer is a small parser in the shape of xterm.js's: the bytes are decoded as UTF-8 the way xterm.js's
// Utf8ToUtf32 decodes them (invalid bytes are dropped), and the code points run through the subset of xterm.js's
// VT500 transition table that decides where an escape sequence starts and ends (ESC, CSI, OSC, DCS, SOS/PM/APC,
// the C1 forms, CAN and SUB). So a mode sequence inside an OSC or DCS payload does not count, and a bare ESC inside an
// OSC title ends it successfully (the title is committed), as in xterm.js. The state machine carries every partial
// sequence and code point across Feed calls, so any split of the stream gives the same result.
//
// Acted on, as xterm.js acts on them:
//   - CSI ? Pm h / l (DECSET / DECRST), no intermediates; Pm as xterm.js's Params: at most 32 params, each clamped to
//     2^31-1, sub-parameters after ':' ignored. Modes 1 (DECCKM), 47/1047/1049 (alternate screen), 2004 (bracketed
//     paste), 9/1000/1002/1003 (mouse tracking) and 1006/1016 (mouse encoding), applied left to right. A reset of
//     any member of a mouse group sets the group to 0 (xterm.js keeps one slot per group).
//   - ESC c (RIS): every mode off, normal buffer. The title is kept.
//   - CSI ! p (DECSTR): DECCKM and bracketed paste off (xterm.js's softReset); mouse and alternate screen kept.
//   - OSC 0 and OSC 2 (id read as a number, so "00" is 0): the payload is the title. BEL, ST and ESC end it
//     successfully; CAN, SUB and the C1 controls abandon it. C0 controls inside the payload are dropped; DEL (0x7F) is
//     kept, as xterm.js keeps it.
//
// Ours, not xterm.js's (memory bounds, the same as the Python reference keeps):
//   - an OSC payload over 4096 bytes is abandoned (no title from it);
//   - a CSI sequence with more than 68 bytes before its final byte is dropped. The Python reference carries at most
//     64 bytes of a split sequence, which depends on where the chunk ends; this is the split-independent bound that
//     agrees with both of its pinned cap cases (a 64-byte carry completed by "1049h" is kept, 65 is dropped);
//   - Title is cut to Options.TitleMax bytes (default 1024) on a rune boundary.
package sniff

import "unicode/utf8"

// Modes is the DEC private-mode state a viewer must re-assert. MouseTrack is 0 (off), 9, 1000, 1002 or 1003;
// MouseEnc is 0 (the default encoding), 1006 or 1016.
type Modes struct {
	Alt, DECCKM, BracketedPaste bool
	MouseTrack, MouseEnc        int
}

// Options configures a Sniffer. TitleMax is the longest title reported, in bytes, cut on a rune boundary;
// 0 or less means DefaultTitleMax.
type Options struct {
	TitleMax int
}

// DefaultTitleMax is the title cap when Options.TitleMax is not set.
const DefaultTitleMax = 1024

const (
	maxPayload = 4096       // OSC payload bytes kept; more abandons the sequence
	maxSeq     = 68         // CSI bytes before the final byte; more drops the sequence
	maxParams  = 32         // xterm.js Params.maxLength
	maxParam   = 0x7FFFFFFF // xterm.js Params MAX_VALUE
	maxOscID   = 1 << 20    // an OSC id is clamped here; only 0 and 2 matter
)

type state uint8

const (
	ground state = iota
	escape
	escInter
	csiEntry
	csiParam
	csiInter
	csiIgnore
	oscString
	dcsEntry
	dcsParam
	dcsInter
	dcsIgnore
	dcsPass
	sosPmApc
)

type oscState uint8

const (
	oscID oscState = iota
	oscPayload
	oscAbort
)

// Sniffer tracks one terminal's modes and title. It is not safe for concurrent use: the caller serialises Feed,
// Modes, Title and Reset.
type Sniffer struct {
	titleMax int

	modes Modes
	title string
	have  bool

	// UTF-8 decoder: a code point in progress.
	cp   rune
	need int // continuation bytes still expected
	ulen int // bytes of the code point so far

	st      state
	collect uint32 // prefix and intermediates, shifted in as xterm.js does
	seqLen  int    // bytes of the current ESC/CSI sequence so far

	params [maxParams]int32
	nparam int
	reject bool // a 33rd param was started: later digits are ignored
	sub    bool // digits go to a sub-parameter (ignored)

	osc     oscState
	oscID   int // -1 until a digit
	capture bool
	hit     bool // the payload passed maxPayload
	plen    int
	payload [maxPayload]byte
}

// New returns a sniffer in the power-on state: every mode off, no title.
func New(o Options) *Sniffer {
	s := &Sniffer{}
	s.titleMax = o.TitleMax
	if s.titleMax <= 0 {
		s.titleMax = DefaultTitleMax
	}
	s.Reset()
	return s
}

// Reset puts the sniffer back in the power-on state (a new PTY): every mode off, no title, no partial sequence.
func (s *Sniffer) Reset() {
	*s = Sniffer{titleMax: s.titleMax, oscID: -1}
}

// Modes returns the current modes.
func (s *Sniffer) Modes() Modes { return s.modes }

// Title returns the last OSC 0 or 2 title, cut to TitleMax bytes on a rune boundary; ok is false before the first.
func (s *Sniffer) Title() (title string, ok bool) { return s.title, s.have }

// Feed consumes the next bytes of the stream. p is not retained. Memory does not grow with the input: a partial
// sequence keeps at most maxPayload bytes.
func (s *Sniffer) Feed(p []byte) {
	for _, b := range p {
		s.decode(b)
	}
}

// decode is xterm.js's Utf8ToUtf32 over the concatenated stream: a wrong continuation byte drops the code point in
// progress and starts over at that byte; an overlong form, a surrogate or a value past U+10FFFF is dropped whole;
// any other stray byte is dropped.
func (s *Sniffer) decode(b byte) {
	if s.need > 0 {
		if b&0xC0 == 0x80 {
			s.cp = s.cp<<6 | rune(b&0x3F)
			s.ulen++
			s.need--
			if s.need > 0 {
				return
			}
			c := s.cp
			switch s.ulen {
			case 2:
				if c < 0x80 {
					return
				}
			case 3:
				if c < 0x800 || (c >= 0xD800 && c <= 0xDFFF) {
					return
				}
			case 4:
				if c < 0x10000 || c > 0x10FFFF {
					return
				}
			}
			s.step(c, s.ulen)
			return
		}
		s.need = 0
	}
	switch {
	case b < 0x80:
		s.step(rune(b), 1)
	case b&0xE0 == 0xC0:
		s.cp, s.need, s.ulen = rune(b&0x1F), 1, 1
	case b&0xF0 == 0xE0:
		s.cp, s.need, s.ulen = rune(b&0x0F), 2, 1
	case b&0xF8 == 0xF0:
		s.cp, s.need, s.ulen = rune(b&0x07), 3, 1
	}
}

// executable reports xterm.js's EXECUTABLES: C0 but CAN, SUB and ESC.
func executable(c rune) bool {
	return c <= 0x17 || c == 0x19 || (c >= 0x1C && c <= 0x1F)
}

// clear is xterm.js's CLEAR action: params and collect start over.
func (s *Sniffer) clear() {
	s.collect = 0
	s.params[0] = 0
	s.nparam = 1
	s.reject = false
	s.sub = false
}

// step runs one code point through the transition table. n is its UTF-8 length.
func (s *Sniffer) step(c rune, n int) {
	// xterm.js's "anywhere" rules, with OSC_STRING's and DCS_PASSTHROUGH's own entries for the same codes.
	switch {
	case c == 0x18 || c == 0x1A || c == 0x99 || c == 0x9A || (c >= 0x80 && c <= 0x8F) || (c >= 0x91 && c <= 0x97):
		s.st = ground
		return
	case c == 0x9C:
		if s.st == oscString {
			s.oscEnd()
		}
		s.st = ground
		return
	case c == 0x1B:
		if s.st == oscString {
			s.oscEnd() // xterm.js's OSC_END with success; DCS_UNHOOK, OSC_END and CLEAR all clear
		}
		s.clear()
		s.st = escape
		s.seqLen = n
		return
	case c == 0x9D:
		s.oscStart()
		return
	case c == 0x98 || c == 0x9E || c == 0x9F:
		s.st = sosPmApc
		return
	case c == 0x9B:
		s.clear()
		s.st = csiEntry
		s.seqLen = n
		return
	case c == 0x90:
		s.clear()
		s.st = dcsEntry
		return
	}
	na := c >= 0xA0 // xterm.js's NON_ASCII_PRINTABLE class
	switch s.st {
	case ground:
	case escape:
		s.seqLen += n
		switch {
		case na:
			s.st = ground
		case executable(c), c == 0x7F:
		case c >= 0x20 && c <= 0x2F:
			s.collect = s.collect<<8 | uint32(c)
			s.st = escInter
		case c == '[':
			s.clear()
			s.st = csiEntry
		case c == ']':
			s.oscStart()
		case c == 'P':
			s.clear()
			s.st = dcsEntry
		case c == 'X' || c == '^' || c == '_':
			s.st = sosPmApc
		default: // 0x30-0x7E but the above: ESC_DISPATCH
			s.escDispatch(c)
			s.st = ground
		}
	case escInter:
		switch {
		case na:
			s.st = ground
		case executable(c), c == 0x7F:
		case c >= 0x20 && c <= 0x2F:
			s.collect = s.collect<<8 | uint32(c)
		default:
			s.escDispatch(c)
			s.st = ground
		}
	case csiEntry, csiParam, csiInter:
		s.seqLen += n
		if na {
			s.st = ground
			return
		}
		if c >= 0x40 && c <= 0x7E {
			s.csiDispatch(c)
			s.st = ground
			return
		}
		if s.seqLen > maxSeq {
			s.st = csiIgnore
			return
		}
		switch {
		case executable(c), c == 0x7F:
		case c >= 0x20 && c <= 0x2F:
			s.collect = s.collect<<8 | uint32(c)
			s.st = csiInter
		case s.st == csiInter: // 0x30-0x3F after an intermediate
			s.st = csiIgnore
		case c <= 0x3B:
			s.param(c)
			s.st = csiParam
		case s.st == csiEntry: // 0x3C-0x3F: a prefix
			s.collect = s.collect<<8 | uint32(c)
			s.st = csiParam
		default: // a prefix after a param
			s.st = csiIgnore
		}
	case csiIgnore:
		if !na && c >= 0x40 && c <= 0x7E {
			s.st = ground
		}
	case oscString:
		switch {
		case c == 0x07:
			s.oscEnd()
			s.st = ground
		case na, c >= 0x20:
			s.oscPut(c)
		}
	case dcsEntry, dcsParam, dcsInter:
		switch {
		case na:
			s.st = ground
		case executable(c), c == 0x7F:
		case c >= 0x40:
			s.st = dcsPass
		case c >= 0x20 && c <= 0x2F:
			s.collect = s.collect<<8 | uint32(c)
			s.st = dcsInter
		case s.st == dcsInter:
			s.st = dcsIgnore
		case c <= 0x3B:
			s.param(c)
			s.st = dcsParam
		case s.st == dcsEntry:
			s.collect = s.collect<<8 | uint32(c)
			s.st = dcsParam
		default:
			s.st = dcsIgnore
		}
	case dcsIgnore, dcsPass:
	case sosPmApc:
		if na {
			s.st = ground
		}
	}
}

// param is xterm.js's PARAM action on one of 0-9 : ;.
func (s *Sniffer) param(c rune) {
	switch c {
	case ';':
		s.sub = false
		if s.nparam >= maxParams {
			s.reject = true
			return
		}
		s.params[s.nparam] = 0
		s.nparam++
	case ':':
		s.sub = true
	default:
		if s.reject || s.sub || s.nparam == 0 {
			return
		}
		v := int64(s.params[s.nparam-1])*10 + int64(c-'0')
		if v > maxParam {
			v = maxParam
		}
		s.params[s.nparam-1] = int32(v)
	}
}

func (s *Sniffer) escDispatch(final rune) {
	if s.collect == 0 && final == 'c' { // RIS
		s.modes = Modes{}
	}
}

func (s *Sniffer) csiDispatch(final rune) {
	switch {
	case s.collect == '?' && (final == 'h' || final == 'l'):
		on := final == 'h'
		for _, p := range s.params[:s.nparam] {
			switch p {
			case 1:
				s.modes.DECCKM = on
			case 47, 1047, 1049:
				s.modes.Alt = on
			case 2004:
				s.modes.BracketedPaste = on
			case 9, 1000, 1002, 1003:
				s.modes.MouseTrack = 0
				if on {
					s.modes.MouseTrack = int(p)
				}
			case 1006, 1016:
				s.modes.MouseEnc = 0
				if on {
					s.modes.MouseEnc = int(p)
				}
			}
		}
	case s.collect == '!' && final == 'p': // DECSTR
		s.modes.DECCKM = false
		s.modes.BracketedPaste = false
	}
}

// oscStart is xterm.js's OSC_START: any earlier OSC is dropped.
func (s *Sniffer) oscStart() {
	s.st = oscString
	s.osc = oscID
	s.oscID = -1
	s.capture = false
	s.hit = false
	s.plen = 0
}

func (s *Sniffer) oscPut(c rune) {
	switch s.osc {
	case oscAbort:
		return
	case oscID:
		switch {
		case c == ';':
			s.osc = oscPayload
			s.capture = s.oscID == 0 || s.oscID == 2
		case c >= '0' && c <= '9':
			if s.oscID < 0 {
				s.oscID = 0
			}
			s.oscID = min(s.oscID*10+int(c-'0'), maxOscID)
		default:
			s.osc = oscAbort
		}
		return
	}
	if !s.capture || s.hit {
		return
	}
	var buf [utf8.UTFMax]byte
	k := utf8.EncodeRune(buf[:], c)
	if s.plen+k > maxPayload {
		s.hit = true
		return
	}
	s.plen += copy(s.payload[s.plen:], buf[:k])
}

// oscEnd is a successful OSC_END. An OSC still reading its id ends with an empty payload, as in xterm.js.
func (s *Sniffer) oscEnd() {
	if s.osc == oscID {
		s.capture = s.oscID == 0 || s.oscID == 2
	} else if s.osc == oscAbort {
		return
	}
	if !s.capture || s.hit {
		return
	}
	p := s.payload[:s.plen]
	if len(p) > s.titleMax {
		k := s.titleMax
		for k > 0 && !utf8.RuneStart(p[k]) {
			k--
		}
		p = p[:k]
	}
	if !s.have || s.title != string(p) { // an unchanged title (a prompt that re-sets it) allocates nothing
		s.title = string(p)
		s.have = true
	}
}
