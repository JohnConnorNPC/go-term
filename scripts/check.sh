#!/usr/bin/env bash
# check.sh: the root module's local checks, in order: gofmt -l (must print nothing), deps (go.mod requires exactly
# scripts/deps.allow and replaces nothing), modverify (go mod download leaves go.mod and go.sum unchanged, go mod
# verify passes, every module is in the cache), readonly (go build with -mod=readonly), go vet, go test, go test -race, cross builds (linux/amd64, linux/arm64, darwin/arm64,
# windows/amd64; 386 and armv7 reported), govulncheck (skipped when vuln.go.dev does not answer, unless GOVULNCHECK=require), and two 60 s
# fuzz runs: fuzz (FuzzFeed on ./sniff) and fuzz-instance (FuzzInstance on ./screen). Every step runs even when an
# earlier one fails; the exit code is 1 when any step failed. Uses the go on PATH (GO=<path> overrides).
# A step's result is its "== <step> exit N" line; a skipped step prints "== <step> skipped: ... (not a pass)" instead.
# A step that changes git status --porcelain (writes a tracked or an untracked, not ignored file) is red too.
# The conformance/candidates module is its own module and is not covered here (conformance/scripts/go.sh runs
# it in the pinned container).
# The race step on conformance/runner takes far longer than go test's 10-minute default, hence the timeouts.
# Usage: bash scripts/check.sh [fuzztime]   (default 60s for each fuzz step; 0 skips both)
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
step() { # step <name> <cmd...>: red when cmd fails, and when it changes git status --porcelain (#12 S2)
	local name=$1 before after rc; shift
	echo "== $name: $*"
	before=$(git status --porcelain 2>&1)
	"$@"; rc=$?
	after=$(git status --porcelain 2>&1)
	if [ "$before" != "$after" ]; then
		echo "$name changed the tree (git status --porcelain, before and after):"
		diff <(printf '%s\n' "$before") <(printf '%s\n' "$after")
		[ $rc -ne 0 ] || rc=1
	fi
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
# x/sys it imports, #6) and replaces nothing. go.mod is read through go's own parser
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
}
step deps allowed_deps
# modverify: go mod verify checks only the modules already in the module cache, so on an empty cache it passes
# having checked no module zip (#12). The step downloads first, verifies, and requires every module that provides
# a package to the build or its tests (go list -deps -test) to be in the cache. go mod download also writes missing go.sum lines back (#12 M1), so go.mod and go.sum must
# come out unchanged: a change means the committed files are incomplete (red, and the committed copies are put
# back). With a warm cache the step is green offline; an empty cache offline fails at the download.
modverify() {
	local t mods all have rc=0
	t=$(mktemp -d) || return 2
	cp go.mod go.sum "$t/" || { rm -rf "$t"; return 2; }
	"$GO" mod download || rc=1
	if ! cmp -s go.mod "$t/go.mod" || ! cmp -s go.sum "$t/go.sum"; then
		echo "modverify: go mod download changed go.mod or go.sum (the committed files are incomplete):"
		diff "$t/go.mod" go.mod; diff "$t/go.sum" go.sum
		cp "$t/go.mod" "$t/go.sum" . || rc=2
		[ $rc = 2 ] || rc=1
	fi
	rm -rf "$t"
	[ $rc = 0 ] || return $rc
	"$GO" mod verify || return 1
	mods=$("$GO" list -e -deps -test -f '{{with .Module}}{{if not .Main}}{{.Path}} {{.Dir}}{{end}}{{end}}' ./...) ||
		{ echo "modverify: go list failed (unmeasured)"; return 2; }
	mods=$(printf '%s\n' "$mods" | awk 'NF && !seen[$0]++')
	all=$(printf '%s\n' "$mods" | awk 'NF >= 1' | wc -l)
	have=$(printf '%s\n' "$mods" | awk 'NF >= 2' | wc -l)
	[ "$all" -gt 0 ] || { echo "modverify: go list lists no modules (unmeasured)"; return 2; }
	[ "$have" = "$all" ] || { echo "modverify: $((all - have)) of $all modules not in GOMODCACHE:"; printf '%s\n' "$mods" | awk 'NF == 1'; return 1; }
	echo "REPORT modverify: $all of $all modules in GOMODCACHE, verified"
}
step modverify modverify
# readonly: the build needs no go.mod or go.sum update. GOFLAGS is exported above too; the step sets it itself so
# that a later change to the export cannot drop it.
step readonly env GOFLAGS=-mod=readonly "$GO" build ./...
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
# govulncheck: the pinned govulncheck, built by $GO, scans ./... against vuln.go.dev at symbol level (a vulnerable
# module that is required but never called is not a finding), for this host's GOOS/GOARCH only. It judges the stdlib
# of the go on its PATH, so PATH starts with "$GO" env GOROOT's bin, and the Go its own report names ("Go: ...") must
# equal "$GO" env GOVERSION. The DB is live: a green commit can turn red later. vuln.go.dev is probed first: with
# GOVULNCHECK=auto (default) the step is skipped, with no exit line (neither a pass nor a failure), only when the probe
# fails (or there is no curl); GOVULNCHECK=require (CI) makes that rc 2. With the DB up, a failed install is rc 2 in
# both modes. Green is rc 0 with "No vulnerabilities found."; rc 3 is "vulnerabilities found"; anything else is the
# tool failing.
vulnpin=golang.org/x/vuln/cmd/govulncheck@v1.9.0
vulnmode=${GOVULNCHECK:-auto} vulnwhy="" vulnrc=0
vulncheck_probe() { # rc 0 DB up, 1 unavailable ($vulnwhy), 2 bad mode
	case $vulnmode in auto | require) ;; *) echo "govulncheck: GOVULNCHECK=$vulnmode, want auto or require"; return 2 ;; esac
	command -v curl > /dev/null || { vulnwhy="no curl to probe vuln.go.dev"; return 1; }
	curl -fsS --max-time 10 -o /dev/null https://vuln.go.dev/index/db.json || { vulnwhy="vuln.go.dev did not answer"; return 1; }
}
vulncheck() { # reads $vulnrc, vulncheck_probe's rc
	local root want d rc got
	case $vulnrc in 0) ;; 1) echo "govulncheck: $vulnwhy (GOVULNCHECK=require)"; return 2 ;; *) return 2 ;; esac
	{ root=$("$GO" env GOROOT) && want=$("$GO" env GOVERSION); } || { echo "govulncheck: GO=$GO does not run"; return 2; }
	d=$(mktemp -d) || return 2
	if ! GOBIN=$d "$GO" install "$vulnpin"; then
		echo "govulncheck: $vulnpin did not install, with vuln.go.dev up (tool failure)"; rm -rf "$d"; return 2
	fi
	PATH="$root/bin:$PATH" "$d/govulncheck" -show version ./... > "$d/out" 2>&1
	rc=$?
	cat "$d/out"
	got=$(grep -m1 '^Go: ' "$d/out")
	if [ $rc = 0 ] && ! grep -q '^No vulnerabilities found\.$' "$d/out"; then
		echo "govulncheck: rc 0 without \"No vulnerabilities found.\""; rc=2
	fi
	[ "$got" = "Go: $want" ] || { echo "govulncheck: judged \"${got:-no Go: line}\", want \"Go: $want\" (\$GO)"; rc=2; }
	case $rc in 0) ;; 3) echo "govulncheck: vulnerabilities found (rc 3)" ;; *) echo "govulncheck: tool failed (rc $rc)" ;; esac
	rm -rf "$d"
	return $rc
}
echo "govulncheck mode: $vulnmode"
vulncheck_probe
vulnrc=$?
if [ $vulnrc = 1 ] && [ "$vulnmode" = auto ]; then
	echo "== govulncheck skipped: $vulnwhy (not a pass)"
else
	step govulncheck vulncheck
fi
if [ "$FUZZTIME" != 0 ] && [ -d sniff ]; then
	step fuzz "$GO" test -run '^$' -fuzz '^FuzzFeed$' -fuzztime "$FUZZTIME" ./sniff
fi
# fuzz-instance: FuzzInstance (screen/fuzz_test.go): each input stops past an 8,000,000-cell work budget, and one still
# running after 60 s fails as a hang. A crasher is written to screen/testdata/fuzz/FuzzInstance (porcelain-after shows it).
if [ "$FUZZTIME" != 0 ]; then
	step fuzz-instance "$GO" test -run '^$' -fuzz '^FuzzInstance$' -fuzztime "$FUZZTIME" ./screen
fi
echo "porcelain-after-begin"; git status --porcelain; echo "porcelain-after-end"
echo "FINAL EXIT $fail"
exit $fail
