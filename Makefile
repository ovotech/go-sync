SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

lint/golangci_lint: ## Lint using golangci-lint.
	golangci-lint run ./...

lint: ## Lint Go Sync.
	make lint/golangci_lint
.PHONY: lint lint/*

fix/gofmt: ## Fix formatting with gofmt.
	gofmt -w ./..

fix/goimports: ## Fix imports.
	goimports -w ./..

fix/golangci_lint: ## Fix golangci-lint errors.
	golangci-lint run --fix ./...

fix: ## Fix common linter errors.
	make fix/gofmt
	make fix/goimports
	make fix/golangci_lint
.PHONY: fix fix/*

generate/mockery: ## Generate Mockery mocks.
	rm -rf ./internal/mocks
	mockery --all --exported --with-expecter --output ./internal/mocks

generate: ## Generate automated code.
	make generate/mockery
.PHONY: generate generate/*

.DEFAULT_GOAL := help
help: Makefile ## Display list of available commands.
	@grep -E '(^[a-zA-Z_-]+:.*?##.*$$)|(^##)' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}' | sed -e 's/\[32m##/[33m/'
