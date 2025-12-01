# Conformance

## Overview

SECA Conformance is a comprehensive testing tool for validating Cloud Service Provider (CSP) implementations against the [SECA API specification](https://spec.secapi.cloud). The tool validates API endpoints, resource lifecycle management, and compliance with SECA standards across for a specific cloud provider.

The project is built as a go test binary (`secatest`) that compiles test scenarios into an executable.
It uses It uses Cobra framework to provied the command line base strucute, uses the [SECA GO SDK](https://github.com/eu-sovereign-cloud/go-sdk) to comunicate with the SECA API implementation and uses Allure Reportin tool, trought the [AllureGo framework](https://github.com/ozontech/allure-go) for rich test reporting.

## Build and Development Commands

### Building the Binary
```bash
make                    # Build secatest binary to ./dist/secatest
make clean              # Remove binaries and reports
```

### Installation
```bash
make install            # Copy binary to /usr/local/bin (requires sudo)
```

### Code Quality
```bash
make fmt                # Format Go code with gofumpt and WireMock JSON mappings
make lint               # Run golangci-lint with .golangci.yml config
```

### Running Tests

#### Using the Built Binary
```bash
# Basic run (requires all configuration flags)
./dist/secatest run \
  --provider.region.v1=<URL> \
  --provider.authorization.v1=<URL> \
  --client.auth.token=<TOKEN> \
  --client.region=<REGION> \
  --client.tenant=<TENANT> \
  --scenarios.users=<USERS> \
  --scenarios.cidr=<CIDR> \
  --scenarios.public.ips=<IPS> \
  --retry.base.delay=<SECONDS> \
  --retry.base.interval=<SECONDS> \
  --retry.max.attempts=<COUNT>

# Filter specific scenarios
./dist/secatest run --scenarios.filter="Compute.V1.LifeCycle" [other flags...]
```

#### Using Make (with mock server)
```bash
make mock-start         # Start WireMock server in background
make run                # Run tests against mock server
make mock-stop          # Stop WireMock server

# Alternative: run mock in foreground
make mock-run           # Runs docker compose up (no -d flag)
```

#### Using Go Test Directly
```bash
make test               # Run go test with all required flags
```

### Test Management
```bash
./dist/secatest list    # List available test scenarios
./dist/secatest report <PATH>  # View Allure report (opens browser)
make report             # View report from default path (./reports/results)
```

## Architecture

### Test Entry Point and Command Structure

The project uses Go's `TestMain` in `secatest/main_test.go` to create a CLI-driven test suite. The binary accepts three subcommands:
- `run` - Execute conformance tests
- `list` - Show available test scenarios
- `report` - Generate and view Allure reports

Configuration is handled via CLI flags (e.g., `--client.auth.token`) which map to environment variables (e.g., `CLIENT_AUTH_TOKEN`).

### Client Architecture (Global vs Regional)

The SDK provides two client types in `secatest/clients.go`:

1. **GlobalClient** - Accesses global resources (regions, authorization)
   - Endpoints: Region V1, Authorization V1
   - Used for: Regions, Roles, RoleAssignments

2. **RegionalClient** - Accesses region-scoped resources
   - Created from GlobalClient with a specific region
   - Endpoints: Workspace V1, Storage V1, Compute V1, Network V1
   - Used for: Workspaces, BlockStorage, Images, Instances, Networks, etc.

### Test Suite Organization

Test suites follow a hierarchical structure defined in `secatest/suites.go`:

- **testSuite** - Base suite with common functionality (mock client, retry config, tenant)
  - **globalTestSuite** - For global resources (uses GlobalClient)
  - **regionalTestSuite** - For regional resources (uses RegionalClient)
  - **mixedTestSuite** - For tests needing both clients

Each API version has its own suite file and test file:
- `*_v1_suite.go` - Suite definition and initialization
- `*_v1_test.go` - Actual test scenarios
- `*_v1_steps.go` - Reusable test steps and helpers

### Step-Based Testing Pattern

Tests use a step-based approach defined in `secatest/steps.go` with generics:

1. **Generic Step Functions** - Parameterized by resource type `[R, M, E]`:
   - `R` - Resource type (e.g., `*schema.Role`)
   - `M` - Metadata type (e.g., `*schema.GlobalTenantResourceMetadata`)
   - `E` - Spec type (e.g., `schema.RoleSpec`)

2. **Step Categories**:
   - **CreateOrUpdate Steps** - Create/update resources and verify response
   - **Get Steps** - Retrieve resources with retry/wait logic
   - **Verify Steps** - Assert metadata, spec, labels, and status

3. **Resource Scopes**:
   - Tenant-scoped: `createOrUpdateTenantResourceStep`
   - Workspace-scoped: `createOrUpdateWorkspaceResourceStep`
   - Network-scoped: `createOrUpdateNetworkResourceStep`

The pattern allows type-safe, reusable step definitions across different resource types.

### Resource Lifecycle and Retry Logic

Resources follow an async creation pattern:
1. Create/update returns immediately
2. Resource enters transitional state (e.g., `CREATING`)
3. Tests poll using `ResourceObserverConfig` with:
   - `baseDelay` - Initial wait before first check
   - `baseInterval` - Time between retry attempts
   - `maxAttempts` - Maximum number of retries

This is implemented in `getResourceWithObserver` in `secatest/steps.go`.

### Mock Testing with WireMock

The `internal/mock` package provides WireMock-based mocking:
- Mock server configured via Docker Compose (`wiremock/docker-compose.yml`)
- Mappings in `wiremock/config/mappings/`
- Scenario configuration in `internal/mock/*_scenarios*.go`
- Mock mode enabled with `--mock.enabled=true` flag

Mock scenarios can be reset between tests using `suite.resetAllScenarios()`.

### Code Generation Helpers

`secalib/generators.go` provides helper functions for generating:
- Random resource names (e.g., `GenerateRoleName()`)
- Resource paths (e.g., `GenerateRoleResource(tenant, role)`)
- References (e.g., `GenerateSkuRef(name)`)
- Network addresses (subnet CIDR calculations, IP addressing)

Use these instead of hardcoding names to avoid conflicts between test runs.

## Key Concepts

### Test Scenarios

Available scenarios (see `secatest list`):
- `Authorization.V1.LifeCycle` - Roles and role assignments
- `Region.V1.LifeCycle` - Regional resources
- `Workspace.V1.LifeCycle` - Workspace management
- `Storage.V1.LifeCycle` - Block storage and images
- `Compute.V1.LifeCycle` - VM instances
- `Network.V1.LifeCycle` - Networks, subnets, gateways, NICs, etc.
- `Foundation.V1.Usage` - Usage metrics

### Resource Hierarchies

Understanding the resource hierarchy is critical:

```
Tenant
├── Role (global)
├── RoleAssignment (global)
└── Region
    └── Workspace
        ├── BlockStorage
        │   └── Image
        ├── Instance
        └── Network
            ├── InternetGateway
            ├── RouteTable
            ├── Subnet
            ├── PublicIp
            ├── Nic
            └── SecurityGroup
```

### Provider Structure

The SECA API is organized into providers:
- **seca.region** (Region V1) - Global/regional resource management
- **seca.authorization** (Authorization V1) - RBAC
- **seca.workspace** (Workspace V1) - Workspace isolation
- **seca.storage** (Storage V1) - Block storage and images
- **seca.compute** (Compute V1) - VM instances and SKUs
- **seca.network** (Network V1) - Networking resources

Each provider has its own base URL endpoint.

## Development Notes

### Adding New Test Scenarios

1. Add scenario name constant to the appropriate `*_v1_suite.go`
2. Create test function in `*_v1_test.go` (must start with `Test`)
3. Use step functions from `steps.go` for consistency
4. Add to list command in `main_test.go`

### Test Isolation

Tests should be isolated and not depend on external state:
- Use random resource names (via generators)
- Clean up resources after tests (delete steps)
- Reset mock scenarios between tests

### Adding Mock Scenarios

1. Create scenario configuration in `internal/mock/*_scenarios_v1.go`
2. Add WireMock mappings to `wiremock/config/mappings/`
3. Format JSON with `make fmt`
4. Call scenario setup in test suite initialization
