# 0006: snapshot and restore by re-feed (gh #8)

## Decision
1. **Restore re-feeds bytes** (R15, Studio SA1 Q4): a new instance is fed `Preamble` (Studio's raw snapshot preamble,
   `ESC[0m ESC[2J ESC[H`), the ring as Studio holds it, CAN, then `sniff.Postamble` of the agent-reported modes. The
   reported state is the sniffer's view of the whole stream, so a mode set that scrolled out of the ring survives.
   CAN (review round 1, S2) abandons a sequence the ring was cut in, as xterm.js does: without it the postamble's ESC
   would end a cut OSC successfully and commit a fragment title. No handler takes CAN, so a whole ring is unchanged.
2. **No trim, no strip** (rulings Q1): Studio owns the ring and its raw snapshot (`raw.py`'s replay from the last clear,
   resync after an evicted head, the DA/DSR/OSC 52 strip). Trimming here would drop history and break the round trip.
   A restored instance sets no callbacks (Instance sets none, 0005 item 5), so queries in the ring are answered to
   nobody. The replay is Studio's Python snapshot plus that one CAN byte (0007).
3. **A versioned Go value, no byte encoding** (Q2): `Snapshot{Version, Cols, Rows, Modes, Title, TitleSet, Ring}` with
   `SnapshotVersion = 1`; any other version is an error. No xterm.js serialize addon format is read or written; Studio
   owns the JSON or wire form.
4. **Modes and title in the read** (Q3): `Frame` gains `BracketedPaste`, `MouseTrack`, `MouseEnc` (numbered as
   `sniff.Modes`), `Title` and `Modes()`, read from vt's state, so the restored state is observable.
5. **Alt screen aligned only on a mismatch** (GO ruling): xterm.js's `CSI ? 1049 h` on the alt screen does not clear
   it (`activateAltBuffer` returns early), so prefixing it to every replay would leave a ring's pre-switch output on
   the alt grid. Restore feeds once; reported but not reached, it feeds again after `CSI ? 1049 h` (no switch in the
   ring means the whole ring was drawn on the alt screen); reached but not reported, `CSI ? 1049 l` follows. Both
   directions and the agreeing case are tested (`TestRestoreBasic`).
6. **The injected title is set last**, through vt's own `SetTitle`.

## Known difference kept
The sniffer cuts a title to `TitleMax` (1024 bytes by default), so a restored instance has a longer title cut where an
instance fed the stream has it whole.

## Tests
`TestReadModes` (every mode against the sniffer, step by step), `TestRestoreBasic`, `TestReplayBytes` (all 60 rows of
`sniff/testdata/postamble.json`), `TestRestoreRoundTrip` (its streams; empty rings for all 120 mode sets; every head and
tail cut, with and without an injected title), `TestRestoreLongTitle` (the TitleMax cut, 1024 and 1023 bytes).
