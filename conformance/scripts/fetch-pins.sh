#!/usr/bin/env bash
# Host-side: download every tarball in gen/pins.lst into gen/vendor/ (if absent) and verify its sha512 against the pinned integrity.
# Plain file fetches only (docker-rules-3: containers have no network). Any mismatch exits 1.
set -euo pipefail
gen=$(cd "$(dirname "$0")/../gen" && pwd); vendor="${gen:?}/vendor"; mkdir -p "$vendor"
ok=0; total=0
while read -r name ver file url integ; do
  case "$name" in ''|'#'*) continue;; esac
  total=$((total+1)); f="$vendor/${file:?}"
  [ -s "$f" ] || curl -sSf --max-time 120 -o "$f" "${url:?}"
  got="sha512-$(openssl dgst -sha512 -binary "$f" | base64 -w0)"
  if [ "$got" = "${integ:?}" ]; then ok=$((ok+1)); else echo "FETCH FAIL integrity $name@$ver"; fi
done < "$gen/pins.lst"
[ "$total" -gt 0 ] && [ "$ok" -eq "$total" ] || { echo "FETCH FAIL $ok/$total"; exit 1; }
echo "FETCH OK $ok/$total"
