#!/usr/bin/env sh

set -e

export GOOS="linux"
export GOARCH="amd64"
export CGO_ENABLED="0"

for directory in ./lambda/*; do
  echo "Building $directory..."
  go build -o "./bin/$directory/main" "$directory"
  ../../../../bin/build-lambda-zip -o "./bin/$directory/main.zip" "./bin/$directory/main"
done

for i in "$@"; do
  if [ "$i" = "--push" ]; then
    (export GOOS="windows" && export GOARCH="amd64" && cd ./pulumi && pulumi up)
    break
  fi
done
