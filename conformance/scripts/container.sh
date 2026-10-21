# Sourced by node.sh and go.sh: run_container <image> <workdir relative to the repo> '<sh command>' [docker -e args via ENVS].
# One foreground container: no network, repo read-only except conformance/, caches under CACHE.
# GOTERM_NAME (default goterm-<user>-<pid>) names it; CACHE (default conformance/.cache, gitignored) holds HOME and caches.
run_container() {
  local image=${1:?image} wd=${2:?workdir relative to the repo} cmd=${3:?command}
  local repo; repo=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)
  local cache=${CACHE:-$repo/conformance/.cache} name=${GOTERM_NAME:-goterm-$(id -un)-$$}
  mkdir -p "${cache:?}/home"
  docker run --rm --cpus 4 --cpu-shares 256 --memory 8g --network none \
    --label gm.project=goterm --name "${name:?}" \
    --user "$(id -u):$(id -g)" -e HOME=/cache/home "${ENVS[@]}" \
    -v "${repo:?}:/repo:ro" -v "${repo:?}/conformance:/repo/conformance" -v "${cache:?}:/cache" \
    -w "/repo/${wd}" "$image" nice -n 10 sh -c "$cmd"
}
