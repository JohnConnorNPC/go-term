#!/usr/bin/env bash
# sumdb-check.sh [go.sum]: host curl of sum.golang.org for each go.sum module@version; every go.sum line (zip and
# /go.mod hashes) must equal the checksum database's line. Prints SUMDB FAIL <mod>@<ver> per mismatch, then a total.
set -uo pipefail
cd "$(dirname "$0")/.." || exit 2
sum=${1:-go.sum}
esc() { sed -E 's/[A-Z]/!\L&/g' <<<"$1"; }
ok=0 n=0
while read -r mod ver; do
  db=$(nice -n 19 ionice -c3 curl -fsS --retry 2 "https://sum.golang.org/lookup/$(esc "$mod")@$(esc "$ver")") || db=""
  while read -r line; do
    n=$((n + 1))
    if grep -qxF -- "$line" <<<"$db"; then ok=$((ok + 1)); else echo "SUMDB FAIL $mod@$ver ${line##* }"; fi
  done < <(awk -v m="$mod" -v v="$ver" '$1 == m && ($2 == v || $2 == v "/go.mod")' "$sum")
done < <(awk '{sub("/go.mod$", "", $2); print $1, $2}' "$sum" | sort -u)
echo "SUMDB $([ $ok = $n ] && [ $n -gt 0 ] && echo OK || echo FAIL) $ok/$n"
[ $ok = $n ] && [ $n -gt 0 ]
