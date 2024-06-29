.PHONY: deps build

TAG=0.2.0

deps:
	go mod tidy

test:
	./_bin/test.sh

tag:
	./_bin/tag.sh ${TAG}
