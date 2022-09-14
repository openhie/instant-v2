#!/bin/bash

cp -r ../../cli ../../cli-tmp
cp ./features/test-conf/* ../../cli-tmp/cmd

cd ../../cli-tmp/cmd

GOOS=darwin GOARCH=amd64 go build -o ../../cli/cmd/features/test-platform-macos
GOOS=linux GOARCH=amd64 go build -o ../../cli/cmd/features/test-platform-linux
GOOS=windows GOARCH=amd64 go build -o ../../cli/cmd/features/test-platform.exe
go clean

cd ../..
rm -rf cli-tmp
