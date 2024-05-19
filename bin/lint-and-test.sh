#!/bin/bash

for dir in $(find . -maxdepth 1 -type d); do
  dirname=$(basename "$dir")
  if [ "$dirname" != "." ] && [[ "$dirname" != .* ]] && [[ "$dirname" != bin ]] && [[ "$dirname" != dev ]]; then
    cd "$dirname"
    goimports -w .
    gofmt -s -w .
    golangci-lint run ./...

    go test -coverprofile=coverage.out -timeout 30s ./...
    go tool cover -func=coverage.out
    go tool cover -html=coverage.out -o coverage.html
    rm coverage.out coverage.html
    cd ..
  fi
done
