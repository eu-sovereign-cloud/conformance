# HOWTO: Add a New Use-Case Scenario

This guide walks through adding a new conformance test scenario end-to-end, using the
`Compute.V1.ProviderLifeCycle` suite
([internal/conformance/suites/compute/provider_lifecycle_v1.go](../internal/conformance/suites/compute/provider_lifecycle_v1.go))
as the worked example — it drives a `compute.v1.Instance` through its full lifecycle: create,
get, update, start, restart, stop, delete. See [STRUCTURE.md](STRUCTURE.md) for a map of the
repository layout referenced below.

A scenario is a single `go test` case that:

1. Builds the resources it needs (via `pkg/builders`).
2. Registers the mock HTTP responses those resources will trigger (via `internal/mock/scenarios`),
   used when running against WireMock instead of a real provider.
3. Drives the actual test steps against the SDK client (via `internal/conformance/steps`) and
   asserts on the responses.
4. Is wired into a suite (a group of related scenarios) and registered with a CLI test entrypoint.

## 1. Name the suite

Add a constant in [internal/constants/suites_v1.go](../internal/constants/suites_v1.go), following
the `<Domain>.V1.<Name>` naming convention already used there:

```go
ComputeProviderLifeCycleV1SuiteName SuiteName = "Compute.V1.ProviderLifeCycle"
```

Add it to the `AllSuiteNames` slice in the same file — this is what powers the `conformance list`
CLI command and the `--scenarios` regexp filter.

## 2. Define the suite's parameters

If the suite needs to share built resources between its setup phase and its test steps, add a
params struct in [internal/conformance/params/params_v1.go](../internal/conformance/params/params_v1.go):

```go
type ComputeProviderLifeCycleV1Params struct {
    Workspace       *schema.Workspace
    BlockStorage    *schema.BlockStorage
    InitialInstance *schema.Instance
    UpdatedInstance *schema.Instance
}
```

## 3. Build the resources (`BeforeAll`)

Create the suite type under `internal/conformance/suites/<domain>/`, embedding
`suites.RegionalTestSuite`. In `BeforeAll`, use the fluent builders in `pkg/builders` to construct
each resource, and `pkg/generators` for names, SKU refs, and other test data. An `Instance`
requires a `Workspace` and a `BlockStorage` (as its boot volume) to already exist:

```go
type ProviderLifeCycleV1TestSuite struct {
    suites.RegionalTestSuite
    config *ProviderLifeCycleV1Config
    params *params.ComputeProviderLifeCycleV1Params
}

func CreateProviderLifeCycleV1TestSuite(regionalTestSuite suites.RegionalTestSuite, config *ProviderLifeCycleV1Config) *ProviderLifeCycleV1TestSuite {
    suite := &ProviderLifeCycleV1TestSuite{RegionalTestSuite: regionalTestSuite, config: config}
    suite.ScenarioName = constants.ComputeProviderLifeCycleV1SuiteName.String()
    return suite
}

func (suite *ProviderLifeCycleV1TestSuite) BeforeAll(t provider.T) {
    t.AddParentSuite(suites.ComputeParentSuite)

    instanceSkuName := suite.config.InstanceSkus[rand.Intn(len(suite.config.InstanceSkus))]
    storageSkuName := suite.config.StorageSkus[rand.Intn(len(suite.config.StorageSkus))]
    initialZone := suite.config.AvailableZones[rand.Intn(len(suite.config.AvailableZones))]

    workspaceName := generators.GenerateWorkspaceName()
    blockStorageName := generators.GenerateBlockStorageName()
    instanceName := generators.GenerateInstanceName()

    storageSkuRef := generators.GenerateSkuRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, storageSkuName)
    instanceSkuRef := generators.GenerateSkuRefObject(sdkconsts.ComputeProviderV1Name, suite.Tenant, instanceSkuName)
    blockStorageRef := generators.GenerateBlockStorageRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, workspaceName, blockStorageName)

    workspace, err := builders.NewWorkspaceBuilder().
        Name(workspaceName).
        Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
        Tenant(suite.Tenant).Region(suite.Region).
        Labels(schema.Labels{constants.EnvLabel: constants.EnvDevelopmentLabel}).
        Build()
    if err != nil {
        t.Fatalf("Failed to build Workspace: %v", err)
    }

    blockStorage, err := builders.NewBlockStorageBuilder().
        Name(blockStorageName).
        Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
        Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
        Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
        Spec(&schema.BlockStorageSpec{SkuRef: *storageSkuRef, SizeGB: constants.BlockStorageInitialSize}).
        Build()
    if err != nil {
        t.Fatalf("Failed to build BlockStorage: %v", err)
    }

    initialInstance, err := builders.NewInstanceBuilder().
        Name(instanceName).
        Provider(sdkconsts.ComputeProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
        Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
        Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
        Spec(&schema.InstanceSpec{
            SkuRef:     *instanceSkuRef,
            Zone:       initialZone,
            BootVolume: schema.VolumeReference{DeviceRef: *blockStorageRef},
        }).
        Build()
    if err != nil {
        t.Fatalf("Failed to build Instance: %v", err)
    }

    // updatedInstance mirrors initialInstance but with a different Zone — used later to test update.

    p := &params.ComputeProviderLifeCycleV1Params{
        Workspace:       workspace,
        BlockStorage:    blockStorage,
        InitialInstance: initialInstance,
        UpdatedInstance: updatedInstance,
    }
    suite.params = p

    // Register the mock stubs for this suite (no-op when running against a real provider)
    if err := suites.SetupMockIfEnabled(suite.TestSuite, mockcompute.ConfigureProviderLifecycleScenarioV1, *p); err != nil {
        t.Fatalf("Failed to setup mock: %v", err)
    }
}
```

Building resources here (rather than inline in the test) never issues API calls — it only
constructs and validates the payloads. The actual create/get/update/delete calls happen as
explicit steps in `TestScenario`.

## 4. Configure the mock stubs

Under `internal/mock/scenarios/<domain>/`, add a `Configure<Name>` function that registers one
WireMock stub per HTTP call your scenario will make, **in the order the test steps will call
them**:

```go
func ConfigureProviderLifecycleScenarioV1(scenario *mockscenarios.Scenario, params params.ComputeProviderLifeCycleV1Params) error {
    configurator, err := scenario.StartConfiguration()
    if err != nil {
        return err
    }
    workspace, blockStorage, instance := params.Workspace, params.BlockStorage, params.InitialInstance

    workspaceUrl := generators.GenerateWorkspaceURL(sdkconsts.WorkspaceProviderV1Name, workspace.Metadata.Tenant, workspace.Metadata.Name)
    blockUrl := generators.GenerateBlockStorageURL(sdkconsts.StorageProviderV1Name, blockStorage.Metadata.Tenant, workspace.Metadata.Name, blockStorage.Metadata.Name)
    instanceUrl := generators.GenerateInstanceURL(sdkconsts.ComputeProviderV1Name, instance.Metadata.Tenant, workspace.Metadata.Name, instance.Metadata.Name)
    instanceStartUrl := generators.GenerateInstanceStartURL(sdkconsts.ComputeProviderV1Name, instance.Metadata.Tenant, workspace.Metadata.Name, instance.Metadata.Name)
    instanceStopUrl := generators.GenerateInstanceStopURL(sdkconsts.ComputeProviderV1Name, instance.Metadata.Tenant, workspace.Metadata.Name, instance.Metadata.Name)
    instanceRestartUrl := generators.GenerateInstanceRestartURL(sdkconsts.ComputeProviderV1Name, instance.Metadata.Tenant, workspace.Metadata.Name, instance.Metadata.Name)

    // Workspace: create -> creating -> active
    if err := configurator.ConfigureCreateWorkspaceStub(workspace, workspaceUrl, scenario.MockParams); err != nil {
        return err
    }
    if err := configurator.ConfigureGetCreatingWorkspaceStub(workspace, workspaceUrl, scenario.MockParams); err != nil {
        return err
    }
    if err := configurator.ConfigureGetActiveWorkspaceStub(workspace, workspaceUrl, scenario.MockParams); err != nil {
        return err
    }

    // BlockStorage: create -> creating -> active (same pattern as Workspace)
    if err := configurator.ConfigureCreateBlockStorageStub(blockStorage, blockUrl, scenario.MockParams); err != nil {
        return err
    }
    if err := configurator.ConfigureGetCreatingBlockStorageStub(blockStorage, blockUrl, scenario.MockParams); err != nil {
        return err
    }
    if err := configurator.ConfigureGetActiveBlockStorageStub(blockStorage, blockUrl, scenario.MockParams); err != nil {
        return err
    }

    // Instance: create -> creating -> active
    if err := configurator.ConfigureCreateInstanceStub(instance, instanceUrl, scenario.MockParams); err != nil {
        return err
    }
    if err := configurator.ConfigureGetCreatingInstanceStub(instance, instanceUrl, scenario.MockParams); err != nil {
        return err
    }
    if err := configurator.ConfigureGetActiveInstanceStub(instance, instanceUrl, scenario.MockParams); err != nil {
        return err
    }

    // Instance: update (new Zone) -> updating -> active
    instance = params.UpdatedInstance
    if err := configurator.ConfigureUpdateInstanceStub(instance, instanceUrl, scenario.MockParams); err != nil {
        return err
    }
    if err := configurator.ConfigureGetUpdatingInstanceStub(instance, instanceUrl, scenario.MockParams); err != nil {
        return err
    }
    if err := configurator.ConfigureGetActiveInstanceStub(instance, instanceUrl, scenario.MockParams); err != nil {
        return err
    }

    // Instance: start / restart / stop — each is a POST action, then a GET to observe the new state
    if err := configurator.ConfigureInstanceOperationStub(instance, instanceStartUrl, scenario.MockParams); err != nil {
        return err
    }
    if err := configurator.ConfigureGetActiveInstanceStub(instance, instanceUrl, scenario.MockParams); err != nil {
        return err
    }
    if err := configurator.ConfigureInstanceOperationStub(instance, instanceRestartUrl, scenario.MockParams); err != nil {
        return err
    }
    if err := configurator.ConfigureGetActiveInstanceStub(instance, instanceUrl, scenario.MockParams); err != nil {
        return err
    }
    if err := configurator.ConfigureInstanceOperationStub(instance, instanceStopUrl, scenario.MockParams); err != nil {
        return err
    }
    if err := configurator.ConfigureGetActiveInstanceStub(instance, instanceUrl, scenario.MockParams); err != nil {
        return err
    }

    // Teardown, in reverse dependency order: Instance -> BlockStorage -> Workspace.
    // Each is delete -> get-deleting -> get-not-found.
    if err := configurator.ConfigureDeleteStub(instanceUrl, scenario.MockParams); err != nil {
        return err
    }
    if err := configurator.ConfigureGetDeletingInstanceStub(instance, instanceUrl, scenario.MockParams); err != nil {
        return err
    }
    if err := configurator.ConfigureGetNotFoundStub(instanceUrl, scenario.MockParams); err != nil {
        return err
    }

    // ... same delete/get-deleting/get-not-found pattern for blockUrl, then workspaceUrl ...

    return scenario.FinishConfiguration(configurator)
}
```

The `internal/mock/stubs.Configurator` used above already has helpers for the common
create/get(creating|active|updating|deleting)/delete/not-found patterns per resource type, plus
`ConfigureInstanceOperationStub` for start/stop/restart-style actions — check
`internal/mock/stubs/` before writing a new one.

## 5. Write the test steps (`TestScenario`)

Use `steps.NewStepsConfigurator` to get access to the reusable Gherkin-style step helpers in
`internal/conformance/steps/<domain>_v1.go` (e.g. `CreateOrUpdateWorkspaceV1Step`,
`CreateOrUpdateInstanceV1Step`, `GetInstanceV1Step`, `StartInstanceV1Step`,
`RestartInstanceV1Step`, `StopInstanceV1Step`, `DeleteInstanceV1Step`,
`WatchInstanceUntilDeletedV1Step`). Each `Get*Step` returns the refreshed resource so later steps
can chain off its latest state:

```go
func (suite *ProviderLifeCycleV1TestSuite) TestScenario(t provider.T) {
    suite.StartScenario(t, sdkconsts.ComputeProviderV1Name)
    suite.ConfigureResources(t, string(schema.RegionalResourceMetadataKindResourceKindInstance))
    suite.ConfigureDepends(t,
        string(schema.RegionalResourceMetadataKindResourceKindWorkspace),
        string(schema.RegionalResourceMetadataKindResourceKindBlockStorage),
    )

    stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

    // Workspace
    workspace := suite.params.Workspace
    stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", t, suite.Client.WorkspaceV1, workspace,
        steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
            Metadata:       workspace.Metadata,
            Labels:         workspace.Labels,
            ResourceStates: suites.CreatedResourceExpectedStates,
        },
    )
    workspaceTRef := secapi.TenantReference{Tenant: secapi.TenantID(suite.Tenant), Name: workspace.Metadata.Name}
    stepsBuilder.GetWorkspaceV1Step("Get the created workspace", suite.Client.WorkspaceV1, workspaceTRef,
        steps.ResponseExpectsWithCondition[schema.RegionalResourceMetadata, schema.WorkspaceSpec, schema.WorkspaceStatus]{
            Metadata: workspace.Metadata,
            ResourceStatus: schema.WorkspaceStatus{
                State:      schema.ResourceStateActive,
                Conditions: suites.GetConditionAfterCreating,
            },
        },
    )

    // BlockStorage (same Create -> Get pattern as Workspace) ...
    block := suite.params.BlockStorage
    blockWRef := secapi.WorkspaceReference{Tenant: secapi.TenantID(suite.Tenant), Workspace: secapi.WorkspaceID(workspace.Metadata.Name), Name: block.Metadata.Name}
    // stepsBuilder.CreateOrUpdateBlockStorageV1Step(...) / stepsBuilder.GetBlockStorageV1Step(...)

    // Instance — create
    instance := suite.params.InitialInstance
    expectInstanceMeta := instance.Metadata
    expectInstanceSpec := &instance.Spec
    stepsBuilder.CreateOrUpdateInstanceV1Step("Create an instance", t, suite.Client.ComputeV1, instance,
        steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
            Metadata:       expectInstanceMeta,
            Spec:           expectInstanceSpec,
            ResourceStates: suites.CreatedResourceExpectedStates,
        },
    )
    instanceWRef := secapi.WorkspaceReference{Tenant: secapi.TenantID(suite.Tenant), Workspace: secapi.WorkspaceID(workspace.Metadata.Name), Name: instance.Metadata.Name}
    instance = stepsBuilder.GetInstanceV1Step("Get the created instance", suite.Client.ComputeV1, instanceWRef,
        steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus]{
            Metadata: expectInstanceMeta,
            Spec:     expectInstanceSpec,
            ResourceStatus: schema.InstanceStatus{
                State:      schema.ResourceStateActive,
                Conditions: suites.GetConditionAfterCreating,
            },
        },
    )

    // Instance — update Zone
    instance.Spec.Zone = suite.params.UpdatedInstance.Spec.Zone
    expectInstanceSpec.Zone = instance.Spec.Zone
    stepsBuilder.CreateOrUpdateInstanceV1Step("Update the instance", t, suite.Client.ComputeV1, instance,
        steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
            Metadata:       expectInstanceMeta,
            Spec:           expectInstanceSpec,
            ResourceStates: suites.UpdatedResourceExpectedStates,
        },
    )
    instance = stepsBuilder.GetInstanceV1Step("Get the updated instance", suite.Client.ComputeV1, instanceWRef,
        steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus]{
            Metadata: expectInstanceMeta,
            Spec:     expectInstanceSpec,
            ResourceStatus: schema.InstanceStatus{
                State:      schema.ResourceStateActive,
                Conditions: suites.GetConditionAfterUpdating,
            },
        },
    )

    // Instance — power actions: start, restart, stop, each followed by a Get to confirm state
    stepsBuilder.StartInstanceV1Step("Start the instance", suite.Client.ComputeV1, instance)
    instance = stepsBuilder.GetInstanceV1Step("Get the started instance", suite.Client.ComputeV1, instanceWRef,
        steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus]{
            Metadata:       expectInstanceMeta,
            Spec:           expectInstanceSpec,
            ResourceStatus: schema.InstanceStatus{State: schema.ResourceStateActive, Conditions: suites.GetConditionAfterStarting},
        },
    )
    stepsBuilder.RestartInstanceV1Step("Restart the instance", suite.Client.ComputeV1, instance)
    // ... Get after restart, using suites.GetConditionAfterRestarting ...
    stepsBuilder.StopInstanceV1Step("Stop the instance", suite.Client.ComputeV1, instance)
    // ... Get after stop, using suites.GetConditionAfterStopping ...

    // Teardown, in reverse dependency order.
    stepsBuilder.DeleteInstanceV1Step("Delete the instance", t, suite.Client.ComputeV1, instance)
    stepsBuilder.WatchInstanceUntilDeletedV1Step("Watch the instance deletion", t, suite.Client.ComputeV1, instanceWRef)
    stepsBuilder.DeleteBlockStorageV1Step("Delete the block storage", t, suite.Client.StorageV1, block)
    stepsBuilder.WatchBlockStorageUntilDeletedV1Step("Watch the block storage deletion", t, suite.Client.StorageV1, blockWRef)
    stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", t, suite.Client.WorkspaceV1, workspace)
    stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", t, suite.Client.WorkspaceV1, workspaceTRef)

    suite.FinishScenario()
}

func (suite *ProviderLifeCycleV1TestSuite) AfterAll(t provider.T) {
    suite.ResetAllScenarios()
}
```

Naming convention: each step name is a short, human-readable Gherkin-style phrase — it shows up
directly in the Allure report, so make it describe intent ("Create an instance", "Get the started
instance") rather than the mechanics.

If a step you need doesn't exist yet in `internal/conformance/steps/`, add it there rather than
inlining raw SDK calls in the suite — this keeps assertions and API-call wrapping reusable across
scenarios.

## 6. Register the suite with a CLI entrypoint

In the relevant `cmd/conformance/<domain>_test.go` (e.g. `compute_test.go`), construct and
conditionally run the suite inside the domain's `Test<Domain>V1Suites` function:

```go
providerLifeCycleSuite := compute.CreateProviderLifeCycleV1TestSuite(regionalTestSuite,
    &compute.ProviderLifeCycleV1Config{
        AvailableZones: config.Clients.RegionZones,
        InstanceSkus:   config.Clients.InstanceSkus,
        StorageSkus:    config.Clients.StorageSkus,
    },
)
if providerLifeCycleSuite.CanRun(config.Parameters.ScenariosRegexp) {
    suite.RunSuite(t, providerLifeCycleSuite)
}
```

`CanRun` checks the suite's name against the `--scenarios` regexp flag, so the suite becomes
automatically filterable by name (`Compute.V1.ProviderLifeCycle`) without further wiring.

## 7. Run it

```sh
make mock-run                                       # start local WireMock
make run SCENARIOS=Compute.V1.ProviderLifeCycle      # run just the new suite
make report                                          # view the Allure report
```

Or list all available scenarios first with `dist/secatest list` (see the [README](../README.md)
for full CLI usage).

## Checklist

- [ ] Suite name constant added to `internal/constants/suites_v1.go` and `AllSuiteNames`
- [ ] Params struct added to `internal/conformance/params/params_v1.go` (if state is shared)
- [ ] Suite type + `BeforeAll` built with `pkg/builders` / `pkg/generators`
- [ ] Mock stub configurator added under `internal/mock/scenarios/<domain>/`, stubbing every
      call the scenario makes, in call order, including teardown
- [ ] `TestScenario` written using reusable helpers from `internal/conformance/steps/`
- [ ] `AfterAll` calls `suite.ResetAllScenarios()`
- [ ] Suite registered and gated behind `CanRun(...)` in `cmd/conformance/<domain>_test.go`
- [ ] Verified locally against WireMock with `make mock-run && make run SCENARIOS=<name>`
