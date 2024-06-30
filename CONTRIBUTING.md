# Contributing

## Guidelines for pull requests

- Write tests for any changes.
- Separate unrelated changes into multiple pull requests.
- For bigger changes, make sure you start a discussion first by creating an issue and explaining the intended change.

## Requirements

* [Go](https://golang.org/dl/)

## Release

1. All notable changes comming with the new version should be documented in [CHANGELOG.md](https://raw.githubusercontent.com/zoomio/microbatch/main/CHANGELOG.md).
2. Run tests with `./_bin/test.sh`, make sure everything is passing.
3. Bump the `TAG` variable inside the `Makefile` to the desired version, 
4. Push and trigger new release on GitHub via `make tag`.