#!/bin/sh

# formater
gofmt -w=true -s $(find . -type f -name '*.go' -not -path "./vendor/*")

# linter
go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck ./...

# tests & coverage
go test -coverprofile=_dist/coverage.out -v ./...
go tool cover -func=_dist/coverage.out

# check race conditions
go test -race ./...

# clean after self
go mod tidy