#!/bin/bash
# This script builds Linux, Windows and MacOS binaries into the cli/bin directory
mkdir -p bin

cd src || exit

GOOS=darwin GOARCH=amd64 go build -o ../bin/gocli-macos
GOOS=linux GOARCH=amd64 go build -o ../bin/gocli-linux
GOOS=windows GOARCH=amd64 go build -o ../bin/gocli.exe
go clean
