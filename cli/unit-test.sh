#!/bin/bash

FILE_PATH=$(
  cd "$(dirname "${BASH_SOURCE[0]}")" || exit
  pwd -P
)

cd "$FILE_PATH"/src/core || exit
go test .
