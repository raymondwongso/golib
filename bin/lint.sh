#!/bin/bash

set -euo pipefail

for dir in $(find . -maxdepth 1 -type d); do
  dirname=$(basename "$dir")
  if [ "$dirname" != "." ] && [[ "$dirname" != .* ]] && [[ "$dirname" != bin ]] && [[ "$dirname" != dev ]]; then
    cd "$dirname"

    for file in $(find . -type f -name "*.go" -not -path "*vendor/*"); do
      goimports -w $file
      gofmt -s -w $file
    done

    golangci-lint -v run ./...
  fi
done
