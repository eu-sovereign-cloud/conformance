GO := go
TOOLS_GOMOD := -modfile=./tools/go.mod
GO_TOOL := $(GO) run $(TOOLS_GOMOD)

WIREMOCK_PATH := $(shell pwd)/wiremock

REPORTS_PATH := $(shell pwd)/reports
RESULTS_DIR := results
RESULTS_PATH := $(REPORTS_PATH)/$(RESULTS_DIR)

DIST_DIR := ./dist
DIST_BIN := $(DIST_DIR)/secatest

.PHONY: $(DIST_BIN)
$(DIST_BIN):
	@echo "Building test code..."
	$(GO) test -c -o $(DIST_BIN) ./...

.PHONY: mock
mock:
	@echo "Running mock..."
	docker compose -f $(WIREMOCK_PATH)/docker-compose.yml -p seca-conformance up

.PHONY: run
run:
	@echo "Running test tool..."
	rm -rf $(RESULTS_PATH)
	mkdir -p $(RESULTS_PATH)
	$(DIST_BIN) --region_url=http://localhost:8080/providers/seca.region \
		--authorization_url=http://localhost:8080/providers/seca.authorization \
		--region_name=eu-central-1 \
		--auth_token=test-token \
		--results_path=$(RESULTS_PATH)

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

.PHONY: clean
clean:
	@echo "Cleaning up binaries and reports..."
	rm -rf $(DIST_DIR)
	rm -rf $(REPORTS_PATH)