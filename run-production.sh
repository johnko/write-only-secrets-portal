#!/usr/bin/env bash
set -exuo pipefail

# open "http://localhost:8888/"

# renovate: datasource=golang-version depName=golang packageName=golang
GOLANG_VERSION=1.27.0

WORKDIR=$(dirname "$0")
pushd "$WORKDIR"
cd ./src/
mise exec go@$GOLANG_VERSION -- go clean
if [ ! -e go.mod ]; then
  mise exec go@$GOLANG_VERSION -- go mod init main
  mise exec go@$GOLANG_VERSION -- go get github.com/aws/aws-sdk-go-v2/aws
  mise exec go@$GOLANG_VERSION -- go get github.com/aws/aws-sdk-go-v2/config
  mise exec go@$GOLANG_VERSION -- go get github.com/aws/aws-sdk-go-v2/service/secretsmanager
else
  mise exec go@$GOLANG_VERSION -- go mod download
fi
mise exec go@$GOLANG_VERSION -- go fmt ./main.go
mise exec go@$GOLANG_VERSION -- go mod tidy
mise exec go@$GOLANG_VERSION -- go run ./main.go
popd
