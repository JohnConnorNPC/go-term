# 0003 Engine: fork 0magnet/xterm-go vt/, gate on @xterm/headless 6.0.0 (issue #3)

Status: **Decided** (owner, 2026-11-13, QUEUE-2: Q1 gate 6.0.0, Q2 fork magnet, Q3 go 1.24). Supersedes the open
engine choice of 0002 ("survey-decides"); 0002 stays as the #2 record.

## Why 0002's premise lapsed
0002 gated on x53 because x53 = xterm-headless@5.3.0 was "the version Studio renders against". Studio has since moved
its page and its headless reference vectors to 6.0.0 (browserland-studio gh #68, F4: `tools/headless/package.json`
names `"@xterm/headless": "6.0.0"`, `internal/webui/vendor.go` `@xterm/xterm@6.0.0`). Every candidate failed x53 on one
probe, DECRQM `?2026`, where the ports answer as 6.0.0 does. Studio's Go floor is 1.27 (its go.mod), so the 1.23
concern of 0002's facts no longer applies.

## Decisions
1. **Gate = xcur** (`@xterm/headless@6.0.0`, conformance/golden/xcur): every case, its probe twins and every split
   variant must PASS. x53 is reported, not gated. The conformance CLI's `-ref` defaults to `xcur`.
2. **Engine = fork of github.com/0magnet/xterm-go `vt/`** at `aa0a32cd702b857db0c2807305c271e1fac1f76c`, imported as
   `github.com/JohnConnorNPC/go-term/vt` (a byte copy; only esctest's import path changed; vt/UPSTREAM, vt/LICENSE,
   NOTICE). Why magnet over gitpod-io/xterm-go: xcur 312/314 vs 303/314 in 0002; magnet's two misses were our 5.3.0-era
   equivalence rows (M7/M9 report the DECSCUSR override as `opts.*`; 6.0.0 leaves the options alone), gitpod's were
   nine behaviour bugs (origin-mode CPR, DECRC origin, scrollback 0, soft hyphen, U+FFFD split); magnet ports exactly
   6.0.0 and its licence is the xterm.js licence. Its upstream (one maintainer) matters less once the code is ours.
3. **Go floor = 1.24** (go.mod `go 1.24`; vt needs 1.24).
   Superseded by 0008 (#15, 2026-11-14): the floor is `go 1.27.0`, CI runs 1.27.x.
4. **conformance/candidates stays frozen** as the #2 record: its equivalence table (including M7/M9), scores and
   record-check are not edited or re-scored. The in-tree adapter (conformance/goterm) reads vt natively, with no table.
   Amended by 0008 (#15): its go.mod go line and x/sys line (with go.sum, mods.seed) follow the root's; the record
   itself is unchanged.

## Result
`go test -run TestGateXcur ./conformance/goterm/` → `GATE goterm xcur 314/314`; the CLI
(`go run ./conformance/cmd/conformance -engine goterm`) → `TOTAL 314/314 PASS skipped=7`. No change to vt/ was needed.
x53, reported (`-ref x53`): `TOTAL 0/314 FAIL skipped=0`. 310 cases fail first on the `B.t1.?2026` probe (6.0.0 knows
mode 2026, 5.3.0 does not); the other 4 are the cases where 5.3.0 and 6.0.0 behave differently (golden/diff.json):
`csi/rep-after-wide`, `dcs/decrqss-sp-q`, `decmodes/decscusr-bar`, `decmodes/decscusr-reset`.

## Follow-ups (QUEUE-2)
Studio read API (grid, cursor, attr runs, hashes, deltas), a safe instance (concurrency, panic isolation, memory bounds),
restore by re-feed with sniff's postamble, then Studio's A2a/A2b vectors at 100 % and tag v0.3.0 for Studio #9.
#4 (hardening of the #2 verdict tooling) is closed as superseded; its quoted-replace item is covered for the root
module by check.sh's `deps` step.
