SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules
TARGET ?= .
ADAPTERS:=$$(ls -d adapters/*/ | sed 's/\(.*\)/.\/\1.../')

lint/golangci_lint: ## Lint using golangci-lint.
	golangci-lint run ./... $(ADAPTERS)

lint: ## Lint Go Sync.
	make lint/golangci_lint
.PHONY: lint lint/*

fix/gofmt: ## Fix formatting with gofmt.
	gofmt -s -w .

fix/gci: ## Fix imports.
	gci write .

fix/golangci_lint: ## Fix golangci-lint errors.
	golangci-lint run --fix ./... $(ADAPTERS)

fix: ## Fix common linter errors.
	make fix/gofmt
	make fix/gci
	make fix/golangci_lint
.PHONY: fix fix/*

generate/mockery: ## Generate mocks.
	mockery --all --with-expecter --inpackage --testonly

generate: ## Generate automated code.
	make generate/mockery
	make fix
.PHONY: generate generate/*

test: ## Test Go Sync and included adapters.
	go test ./... $(ADAPTERS)
.PHONY: test

report: ## Test and produce a JUnit report.
	go test -v 2>&1 -count=1 ./... $(ADAPTERS) | go-junit-report -set-exit-code > report.xml
.PHONY: report

ci/tag-adapters: ## Tag all adapters with $RELEASE_VERSION environment variable. For use in CI.
	for adapter in $(shell ls -d adapters/*); do git tag $${adapter}/$$RELEASE_VERSION; done

tidy: ## Run go mod tidy in all adapters.
	go mod tidy
	for adapter in $(shell ls -d adapters/*); do sh -c "cd $${adapter} && go mod tidy"; done


.DEFAULT_GOAL := help
help: Makefile ## Display list of available commands.
	@grep -E '(^[a-zA-Z_-]+:.*?##.*$$)|(^##)' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}' | sed -e 's/\[32m##/[33m/'
