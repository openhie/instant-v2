#!/bin/bash

FILE_PATH=$(
  cd "$(dirname "${BASH_SOURCE[0]}")" || exit
  pwd -P
)

cd "$FILE_PATH"/src/core/parse || exit
go test .

cd "$FILE_PATH"/src/util/slice || exit
go test .

cd "$FILE_PATH"/src/util/file || exit
go test .
