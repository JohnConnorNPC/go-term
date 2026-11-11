#!/usr/bin/env bash
# check.sh: the root module's local checks, in order: gofmt -l (must print nothing), go vet, go test, go test -race,
# and a 60 s FuzzFeed run on ./sniff. Every step runs even when an earlier one fails; the exit code is 1 when any
# step failed. Uses the go on PATH (GO=<path> overrides). The conformance/candidates module is its own module and
# is not covered here (conformance/scripts/go.sh runs it in the pinned container).
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
step vet "$GO" vet ./...
step test "$GO" test -count=1 -timeout 30m ./...
step race "$GO" test -count=1 -race -timeout 90m ./...
if [ "$FUZZTIME" != 0 ] && [ -d sniff ]; then
	step fuzz "$GO" test -run '^$' -fuzz '^FuzzFeed$' -fuzztime "$FUZZTIME" ./sniff
fi
echo "porcelain-after-begin"; git status --porcelain; echo "porcelain-after-end"
echo "FINAL EXIT $fail"
exit $fail
