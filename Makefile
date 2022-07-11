export GOPROXY=$("https://proxy.golang.org")

SHELL = /bin/bash
GOLANGCI_LINT = v1.46.2

.PHONY: clean-test-cache
clean-test-cache:
	@go clean -testcache ./...

.PHONY: test
test: clean-test-cache
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

.PHONY: build-port-service-linux
build-port-service-linux:
	@echo "Building Port Service binary for Linux"
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o cmd/port-service cmd/main.go

.PHONY: docker-build
docker-build:
	docker build -t albertogviana/port-service:latest -f Dockerfile --no-cache .