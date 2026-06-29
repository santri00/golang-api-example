BINARY_NAME=api
GO=go
GOLANGCI_LINT=golangci-lint

.PHONY: help fmt vet lint test check up

help: ## Show available targets
	@echo Available targets:
	@echo   fmt       - Run go fmt on all packages
	@echo   vet       - Run go vet on all packages
	@echo   lint      - Run golangci-lint on all packages
	@echo   test      - Run all tests
	@echo   test-race - Run all tests with race detector //requires CGO and a C compiler//
	@echo   check     - Run all quality gates: fmt, vet, lint, test
	@echo   up        - Run quality gates then start the stack via docker compose
	@echo   down      - Stop the service and reomve the containers

fmt: ## Run go fmt on all packages
	$(GO) fmt ./...

vet: fmt ## Run go vet on all packages
	$(GO) vet ./...

lint: vet ## Run golangci-lint on all packages
	$(GOLANGCI_LINT) run ./...

test: lint ## Run all tests
	$(GO) test -count=1 ./...

test-race: lint ## Run all tests with race detector (requires CGO_ENABLED=1 and a C compiler)
	CGO_ENABLED=1 $(GO) test -race -count=1 ./...

check: test ## Run all quality gates (fmt, vet, lint, test)

up: check ## Run quality gates then start the stack via docker compose
	docker compose up -d

down: ## Stop service and remove service containers
	docker compose down
