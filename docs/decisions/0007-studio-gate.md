# 0007: the Studio gate: vectors ledger and differential fuzz (gh #9)

## Decision
1. **Every vector, one disposition** (rulings-0 R2, R3): Studio's A2b vectors (`vectors/headless`, 555) must all pass:
   rows as xterm gives them (written spaces kept) and as `Read` trims them, cursor, buffer, history, attribute runs,
   modes and title, and every `output_state` (Python's output bytes fed whole). Each A2a case (`vectors/python`, 746)
   is GO, LISTED, STUDIO or UNCOMPARED (`TestStudioLedger`); a case that fits none fails the test.
2. **LISTED means Studio's list**: go-term differs from Python, `vectors/headless/divergence.json` lists the case
   against python with the same fields, and go-term equals headless there. A listed case where go-term equals Python
   is allowed only for the colours (pyte's names) and the sniffer caps (Python's bounds). The list is Studio's and
   pinned by sha256; go-term never edits it. Studio signs it off after the v0.3.0 tag,
   which #15 made after the move to Go 1.27 (owner, ask-4 q2).
3. **Studio's byte producers are not ported** (R2): raw snapshot, pyte re-emit and textgrid's trim are Studio's
   (`internal/snapshot`); their bytes are checked through `output_state` (115 cases). The 52 `pyte._color_params`
   cases give SGR params, not bytes: go-term reads each back and must get the same colour name, except fg and bg of
   `brown`, `brightbrown` and `bfightmagenta`, which Studio's table emits as 39/49 (default). Studio's re-emit thus
   drops SGR 33 (pyte's "brown"); that is Studio's to fix. Every disposition count is asserted, and every
   divergence.json row against python must be reached as GO or LISTED.
4. **One byte more than Python's snapshot**: `ReplayBytes` puts CAN between ring and postamble (0006, #8 S2), so it is
   Python's `raw.render + _mode_reassert` plus one 0x18. `TestRawReplay` checks that exactly on the 6 `raw.render`
   vectors Studio does not trim, for all 60 mode sets. Recorded here and in the README, not in Studio's list
   (rulings-2 Q2).
5. **The fuzz goes through the conformance runner** (R1): `cmd/fuzzgen` writes seeded cases, `gen/fuzz.mjs` records
   xcur goldens for them in the pinned node container (`--network none`, `--pull never`), and the runner scores them
   with `-fuzz`: split variants headless itself gets wrong are scored against its whole-write result (go-term is
   split-safe since #14, 0005), so no cut is skipped. `scripts/fuzz.sh` loops batches in one container; it is outside
   `scripts/check.sh` (rulings-2 Q3).

## Tests
`TestVectorPins`, `TestHeadlessVectors`, `TestRawReplay`, `TestPythonTextgrid`, `TestStudioLedger` (screen);
`cmd/fuzzgen`'s `TestDeterministic` and `TestBounds`; the runner's `TestRunScoreSensitive` and `TestCLIFuzz`.
