# go-term

Headless terminal (VT) engine in Go with xterm.js-compatible state, for Browserland Studio.

Work in progress: v0.x makes no compatibility promise. MIT licensed (see `LICENSE` and `NOTICE`).

## Layout

| path | what |
|---|---|
| `vt/` | the terminal engine: a fork of github.com/0magnet/xterm-go `vt/` (a Go port of xterm.js 6.0.0), see `vt/UPSTREAM` |
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
passes;
`go vet`; `go test` (including the conformance gate, every case against `@xterm/headless@6.0.0`); `go test -race`;
`cross`: `CGO_ENABLED=0` builds for linux/amd64, linux/arm64, darwin/arm64 and windows/amd64 (linux/386 and linux/arm
are REPORT lines only); a `FuzzFeed` run on `./sniff`. The GitHub workflow `check.yml` runs the same script.
The engine choice is `docs/decisions/0003-engine.md`.

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
