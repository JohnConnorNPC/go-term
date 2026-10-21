#!/usr/bin/env bash
# node.sh <workdir relative to the repo> '<sh command>': the pinned node image (1/pins.md); globs expand inside the container.
set -euo pipefail
IMAGE_DIGEST=sha256:7cd3fbc830c75c92256fe1122002add9a1c025831af8770cd0bf8e45688ef661
ENVS=(-e "IMAGE_DIGEST=${IMAGE_DIGEST:?}" -e npm_config_cache=/cache/npm -e npm_config_update_notifier=false)
. "$(dirname "$0")/container.sh"
run_container "node:20.19.2-bookworm-slim@${IMAGE_DIGEST:?}" "$@"
