# Sourced by node.sh and go.sh: run_container <image> <workdir relative to the repo> '<sh command>' [docker -e args via ENVS].
# One foreground container: no network, repo read-only except conformance/, caches under CACHE.
run_container() {
  local image=${1:?image} wd=${2:?workdir relative to the repo} cmd=${3:?command}
  local repo; repo=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)
  local cache=${CACHE:-/home/debian/agentwork/goterm/survey/cache}
  mkdir -p "${cache:?}/home"
  docker run --rm --cpus 4 --memory 8g --network none \
    --label gm.project=goterm --name goterm-w3-1 \
    --user "$(id -u):$(id -g)" -e HOME=/cache/home "${ENVS[@]}" \
    -v "${repo:?}:/repo:ro" -v "${repo:?}/conformance:/repo/conformance" -v "${cache:?}:/cache" \
    -w "/repo/${wd}" "$image" nice -n 10 sh -c "$cmd"
}
