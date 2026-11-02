GOLANGCI_LINT_VERSION := v2.15.0
GOLANGCI_LINT_BIN := $(CURDIR)/.cache/bin/golangci-lint-$(GOLANGCI_LINT_VERSION)
GOLANGCI_LINT := $(GOLANGCI_LINT_BIN)/golangci-lint
GOVULNCHECK_VERSION := v1.9.0

.PHONY: all build test coverage lint fmt check-fmt vuln tidy clean check ci

all: build check

build:
	go build ./...

test:
	go test -v -race ./...

coverage:
	go test -v -covermode=count -coverpkg=./... -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

$(GOLANGCI_LINT):
	GOBIN="$(GOLANGCI_LINT_BIN)" go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run

fmt: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) fmt

check-fmt: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) fmt --diff

vuln:
	go run golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION) ./...

tidy:
	go mod tidy

clean:
	rm -f coverage.out
	@test -x $(GOLANGCI_LINT) && $(GOLANGCI_LINT) cache clean || true
	go clean -testcache

check: test lint check-fmt

ci: check coverage
