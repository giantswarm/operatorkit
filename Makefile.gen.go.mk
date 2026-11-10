# DO NOT EDIT. Generated with:
#
#    devctl
#
#    https://github.com/giantswarm/devctl/blob/30010fbc8809ecb0f1422ee655a8105c299695b1/pkg/gen/input/makefile/internal/file/Makefile.gen.go.mk.template
#

APPLICATION    := $(shell go list -m | cut -d '/' -f 3)
BUILDTIMESTAMP := $(shell date -u '+%FT%TZ')
GITSHA1        := $(shell git rev-parse --verify HEAD)
MODULE         := $(shell go list -m)
# Build the whole main package, including sibling files, at the root or in cmd/.
MAIN_SOURCE    := $(shell if test -e cmd/main.go; then echo ./cmd; else echo .; fi)
OS             := $(shell go env GOOS)
ARCH           := $(shell go env GOARCH)
SOURCES        := $(shell find . -name '*.go')
# Older gitsemver releases cannot open a worktree's .git file. Fall back to
# local version tags; without one, leave VERSION empty so shipping targets fail.
VERSION        := $(shell if version=$$(gitsemver get 2>/dev/null) && [ -n "$$version" ]; then \
  printf '%s\n' "$$version"; \
  else git describe --tags --match 'v[0-9]*.[0-9]*.[0-9]*' 2>/dev/null | sed 's/^v//'; fi)
# Expanded first in every recipe that stamps VERSION into its output, so failed
# lookups fail the build instead of producing an unversioned binary.
# Targets that do not ship a version (test, lint, ...) never expand it.
REQUIRE_VERSION = $(if $(VERSION),,$(error VERSION is empty: `gitsemver get` and `git describe` found no version in $(CURDIR); pass VERSION=<version> to override))
ifeq ($(OS), linux)
EXTLDFLAGS := -static
endif
LDFLAGS        ?= -w -linkmode 'auto' -extldflags '$(EXTLDFLAGS)' \
  -X '$(MODULE)/pkg/project.version=$(VERSION)' \
  -X '$(MODULE)/pkg/project.buildTimestamp=$(BUILDTIMESTAMP)' \
  -X '$(MODULE)/pkg/project.gitSHA=$(GITSHA1)'

# Copy to a unique sibling before renaming so a running binary keeps its inode
# and readers only see complete files. Keep the versioned artifact for packaging.
INSTALL_BINARY = tmp=$$(mktemp "$(2).tmp.XXXXXX") || exit $$?; \
	trap 'rm -f "$$tmp"' EXIT; \
	trap 'exit 1' HUP INT TERM; \
	cp -a "$(1)" "$$tmp" && mv -f "$$tmp" "$(2)"

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
	@$(call INSTALL_BINARY,$<,$@)

$(APPLICATION)-darwin: $(APPLICATION)-v$(VERSION)-darwin-amd64
	@echo "====> $@"
	@$(call INSTALL_BINARY,$<,$@)

$(APPLICATION)-darwin-arm64: $(APPLICATION)-v$(VERSION)-darwin-arm64
	@echo "====> $@"
	@$(call INSTALL_BINARY,$<,$@)

$(APPLICATION)-linux: $(APPLICATION)-v$(VERSION)-linux-amd64
	@echo "====> $@"
	@$(call INSTALL_BINARY,$<,$@)

$(APPLICATION)-linux-amd64: $(APPLICATION)-v$(VERSION)-linux-amd64
	@echo "====> $@"
	@$(call INSTALL_BINARY,$<,$@)

$(APPLICATION)-linux-arm64: $(APPLICATION)-v$(VERSION)-linux-arm64
	@echo "====> $@"
	@$(call INSTALL_BINARY,$<,$@)

$(APPLICATION)-windows-amd64.exe: $(APPLICATION)-v$(VERSION)-windows-amd64.exe
	@echo "====> $@"
	@$(call INSTALL_BINARY,$<,$@)

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
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o "$(shell go env GOPATH)/bin/$(APPLICATION)" "$(MAIN_SOURCE)"

.PHONY: run
run: ## Runs the main package with go run.
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
lint: ## Checks import formatting and runs golangci-lint (installs goimports if missing).
	@echo "====> $@"
	@tool=$$(command -v goimports) || { \
		bindir=$$(go env GOBIN) || exit $$?; \
		if [ -z "$$bindir" ]; then \
			gopath=$$(go env GOPATH) || exit $$?; \
			bindir="$$gopath/bin"; \
		fi; \
		GOBIN="$$bindir" go install golang.org/x/tools/cmd/goimports@v0.50.0 || exit $$?; \
		tool="$$bindir/goimports"; \
	}; \
	files=$$("$$tool" -local "$(MODULE)" -l .) || exit $$?; \
	if [ -n "$$files" ]; then \
		printf '%s\n' "$$files"; \
		echo "Run 'make imports' to fix import formatting."; \
		exit 1; \
	fi
	golangci-lint run --timeout=15m ./...

# Extend lint from Makefile.custom.mk without replacing its generated recipe.
.PHONY: lint-custom
lint: lint-custom

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
build-docker: $(APPLICATION)-linux-$(ARCH) ## Builds docker image to registry for the host architecture.
	$(REQUIRE_VERSION)
	@echo "====> $@"
	@$(call INSTALL_BINARY,$<,$(APPLICATION))
	docker build -t ${APPLICATION}:${VERSION} .
