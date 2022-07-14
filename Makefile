export GOPROXY=$("https://proxy.golang.org")

SHELL = /bin/bash
GOLANGCI_LINT = v1.46.2

.PHONY: clean-test-cache
clean-test-cache:
	@go clean -testcache ./...

.PHONY: test
test: mysql-truncate-tables clean-test-cache
	@go test -v -timeout 30s ./... -race -cover -coverpkg=./... -coverprofile=coverage.out

.PHONY: unit-test
unit-test: clean-test-cache
	@go test -v -timeout 30s ./... -cover -coverprofile=coverage_unit.out -race -run Unit

.PHONY: integration-test
integration-test: mysql-truncate-tables clean-test-cache
	@go test -v -timeout 30s ./... -cover -coverprofile=coverage_integration.out -race -run Integration

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
check: test lint # Run tests and linters

.PHONY: mysql-truncate-tables
mysql-truncate-tables:
	@echo "Truncate MySQL Tables before run Integration Tests"
	@./scripts/mysql-ports-test-truncate.sh

.PHONY: build-port-service-linux
build-port-service-linux:
	@echo "Building Port Service binary for Linux"
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o cmd/port-service cmd/main.go

.PHONY: docker-build
docker-build:
	docker build -t albertogviana/port-service:latest -f Dockerfile --no-cache .