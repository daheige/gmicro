#!/usr/bin/env bash

goci=$(which "golangci-lint")
if [ -z $goci ]; then
    echo "please install golangci-lint use cmd: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
    exit 0
fi

echo "current golangci-lint version"
golangci-lint version

echo "golangci-lint check begin"
go mod tidy
golangci-lint run ./... > golangci.log

echo "golangci check success,log write into golangci.log"

exit 0
