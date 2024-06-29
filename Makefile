.PHONY: deps build

TAG=0.1.0

deps:
	go mod tidy

test:
	./_bin/test.sh

tag:
	./_bin/tag.sh ${TAG}
