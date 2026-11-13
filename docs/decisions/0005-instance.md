# 0005: the safe instance (gh #7)

## Decision
1. **`screen.Instance` in package screen** (rulings Q1): it owns `Config`, `Frame` and `FrameCache` already. The
   lock-free `Read`, `AttrRuns` and hashes stay public for single-goroutine callers; an Instance calls them under its
   lock and keeps one `FrameCache` of its own.
2. **One `sync.Mutex`** (Q2) for `Write`, `Resize` and `Read`. Each `Read` copies the grid into a `Frame` under the
   lock, which is the snapshot. There is no RWMutex, since that would need proof that a read never writes vt
   state, and no copy-on-write clone.
3. **Payload cap per terminal** (Q3): `vt.Options.PayloadLimit` (0 keeps xterm.js's 10,000,000), counted by the OSC
   and DCS dispatch and the APC buffer, so a string past the cap is dropped whole as soon as it passes it. Buffers over
   64 KiB are released when their string ends. An Instance defaults to 1 MiB code units (`DefaultPayloadLimit`):
   a longer title, link, colour or DECRQSS string is dropped where headless would take it, a documented divergence.
   APC is buffered only when a handler is set.
4. **Size limits, no cell budget knob** (Q4, Q5): cols 1..2000 (1 reads as `vt.MinimumCols`, 2), rows 1..1000; outside
   them `New` and `Resize` return an error and the size is unchanged. History stays bounded by `ScrollbackCells`.
   The worst case measured is about 55 MiB per instance.
5. **Panics fail the instance** (R08): `recover` under the lock; `ErrFailed` with the panic and stack is sticky, since
   vt's state may be torn mid-write.
6. **Replies and title are not exposed** (Q6): the sniffer owns modes and title for Studio.

## Known difference kept (Q7, manager ruling)
vt ports xterm.js's `Utf8ToUtf32` as is: a held continuation byte `0x80` reads as the end of the held bytes, so a write
ending right after it loses the character (`E2 80 | 94`, an em dash; ZWSP, ZWJ, U+20000). `@xterm/headless` does the
same: its recorded split-sensitive goldens wide/cjk-ext-b, zwj-family and zwsp. `TestInstanceSplit` names these cuts
(`quirkCut`) and asserts that they still differ. A fix is a separate question for after #7.

## Engine fix found by the fuzz
A zero-width char that joins nothing (after a C0 control) in the pending-wrap column was written one cell past the line
(index out of range). xterm.js's typed array drops that write, and vt now skips it (`TestCombiningAtPendingWrap`,
FuzzInstance crasher `d791907a10fdac2f` kept as a seed).
