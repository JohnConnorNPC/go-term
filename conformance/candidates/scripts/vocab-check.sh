#!/usr/bin/env bash
# vocab-check.sh: the vendored typings excerpt (cmd/facts/testdata, read by TestGoldenVocab) is true to the 5.3.0
# typings in ../gen/node_modules (#1's build output: fetch-pins.sh + npm ci): the whole file's sha256 is the one the
# excerpt's header names and each excerpt line is a line of it. Prints VOCAB OK n lines, VOCAB FAIL lines (exit 1), or
# VOCAB UNMEASURED without the typings (exit 2); the repo bans skipped tests, so this is not a go test.
set -uo pipefail
cd "$(dirname "$0")/.." || exit 2
ex=cmd/facts/testdata/xterm-headless-5.3.0.d.ts real=../gen/node_modules/xterm-headless/typings/xterm-headless.d.ts
[ -r "$real" ] || { echo "VOCAB UNMEASURED $real absent"; exit 2; }
want=$(grep -oE 'sha256 [0-9a-f]{64}' "$ex" | cut -c8-) got=$(sha256sum "$real" | cut -c1-64)
fail=0
[ -n "$want" ] && [ "$want" = "$got" ] || { echo "VOCAB FAIL the typings are sha256 $got, the excerpt names ${want:-none}"; fail=1; }
while IFS= read -r l; do
  grep -Fxq -- "$l" "$real" || { echo "VOCAB FAIL not a line of the typings: $l"; fail=1; }
done < <(grep -v '^//' "$ex")
[ $fail = 0 ] && echo "VOCAB OK $(grep -vc '^//' "$ex") lines"
exit $fail
