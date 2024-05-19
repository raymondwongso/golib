#!/bin/bash

set -euo pipefail

for dir in $(find . -maxdepth 1 -type d); do
  dirname=$(basename "$dir")
  if [ "$dirname" != "." ] && [[ "$dirname" != .* ]] && [[ "$dirname" != bin ]] && [[ "$dirname" != dev ]]; then
    cd "$dirname"
    golangci-lint run ./...
  fi
done
