# go-term

Headless terminal (VT) engine in Go with xterm.js-compatible state, for Browserland Studio.

Work in progress: v0.x makes no compatibility promise. MIT licensed (see `LICENSE` and `NOTICE`).

## Layout

| path | what |
|---|---|
| `vt/` | the terminal engine: a fork of github.com/0magnet/xterm-go `vt/` (a Go port of xterm.js 6.0.0), see `vt/UPSTREAM` |
| `screen/` | Studio's typed read of a `vt` terminal: grid, cursor, modes, title, attribute runs, hashes, deltas, scrollback caps; the safe `Instance` and `Restore` (below) |
| `sniff/` | the mode and title sniffer and the snapshot postamble (below) |
| `conformance/` | the conformance suite: cases, xterm-headless goldens, the Go runner and CLI (see its README) |
| `conformance/goterm/` | the runner adapter for `vt/` and the gate test (`TestGateXcur`) |
| `conformance/candidates/` | the survey of third-party Go emulators (its own module; the frozen #2 record) |
| `docs/decisions/` | decision records |
| `scripts/check.sh` | the checks: `gofmt -l` empty, deps, `go vet`, `go test`, `go test -race`, cross builds, fuzz |
| `.github/workflows/check.yml` | CI: runs `scripts/check.sh` on push and pull request to master |

## Checks

`bash scripts/check.sh` runs the root module's checks with the `go` on `PATH` (`GO=<path>` picks another) and exits
non-zero when any step fails. The steps: `gofmt -l` prints nothing; `deps`: go.mod requires exactly
`scripts/deps.allow` (golang.org/x/crypto for blake2b, golang.org/x/sys indirect), replaces nothing, and `go mod verify`
passes; `go vet`; `go test` (including the conformance gate, every case against `@xterm/headless@6.0.0`); `go test -race`;
`cross`: `CGO_ENABLED=0` builds for linux/amd64, linux/arm64, darwin/arm64 and windows/amd64 (linux/386 and linux/arm
are REPORT lines only); a `FuzzFeed` run on `./sniff`. The GitHub workflow `check.yml` runs the same script.
The engine choice is `docs/decisions/0003-engine.md`; the read API's choices are `docs/decisions/0004-screen.md`,
the safe instance's `docs/decisions/0005-instance.md`.

Widths are Unicode v6, as `@xterm/headless` 6.0.0 without the unicode11 addon (what Studio's page loads); Unicode v11 is
not supported yet. DEC mode 2027 (grapheme clusters), which 6.0.0 does not have, is off unless
`vt.Options.GraphemeMode2027` is set.

## screen

`github.com/JohnConnorNPC/go-term/screen` reads a `*vt.Terminal` into the values of Studio's `screen_text` reply; Studio
builds the JSON. The functions take no lock: the caller serialises them with `Write`, or uses an `Instance`.

```go
t, err := screen.NewTerminal(80, 24, screen.Config{})       // history: min(1000 lines, 100,000 cells / cols)
f := screen.Read(t, screen.Request{View: screen.ViewScreen, Attrs: true})
f.Text(); f.Lines; f.Cols; f.Rows; f.Cursor; f.CursorHidden  // Lines padded to cols; Trimmed() trims them
f.AltScreen; f.AppCursor; f.View; f.HistoryLines; f.AttrRuns; f.ContentHash; f.StableHash
f.BracketedPaste; f.MouseTrack; f.MouseEnc; f.Title; f.Modes() // modes as sniff.Modes; the last OSC 0/2 title
cache := screen.NewFrameCache()                               // one per terminal: the last 16 frames
changed, delta := cache.Delta(&f, since, degraded)            // delta: send changed rows only, text ""
err = screen.Resize(t, cols, rows, screen.Config{})           // re-applies the history cap
```

- Rows are space-padded to the width as Python's renderers give them (a wide glyph is one rune over two cells). The
  cursor is grid-relative and clamped to the grid (the pending-wrap column reads as the last one); `CursorHidden` is
  DECTCEM. `ViewScrollback` with `Lines` prepends history; on the alt screen it reads as `ViewScreen`.
- `ContentHash` is blake2b-128 (`digest_size=16`) of the UTF-8 text, 32 lowercase hex digits. `StableHash` is the same
  over the text with the cursor's glyph replaced by one NUL (all cells and runes of that glyph), equal to `ContentHash`
  when there is no cursor; `screen.StableHash` is Python's `_stable_hash` itself (the column is a rune index).
- `AttrRuns`: maximal runs of (fg, bg, reverse, bold, underscore) other than the default, `len` in cells, at most 2000.
  Colours are pyte's names: `default`; `black red green brown blue magenta cyan white` and `bright*` for SGR 30-37,
  40-47, 90-97, 100-107 (bg 105 is `brightmagenta`); 6-hex from pyte's 256 table for every `38;5;n`/`48;5;n` (n=1 is
  `cd0000`); `%02x%02x%02x` for `38;2`. Dim, italic and the other attributes are not part of a run.
- `Delta`: a delta only when `since` is a remembered hash, the read is not degraded, the view is screen, the row counts
  match and at most `int(rows*0.6)` rows changed (an unchanged screen always qualifies); otherwise the full text.

Tested against Studio's vectors (copies in `screen/testdata`, sha256-pinned): python colours 53/53, pyte attr runs 11/11,
python `_content_hash` 69/69 and `_stable_hash` 68/68 byte for byte, headless textgrid 115/115 and hashes 74/74; every
xcur golden through `Read` (314/314, wide 30/30); pyte's 256 table.

### Instance

`screen.Instance` is a terminal to share between goroutines in a server: `Write`, `Resize` and `Read` hold one lock, so
each `Read` is a consistent snapshot taken between writes.

```go
i, err := screen.New(80, 24, screen.Config{})                 // cols 1..2000, rows 1..1000, else an error
n, err := i.Write(chunk)                                      // io.Writer; a chunk may end anywhere
f, err := i.Read(screen.Request{Attrs: true})                 // a Frame, as screen.Read
changed, delta, err := i.Delta(&f, since, degraded)           // the instance's own FrameCache
err = i.Resize(cols, rows)                                    // out of range: an error, size unchanged
err = i.Err()                                                 // nil, or ErrFailed once a panic failed it
```

- **Chunk splits**: a write may end inside any ESC, CSI, OSC, DCS or APC sequence or UTF-8 character; vt carries the
  rest over to the next write. One exception is inherited from xterm.js's UTF-8 decoder and kept for parity with
  `@xterm/headless`: a write that ends right after a `0x80` continuation byte of a 3- or 4-byte character loses that
  character (`E2 80` + `94` for an em dash; ZWSP, ZWJ and U+20000 split the same way: headless's split-sensitive
  goldens).
- **Panics**: a panic during a call fails that instance only: every later call returns `ErrFailed` wrapping the panic
  (the call that panicked also gets the stack). vt's state may be torn, so the instance is not reused; other instances
  are unaffected.
- **Callbacks**: an Instance sets none of vt's callbacks: replies are not exposed, and modes and title are read from vt's
  state into the `Frame`. Any callback added later
  runs inside `Write`, under the instance's lock, and must not call the Instance (it would deadlock).
- **Repeat counts**: a CSI count is capped at 2^31-1. CHT and CBT stop at a tab stop that no longer moves; IL, DL, SU
  and SD stop at the scroll region's height; REP keeps its count past a full buffer modulo one line. None of these
  changes what is drawn; each only drops turns that would change nothing.
- **Memory**: an OSC, DCS or APC string longer than `Config.PayloadLimit` (default `DefaultPayloadLimit`, 1 MiB code
  units, 4 MiB as runes) is dropped whole, where xterm.js keeps up to 10,000,000. 64 MiB of an unterminated string
  peaks at about 4.4 MiB of heap. At the largest size, 2000x1000 with both screens and the history full, an instance
  holds about 55 MiB.

### Restore

`screen.Restore` rebuilds an instance by re-feeding Studio's ring of raw pty output with the mode and title state the
agent's sniffer reported for the whole stream, so a mode set that scrolled out of the ring is still restored.

```go
i, err := screen.Restore(screen.Snapshot{
	Version: screen.SnapshotVersion,                          // 1; any other version is an error
	Cols: 80, Rows: 24,                                       // New's limits
	Modes: s.Modes(), Title: title, TitleSet: ok,             // the sniffer's report (s.Title())
	Ring: ring,                                               // as Studio holds it
}, screen.Config{})
b := screen.ReplayBytes(m, ring)                              // what Restore writes
```

- **What is fed**: `screen.Preamble` (`ESC[0m ESC[2J ESC[H`, Studio's raw snapshot preamble: no hard reset, no
  `ESC[3J`), the ring, `screen.Cancel` (CAN), then `sniff.Postamble(m)` (bracketed paste, DECCKM, the mouse groups
  reset-then-set). CAN abandons a CSI, OSC, DCS or APC the ring was cut in, as xterm.js does, so the postamble's ESC
  does not complete it (a cut OSC title is not committed); after a whole ring it does nothing. The ring is
  not trimmed and queries are not stripped: that is Studio's, and a restored instance answers nothing (no callbacks).
- **Alt screen**: the postamble does not carry it, so `Restore` aligns it only when the replay disagrees with
  `Modes.Alt`. Reported but not reached (the switch scrolled out: the whole ring was drawn on the alt screen), the replay
  is fed again to a new instance after `CSI ? 1049 h`; reached but not reported, `CSI ? 1049 l` follows. A replay that
  agrees is kept as is: `CSI ? 1049 h` on the alt screen does not clear it, so a blind prefix would leave the ring's
  output from before its own switch on the alt grid. A ring that leaves the alt screen although `Alt` is reported gets
  `CSI ? 1049 h` last, an empty alt grid.
- **Title**: when `TitleSet`, the injected title replaces the ring's. The sniffer cuts a title to its `TitleMax` (1024
  bytes by default) on a rune boundary, so a longer title restores cut where the fed instance has it whole: a
  documented difference (`TestRestoreLongTitle`).
- **Round trip**: a stream fed to one instance and restored from the same bytes with its sniffed state read the same
  (lines with history, cursor, attribute runs, hashes, modes, title); an empty ring and a ring cut at any byte at
  either end restore without error, report the injected modes and, with no injected title, the title of the ring
  alone (no fragment) (`TestRestoreRoundTrip`).
- **Format**: `Snapshot` is a Go value, versioned by `SnapshotVersion`; go-term defines no byte encoding and never uses
  xterm.js's serialize addon format. Studio owns the wire form.

## sniff

`github.com/JohnConnorNPC/go-term/sniff` watches a terminal's output for what a snapshot cannot carry and turns it back
into bytes:

```go
s := sniff.New(sniff.Options{})  // TitleMax 0 = 1024 bytes
s.Feed(chunk)                    // every PTY chunk, in order; chunk-safe, bounded, never retains chunk
m := s.Modes()                   // {Alt, DECCKM, BracketedPaste bool; MouseTrack 0|9|1000|1002|1003; MouseEnc 0|1006|1016}
title, ok := s.Title()           // the last OSC 0/2 title; ok is false before the first; untrusted text
snapshot = append(snapshot, sniff.Postamble(m)...)
```

The title is whatever the program wrote: valid UTF-8 within `TitleMax`, but DEL and format characters (bidi overrides
such as U+202E, U+2028/2029, U+FEFF) pass through. Treat it as untrusted text and escape or sanitise it before showing
it in a terminal, a log or a page.

A `Sniffer` is not safe for concurrent use; `Reset` returns it to the power-on state (a new PTY). It parses like
xterm.js (UTF-8 decoded as xterm.js decodes it, then the escape-sequence states of xterm.js's VT500 table), so a mode
sequence inside an OSC or DCS payload does not count, C1 controls count when UTF-8 encoded, and RIS (`ESC c`) clears
every mode while DECSTR (`CSI ! p`) clears DECCKM and bracketed paste. `Postamble` writes `CSI ? 2004 h|l`,
`CSI ? 1 h|l`, then per mouse group the resets of every inactive member before the set of the active one; the alternate
screen is left to the grid replay.

Caps (ours, not xterm.js's): an OSC payload over 4096 bytes is abandoned; a CSI with more than 68 bytes before its
final byte is dropped; the title is cut to `TitleMax` bytes on a rune boundary. Memory does not grow with the input.

Tested against:
- Browserland Studio's vectors (copies in `sniff/testdata`, sha256-pinned): the Python reference sniffers' 118 altscreen
  and 80 titles cases plus the 4 cap cases, op by op, and the same 202 inputs through `@xterm/headless` 6.0.0 per chunk;
- every conformance case's x53 and xcur golden, step by step (modes, buffer, title);
- the Python reference's postamble for all 60 mode combinations (`sniff/testdata/postamble_oracle.py`);
- splits at every offset, 64 MiB bounds and `FuzzFeed`.

Known differences:
- from the Python reference, where xterm.js wins (owner ruling): a bare ESC inside an OSC title ends the OSC and commits
  the title (`test_st_terminated_survives_any_split[21]#0`, `test_esc_abort_inside_payload#0`). The Python reference
  also counts mode sequences inside OSC/DCS payloads, ignores DECSTR, skips a mode sequence with `:` sub-parameters and
  takes any number of params; no vector shows a difference there, and this package follows xterm.js;
- not a difference, noted: DEL (0x7F) inside an OSC title payload stays in `Title()`, as in xterm.js and Python, while
  the other C0 controls are dropped and C1 controls abandon the title;
- from xterm.js, the caps above (`synthetic/title-payload-cap#1`, `synthetic/mode-carry-cap#1`).

## Releases

- v0.2.0: the `sniff` package and `scripts/check.sh`. v0.x makes no compatibility promise.
- v0.1.0: the conformance harness.
