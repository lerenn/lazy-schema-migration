.DEFAULT_GOAL  := help

DAGGER_COMMAND := _EXPERIMENTAL_DAGGER_INTERACTIVE_TUI=true dagger run go run ./ci

.PHONY: ci
ci: ## Run the CI
	@${DAGGER_COMMAND} all

.PHONY: lint
lint: ## Lint the code
	@${DAGGER_COMMAND} linter

.PHONY: test
test: ## Perform tests
	@${DAGGER_COMMAND} test

.PHONY: help
help: ## Display this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
