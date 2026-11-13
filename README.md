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
Two fuzzers run outside `check.sh` (its `go test` step runs only `FuzzInstance`'s seed corpus):
`go test -run '^$' -fuzz '^FuzzInstance$' ./screen` (each input stops past an 8,000,000-cell work budget and fails as
a hang after 60 s), and the differential fuzz against `@xterm/headless`, `bash conformance/scripts/fuzz.sh` (it needs
docker; see "Studio API").
Decision records in `docs/decisions/`: the engine choice `0003-engine.md` (`0002-engine.md` is the superseded survey
record), the read API `0004-screen.md`, the safe instance `0005-instance.md`, snapshot and restore `0006-restore.md`,
and the Studio gate (vectors ledger, differential fuzz) `0007-studio-gate.md`.

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
f.BracketedPaste; f.MouseTrack; f.MouseEnc; f.Title; f.Modes() // modes as sniff.Modes; the window title
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
- `Title` is vt's window title: the last OSC 0/2 or a title-stack pop (`CSI 23 t`) set it. `sniff`'s `Title()` sees
  OSC 0/2 only, so the two can differ after a pop; a restored instance has the sniffer's title (see Restore).
- `Delta`: a delta only when `since` is a remembered hash, the read is not degraded, the view is screen, the row counts
  match and at most `int(rows*0.6)` rows changed (an unchanged screen always qualifies); otherwise the full text.

Tested against Studio's vectors (copies in `screen/testdata`, sha256-pinned): python colours 53/53, pyte attr runs 11/11,
python `_content_hash` 69/69 and `_stable_hash` 68/68 byte for byte; every A2b (headless) vector, 555 of 555, with
rows, cursor, buffer, history, attribute runs, modes and title (textgrid 115, hashes 74, colours 53, pyte 26, raw 32,
synthetic 57, altscreen 118, titles 80) plus all 115 `output_state`s; every xcur golden through `Read` (314/314, wide
30/30); pyte's 256 table. See "Studio API" for the A2a ledger.

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
  rest over to the next write. This includes a write that ends right after a held `0x80` continuation byte (`E2 80` +
  `94`, an em dash; ZWSP, ZWJ, U+20000), where `@xterm/headless` 6.0.0 loses the character (its split-sensitive
  goldens wide/cjk-ext-b, wide/zwj-family, wide/zwsp): vt keeps it, a deliberate divergence (owner decision
  2026-11-14, #14).
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
  not trimmed and queries are not stripped: that is Studio's, and a restored instance answers nothing (no callbacks:
  Instance sets none, 0005 item 5).
- **One byte more than Studio's Python snapshot**: Python's snapshot is `raw.render(ring)` + `_mode_reassert(sniffer)`
  (`old/webterm/agent/agent.py`); `ReplayBytes` is the same bytes with CAN between ring and postamble. Where Studio
  trims nothing, `ReplayBytes` equals Python's snapshot plus that CAN exactly (`TestRawReplay`: 6 of 31 `raw.render`
  vectors, all 60 mode sets; the 24 Studio trims or strips and 1 digest are Studio's).
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

## Studio API

What Browserland Studio uses, by package; the sections above give the details.
- `screen`: `NewTerminal`, `Read(t, Request{View, Lines, Attrs})` -> `Frame` (`Lines`, `Text()`, `Trimmed()`, `Cols`,
  `Rows`, `View`, `Cursor`, `CursorHidden`, `AltScreen`, `AppCursor`, `BracketedPaste`, `MouseTrack`, `MouseEnc`,
  `Title`, `Modes()`, `HistoryLines`, `AttrRuns`, `ContentHash`, `StableHash`), `AttrRun`, `AttrRuns`/`MaxAttrRuns`,
  `Colour`/`Default`, `ContentHash`/`StableHash`/`HashSize`, `NewFrameCache`/`FrameCacheSize` with
  `FrameCache.Delta(&f, since, degraded) (changed []ChangedRow, delta)`, `Resize`, `Config` (`ScrollbackLines`, `ScrollbackCells`,
  `PayloadLimit`, `Unicode`; `HistoryCap`, `DefaultScrollbackLines`), `MaxCols`/`MaxRows`.
- `screen.Instance` (safe for concurrent use): `New`, `Write`, `Resize`, `Read(r) (Frame, error)`,
  `Delta(&f, since, degraded) (changed, delta, err)` (the instance's own cache; an error once failed), `Err`,
  `ErrFailed`. Restore: `Snapshot`, `SnapshotVersion`, `Restore` (an error for another version, a size out of `New`'s
  limits or a failed instance), `ReplayBytes`, `Preamble`, `Cancel`: see "Restore".
- `sniff`: `New(Options{TitleMax})`/`DefaultTitleMax`, `Feed`, `Reset`, `Modes()`, `Title()`, `Modes`, `Postamble`,
  `MouseTrackModes`/`MouseEncModes`: see "sniff".
- Vectors (Studio @ed3249b, sha256-pinned, `go test ./screen -run 'TestHeadlessVectors|TestStudioLedger' -v`): every
  A2b vector passes (555/555); every A2a case has one disposition: GO 546 (go-term gives Python's value), LISTED 27
  (Python differs from headless and go-term follows headless: listed in Studio's `vectors/headless/divergence.json`,
  pinned at `screen/testdata/headless/divergence.json`), STUDIO 167 (Studio's own producers: raw snapshot, pyte
  re-emit, textgrid trim; the 115 that give bytes are fed and equal headless's `output_state`; the 52 `_color_params`
  cases give SGR params, which go-term reads back as the same colour for 46 and as default for the 6 fg/bg of `brown`,
  `brightbrown` and `bfightmagenta`, the names Studio's table emits as 39/49), UNCOMPARED 6 (Studio's own list, plus
  `pyte.available`). Each count is asserted, not only the sum. `38;5;n` with n < 16 reads as pyte's hex (`AttrRuns` above), not headless's palette index.
- Differential fuzz against `@xterm/headless` 6.0.0: `bash conformance/scripts/fuzz.sh 600 <seed>` (conformance README
  "Fuzz").

## Releases

- v0.3.0: `screen` (Read, Instance, Restore) and `sniff`'s mode and title fixes; vt keeps a character a write ended
  inside (#14); Studio's vectors and the differential fuzz (#9).
- v0.2.0: the `sniff` package and `scripts/check.sh` (no `screen` yet). v0.x makes no compatibility promise.
- v0.1.0: the conformance harness.
