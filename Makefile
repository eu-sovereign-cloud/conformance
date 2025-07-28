GO := go
TOOLS_GOMOD := -modfile=./tools/go.mod
GO_TOOL := $(GO) run $(TOOLS_GOMOD)
WIREMOCK_PATH := $(shell pwd)/wiremock
REPORTS_PATH := $(shell pwd)/reports
RESULTS_DIR := results
RESULTS_PATH := $(REPORTS_PATH)/$(RESULTS_DIR)

.PHONY: build
build:
	@echo "Building code..."
	$(GO) build ./...

.PHONY: mock
mock:
	@echo "Running mock..."
	docker compose -f $(WIREMOCK_PATH)/docker-compose.yml -p seca-conformance up

.PHONY: run
run:
	@echo "Running tool..."
	rm -rf $(RESULTS_PATH)
	mkdir -p $(RESULTS_PATH)
	cp reports/categories.json $(RESULTS_PATH)/ 2>/dev/null || true
	ALLURE_OUTPUT_PATH=$(REPORTS_PATH) ALLURE_OUTPUT_FOLDER=$(RESULTS_DIR) $(GO) test -count=1 ./secapi/...

.PHONY: report
report:
	@echo "Running report..."
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

.PHONY: sec
sec:
	@echo "Running gosec..."
	$(GO_TOOL) github.com/securego/gosec/v2/cmd/gosec ./...

.PHONY: dev
dev: fmt lint vet sec
