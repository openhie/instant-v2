#!/bin/bash

FILE_PATH=$(
  cd "$(dirname "${BASH_SOURCE[0]}")" || exit
  pwd -P
)
export FILE_PATH

mkdir -p "$FILE_PATH"/bin

cd "$FILE_PATH"/src || exit

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o "$FILE_PATH"/bin/instant-macos
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "$FILE_PATH"/bin/instant-linux
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o "$FILE_PATH"/bin/instant-win.exe

go clean
