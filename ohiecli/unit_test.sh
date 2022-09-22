#!/bin/bash

FILE_PATH=$(
  cd "$(dirname "${BASH_SOURCE[0]}")" || exit
  pwd -P
)

cd "$FILE_PATH"/cmd/docker || exit
go test .

cd "$FILE_PATH"/cmd/prompts || exit
go test .

cd "$FILE_PATH"/cmd/utils || exit
go test .
