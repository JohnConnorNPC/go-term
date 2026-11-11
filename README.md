# go-term

Headless terminal (VT) engine in Go with xterm.js-compatible state, for Browserland Studio.

Private, work in progress.

## Layout

| path | what |
|---|---|
| `conformance/` | the conformance suite: cases, xterm-headless goldens, the Go runner and CLI (see its README) |
| `conformance/candidates/` | the survey of third-party Go emulators (its own module) |
| `docs/decisions/` | decision records |
| `scripts/check.sh` | the local checks: `gofmt -l` empty, `go vet`, `go test`, `go test -race` |

## Checks

`bash scripts/check.sh` runs the root module's checks with the `go` on `PATH` (`GO=<path>` picks another) and exits
non-zero when any step fails. There is no CI workflow yet (#3).
