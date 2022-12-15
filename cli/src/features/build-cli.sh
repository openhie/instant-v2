#!/bin/bash

cp -r ../../cli ../../cli-tmp
cp ./features/test-conf/* ../../cli-tmp/src

cd ../../cli-tmp/src || exit

GOOS=darwin GOARCH=amd64 go build -o ../../cli/src/features/test-platform-macos
GOOS=linux GOARCH=amd64 go build -o ../../cli/src/features/test-platform-linux
GOOS=windows GOARCH=amd64 go build -o ../../cli/src/features/test-platform.exe
go clean

cd ../.. || exit
rm -rf cli-tmp
