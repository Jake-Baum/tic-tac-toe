#!/usr/bin/env sh

export GOOS="linux"
export GOARCH="amd64"
export CGO_ENABLED="0"

for directory in ./lambda/*/; do
  echo "Building $directory..."
  go build -o "./bin/$directory/main" "$directory"
  ../../../../bin/build-lambda-zip -o "./bin/$directory/main.zip" "./bin/$directory/main"
done
