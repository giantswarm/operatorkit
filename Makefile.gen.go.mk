# DO NOT EDIT. Generated with:
#
#    devctl@4.3.0
#

APPLICATION    := $(shell go list -m | cut -d '/' -f 3)
BUILDTIMESTAMP := $(shell date -u '+%FT%TZ')
GITSHA1        := $(shell git rev-parse --verify HEAD)
MODULE         := $(shell go list -m)
OS             := $(shell go env GOOS)
SOURCES        := $(shell find . -name '*.go')
VERSION        := $(shell architect project version)
ifeq ($(OS), linux)
EXTLDFLAGS := -static
endif
LDFLAGS        ?= -w -linkmode 'auto' -extldflags '$(EXTLDFLAGS)' \
  -X '$(shell go list -m)/pkg/project.buildTimestamp=${BUILDTIMESTAMP}' \
  -X '$(shell go list -m)/pkg/project.gitSHA=${GITSHA1}'

.DEFAULT_GOAL := build

.PHONY: build build-darwin build-linux
## build: builds a local binary
build: $(APPLICATION)
	@echo "====> $@"
## build-darwin: builds a local binary for darwin/amd64
build-darwin: $(APPLICATION)-darwin
	@echo "====> $@"
## build-linux: builds a local binary for linux/amd64
build-linux: $(APPLICATION)-linux
	@echo "====> $@"

$(APPLICATION): $(APPLICATION)-v$(VERSION)-$(OS)-amd64
	@echo "====> $@"
	cp -a $< $@

$(APPLICATION)-darwin: $(APPLICATION)-v$(VERSION)-darwin-amd64
	@echo "====> $@"
	cp -a $< $@

$(APPLICATION)-linux: $(APPLICATION)-v$(VERSION)-linux-amd64
	@echo "====> $@"
	cp -a $< $@

$(APPLICATION)-v$(VERSION)-%-amd64: $(SOURCES)
	@echo "====> $@"
	CGO_ENABLED=0 GOOS=$* GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $@ .

.PHONY: install
## install: install the application
install:
	@echo "====> $@"
	go install -ldflags "$(LDFLAGS)" .

.PHONY: run
## run: runs go run main.go
run:
	@echo "====> $@"
	go run -ldflags "$(LDFLAGS)" -race .

.PHONY: clean
## clean: cleans the binary
clean:
	@echo "====> $@"
	rm -f $(APPLICATION)*
	go clean

.PHONY: imports
## imports: runs goimports
imports:
	@echo "====> $@"
	goimports -local $(MODULE) -w .

.PHONY: lint
## lint: runs golangci-lint
lint:
	@echo "====> $@"
	golangci-lint run -E gosec -E goconst --timeout=15m ./...

.PHONY: test
## test: runs go test with default values
test:
	@echo "====> $@"
	go test -ldflags "$(LDFLAGS)" -race ./...

.PHONY: build-docker
## build-docker: builds docker image to registry
build-docker: build-linux
	@echo "====> $@"
	cp -a $(APPLICATION)-linux $(APPLICATION)
	docker build -t ${APPLICATION}:${VERSION} .
