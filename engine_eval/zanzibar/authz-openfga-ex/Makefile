MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
PROJECT_PATH := $(patsubst %/,%,$(dir $(MKFILE_PATH)))
DOCS_DIR := $(PROJECT_PATH)/docs

.DEFAULT_GOAL := help
SHELL = bash
CURL ?= curl
LOCAL_BIN_PATH := ${PROJECT_PATH}/bin
# Add the project-level bin directory into PATH. Needed in order
# for `go generate` to use project-level bin directory binaries first
export PATH := ${LOCAL_BIN_PATH}:$(PATH)
# The details of the application:
binary:=authz-openfga-ex

# The version needs to be different for each deployment because otherwise the
# cluster will not pull the new image from the internal registry:
version:=$(shell date +%s)

GO := go
GOFMT := gofmt
# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell $(GO) env GOBIN))
GOBIN=$(shell $(GO) env GOPATH)/bin
else
GOBIN=$(shell $(GO) env GOBIN)
endif

LOCAL_BIN_PATH := ${PROJECT_PATH}/bin

export PATH := ${LOCAL_BIN_PATH}:$(PATH)


binary:
	$(GO) build -o bin/$(binary) ./authz-openfga/main.go
.PHONY: binary

test:
	$(GO) test ./test -test.v
.PHONY: test

run:
	$(GO) run ./authz-openfga/main.go
.PHONY: run

# Examples:
#   make test TESTFLAGS="-run TestSomething"

# Prints a list of useful targets.
help:
	@echo "Kafka Service Fleet Manager make targets"
	@echo ""
	@echo "make binary                              compile binaries"
	@echo "make test                                run test"
	@echo "make run                                 run the application"
.PHONY: help