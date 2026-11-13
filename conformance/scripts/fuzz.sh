#!/usr/bin/env bash
# fuzz.sh [seconds=600] [seed=epoch]: the differential fuzz against @xterm/headless 6.0.0 (gh #9 A2).
# Host: builds cmd/fuzzgen and cmd/conformance (static, CGO_ENABLED=0) into conformance/.cache/fuzz/bin, under nice.
# Then ONE node container (node.sh: --network none, --cpu-shares 256) runs scripts/fuzz-loop.sh until the time is
# spent: fuzzgen batch of 100 -> FUZZ_SHARDS (default 4) x gen/fuzz.mjs --shard -> --merge -> conformance -fuzz.
# Needs conformance/gen/node_modules (README "Running it"). Last line: FUZZ seconds= batches= cases= variants= hsens=
# diffs=; exit 0 no diff, 1 a diff (batches kept as conformance/.cache/fuzz/fail-<seed>), 2 nothing ran or unmeasured
# (a usage error, a failed host build, a docker error such as a missing image under --pull never (125), or no FUZZ line).
# The container's output is also kept in conformance/.cache/fuzz/fuzz-<seed>.log.
set -euo pipefail
secs=${1-600} seed=${2-$(date +%s)} shards=${FUZZ_SHARDS-4}   # an empty argument is a usage error, not the default
case "$secs:$seed:$shards" in *[!0-9:]* | :* | *::* | *:) bad=1 ;; *) bad=0 ;; esac
if [ "$bad" = 1 ] || [ "$shards" -lt 1 ]; then echo "usage: fuzz.sh [seconds] [seed] (integers; FUZZ_SHARDS >= 1)"; exit 2; fi
repo=$(cd "$(dirname "$0")/../.." && pwd)
bin=$repo/conformance/.cache/fuzz/bin
mkdir -p "$bin"
if ! (cd "$repo" && CGO_ENABLED=0 nice -n 19 ionice -c3 go build -o "$bin/" ./conformance/cmd/fuzzgen ./conformance/cmd/conformance); then
  echo "FUZZ unmeasured: host go build failed"; exit 2
fi
log=$repo/conformance/.cache/fuzz/fuzz-$seed.log
rc=0
bash "$repo/conformance/scripts/node.sh" conformance "sh scripts/fuzz-loop.sh $secs $seed $shards" 2>&1 | tee "$log" || rc=$?
# Only the loop's own 0/1/2 with its FUZZ line count: docker's 125-127, a kill (137) or a missing line are unmeasured.
last=$(grep '^FUZZ seconds=' "$log" | tail -n 1 || true)
case "$rc:$last" in
  0:*diffs=0 | 1:*diffs=[1-9]* | 2:FUZZ*) exit "$rc" ;;
esac
echo "FUZZ unmeasured: container exit $rc, last FUZZ line '${last}'"
exit 2
