#!/usr/bin/env bash
# go.sh <workdir relative to the repo> '<sh command>': the pinned golang image (1/pins.md); no module downloads.
set -euo pipefail
IMAGE_DIGEST=sha256:21ade014faf952365b9ecbd8dbdb87176fbacb82fac9fb132d9f2a7b8cc3f132
ENVS=(-e "IMAGE_DIGEST=${IMAGE_DIGEST:?}" -e GOTOOLCHAIN=local -e GOPROXY=off -e GOFLAGS=-mod=readonly
      -e GOCACHE=/cache/go-build -e GOPATH=/cache/gopath)
. "$(dirname "$0")/container.sh"
run_container "golang:1.27.3-bookworm@${IMAGE_DIGEST:?}" "$@"
