TAGS	?= sqlite,integration
GO_BIN	?= go

LINT_TARGET	?= ./...
LINT_SKIP_DIRS	?= "internal"

LINTERS_B	= govet,revive,goimports,unused,typecheck,ineffassign,misspell,nakedret,unconvert,megacheck
LINTERS		= $(LINTERS_B),gocognit,gocyclo,misspell
MORE_LINTERS	= $(LINTERS),errorlint,goerr113,forcetypeassert,gosec

DEPRECATED_LINTERS	= deadcode,golint,maligned,ifshort,scopelint,interfacer,varcheck,nosnakecase,structcheck
NEVER_LINTERS		= godox,forbidigo,varnamelen,paralleltest,testpackage
DISABLED_LINTERS	= $(DEPRECATED_LINTERS),$(NEVER_LINTERS),exhaustivestruct,ireturn,gci,gochecknoglobals,wrapcheck
DISABLED_LINTERS_TEST	= $(DISABLED_LINTERS),nolintlint,wsl,goerr113,forcetypeassert,errcheck,goconst,scopelint

test:
	$(GO_BIN) test -tags ${TAGS} -cover -race ./...

lint:
	# go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run --disable-all --enable $(LINTERS) \
		--skip-dirs $(LINT_SKIP_DIRS) \
		$(LINT_TARGET)
	golangci-lint run --disable-all --enable $(MORE_LINTERS) \
		--skip-dirs $(LINT_SKIP_DIRS) \
		--tests=false $(LINT_TARGET)

lint-harder:
	golangci-lint run --enable-all --disable $(DISABLED_LINTERS) \
		--skip-dirs $(LINT_SKIP_DIRS) \
		--exclude-use-default=false \
		--tests=false $(LINT_TARGET)
	golangci-lint run --enable-all --disable $(DISABLED_LINTERS_TEST) \
		--skip-dirs $(LINT_SKIP_DIRS) \
		--exclude-use-default=false \
		$(LINT_TARGET)

