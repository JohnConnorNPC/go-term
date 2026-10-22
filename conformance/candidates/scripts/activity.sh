#!/usr/bin/env bash
# activity.sh [repo=engine ...]: report-only upstream activity (rulings-0; no gate): five ACT lines per engine, the last
# commit date (x/vt: the last one touching vt/; the other counts are charmbracelet/x repo-wide), open issues, open
# PRs, contributors and archived. A value the API doesn't give is UNREADABLE.
set -uo pipefail
[ $# -gt 0 ] || set -- gitpod-io/xterm-go=gitpod 0magnet/xterm-go=magnet charmbracelet/x=xvt
v() { local x; x=$(gh api "$@" 2>/dev/null) && [ -n "$x" ] && echo "$x" || echo UNREADABLE; }
for r in "$@"; do
  repo=${r%=*} e=${r#*=} path=; [ "$e" = xvt ] && path='&path=vt'
  prs=$(v "repos/$repo/pulls?state=open&per_page=100" --jq length) all=$(v "repos/$repo" --jq .open_issues_count)
  issues=UNREADABLE; [ "$prs" != UNREADABLE ] && [ "$all" != UNREADABLE ] && issues=$((all - prs))
  echo "ACT $e last-commit $(v "repos/$repo/commits?per_page=1$path" --jq '.[0].commit.committer.date')"
  echo "ACT $e open-issues $issues"
  echo "ACT $e open-prs $prs"
  echo "ACT $e contributors $(v "repos/$repo/contributors?per_page=100" --jq length)"
  echo "ACT $e archived $(v "repos/$repo" --jq .archived)"
done
