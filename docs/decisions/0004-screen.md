# 0004: the Studio read API (gh #6)

## Decision
1. **A package over vt, not a fork of it.** `screen` reads a `*vt.Terminal` through vt's exported buffer, cell and
   mode fields and takes no lock; #7's safe instance calls it under its own lock. vt changed in one place only (2).
2. **DEC mode 2027 off by default.** `@xterm/headless` 6.0.0 has no mode 2027 (no `case 2027` in its bundle) and no
   golden uses it; vt (a later upstream) has it. `vt.Options.GraphemeMode2027` (default false) keeps it for callers that
   want it; off, SM/RM ?2027 are ignored and DECRQM ?2027 answers 0 (rulings Q7).
3. **Rows padded to the width** (rulings Q8): the text and hashes match what Python's renderers return, and the stable
   hash masks a cursor on a blank cell. Headless vectors are compared trimmed (`Frame.Trimmed`).
4. **Stable hash over glyphs** (rulings Q1): the cursor cell's glyph is replaced by one NUL, all of a wide glyph or a
   combining sequence; for one-cell one-rune glyphs this is Python's `_stable_hash`, which `screen.StableHash` is.
5. **Cursor kept when hidden** (Q3) with `CursorHidden`; **attr runs on grid rows only** (Q6); **the frame cache is per
   terminal in go-term**, `degraded` is the caller's flag (Q5).
6. **History cap** min(lines, max(1, cells/cols)), default 1000 lines and 100,000 cells, re-applied on resize (Q4).
   **Unicode v6 only**; `Config.Unicode` accepts "6" and v11 waits for its own issue (Q2).
7. **First dependencies**: golang.org/x/crypto v0.48.0 (blake2b) and golang.org/x/sys v0.41.0 (indirect, blake2b's
   cpu detection): the newest x/crypto whose go.mod says go <= 1.24 (v0.49.0 needs 1.25). check.sh's deps step now
   requires exactly `scripts/deps.allow`, no replace, and `go mod verify`. `go get` wrote the go directive as
   `go 1.24.0`, the same floor as 0003's `go 1.24`.

## Known difference (fixed in #10)
vt collected APC as an OSC string, so BEL ended it; in 6.0.0 APC is the SOS/PM/APC ignore state that only ST ends.
Studio's headless textgrid vector `test_control_string_families_are_skipped#2` showed it (114/115). #10 keeps BEL from
ending an APC in vt, so headless textgrid is 115/115 and `knownHeadless` is empty.
