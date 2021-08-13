# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

.PHONY: test build
.DEFAULT_GOAL := help
VERSION ?= $$(git describe --tags --dirty)

test: ## Run the tests
	mkdir -p cover
	buffalo test -coverprofile=cover/c.out ./...
	go tool cover -func=cover/c.out

build: ## Build the project
	 buffalo build -v --environment production --ldflags "-X auri/config.prodBuild=yes" --ldflags "-X auri/config.version=$(VERSION)" --tags prodBuild

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
