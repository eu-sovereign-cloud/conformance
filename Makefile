GO := go
TOOLS_GOMOD := -modfile=./tools/go.mod
GO_TOOL := $(GO) run $(TOOLS_GOMOD)

WIREMOCK_PATH := $(shell pwd)/wiremock
WIREMOCK_MAPPINGS_PATH := $(WIREMOCK_PATH)/config/mappings

REPORTS_PATH := $(shell pwd)/reports
RESULTS_DIR := results
RESULTS_PATH := $(REPORTS_PATH)/$(RESULTS_DIR)

DIST_DIR := ./dist
DIST_BIN := $(DIST_DIR)/secatest

.PHONY: default
default: $(DIST_BIN)

.PHONY: $(DIST_BIN)
$(DIST_BIN):
	@echo "Building conformance tool..."
	$(GO) test -c -o $(DIST_BIN) ./secatest

.PHONY: install
install:
	@echo "Installing conformance tool..."
	sudo cp $(DIST_BIN) /usr/local/bin

.PHONY: mock-run
mock-run:
	@echo "Running mock..."
	docker compose -f $(WIREMOCK_PATH)/docker-compose.yml -p seca-conformance up

.PHONY: mock-start
mock-start:
	@echo "Starting mock..."
	docker compose -f $(WIREMOCK_PATH)/docker-compose.yml -p seca-conformance up -d

.PHONY: mock-stop
mock-stop:
	@echo "Stopping mock..."
	docker compose -f $(WIREMOCK_PATH)/docker-compose.yml -p seca-conformance down

.PHONY: run
run:
	@echo "Running conformance tests..."
	rm -rf $(RESULTS_PATH)
	mkdir -p $(RESULTS_PATH)
	$(DIST_BIN) run \
	  --provider.region.v1=http://localhost:8080/providers/seca.region \
	  --provider.authorization.v1=http://localhost:8080/providers/seca.authorization \
	  --client.auth.token=test-token \
	  --client.region=region-1 \
	  --client.tenant=tenant-1 \
	  --scenarios.users=user1@secapi.com,user2@secapi.com \
	  --scenarios.cidr=10.1.0.0/16 \
	  --scenarios.public.ips=52.93.126.1/26 \
	  --report.results.path=$(RESULTS_PATH) \
	  --mock.enabled=true \
	  --mock.server.url=http://localhost:8080

.PHONY: report
report:
	@echo "Viewing report..."
	$(DIST_BIN) report $(RESULTS_PATH)

.PHONY: list
list:
	@echo "Listing scenarios..."
	$(DIST_BIN) list

.PHONY: test
test:
	@echo "Running tests..."
	$(GO) test -count=1 -v ./secatest -args run \
	  --provider.region.v1=http://localhost:8080/providers/seca.region \
	  --provider.authorization.v1=http://localhost:8080/providers/seca.authorization \
	  --client.auth.token=test-token \
	  --client.region=region-1 \
	  --client.tenant=tenant-1 \
	  --scenarios.users=user1@secapi.com,user2@secapi.com \
	  --scenarios.cidr=10.1.0.0/16 \
	  --scenarios.public.ips=52.93.126.1/26 \
	  --report.results.path=$(RESULTS_PATH) \
	  --mock.enabled=true \
	  --mock.server.url=http://localhost:8080

.PHONY: fmt
fmt:
	@echo "Formating code..."
	$(GO_TOOL) mvdan.cc/gofumpt -w .
	@echo "Formatting mock mappings..."
	find $(WIREMOCK_MAPPINGS_PATH) -name "*.json" -type f | while read -r file; do \
      jq '.' "$$file" > "$$file.tmp" && mv "$$file.tmp" "$$file"; \
	done

.PHONY: golint
golint:
	@echo "Linting code..."
	$(GO_TOOL) github.com/golangci/golangci-lint/v2/cmd/golangci-lint run --timeout 5m

.PHONY: vet
vet:
	@echo "Running vet..."
	$(GO) vet ./...

.PHONY: sec
sec:
	@echo "Running gosec..."
	$(GO_TOOL) github.com/securego/gosec/v2/cmd/gosec -exclude=G101,G404 ./...

.PHONY: lint
lint: fmt golint vet sec

.PHONY: clean
clean:
	@echo "Cleaning up binaries and reports..."
	rm -rf $(DIST_DIR)
	rm -rf $(REPORTS_PATH)