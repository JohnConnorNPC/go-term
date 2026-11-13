#!/usr/bin/env bash
# check.sh: the root module's local checks, in order: gofmt -l (must print nothing), deps (go.mod requires exactly
# scripts/deps.allow, replaces nothing, go mod verify passes), go vet, go test, go test -race, cross builds (linux/amd64, linux/arm64, darwin/arm64,
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
# deps: the root module requires exactly scripts/deps.allow (path version [indirect]; x/crypto for blake2b and the
# x/sys it imports, #6) and replaces nothing, and go mod verify passes. go.mod is read through go's own parser
# (go mod edit -json), so a quoted path such as replace "example.com/m" => ../m is seen too (#4 R2).
allowed_deps() {
	local j got want
	j=$("$GO" mod edit -json) || return 2
	printf '%s' "$j" | grep -q '"Replace": \[' && { echo "go.mod has a replace:"; printf '%s\n' "$j"; return 1; }
	got=$(printf '%s\n' "$j" | awk '
		/"Require": \[/ { in_req = 1; next }
		in_req && /^\t\]/ { in_req = 0 }
		in_req && /"Path":/ { gsub(/[",]/, "", $2); path = $2; ver = ""; ind = "" }
		in_req && /"Version":/ { gsub(/[",]/, "", $2); ver = $2 }
		in_req && /"Indirect": true/ { ind = " indirect" }
		in_req && /^\t\t\}/ { print path " " ver ind }' | LC_ALL=C sort)
	want=$(grep -v '^#' scripts/deps.allow | grep -v '^$' | LC_ALL=C sort) || return 2
	[ -n "$want" ] || { echo "scripts/deps.allow lists nothing (unmeasured)"; return 2; }
	[ "$got" = "$want" ] || { echo "go.mod requires:"; printf '%s\n' "$got"; echo "scripts/deps.allow:"; printf '%s\n' "$want"; return 1; }
	echo "deps: $(printf '%s\n' "$got" | wc -l) requires match scripts/deps.allow"
	"$GO" mod verify
}
step deps allowed_deps
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
