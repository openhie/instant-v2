#!/bin/bash

cp -r ../cli ../cli-tmp
cp ./features/test-conf/* ../cli-tmp
cd ../cli-tmp

GOOS=darwin GOARCH=amd64 go build -o ../cli/features/test-platform-macos
GOOS=linux GOARCH=amd64 go build -o ../cli/features/test-platform-linux
GOOS=windows GOARCH=amd64 go build -o ../cli/features/test-platform.exe
go clean

cd ..
rm -rf cli-tmp
