GO := go
TOOLS_GOMOD := -modfile=./tools/go.mod
GO_TOOL := $(GO) run $(TOOLS_GOMOD)

.PHONY: mock
mock: generate
	$(GO_TOOL) github.com/vektra/mockery/v2

.PHONY: build
build:
	$(GO) build ./...

.PHONY: test
test:
	$(GO) test -count=1 -cover -coverprofile=coverage.out -v ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	rm coverage.out

.PHONY: fmt
fmt:
	$(GO_TOOL) mvdan.cc/gofumpt -w .

lint: vet golint

.PHONY: vet
vet:
	$(GO) vet ./...

.PHONY: golint
golint:
	$(GO_TOOL) github.com/golangci/golangci-lint/cmd/golangci-lint run --timeout 5m -c .golangci.yml

.PHONY: clean
clean:
	rm -rf $(SCHEMAS_FINAL) $(SPEC_DIST) mock
