#!/usr/bin/env bash
# gh-spdx.sh: GitHub's licence detection per candidate repo, appended to lic.txt as INFO (`LIC <e> gh-spdx INFO <id>`).
set -uo pipefail
for r in gitpod-io/xterm-go=gitpod 0magnet/xterm-go=magnet charmbracelet/x=xvt; do
  id=$(gh api "repos/${r%=*}/license" --jq .license.spdx_id 2>/dev/null) || id=UNREADABLE
  echo "LIC ${r#*=} gh-spdx INFO ${id:-UNREADABLE}"
done
