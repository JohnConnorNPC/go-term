#!/usr/bin/env bash
# record-check.sh [record]: every line of the six evidence files (scores/verdict.txt, equiv.txt, lic.txt, pin.txt,
# facts.txt, activity.txt) occurs as a whole line in the decision record, and scores/verdict.txt is byte for byte the
# verdict the committed summaries give (`go run ./cmd/verdict`, re-run in the pinned image by cand.sh: CACHE as there).
# Prints RECORD OK n/n, or one RECORD FAIL line per absent line or a verdict that differs.
set -uo pipefail
cd "$(dirname "$0")/.." || exit 2
rec=${1:-../../docs/decisions/0002-engine.md}
ok=0 n=0
for f in verdict equiv lic pin facts activity; do
  [ -s "scores/$f.txt" ] || { echo "RECORD FAIL unreadable scores/$f.txt"; n=$((n + 1)); continue; }
  while IFS= read -r l; do
    n=$((n + 1))
    if grep -Fxq -- "$l" "$rec"; then ok=$((ok + 1)); else echo "RECORD FAIL missing: $l"; fi
  done < "scores/$f.txt"
done
tmp=$(mktemp)
n=$((n + 1))
if REPO=$(cd ../.. && pwd) TMO=600 bash scripts/cand.sh run 'go run ./cmd/verdict' > "$tmp" && cmp -s "$tmp" scores/verdict.txt; then
  ok=$((ok + 1))
else
  echo "RECORD FAIL scores/verdict.txt is not the verdict of scores/*.json: $(diff "$tmp" scores/verdict.txt | head -3 | tr '\n' ' ')"
fi
rm -f "$tmp"
echo "RECORD $([ $ok = $n ] && [ $n -gt 0 ] && echo OK || echo FAIL) $ok/$n"
[ $ok = $n ] && [ $n -gt 0 ]
