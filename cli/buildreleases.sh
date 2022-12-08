#!/bin/bash
# This script builds Linux, Windows and MacOS binaries into the cli/bin directory
mkdir -p bin

cd src || exit

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ../bin/cli
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../bin/cli
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ../bin/cli
go clean
