# info
VERSION=v0.0.1

CGO_ENABLE=0
PROJECT_NAME="flowctl"
GO_VERSION=$(go version | awk '{print $3}')
BUILD_TIME=$(date +%FT%T%z)

all: linux-install

linux-install:
    GOOS=linux GOARCH=amd64 go install