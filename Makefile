ORG_PATH=go_starter
PROJ=src
REPO_PATH=$(ORG_PATH)/$(PROJ)
export PATH := $(PWD)/bin:$(PATH)
THIS_DIRECTORY:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

VERSION ?= $(shell ./scripts/git-version)

$( shell mkdir -p bin )

export GOBIN=$(PWD)/bin

LD_FLAGS="-w -X $(ORG_PATH)/version.Version=$(VERSION)"

build: bin/go_starter

bin/go_starter:
	@go install -v -ldflags $(LD_FLAGS) $(REPO_PATH)/cmd/go_starter

.PHONY: release-binary
release-binary:
	@go build -o /go/bin/go_starter -v -ldflags $(LD_FLAGS) $(REPO_PATH)/cmd/go_starter

fmt:
	@./scripts/gofmt ./...

clean:
	@rm -rf bin/

testall: fmt

FORCE:

.PHONY: fmt testall
