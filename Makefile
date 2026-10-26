# DO NOT EDIT. Generated with:
#
#    devctl
#
#    https://github.com/giantswarm/devctl/blob/c0a255e412bf450e122f71e563d74a9bd9f9cddf/pkg/gen/input/makefile/internal/file/Makefile.template
#

include Makefile.*.mk

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; n = 0; w = 0} /^[a-zA-Z%\\\/_0-9-]+:.*?##/ { order[n] = "t"; targets[n] = $$1; descs[n] = $$2; if (length($$1) > w) w = length($$1); n++ } /^##@/ { order[n] = "c"; cats[n] = substr($$0, 5); n++ } END { printf "\nUsage:\n  make \033[36m<target>\033[0m\n"; for (i = 0; i < n; i++) { if (order[i] == "c") printf "\n\033[1m%s\033[0m\n", cats[i]; else printf "  \033[36m%-*s\033[0m  %s\n", w, targets[i], descs[i] } }' $(MAKEFILE_LIST)
