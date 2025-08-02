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

##@ Ansible

.PHONY: ansible
ansible: deps ## Runs ansible.
	ansible-playbook -e @./ansible/variables.var -e @./ansible/secrets.encrypted --ask-vault-pass -i ./ansible/inventory.yaml ./ansible/site.yml

.PHONY: ansible-dry-run
ansible-dry-run: deps ## Runs ansible in dry-run mode.
	ansible-playbook -e @./ansible/variables.var -e @./ansible/secrets.encrypted --ask-vault-pass -i ./ansible/inventory.yaml ./ansible/site.yml --check --diff

##@ Grafana

.PHONY: generate-dashboards
generate-dashboards: ## Generates Grafana dashboards.
	cd grafana/dashboards && go mod vendor && go run . ../resources/Dashboard

##@ Dependencies

.PHONY: deps
deps: ansible-deps ## Installs the dependencies.

.PHONY: ansible-deps
ansible-deps: check-binaries ## Installs ansible-related dependencies.
	ansible-galaxy role install -r ./ansible/requirements.yml
	ansible-galaxy collection install -r ./ansible/requirements.yml

.PHONY: check-binaries
check-binaries: ## Check that the required binaries are present.
	@ansible --version >/dev/null 2>&1 || (echo "ERROR: ansible is required."; exit 1)
	@ansible-galaxy --version >/dev/null 2>&1 || (echo "ERROR: ansible-galaxy is required."; exit 1)
