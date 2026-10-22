#!/usr/bin/env bash
# fetch-mods.sh: host curl (no go process) of every module@version in go.sum and mods.seed (the candidates' own go.sum
# lines, plus each `missing` module a `cand.sh tidy` run names) into the file proxy ${CACHE}/goproxy (GOPROXY layout).
# A `/go.mod` line needs .info + .mod; any other line also needs the .zip. Files already present are re-checked
# (non-empty, and `unzip -tq` for zips); a bad one is fetched again. go.sum, in the container, verifies every hash.
set -uo pipefail
: "${CACHE:?cache dir (outside the repo)}"
cd "$(dirname "$0")/.." || exit 2
esc() { sed -E 's/[A-Z]/!\L&/g' <<<"$1"; }
ok=0 n=0
while read -r mod ver; do
  v=${ver%/go.mod} exts="info mod"; [ "$v" = "$ver" ] && exts="info mod zip"
  d=$CACHE/goproxy/$(esc "$mod")/@v n=$((n + 1)) bad=0
  mkdir -p "$d"
  for e in $exts; do
    f=$d/$(esc "$v").$e
    { [ -s "$f" ] && { [ "$e" != zip ] || unzip -tq "$f" >/dev/null 2>&1; }; } && continue
    nice -n 19 ionice -c3 curl -fsS --retry 2 -o "$f" "https://proxy.golang.org/$(esc "$mod")/@v/$(esc "$v").$e" || { rm -f "$f"; bad=1; echo "FETCH FAIL $mod@$v.$e"; }
  done
  [ $bad = 0 ] && ok=$((ok + 1))
done < <(cat go.sum mods.seed 2>/dev/null | awk 'NF >= 2 {print $1, $2}' | sort -u)
echo "FETCH $([ $ok = $n ] && [ $n -gt 0 ] && echo OK || echo FAIL) $ok/$n"
[ $ok = $n ] && [ $n -gt 0 ]
