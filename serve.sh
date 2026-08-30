#!/usr/bin/env bash
set -euo pipefail

# open "http://localhost:8888/"

# renovate: datasource=golang-version depName=golang packageName=golang
GOLANG_VERSION=1.27.0

WORKDIR=$(dirname "$0")
pushd "$WORKDIR"
cd ./src/
mise exec go@$GOLANG_VERSION -- go clean
mise exec go@$GOLANG_VERSION -- go fmt ./main.go
mise exec go@$GOLANG_VERSION -- go run ./main.go
popd
