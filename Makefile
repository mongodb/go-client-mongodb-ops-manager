# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

SOURCE_FILES?=./...
PKG_NAME=opsmngr
GOLANGCI_VERSION=v1.23.6
COVERAGE=coverage.out

export PATH := ./bin:$(PATH)
export GO111MODULE := on

.PHONY: setup
setup:  ## Install dev tools
	@echo "==> Installing dependencies..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s $(GOLANGCI_VERSION)

.PHONY: link-git-hooks
link-git-hooks: ## Install git hooks
	@echo "==> Installing all git hooks..."
	find .git/hooks -type l -exec rm {} \;
	find .githooks -type f -exec ln -sf ../../{} .git/hooks/ \;

.PHONY: fmt
fmt: ## Format code
	@echo "==> Formatting all files..."
	gofmt -w -s ${PKG_NAME}
	goimports -w ${PKG_NAME}

.PHONY: test
test: ## Run tests
	@echo "==> Running tests..."
	go test -race -cover -count=1 -coverprofile ${COVERAGE} ${SOURCE_FILES}

.PHONY: lint
lint: ## Run linter
	@echo "==> Linting all packages..."
	golangci-lint run $(SOURCE_FILES) -E goimports -E golint -E misspell -E unconvert -E maligned -E bodyclose -E gosec

.PHONY: check
check: test lint ## Run tests and linters

.PHONY: addlicense
addlicense:
	find . -name '*.go' | while read -r file; do addlicense -c "MongoDB Inc" "$$file"; done

.PHONY: list
list: ## List all make targets
	@${MAKE} -pRrn : -f $(MAKEFILE_LIST) 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | sort

.PHONY: help
.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
