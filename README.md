# go-term

Headless terminal (VT) engine in Go with xterm.js-compatible state, for Browserland Studio.

Work in progress: v0.x makes no compatibility promise. MIT licensed (see `LICENSE` and `NOTICE`).

## Layout

| path | what |
|---|---|
| `vt/` | the terminal engine: a fork of github.com/0magnet/xterm-go `vt/` (a Go port of xterm.js 6.0.0), see `vt/UPSTREAM` |
| `screen/` | Studio's typed read of a `vt` terminal: grid, cursor, modes, attribute runs, hashes, deltas, scrollback caps (below) |
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
The engine choice is `docs/decisions/0003-engine.md`; the read API's choices are `docs/decisions/0004-screen.md`.

Widths are Unicode v6, as `@xterm/headless` 6.0.0 without the unicode11 addon (what Studio's page loads); Unicode v11 is
not supported yet. DEC mode 2027 (grapheme clusters), which 6.0.0 does not have, is off unless
`vt.Options.GraphemeMode2027` is set.

## screen

`github.com/JohnConnorNPC/go-term/screen` reads a `*vt.Terminal` into the values of Studio's `screen_text` reply; Studio
builds the JSON. The functions take no lock: the caller serialises them with `Write`.

```go
t, err := screen.NewTerminal(80, 24, screen.Config{})       // history: min(1000 lines, 100,000 cells / cols)
f := screen.Read(t, screen.Request{View: screen.ViewScreen, Attrs: true})
f.Text(); f.Lines; f.Cols; f.Rows; f.Cursor; f.CursorHidden  // Lines padded to cols; Trimmed() trims them
f.AltScreen; f.AppCursor; f.View; f.HistoryLines; f.AttrRuns; f.ContentHash; f.StableHash
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
python `_content_hash` 69/69 and `_stable_hash` 68/68 byte for byte, headless textgrid 114/115 and hashes 74/74 (the
one miss is an engine difference, below); every xcur golden through `Read` (314/314, wide 30/30); pyte's 256 table.

Known difference from `@xterm/headless` 6.0.0: `vt` ends an APC string (`ESC _`) at BEL, so text after
`ESC _ ... BEL` prints; 6.0.0 ignores APC until ST (`test_control_string_families_are_skipped#2`).

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
