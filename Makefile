GO := go
TOOLS_GOMOD := -modfile=./tools/go.mod
GO_TOOL := $(GO) run $(TOOLS_GOMOD)
REPORTS_PATH := $(shell pwd)/reports
RESULTS_DIR := results
RESULTS_PATH := $(REPORTS_PATH)/$(RESULTS_DIR)

.PHONY: build
build:
	@echo "Building code..."
	$(GO) build ./...

.PHONY: run
run:
	@echo "Running tool..."
	rm -rf $(RESULTS_PATH)
	ALLURE_OUTPUT_PATH=$(REPORTS_PATH) ALLURE_OUTPUT_FOLDER=$(RESULTS_DIR) $(GO) test -count=1 ./secapi/...
	allure serve $(RESULTS_PATH)

.PHONY: fmt
fmt:
	@echo "Formating code..."
	$(GO_TOOL) mvdan.cc/gofumpt -w .

.PHONY: lint
lint:
	@echo "Linting code..."
	$(GO_TOOL) github.com/golangci/golangci-lint/cmd/golangci-lint run --timeout 5m

.PHONY: vet
vet:
	@echo "Running vet..."
	$(GO) vet ./...

dev: fmt lint vet
