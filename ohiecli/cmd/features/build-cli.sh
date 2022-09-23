#!/bin/bash

cp -r ../../ohiecli ../../cli-tmp
cp ./features/test-conf/* ../../cli-tmp/cmd

cd ../../cli-tmp/cmd || exit

GOOS=darwin GOARCH=amd64 go build -o ../../ohiecli/cmd/features/test-platform-macos
GOOS=linux GOARCH=amd64 go build -o ../../ohiecli/cmd/features/test-platform-linux
GOOS=windows GOARCH=amd64 go build -o ../../ohiecli/cmd/features/test-platform.exe
go clean

cd ../.. || exit
rm -rf cli-tmp