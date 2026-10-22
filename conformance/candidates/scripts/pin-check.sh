#!/usr/bin/env bash
# pin-check.sh [go.mod]: is each pin still upstream's latest (rulings-2 #2)? gitpod and magnet: the repo's HEAD sha12
# must equal the pseudo-version's; x/vt: the last commit touching vt/ must be an ancestor of the pin, or the pin
# itself. Prints one PIN line per engine (the pin.txt format); exit 1 when any is MOVED. A `replace` of a candidate
# module in go.mod (any form: whole module, one version, a local path, in a replace block) makes it REPLACED: the code
# scored is then not the pin, so exit 1 too.
set -uo pipefail
cd "$(dirname "$0")/.." || exit 2
gomod=${1:-go.mod}
pin() { grep -E "^[[:space:]]*$1 v" "$gomod" | grep -oE '[0-9a-f]{12}$'; }
replaced=$(awk '/^replace[ \t]*\(/ {b = 1; next} b && /^\)/ {b = 0; next} b {print $1} /^replace[ \t]+[^( \t]/ {print $2}' "$gomod")
rep() { grep -qxF "$1" <<<"$replaced" && echo REPLACED; }
out=0
for r in gitpod-io/xterm-go=gitpod 0magnet/xterm-go=magnet; do
  repo=${r%=*} e=${r#*=} p=$(pin "github.com/$repo")
  head=$(gh api "repos/$repo/commits/HEAD" --jq .sha | cut -c1-12)
  s=MOVED; [ -n "$p" ] && [ "$head" = "$p" ] && s=SAME
  s=$(rep "github.com/$repo" || echo $s)
  [ $s = SAME ] || out=1
  echo "PIN $e pinned=${p:-NONE} head=${head:-UNREADABLE} $s"
done
p=$(pin github.com/charmbracelet/x/vt)
last=$(gh api 'repos/charmbracelet/x/commits?path=vt&per_page=1' --jq '.[0].sha' | cut -c1-12)
s=MOVED; case $(gh api "repos/charmbracelet/x/compare/$last...$p" --jq .status) in ahead | identical) [ -n "$p" ] && s=SAME ;; esac
s=$(rep github.com/charmbracelet/x/vt || echo $s)
[ $s = SAME ] || out=1
echo "PIN xvt pinned=${p:-NONE} head=${last:-UNREADABLE} $s"
exit $out
