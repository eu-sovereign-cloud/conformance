# Coding Guidelines

## Formatting

- Formatter: `gofumpt` (stricter than `gofmt`; run via `make fmt`, which also pretty-prints the
  WireMock mapping JSON under `wiremock/config/mappings` via `jq`)
- Linter: `golangci-lint v2` (run via `make lint`, config in `.golangci.yml`)
- A separate, opt-in duplicate-code check exists (`make dupl`, config in `.golangci-dupl.yml`,
  `dupl` threshold `100`) — not part of `make lint`, run it manually before adding a large new
  suite file that closely mirrors an existing one
- CI (`.github/workflows/linting.yml`) runs `make lint` on every push/PR (skipped for docs-only
  changes) — a failing lint blocks the PR

## Imports

Group imports into two blocks, separated by a blank line: everything except the allure-go
framework packages, then allure-go last. Within the first block, stdlib comes first, then this
module's own packages and the `go-sdk` packages together:

```go
import (
	"log/slog"
	"math/rand"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockUsage "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/usage"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)
```

- `github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema` and `.../secapi` are imported directly,
  **unaliased** — do not alias them as `sdk` or similar.
- `github.com/eu-sovereign-cloud/go-sdk/pkg/constants` is always aliased `sdkconsts` (it would
  otherwise collide with this module's own `internal/constants`).
- Per-domain mock scenario packages (`internal/mock/scenarios/<domain>`, package name `mockxxx`)
  are aliased to a capitalized `mockXxx` on import, e.g. `mockUsage`, `mockCompute`.
- `goimports` is enforced by the linter — run `make fmt` rather than hand-ordering imports.

## Package Structure

- Conformance test suites live under `internal/conformance/suites/<domain>/`
  (`authorization/`, `compute/`, `network/`, `region/`, `storage/`, `usage/`, `workspace/`), one
  file per suite, `package <domain>`.
- Matching mock configuration lives under `internal/mock/scenarios/<domain>/`, `package mock<domain>`.
- Shared/reusable fixture code (`pkg/builders`, `pkg/generators`, `pkg/wrappers`) is public and
  domain-agnostic — never put suite-specific logic there.
- Suite/scenario/mock config constants live in `internal/constants/` (`suites_v1.go`,
  `conditions.go`, `operations.go`, `test.go`), never scattered as string/int literals in suite
  files.
- Never inline raw SDK calls in a suite — add a reusable helper in
  `internal/conformance/steps/<domain>_v1.go` instead (see [HOWTO.md](HOWTO.md)).

## Naming

- Suite struct: `XxxV1TestSuite`; its config: `XxxV1Config`; its constructor:
  `CreateXxxV1TestSuite(baseSuite, config)`.
- Suite name constant: `XxxV1SuiteName` in `internal/constants/suites_v1.go`, value following
  `"<Domain>.V1.<Name>"` (e.g. `"Compute.V1.ProviderLifeCycle"`).
- Shared params struct (when a suite needs to carry fixtures from `BeforeAll` to `TestScenario`):
  `XxxV1Params` in `internal/conformance/params/params_v1.go`.
- Mock config function: `Configure<Name>V1(scenario *mockscenarios.Scenario, params XxxV1Params) error`.
- Step methods: `<Verb><Resource>V1Step`, e.g. `CreateOrUpdateInstanceV1Step`, `GetInstanceV1Step`,
  `DeleteInstanceV1Step`, `StartInstanceV1Step`, `ListInstanceV1Step`,
  `WatchInstanceUntilDeletedV1Step`.
- Builder chains read `Name(...).Provider(...).ApiVersion(...).Tenant(...)[.Workspace(...)][.Region(...)][.Network(...)].Labels(...).Spec(...).Build()`
  — always end with `.Build()`, which returns `(*schema.X, error)`.
- The linter has `revive`'s `var-naming` rule **disabled**: the codebase is inconsistent between
  `Ip`/`IP` (`GeneratePublicIpURL` vs. `schema.RegionalResourceMetadataKindResourceKindPublicIP`)
  because it mirrors whatever casing the generated SDK uses — match the casing of the SDK type or
  existing sibling functions you're extending, don't "fix" it to be more idiomatic Go.
- Avoid abbreviations unless they are domain terms already used throughout the SDK (`sku`, `nic`,
  `cidr`).

## Error Handling

- Every `Builder.Build()` call is followed immediately by an error check that fails the test:

  ```go
  workspace, err := builders.NewWorkspaceBuilder().
      Name(workspaceName).
      Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
      Tenant(suite.Tenant).Region(suite.Region).
      Build()
  if err != nil {
      t.Fatalf("Failed to build Workspace: %v", err)
  }
  ```

- Prefer `t.Fatalf("Failed to build X: %v", err)` in `BeforeAll` (the pattern used by every current
  usage/compute/network suite). An older `slog.Error("Failed to build X", "error", err); t.FailNow()`
  pattern exists in a few files — don't introduce more of it, `t.Fatalf` is the one to copy.
- Mock setup errors follow the same rule: `if err := suites.SetupMockIfEnabled(...); err != nil { t.Fatalf("Failed to setup mock: %v", err) }`.
- Inside a mock `Configure<Name>V1` function, every `configurator.Configure*Stub(...)` call is
  wrapped and its error returned immediately — never ignored, never collected and checked later:

  ```go
  if err := configurator.ConfigureCreateWorkspaceStub(workspace, workspaceUrl, scenario.MockParams); err != nil {
      return err
  }
  ```

- `errcheck`, `ineffassign`, and `unused` are enforced by the linter — never assign an error (or
  anything else) to `_` to silence them; fix the underlying issue instead.

## Scenario Structure

Every suite implements the same three allure-go lifecycle hooks, in this order and for this
purpose — see [HOWTO.md](HOWTO.md) for a full worked example:

1. **`BeforeAll(t)`** — build every fixture with `pkg/builders`/`pkg/generators` (no API calls),
   store them in a `*params.XxxV1Params`, then call `suites.SetupMockIfEnabled(...)`.
2. **`TestScenario(t)`** — call `suite.StartScenario(t, <providers...>)`, then
   `suite.ConfigureResources(t, <kinds under test...>)` (plus `suite.ConfigureDepends(t, ...)` for
   error/constraint suites), get a `steps.NewStepsConfigurator(...)`, then drive create → get →
   (update/action) → get → ... → delete, in dependency order for creation and **reverse**
   dependency order for teardown. Always finish with `suite.FinishScenario()`.
3. **`AfterAll(t)`** — always just `suite.ResetAllScenarios()`.

Step names passed as the first argument to every `Step` call are short, human-readable, Gherkin-style
phrases ("Create a workspace", "Get the started instance") — they appear verbatim in the Allure
report, so describe intent, not mechanics.

## Mock Stub Ordering

WireMock scenario state is a **single sequential chain per suite**: each registered stub requires
the mock's "current state" to match, in addition to URL/verb. This means a `Configure<Name>V1`
function must register stubs in **exactly** the order `TestScenario` will call them at runtime —
across every resource involved, not just within one resource's own create/get/delete sequence.

If a scenario provisions two independent resource stacks (e.g. two workspaces), do **not** register
one stack's full create+delete sequence before starting the next stack's — that only matches
runtime order if the test itself fully tears down stack A before touching stack B. If the test
instead does `create A, create B, list, delete A, delete B` (the common pattern for
isolation/parallel-stack scenarios), the mock config must be split into phase-specific helper
functions (`configureXCreateV1`, `configureXListV1`, `configureXDeleteV1`) called in that same
interleaved order — one combined per-stack function will desync the state machine and produce
either an outright match failure or, worse, a *wrong but valid-looking* stubbed response.

## Status Conditions & State Assertions

- `schema.ResourceState` (`Pending`/`Creating`/`Active`/`Updating`/`Deleting`) is the resource's
  current state; compare it against `constants.CreatedResourceExpectedStates` /
  `UpdatedResourceExpectedStates` (a *set* of allowed values for a create/update response).
- `schema.StatusCondition` entries **accumulate** — each lifecycle transition appends one more to
  the history, it does not replace the list. The `constants.GetConditionAfterX` variables
  (`GetConditionAfterCreating`, `GetConditionAfterUpdating`, `GetConditionAfterStarting`,
  `GetConditionAfterStartingWithoutUpdate`, `GetConditionAfterRestarting`,
  `GetConditionAfterStopping`, `GetConditionAfterDeleting`) are literal, fixed-length sequences
  calibrated to *specific* canonical call orders (e.g. `GetConditionAfterStarting` assumes
  Create → Update → Start; a scenario that goes straight from Create to Start must use
  `GetConditionAfterStartingWithoutUpdate` instead). Picking the wrong constant for your scenario's
  actual call sequence is the most common cause of a `"Status conditions length should match
  expected"` assertion failure — match the constant to your scenario's real sequence, don't assume
  the name closest to your step's description is correct without checking its length against the
  mock's accumulated history.

## Lint Rules to Design Around

- `revive`'s `argument-limit` is set to **5** — a function/method with 6+ explicit parameters (the
  receiver doesn't count) fails lint. When a helper genuinely needs more inputs, bundle related
  ones into an existing or new struct (e.g. pass `stack params.WorkspaceStackV1` instead of five
  separate resource pointers) rather than requesting an exception.
- `goconst` flags repeated string literals — pull anything reused more than once into
  `internal/constants` rather than repeating the literal.
- `gosec` is enabled with `G101` (hardcoded credentials) and `G404` (weak random source) excluded —
  `math/rand` is deliberately fine to use for generating test data (names, zone/SKU selection,
  IPs); don't switch to `crypto/rand` for that.
- `unconvert` flags unnecessary type conversions — when casting an SDK enum, only cast where the
  underlying type actually differs (`schema.ImageSpecCpuArchitecture(...)`, `string(...)`), not
  defensively.

## Comments

- Default to no comments.
- One-line comments only, and only when the WHY is non-obvious (a mock-ordering constraint, a
  schema quirk, a workaround) — never restate what the following code already says.
- Don't add `// TODO` comments.

## Dependencies

- Never add a new dependency to `go.mod` without confirming it's actually needed — check
  `pkg/builders`, `pkg/generators`, `pkg/wrappers`, and `internal/conformance/steps` first, since
  most cross-cutting behavior (name generation, response assertion, retry/watch polling) already
  has a helper.
- Developer tooling (`golangci-lint`, `gofumpt`) is pinned in the isolated `tools/go.mod`, invoked
  via `$(GO) run -modfile=./tools/go.mod -mod=mod ...` from the Makefile — don't add tooling
  dependencies to the main `go.mod`.
- Never use `replace` directives for published modules.
