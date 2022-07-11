export GOPROXY=$("https://proxy.golang.org")

SHELL = /bin/bash
GOLANGCI_LINT = v1.46.2

.PHONY: unit-test
unit-test: clean-test-cache
	@go test -v -timeout 30s ./... -cover -coverprofile=coverage_unit.out -race -run Unit

.PHONY: integration-test
integration-test: clean-test-cache
	@go test -v -timeout 30s ./... -cover -coverprofile=coverage_integration.out -race -run Integration

.PHONY: clean-test-cache
clean-test-cache:
	@go clean -testcache ./...

.PHONY: test-all
test-all: clean-test-cache
	@go test -v -timeout 30s ./... -race -cover -coverpkg=./... -coverprofile=coverage.out

.PHONY: coverage
coverage:
	@go tool cover -html=coverage.out

.PHONY: install-lint
install-lint:
	@test -f ./bin/golangci-lint || curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s ${GOLANGCI_LINT}

.PHONY: lint
lint: install-lint
	@echo "Running golangci-lint"
	@bin/golangci-lint run

.PHONY: check
check: test-all lint # Run tests and linters
