#!/bin/sh
# Inside the node container with -w /repo/conformance: reinstall from the lockfile, verify the pins, generate twice into
# temp dirs and compare both runs with each other and with the committed golden/, byte for byte.
set -eu
TMP=$(mktemp -d); : "${TMP:?}"
fail() { echo "REGEN FAIL $*"; exit 1; }
sums() { (cd "$1" && find . -type f | sort | xargs sha256sum); }
cd gen
grep -v '^#' pins.lst | while read -r name ver file url integ; do
  [ -s "vendor/${file:?}" ] || { echo "REGEN FAIL vendor missing: run fetch-pins.sh ($file)"; exit 1; }
done || exit 1
pinned=$(node -e 'console.log(JSON.parse(require("fs").readFileSync("pins.json", "utf8")).lockSha256)')
h0=$(sha256sum package-lock.json | cut -d' ' -f1)
[ "$h0" = "$pinned" ] || fail "lockfile H0 != pins"
rm -rf "${PWD:?}/node_modules"
npm ci --offline --ignore-scripts --no-audit --no-fund > "$TMP/ci.log" 2>&1 || fail "npm ci $(tail -1 "$TMP/ci.log")"
h1=$(sha256sum package-lock.json | cut -d' ' -f1)
[ "$h1" = "$h0" ] || fail "lockfile H1 != H0"
node verify-pins.mjs > "$TMP/pins.log" || fail "$(cat "$TMP/pins.log")"
node gen.mjs --out "$TMP/A" > "$TMP/a.log" || fail "gen A: $(tail -1 "$TMP/a.log")"
node gen.mjs --out "$TMP/B" > "$TMP/b.log" || fail "gen B: $(tail -1 "$TMP/b.log")"
sums "$TMP/A" > "$TMP/a.sum"; sums "$TMP/B" > "$TMP/b.sum"; sums ../golden > "$TMP/c.sum"
n=$(wc -l < "$TMP/a.sum"); cases=$(find ../cases -name '*.json' | wc -l)
[ "$cases" -gt 0 ] && [ "$n" -eq $((2 * cases + 2)) ] || fail "files=$n cases=$cases"
cmp -s "$TMP/a.sum" "$TMP/b.sum" || fail "A!=B $(diff "$TMP/a.sum" "$TMP/b.sum" | grep '^<' | head -1 | awk '{print $3}')"
cmp -s "$TMP/a.sum" "$TMP/c.sum" || fail "committed $(diff "$TMP/a.sum" "$TMP/c.sum" | grep '^[<>]' | head -1 | awk '{print $3}')"
rm -rf "${TMP:?}"
echo "REGEN OK files=$n"
