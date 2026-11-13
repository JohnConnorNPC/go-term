# fuzz-loop.sh <seconds> <seed> <shards>: fuzz.sh's loop, inside the node container, cwd conformance/ (POSIX sh).
# The time is checked between batches, so a run takes at least <seconds> plus the last batch (~90 s for 100 cases).
# Batch b uses seed+b. A batch that scores clean is deleted; one with a diff is moved to .cache/fuzz/fail-<seed+b>;
# a batch whose goldens or score are unmeasured stops the run (exit 2), kept as .cache/fuzz/bad-<seed+b>.
set -u
secs=$1 seed=$2 shards=$3
bin=.cache/fuzz/bin run=.cache/fuzz/run-$seed
rm -rf "$run" && mkdir -p "$run"
start=$(date +%s) b=0 cases=0 variants=0 hsens=0 diffs=0 bad=0
while [ $(($(date +%s) - start)) -lt "$secs" ]; do
  s=$((seed + b)) dir=$run/b$b
  if ! "$bin/fuzzgen" -seed "$s" -n 100 -root "$dir" -src . >"$run/fuzzgen.log" 2>&1; then
    echo "FUZZ batch $s: fuzzgen failed: $(tail -n 1 "$run/fuzzgen.log")"; bad=1; break
  fi
  i=0 pids=''
  while [ "$i" -lt "$shards" ]; do
    (cd gen && exec node fuzz.mjs --root "../$dir" --shard "$i/$shards") >"$dir/shard-$i.log" 2>&1 &
    pids="$pids $!" i=$((i + 1))
  done
  ok=1
  for p in $pids; do wait "$p" || ok=0; done
  if [ "$shards" = 1 ]; then
    cp "$dir/shard-0.log" "$dir/merge.log"                 # one shard writes manifest.json itself
  elif [ "$ok" = 1 ]; then
    (cd gen && exec node fuzz.mjs --root "../$dir" --merge "$shards") >"$dir/merge.log" 2>&1 || ok=0
  fi
  if [ "$ok" = 0 ]; then
    echo "FUZZ batch $s: goldens failed: $(cat "$dir"/shard-*.log "$dir/merge.log" 2>/dev/null | grep -m 1 FAIL)"
    mv "$dir" ".cache/fuzz/bad-$s"; bad=1; break
  fi
  rc=0
  "$bin/conformance" -root "$dir" -engine goterm -fuzz >"$dir/score.log" 2>&1 || rc=$?
  total=$(tail -n 1 "$dir/score.log")                      # TOTAL <pass>/<n> PASS|FAIL hsens=<k>
  p=$(echo "$total" | sed -n 's/^TOTAL \([0-9]*\)\/[0-9]* [A-Z]* hsens=[0-9]*$/\1/p')
  n=$(echo "$total" | sed -n 's/^TOTAL [0-9]*\/\([0-9]*\) [A-Z]* hsens=[0-9]*$/\1/p')
  k=$(echo "$total" | sed -n 's/^TOTAL .* hsens=\([0-9]*\)$/\1/p')
  v=$(sed -n 's/^fuzz-gen ok .* variants=\([0-9]*\) .*$/\1/p' "$dir/merge.log")
  if [ "$rc" -gt 1 ] || [ -z "$p" ] || [ -z "$n" ] || [ -z "$v" ] || [ "$n" != 100 ] || [ "$v" = 0 ]; then
    echo "FUZZ batch $s: unmeasured (exit $rc): $total"; mv "$dir" ".cache/fuzz/bad-$s"; bad=1; break
  fi
  cases=$((cases + n)) variants=$((variants + v)) hsens=$((hsens + k)) b=$((b + 1))
  if [ "$rc" = 0 ]; then
    rm -rf "$dir"
  else
    diffs=$((diffs + n - p))
    grep -v '^PASS ' "$dir/score.log" | head -n 20
    rm -rf ".cache/fuzz/fail-$s" && mv "$dir" ".cache/fuzz/fail-$s"
  fi
done
rm -rf "$run"                                              # kept batches were moved out
echo "FUZZ seconds=$(($(date +%s) - start)) batches=$b cases=$cases variants=$variants hsens=$hsens diffs=$diffs"
if [ "$bad" = 1 ] || [ "$b" = 0 ]; then exit 2; fi
[ "$diffs" = 0 ]
