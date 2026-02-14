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
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Infra

.PHONY: infra
infra: deps ## Runs pyinfra to setup the infrastructure.
	uv run pyinfra -y ./infra/inventory.py ./infra/deploy.py

.PHONY: infra-dry-run
infra-dry-run: deps ## Runs `make infra` in dry-run mode.
	uv run pyinfra -vvv --dry --diff ./infra/inventory.py ./infra/deploy.py

.PHONY: infra-upgrade
infra-upgrade: deps ## Runs pyinfra to upgrade the infrastructure (apt packages).
	uv run pyinfra -vvv -y ./infra/inventory.py ./infra/upgrade.py

.PHONY: infra-upgrade-dry-run
infra-upgrade-dry-run: deps ## Runs `make infra-upgrade-dry-run` in dry-run mode.
	uv run pyinfra -vvv --dry --diff ./infra/inventory.py ./infra/upgrade.py

##@ Docker

.PHONY: build-forgejo-runners-images
build-forgejo-runners-images:
	docker build \
		-f ./forgejo-runners/Dockerfile.node-trixie \
		-t kphoen/node-trixie .

.PHONY: publish-forgejo-runners-images
publish-forgejo-runners-images: build-forgejo-runners-images
	docker push kphoen/node-trixie

##@ Grafana

.PHONY: generate-dashboards
generate-dashboards: ## Generates Grafana dashboards.
	cd grafana/dashboards && go mod vendor && go run . ../resources/Dashboard

##@ Dependencies

.PHONY: deps
deps: check-binaries ## Verifies and installs the dependencies.

.PHONY: check-binaries
check-binaries: ## Check that the required binaries are present.
	@uv --version >/dev/null 2>&1 || (echo "ERROR: uv is required. See https://docs.astral.sh/uv/getting-started/installation/"; exit 1)
