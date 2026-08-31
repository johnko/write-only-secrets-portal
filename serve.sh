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
# Point to our simulaated AWS secretsmanager service for testing
export AWS_ENDPOINT_URL="http://127.0.0.1:3000/"
export AWS_ACCESS_KEY_ID="test"
export AWS_SECRET_ACCESS_KEY="test"
export AWS_REGION="ca-central-1"
mise exec go@$GOLANG_VERSION -- go run ./main.go
popd
