#!/usr/bin/env bash
# pin-check.sh [go.mod]: is each pin still upstream's latest (rulings-2 #2)? gitpod and magnet: the repo's HEAD sha12
# must equal the pseudo-version's; x/vt: the last commit touching vt/ must be an ancestor of the pin, or the pin
# itself. Prints one PIN line per engine (the pin.txt format); exit 1 when any is MOVED.
set -uo pipefail
cd "$(dirname "$0")/.." || exit 2
gomod=${1:-go.mod}
pin() { grep -E "^[[:space:]]*$1 v" "$gomod" | grep -oE '[0-9a-f]{12}$'; }
out=0
for r in gitpod-io/xterm-go=gitpod 0magnet/xterm-go=magnet; do
  repo=${r%=*} e=${r#*=} p=$(pin "github.com/$repo")
  head=$(gh api "repos/$repo/commits/HEAD" --jq .sha | cut -c1-12)
  s=MOVED; [ -n "$p" ] && [ "$head" = "$p" ] && s=SAME
  [ $s = SAME ] || out=1
  echo "PIN $e pinned=${p:-NONE} head=${head:-UNREADABLE} $s"
done
p=$(pin github.com/charmbracelet/x/vt)
last=$(gh api 'repos/charmbracelet/x/commits?path=vt&per_page=1' --jq '.[0].sha' | cut -c1-12)
s=MOVED; case $(gh api "repos/charmbracelet/x/compare/$last...$p" --jq .status) in ahead | identical) [ -n "$p" ] && s=SAME ;; esac
[ $s = SAME ] || out=1
echo "PIN xvt pinned=${p:-NONE} head=${last:-UNREADABLE} $s"
exit $out
