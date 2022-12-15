#!/bin/bash

cp -r ../../cli ../../cli-tmp
cp ./features/test-conf/* ../../cli-tmp/src

cd ../../cli-tmp/src || exit

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ../../cli/src/features/test-platform-macos
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../../cli/src/features/test-platform-linux
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ../../cli/src/features/test-platform.exe
go clean

cd ../.. || exit
rm -rf cli-tmp
