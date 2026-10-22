#!/usr/bin/env bash
# closure-check.sh: the adapters' non-test Go files hold (a) no golden vocabulary literal and (b) no case folding or
# remapping, and (c) call equiv.Map("<engine>", "<field>" for every (engine, field) pair of equiv/table.json (each
# row's engine line, then its field line; a count that differs from the row count is a FAIL).
set -uo pipefail
cd "$(dirname "$0")/.." || exit 2
src() { find gitpod magnet xvt -name '*.go' ! -name '*_test.go'; }
fail=0
say() { echo "CLOSURE FAIL $*"; fail=1; }
for f in $(src | xargs grep -lwF -e '"none"' -e '"x10"' -e '"vt200"' -e '"drag"' -e '"any"' -e '"block"' -e '"underline"' -e '"bar"'); do say vocab "$f"; done
for f in $(src | xargs grep -lE 'ToLower|ToUpper|EqualFold|strings\.Map|unicode\.To'); do say fold "$f"; done
pairs=$(grep -oE '"(engine|field)": *"[^"]+"' equiv/table.json | sed -E 's/.*: *"([^"]+)"/\1/' | paste -d' ' - -)
[ "$(wc -l <<<"$pairs")" = "$(grep -c '"id":' equiv/table.json)" ] || say unparsed equiv/table.json
while read -r e f; do
  find "$e" -name '*.go' ! -name '*_test.go' -print0 | xargs -0 grep -qF "equiv.Map(\"$e\", \"$f\"" || say map "$e" "$f"
done < <(sort -u <<<"$pairs")
[ $fail = 0 ] && echo "CLOSURE OK $(sort -u <<<"$pairs" | wc -l) pairs"
exit $fail
