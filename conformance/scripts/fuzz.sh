#!/usr/bin/env bash
# fuzz.sh [seconds=600] [seed=epoch]: the differential fuzz against @xterm/headless 6.0.0 (gh #9 A2).
# Host: builds cmd/fuzzgen and cmd/conformance (static, CGO_ENABLED=0) into conformance/.cache/fuzz/bin, under nice.
# Then ONE node container (node.sh: --network none, --cpu-shares 256) runs scripts/fuzz-loop.sh until the time is
# spent: fuzzgen batch of 100 -> FUZZ_SHARDS (default 4) x gen/fuzz.mjs --shard -> --merge -> conformance -fuzz.
# Needs conformance/gen/node_modules (README "Running it"). Last line: FUZZ seconds= batches= cases= variants= hsens=
# diffs=; exit 0 no diff, 1 a diff (batches kept as conformance/.cache/fuzz/fail-<seed>), 2 nothing ran or unmeasured.
set -euo pipefail
secs=${1:-600} seed=${2:-$(date +%s)} shards=${FUZZ_SHARDS:-4}
case "$secs$seed$shards" in *[!0-9]*) echo "usage: fuzz.sh [seconds] [seed] (integers; FUZZ_SHARDS too)"; exit 2 ;; esac
repo=$(cd "$(dirname "$0")/../.." && pwd)
bin=$repo/conformance/.cache/fuzz/bin
mkdir -p "$bin"
(cd "$repo" && CGO_ENABLED=0 nice -n 19 ionice -c3 go build -o "$bin/" ./conformance/cmd/fuzzgen ./conformance/cmd/conformance)
rc=0
bash "$repo/conformance/scripts/node.sh" conformance "sh scripts/fuzz-loop.sh $secs $seed $shards" || rc=$?
exit "$rc"
