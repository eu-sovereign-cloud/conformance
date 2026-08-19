# Glossary

Domain terminology, SDK types, and this repo's own test-framework types. See
[README.md](../README.md) for installation/usage, [STRUCTURE.md](STRUCTURE.md) for a map of the
repository layout, and [HOWTO.md](HOWTO.md) for a worked example of most of the concepts below.

## Domain Terms

| Term | Meaning |
|---|---|
| **SECA** | EU Sovereign Cloud — the [API specification](https://spec.secapi.cloud) this tool validates a CSP's implementation against. |
| **CSP** | Cloud Service Provider — the system under test. |
| **Tenant** | The top-level isolation boundary. All resources belong to a tenant (`--client.tenant`, `suite.Tenant`). |
| **Region** | A geographic deployment (`--client.region`, `suite.Region`). Most resources are region-scoped; `Role`/`RoleAssignment` are tenant-scoped only. |
| **Zone** | A failure domain within a region (e.g. `"a"`, `"b"`), set on `InstanceSpec.Zone` / `SubnetSpec.Zone`. |
| **Workspace** | A logical grouping of resources within a tenant+region — most resources (BlockStorage, Instance, Network, ...) live inside one. |
| **SKU (Sku / SkuRef)** | A named performance/feature tier (`InstanceSku`, `StorageSku`, `NetworkSku`), resolved to a `Reference` via `generators.GenerateSkuRefObject(...)`. |
| **Ref / Reference** | `schema.Reference` — a typed pointer from one resource to another (e.g. `Instance.BootVolume.DeviceRef` → a `BlockStorage`). Built via `generators.GenerateXRefObject(...)`, dereferenced with `*` when assigned into a spec. |
| **BlockStorage** | A persistent volume, referenced by an `Instance` as its `BootVolume` or a `DataVolume`. |
| **Image** | A bootable disk image wrapping a `BlockStorageRef`. |
| **BootVolume / DataVolume** | `InstanceSpec` fields: `BootVolume` (required, single) is the disk the instance boots from; `DataVolumes` (optional, list) are extra attached disks. Both wrap a `VolumeReference{DeviceRef}`. |
| **NIC (Nic)** | Network Interface Card — a `Network`-domain resource attaching an `Instance`'s workspace to a `Subnet`, optionally with `PublicIpRefs`. |
| **AntiAffinityGroup** | Opaque string on `InstanceSpec` hinting that same-valued instances shouldn't share a host — the basis for HA placement. |
| **Labels / Annotations / Extensions** | Resource metadata maps: `Labels` (structured, filterable via `ListOptions.WithLabels(...)`), `Annotations` (free-form, e.g. descriptions), `Extensions` — all asserted for round-trip equality in response checks. |
| **Scope (RoleAssignmentScope)** | The blast radius of a `RoleAssignment`: `Tenants`, `Workspaces`, and/or `Regions` name lists, at least one required. |

## SDK Types (`go-sdk`)

### Client Types (`secapi` package)

| Type | Purpose |
|---|---|
| `secapi.GlobalClient` | Global/tenant-scoped domains: `.AuthorizationV1`, `.RegionV1`. |
| `secapi.RegionalClient` | Region/workspace-scoped domains: `.WorkspaceV1`, `.StorageV1`, `.ComputeV1`, `.NetworkV1`. |

### Reference & Path Types (`secapi` package)

| Type | Fields | Used for |
|---|---|---|
| `secapi.TenantReference` | `Tenant TenantID`, `Name string` | Tenant-scoped resources (Role, RoleAssignment, Workspace) |
| `secapi.WorkspaceReference` | `Tenant`, `Workspace WorkspaceID`, `Name string` | Workspace-scoped resources (BlockStorage, Instance, Nic, ...) |
| `secapi.NetworkReference` | `Tenant`, `Workspace`, `Network NetworkID`, `Name string` | Network-scoped resources (RouteTable, Subnet) |
| `secapi.TenantPath` / `WorkspacePath` / `NetworkPath` | scope IDs only, no resource `Name` | Passed to `List*` steps (e.g. `ListInstanceV1Step`) |
| `secapi.TenantID` / `WorkspaceID` / `NetworkID` | `string` alias | Type-safe scope identifiers |
| `secapi.Reference` | `{Resource string}` | Generic cross-resource reference payload |

### Query & Polling Types

| Type | Purpose |
|---|---|
| `secapi.ListOptions` / `NewListOptions()` | Pagination + label filtering for `List*` calls (`.WithLimit(n)`, `.WithLabels(...)`) |
| `secapi.Iterator[T]` | Paged result set returned by `List*` SDK calls |
| `secapi.ResourceObserverConfig` / `ResourceObserverUntilValueConfig[V]` | Poll config (`BaseDelay`, `BaseInterval`, `MaxAttempts`) backing `Get*UntilXStep` / `Watch*UntilDeletedV1Step` |
| `secapi.ErrRequestPreconditionFailed` | Sentinel error compared against for 409/412/422-style rejections (`requirePreConditionFailedError`) |

### Schema & Metadata Types (`go-sdk/pkg/spec/schema`)

| Type | Scope | Notes |
|---|---|---|
| `schema.GlobalResourceMetadata` | Global | e.g. Region |
| `schema.GlobalTenantResourceMetadata` | Tenant, no region | Role, RoleAssignment |
| `schema.RegionalResourceMetadata` | Tenant + Region | Workspace, Image |
| `schema.RegionalWorkspaceResourceMetadata` | + Workspace | BlockStorage, Instance, Network, InternetGateway, Nic, PublicIp, SecurityGroup |
| `schema.RegionalNetworkResourceMetadata` | + Network | RouteTable, Subnet |
| `schema.SkuResourceMetadata` | Tenant (SKUs) | InstanceSku, StorageSku, NetworkSku |

### State Types

| Type | Values | Notes |
|---|---|---|
| `schema.ResourceState` | `Pending`, `Creating`, `Active`, `Updating`, `Deleting` | Current lifecycle state; `constants.CreatedResourceExpectedStates` / `UpdatedResourceExpectedStates` are the allowed *sets* a create/update response may land in. |
| `schema.StatusCondition` | `{State, LastTransitionAt}` | One entry in `Status.Conditions` — unlike `ResourceState`, conditions **accumulate** one per lifecycle transition. `constants.GetConditionAfterX` vars are the literal expected sequences for specific canonical call orders (mismatching the call order to the wrong `GetConditionAfterX` constant is the most common assertion bug when writing a new scenario). |
| `schema.SecurityGroupRuleDirection` | `Ingress`, `Egress` | `SecurityGroupRuleSpec.Direction` |

## This Repo's Test-Framework Types

### Suite Types (`internal/conformance/suites`)

| Type | Purpose |
|---|---|
| `TestSuite` | Base type embedding allure-go's `suite.Suite`; holds `Tenant`, `AuthToken`, `MockEnabled`, retry config, `ScenarioName`. |
| `GlobalTestSuite` | `TestSuite` + `Client *secapi.GlobalClient` |
| `RegionalTestSuite` | `TestSuite` + `Region string` + `Client *secapi.RegionalClient` |
| `MixedTestSuite` | `TestSuite` + both a `GlobalClient` and a `RegionalClient` — used when a scenario spans global and regional domains (e.g. creating a `Role` *and* a `Workspace`) |

### Suite Lifecycle & Naming

| Term | Meaning |
|---|---|
| `SuiteName` | Typed string (`<Domain>.V1.<Name>`, e.g. `"Usage.V1.FoundationProviders"`) declared in `internal/constants/suites_v1.go`, appended to `AllSuiteNames`. Drives `secatest list` and `--scenarios.filter`. |
| `BeforeAll` / `TestScenario` / `AfterAll` | allure-go lifecycle hooks: build fixtures + mock setup / run the actual test / `suite.ResetAllScenarios()`. |
| `ConfigureResources` / `ConfigureDepends` | Tag which resource `Kind`s a scenario exercises vs. merely depends on — reporting/filtering aid only, not a runtime check. |
| `CanRun(regexp)` | Checks `ScenarioName` against `--scenarios.filter` before a suite is run from `cmd/conformance/<domain>_test.go`. |

### Fixture & Assertion Helpers

| Type/Package | Purpose |
|---|---|
| `pkg/builders` | Fluent, generics-based constructors per resource (`NewWorkspaceBuilder()`, `NewInstanceBuilder()`, ...) — validate required fields, never make API calls. |
| `pkg/generators` | Test-data helpers: random names, `Reference`/URL builders, derived values (subnet CIDRs, NIC addresses). |
| `pkg/wrappers` | `ResourceWrapper[R,M,E,S]` / `GlobalResourceWrapper[R,M,E]` — normalize the differently-shaped SDK response types behind a common `GetResource()`/`GetMetadata()`/`GetSpec()`/`GetStatus()` interface (e.g. `wrappers.NewInstanceWrapper(resp)`). |
| `internal/conformance/params` | Per-suite structs (e.g. `FoundationUsageV1Params`) carrying resources built in `BeforeAll` through to `TestScenario` and the mock configurator. |
| `internal/conformance/steps.StepsConfigurator` | Obtained via `steps.NewStepsConfigurator(...)`; exposes one `Step` method per operation per resource (`CreateOrUpdateInstanceV1Step`, `GetInstanceV1Step`, `StartInstanceV1Step`, `DeleteInstanceV1Step`, `ListInstanceV1Step`, `WatchInstanceUntilDeletedV1Step`, ...). |
| `steps.ResponseExpects[M,E]` / `ResponseExpectsWithCondition[M,E,S]` | Expected `Labels`/`Annotations`/`Extensions`/`Metadata`/`Spec`/`ResourceStates` (or `ResourceStatus`) passed into a create/get step. |
| `suites.VerifyXStep(...)` | Field-by-field assertion helpers in `internal/conformance/suites/asserts*.go` (`VerifyStatusConditionsStep`, `VerifyLabelsStep`, ...), used internally by the generic `Step` implementations. |

## Mocking Types (`internal/mock`)

| Type | Purpose |
|---|---|
| `mockscenarios.Scenario` | A WireMock-backed stand-in for one conformance suite run (`mockscenarios.NewScenario(suiteName, mockParams)`), configured by a `Configure<SuiteName>V1(scenario, params)` function. Lifecycle: `StartConfiguration()` → register stubs → `FinishConfiguration()` → `ResetAllScenarios()` after the test. |
| `stubs.Configurator` | Registers individual WireMock stub rules per resource type (`ConfigureCreateInstanceStub`, `ConfigureGetActiveInstanceStub`, `ConfigureInstanceOperationStub`, `ConfigureDeleteStub`, `ConfigureListInstanceStub`, ...), obtained via `scenario.StartConfiguration()`. |
| WireMock scenario state | WireMock's own stateful sequencing: stubs match by required "current state" as well as URL/verb, scoped one-to-one with the conformance `SuiteName`. **Mock stub registration order must exactly match the runtime call order of `TestScenario`**, across every resource involved — the most common source of mock-only test failures when adding a new scenario. |

## CLI / Infrastructure Terms

| Term | Meaning |
|---|---|
| `secatest` | The compiled CLI binary (`dist/secatest`, built via `go test -c -o dist/secatest ./cmd/conformance`), with subcommands `run`, `list`, `report`, `summary`. |
| `--provider.region.v1` / `--provider.authorization.v1` | Base URL of the CSP's actual implementation of one SECA API domain to test against — distinct from `sdkconsts.XProviderV1Name` (e.g. `"seca.compute"`), which is the domain identifier string embedded in resource metadata and `Role` permissions. |
| `--scenarios.filter` | Regexp matched against `SuiteName` to select which scenarios `run` executes. |
| `--mock.enabled` / `--mock.server.url` | Run against WireMock instead of a real provider. |
| WireMock | The local mock HTTP server (`wiremock/docker-compose.yml`) that suites can be validated against without a real CSP backend. |
| Allure / Allure Report | The test-reporting framework (`ozontech/allure-go` for authoring, Allure Report V2 for viewing) — every `Step` shows up as a named entry in the report opened by `secatest report`. |
