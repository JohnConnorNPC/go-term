#!/usr/bin/env bash
# check.sh: the root module's local checks, in order: gofmt -l (must print nothing), deps (go.mod requires and
# replaces nothing), go vet, go test, go test -race, cross builds (linux/amd64, linux/arm64, darwin/arm64,
# windows/amd64; 386 and armv7 reported), and a 60 s FuzzFeed run on ./sniff. Every step runs even when an
# earlier one fails; the exit code is 1 when any step failed. Uses the go on PATH (GO=<path> overrides).
# The conformance/candidates module is its own module and is not covered here (conformance/scripts/go.sh runs
# it in the pinned container).
# The race step on conformance/runner takes far longer than go test's 10-minute default, hence the timeouts.
# Usage: bash scripts/check.sh [fuzztime]   (default 60s; 0 skips the fuzz step)
set -u
cd "$(dirname "$0")/.."
GO=${GO:-go}
FUZZTIME=${1:-60s}
export GOFLAGS=-mod=readonly
echo "HEAD $(git rev-parse HEAD)"
echo "porcelain-begin"; git status --porcelain; echo "porcelain-end"
echo "go $($GO version)"
echo "command: bash scripts/check.sh $FUZZTIME"
fail=0
step() { # step <name> <cmd...>
	local name=$1; shift
	echo "== $name: $*"
	"$@"; local rc=$?
	echo "== $name exit $rc"
	[ $rc -eq 0 ] || fail=1
}
gofmt_empty() {
	local out
	out=$("${GO%go}gofmt" -l . 2>&1) || { echo "$out"; return 2; }
	[ -z "$out" ] || { echo "gofmt -l lists:"; echo "$out"; return 1; }
}
step gofmt gofmt_empty
# deps: the root module requires and replaces nothing (vt/ is in-tree, stdlib only). Read through go's own go.mod
# parser, so a quoted path such as replace "example.com/m" => ../m is seen too (#4 R2).
no_deps() {
	local j
	j=$("$GO" mod edit -json) || return 2
	printf '%s' "$j" | grep -Eq '"(Require|Replace)": \[' && { echo "go.mod has a require or replace:"; printf '%s\n' "$j"; return 1; }
	return 0
}
step deps no_deps
step vet "$GO" vet ./...
step test "$GO" test -count=1 -timeout 30m ./...
step race "$GO" test -count=1 -race -timeout 90m ./...
# cross: the release targets must build without cgo; 386 and armv7 are reported only (#3 rulings Q5).
cross() {
	local t rc=0
	for t in linux/amd64 linux/arm64 darwin/arm64 windows/amd64; do
		CGO_ENABLED=0 GOOS=${t%/*} GOARCH=${t#*/} "$GO" build ./... || { echo "cross build failed: $t"; rc=1; }
	done
	for t in linux/386 linux/arm; do
		if CGO_ENABLED=0 GOOS=${t%/*} GOARCH=${t#*/} GOARM=7 "$GO" build ./... >/dev/null 2>&1; then
			echo "REPORT $t build ok"; else echo "REPORT $t build failed (not gating)"; fi
	done
	return $rc
}
step cross cross
if [ "$FUZZTIME" != 0 ] && [ -d sniff ]; then
	step fuzz "$GO" test -run '^$' -fuzz '^FuzzFeed$' -fuzztime "$FUZZTIME" ./sniff
fi
echo "porcelain-after-begin"; git status --porcelain; echo "porcelain-after-end"
echo "FINAL EXIT $fail"
exit $fail
