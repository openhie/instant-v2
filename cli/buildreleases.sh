#!/bin/bash
# This script builds Linux, Windows and MacOS binaries into the cli/bin directory
mkdir -p bin
GOOS=darwin GOARCH=amd64 go build && mv ./cli ./bin/gocli-macos &&
    GOOS=linux GOARCH=amd64 go build && mv ./cli ./bin/gocli-linux &&
    GOOS=windows GOARCH=amd64 go build && mv ./cli.exe ./bin/gocli.exe && go clean
