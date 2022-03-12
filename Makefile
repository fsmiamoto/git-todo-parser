GO := $(shell which go)

all: build

test:
	$(GO) test -v ./...


.PHONY: all test
