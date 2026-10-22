#!/usr/bin/env bash
# cand.sh <run|tidy|run123> '<sh command>': the candidates module (conformance/candidates) in the pinned golang image, offline.
# Modules come only from the host-fetched file proxy ${CACHE}/goproxy (fetch-mods.sh); go.sum verifies every hash.
# run123 is run on the pinned golang:1.23 image (rulings-4). TMO (default 1800 s) bounds the command: exit 124 = UNMEASURED.
# The container is ${GOTERM_NAME:-goterm-$USER}-<pid>, labelled gm.project=goterm (container.sh).
set -euo pipefail
: "${REPO:?repo root}" "${CACHE:?cache dir (outside the repo)}"
case ${1:?mode} in
  run | tidy) IMAGE=golang:1.27.3-bookworm@sha256:21ade014faf952365b9ecbd8dbdb87176fbacb82fac9fb132d9f2a7b8cc3f132 ;;
  run123) IMAGE=golang:1.23.12-bookworm@sha256:167053a2bb901972bf2c1611f8f52c44d5fe7e762e5cab213708d82c421614db ;;
  *) echo "cand.sh: mode $1 is not run|tidy|run123" >&2; exit 2 ;;
esac
mod=readonly; [ "$1" = tidy ] && mod=mod
ENVS=(-e GOTOOLCHAIN=local -e GOPROXY=file:///cache/goproxy -e GOFLAGS=-mod=$mod -e GOSUMDB=off -e GOCACHE=/cache/go-build
      -e GOPATH=/cache/gopath -e "TMO=${TMO:-1800}" -e "CAND_CMD=${2:?command}")
. "${REPO:?}/conformance/scripts/container.sh"
CACHE=${CACHE:?} GOTERM_NAME=${GOTERM_NAME:-goterm-$USER}-$$ run_container "$IMAGE" conformance/candidates 'exec timeout "$TMO" sh -c "$CAND_CMD"'
