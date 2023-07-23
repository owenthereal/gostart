SHELL=/bin/bash -o pipefail

TOOLS_DIR ?= $(CURDIR)/bin
export PATH := $(TOOLS_DIR):$(PATH)
.PHONY: tools
tools:
	rm -rf $(TOOLS_DIR) && mkdir -p $(TOOLS_DIR)
	# go tools
	GOBIN=$(TOOLS_DIR) go generate -tags tools tools.go

.PHONY: build
build: gen
	go build -o bin/ ./cmd/...

.PHONY: gen
gen:
	go generate ./...
