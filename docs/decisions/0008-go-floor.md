# 0008: the Go floor is 1.27 (gh #15)

Status: **Decided** (owner, ask-5, 2026-11-14: move off Go 1.24 before the v0.3.0 tag; manager rulings on #15 Q1-Q4).
Reverses QUEUE-2 Q3 (owner, 2026-11-13, "go 1.24"), recorded as 0003 item 3 and 0004 item 7.

## Decision
1. **go.mod says `go 1.27.0`, with no toolchain line; CI runs `1.27.x`** (one job). go-term needs Go 1.27 or newer to
   build. The line is 1.27.0, not the current patch, because conformance/scripts/go.sh and cand.sh run the pinned
   golang:1.27.3 image with `GOTOOLCHAIN=local`, which refuses a go line above 1.27.3.
2. **Why move**: Go 1.24 left support when go1.26.0 shipped (2026-02-10; its last patch is go1.24.13). govulncheck
   judges the standard library of the Go that runs it, so at any 1.24.x it reports stdlib advisories that will never be
   fixed there (#15: 25 reachable at go1.24.13). Go supports the two newest majors: 1.26 ends when 1.28 ships (about
   February 2027), 1.27 lasts until 1.29. Studio, the one consumer, is on go 1.27 already (0003).
3. **Dependencies**: the newest golang.org/x/crypto whose go.mod fits the go line and that is at least 14 days old on
   the day of the bump, and the golang.org/x/sys it requires. On 2026-11-13 that was x/crypto v0.58.0 (2026-10-09) and
   x/sys v0.49.0 (2026-10-06); v0.59.0 (2026-11-11) was too new. `scripts/deps.allow` carries the rule and the lines.
4. **conformance/candidates follows the root**: its `replace github.com/JohnConnorNPC/go-term => ../..` puts the root's
   requires into its module graph, so its go line (`go 1.25`) and its x/sys line move with the root's, or cand.sh
   (`-mod=readonly`, offline proxy) stops building it. Only go.mod, go.sum and mods.seed change there; the #2 record
   (scores, verdict, pins, licences) is not edited or re-scored, and record-check still says RECORD OK.

## What changes for consumers
A module that requires go-term needs Go 1.27 or newer; with `GOTOOLCHAIN=auto` (the default) an older go downloads
go1.27.0, and `go get` raises the consumer's own go line to 1.27.0 if it is lower. MVS picks x/crypto >= v0.58.0 and
x/sys >= v0.49.0. GODEBUG defaults follow the consumer's main module, not go-term's go line. No API, behaviour or
output change: same code, goldens and vectors.

## Checking it
On a host whose go is older than 1.27, run the checks with a pinned toolchain, e.g.
`GOTOOLCHAIN=go1.27.4 bash scripts/check.sh`; without the pin, auto-switch picks go1.27.0, an old patch, and govulncheck
would judge its standard library. `GOTOOLCHAIN=go1.27.4 GOFLAGS= go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...`
prints `No vulnerabilities found.` at v0.3.0.
