# Directories.
API_DIR := api
CRD_DIR := config/crd
SCRIPTS_DIR := hack
TOOLS_DIR := $(SCRIPTS_DIR)/tools
TOOLS_BIN_DIR := $(abspath $(TOOLS_DIR)/bin)

# Binaries.
# Need to use abspath so we can invoke these from subdirectories
CONTROLLER_GEN := $(abspath $(TOOLS_BIN_DIR)/controller-gen)

BUILD_COLOR = ""
GEN_COLOR = ""
NO_COLOR = ""

ifneq (, $(shell command -v tput))
ifeq ($(shell test `tput colors` -ge 8 && echo "yes"), yes)
BUILD_COLOR = \033[0;34m
GEN_COLOR = \033[0;32m
NO_COLOR = \033[0m
endif
endif

DEEPCOPY_BASE = zz_generated.deepcopy
MODULE = $(shell go list -m)
BOILERPLATE = $(SCRIPTS_DIR)/boilerplate.go.txt
YEAR = $(shell date +'%Y')

DEEPCOPY_FILES := $(shell find $(API_DIR) -name $(DEEPCOPY_BASE).go)

all: generate

$(CONTROLLER_GEN): $(TOOLS_DIR)/controller-gen/go.mod
	@echo "$(BUILD_COLOR)Building controller-gen$(NO_COLOR)"
	cd $(TOOLS_DIR)/controller-gen \
 	&& go build -tags=tools -o $(CONTROLLER_GEN) sigs.k8s.io/controller-tools/cmd/controller-gen

.PHONY: generate
generate:
	@$(MAKE) generate-deepcopy
	@$(MAKE) generate-manifests

.PHONY: verify
verify:
	@$(MAKE) clean-generated
	@$(MAKE) generate
	git diff --exit-code

.PHONY: generate-deepcopy
generate-deepcopy: $(CONTROLLER_GEN)
	@echo "$(GEN_COLOR)Generating deepcopy$(NO_COLOR)"
	$(CONTROLLER_GEN) \
	object:headerFile=$(BOILERPLATE),year=$(YEAR) \
	paths=./$(API_DIR)/...

.PHONY: generate-manifests
generate-manifests: $(CONTROLLER_GEN)
	@echo "$(GEN_COLOR)Generating CRDs$(NO_COLOR)"
	$(CONTROLLER_GEN) \
	crd \
	paths=./$(API_DIR)/... \
	output:dir="./$(CRD_DIR)"

.PHONY: clean-generated
clean-generated:
	@echo "$(GEN_COLOR)Cleaning generated files$(NO_COLOR)"
	rm -rf $(CRD_DIR) $(DEEPCOPY_FILES)

.PHONY: clean-tools
clean-tools:
	@echo "$(GEN_COLOR)Cleaning tools$(NO_COLOR)"
	rm -rf $(TOOLS_BIN_DIR)
