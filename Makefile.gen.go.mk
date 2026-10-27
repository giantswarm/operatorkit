# DO NOT EDIT. Generated with:
#
#    devctl
#
#    https://github.com/giantswarm/devctl/blob/c64a781d7ee888be4af92da9cce27c321300a589/pkg/gen/input/makefile/internal/file/Makefile.gen.go.mk.template
#

APPLICATION    := $(shell go list -m | cut -d '/' -f 3)
BUILDTIMESTAMP := $(shell date -u '+%FT%TZ')
GITSHA1        := $(shell git rev-parse --verify HEAD)
MODULE         := $(shell go list -m)
# main() is usually in `main.go`, but sometimes in `cmd/main.go` (for example in newer kubebuilder projects)
MAIN_SOURCE    := $(shell if test -e cmd/main.go; then echo cmd/main.go; else echo main.go; fi)
OS             := $(shell go env GOOS)
SOURCES        := $(shell find . -name '*.go')
VERSION        := $(shell gitsemver get)
# Expanded first in every recipe that stamps VERSION into its output, so a failed
# `gitsemver get` fails the build instead of producing an unversioned binary.
# Targets that do not ship a version (test, lint, ...) never expand it.
REQUIRE_VERSION = $(if $(VERSION),,$(error VERSION is empty: `gitsemver get` failed in $(CURDIR); pass VERSION=<version> to override))
ifeq ($(OS), linux)
EXTLDFLAGS := -static
endif
LDFLAGS        ?= -w -linkmode 'auto' -extldflags '$(EXTLDFLAGS)' \
  -X '$(MODULE)/pkg/project.version=$(VERSION)' \
  -X '$(MODULE)/pkg/project.buildTimestamp=$(BUILDTIMESTAMP)' \
  -X '$(MODULE)/pkg/project.gitSHA=$(GITSHA1)'

.DEFAULT_GOAL := build

##@ Go

.PHONY: build build-darwin build-darwin-64 build-linux build-linux-arm64 build-windows-amd64
build: $(APPLICATION) ## Builds a local binary.
	@echo "====> $@"
build-darwin: $(APPLICATION)-darwin ## Builds a local binary for darwin/amd64.
	@echo "====> $@"
build-darwin-arm64: $(APPLICATION)-darwin-arm64 ## Builds a local binary for darwin/arm64.
	@echo "====> $@"
build-linux: $(APPLICATION)-linux ## Builds a local binary for linux/amd64.
	@echo "====> $@"
build-linux-arm64: $(APPLICATION)-linux-arm64 ## Builds a local binary for linux/arm64.
	@echo "====> $@"
build-windows-amd64: $(APPLICATION)-windows-amd64.exe ## Builds a local binary for windows/amd64.
	@echo "====> $@"

$(APPLICATION): $(APPLICATION)-v$(VERSION)-$(OS)-amd64
	@echo "====> $@"
	cp -a $< $@

$(APPLICATION)-darwin: $(APPLICATION)-v$(VERSION)-darwin-amd64
	@echo "====> $@"
	cp -a $< $@

$(APPLICATION)-darwin-arm64: $(APPLICATION)-v$(VERSION)-darwin-arm64
	@echo "====> $@"
	cp -a $< $@

$(APPLICATION)-linux: $(APPLICATION)-v$(VERSION)-linux-amd64
	@echo "====> $@"
	cp -a $< $@

$(APPLICATION)-linux-arm64: $(APPLICATION)-v$(VERSION)-linux-arm64
	@echo "====> $@"
	cp -a $< $@

$(APPLICATION)-windows-amd64.exe: $(APPLICATION)-v$(VERSION)-windows-amd64.exe
	@echo "====> $@"
	cp -a $< $@

$(APPLICATION)-v$(VERSION)-%-amd64: $(SOURCES)
	$(REQUIRE_VERSION)
	@echo "====> $@"
	CGO_ENABLED=0 GOOS=$* GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS)" -o $@ "$(MAIN_SOURCE)"

$(APPLICATION)-v$(VERSION)-%-arm64: $(SOURCES)
	$(REQUIRE_VERSION)
	@echo "====> $@"
	CGO_ENABLED=0 GOOS=$* GOARCH=arm64 go build -trimpath -ldflags "$(LDFLAGS)" -o $@ "$(MAIN_SOURCE)"

$(APPLICATION)-v$(VERSION)-windows-amd64.exe: $(SOURCES)
	$(REQUIRE_VERSION)
	@echo "====> $@"
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS)" -o $@ "$(MAIN_SOURCE)"

.PHONY: install
install: ## Install the application.
	$(REQUIRE_VERSION)
	@echo "====> $@"
	go build -ldflags "$(LDFLAGS)" -o "$(shell go env GOPATH)/bin/$(APPLICATION)" "$(MAIN_SOURCE)"

.PHONY: run
run: ## Runs go run main.go.
	@echo "====> $@"
	go run -ldflags "$(LDFLAGS)" -race "$(MAIN_SOURCE)"

.PHONY: clean
clean: ## Cleans the binary.
	@echo "====> $@"
	rm -f $(APPLICATION)*
	go clean

.PHONY: imports
imports: ## Runs goimports.
	@echo "====> $@"
	goimports -local $(MODULE) -w .

.PHONY: lint
lint: ## Runs golangci-lint.
	@echo "====> $@"
	golangci-lint run --timeout=15m ./...

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

# `-deps ./...` lists only the packages actually compiled into this module, so nancy
# reports vulnerabilities in code we really ship. `-m all` walks the entire module
# graph instead and flags modules that are never built: on team-stamper it reported 7
# vulnerable modules of which 5 are not in the build at all, including golang.org/x/crypto
# and its 13 CVEs. #680 moved to `-deps ./...` for exactly this reason; #1964 moved back
# to `-m all` only to dodge nancy's 10 MB stdin cap, which nancy made configurable and
# defaulted to 100 MB in v2.1.0 -- so the workaround costs accuracy for nothing.
.PHONY: nancy
nancy: ## Runs nancy (requires v2.1.0 or newer).
	@echo "====> $@"
	CGO_ENABLED=0 go list -json -deps ./... | nancy sleuth --skip-update-check --quiet --exclude-vulnerability-file ./.nancy-ignore --additional-exclude-vulnerability-files ./.nancy-ignore.generated

# Race detector needs a C toolchain. The architect CI image has none and runs
# with CGO_ENABLED=0, so degrade to cgo-free there; everywhere a compiler exists
# (laptops, coding agents, GitHub Actions, any cgo-capable runner) keeps -race.
RACE := $(shell { [ "$${CGO_ENABLED:-1}" != "0" ] && { command -v gcc || command -v clang; } >/dev/null 2>&1; } && echo -race)

.PHONY: test
test: ## Runs go test with default values (race detector when a C toolchain is available).
	@echo "====> $@"
	go test -ldflags "$(LDFLAGS)" $(RACE) ./...

.PHONY: build-docker
build-docker: build-linux ## Builds docker image to registry.
	$(REQUIRE_VERSION)
	@echo "====> $@"
	cp -a $(APPLICATION)-linux $(APPLICATION)
	docker build -t ${APPLICATION}:${VERSION} .
