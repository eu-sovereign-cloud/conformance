package compute

import (
	"context"
	"math/rand"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockcompute "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/compute"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// InstanceErrorV1TestSuite verifies that Instance resources with invalid references
// are rejected by the API with 422 Unprocessable Entity.
//
// Scenarios tested:
//   - Create instance with non-existent SKU
//   - Create instance with non-existent workspace
//   - Create instance with non-existent boot volume ref
//   - Create instance with invalid zone
type InstanceErrorV1TestSuite struct {
	suites.RegionalTestSuite

	config *InstanceErrorV1Config
	params *params.InstanceErrorV1Params
}

type InstanceErrorV1Config struct {
	AvailableZones []string
	InstanceSkus   []string
	StorageSkus    []string
}

func CreateInstanceErrorV1TestSuite(regionalTestSuite suites.RegionalTestSuite, config *InstanceErrorV1Config) *InstanceErrorV1TestSuite {
	suite := &InstanceErrorV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		config:            config,
	}
	suite.ScenarioName = constants.InstanceErrorV1SuiteName.String()
	return suite
}

func (suite *InstanceErrorV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite(suites.ComputeParentSuite)

	// Select valid skus and zone for setup
	storageSkuName := suite.config.StorageSkus[rand.Intn(len(suite.config.StorageSkus))]
	instanceSkuName := suite.config.InstanceSkus[rand.Intn(len(suite.config.InstanceSkus))]
	zone := suite.config.AvailableZones[rand.Intn(len(suite.config.AvailableZones))]

	workspaceName := generators.GenerateWorkspaceName()
	blockStorageName := generators.GenerateBlockStorageName()

	storageSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, storageSkuName)
	instanceSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.ComputeProviderV1Name, suite.Tenant, instanceSkuName)
	invalidSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.ComputeProviderV1Name, suite.Tenant, "non-existent-sku")
	blockStorageRefObj := generators.GenerateBlockStorageRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, workspaceName, blockStorageName)
	nonExistentBlockStorageRefObj := generators.GenerateBlockStorageRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, workspaceName, "non-existent-block-storage")
	baseLabels := schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}

	// Valid workspace and block storage for setup
	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(baseLabels).
		Annotations(schema.Annotations{"description": "Workspace for instance error scenarios testing"}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	blockStorage, err := builders.NewBlockStorageBuilder().
		Name(blockStorageName).
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(baseLabels).
		Annotations(schema.Annotations{"description": "BlockStorage for instance error scenarios testing"}).
		Spec(&schema.BlockStorageSpec{
			SkuRef: *storageSkuRefObj,
			SizeGB: constants.BlockStorageInitialSize,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build BlockStorage: %v", err)
	}

	buildInstance := func(name string, workspaceRef string, skuRef schema.Reference, bootVolumeRef schema.Reference, instanceZone string) *schema.Instance {
		instance, err := builders.NewInstanceBuilder().
			Name(name).
			Provider(sdkconsts.ComputeProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(suite.Tenant).Workspace(workspaceRef).Region(suite.Region).
			Labels(baseLabels).
			Annotations(schema.Annotations{"description": "Instance for error scenario testing"}).
			Spec(&schema.InstanceSpec{
				SkuRef: skuRef,
				Zone:   instanceZone,
				BootVolume: schema.VolumeReference{
					DeviceRef: bootVolumeRef,
				},
			}).Build()
		if err != nil {
			t.Fatalf("Failed to build Instance: %v", err)
		}
		return instance
	}

	p := &params.InstanceErrorV1Params{
		Workspace:    workspace,
		BlockStorage: blockStorage,

		// Non-existent SKU — valid workspace + zone + boot volume, SKU does not exist
		InvalidSkuInstance: buildInstance(
			generators.GenerateInstanceName(),
			workspaceName,
			*invalidSkuRefObj,
			*blockStorageRefObj,
			zone,
		),

		// Non-existent workspace — workspace was never created
		NonExistentWorkspaceInstance: buildInstance(
			generators.GenerateInstanceName(),
			"non-existent-workspace",
			*instanceSkuRefObj,
			*blockStorageRefObj,
			zone,
		),

		// Non-existent boot volume — valid workspace + SKU + zone, block storage does not exist
		NonExistentBootVolumeInstance: buildInstance(
			generators.GenerateInstanceName(),
			workspaceName,
			*instanceSkuRefObj,
			*nonExistentBlockStorageRefObj,
			zone,
		),

		// Invalid zone — zone string not in provider available zones
		InvalidZoneInstance: buildInstance(
			generators.GenerateInstanceName(),
			workspaceName,
			*instanceSkuRefObj,
			*blockStorageRefObj,
			"invalid-zone",
		),
		PowerInstance: buildInstance(
			generators.GenerateInstanceName(),
			workspaceName,
			*instanceSkuRefObj,
			*blockStorageRefObj,
			zone,
		),
	}

	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mockcompute.ConfigureInstanceErrorV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *InstanceErrorV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t, sdkconsts.ComputeProviderV1Name)
	suite.ConfigureResources(t, string(schema.RegionalResourceMetadataKindResourceKindInstance))
	suite.ConfigureDepends(t,
		string(schema.RegionalResourceMetadataKindResourceKindWorkspace),
		string(schema.RegionalResourceMetadataKindResourceKindBlockStorage),
	)

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Workspace setup
	workspace := suite.params.Workspace
	expectWorkspaceMeta := workspace.Metadata
	expectWorkspaceLabels := workspace.Labels
	expectWorkspaceAnnotations := workspace.Annotations
	expectWorkspaceExtensions := workspace.Extensions

	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", t, suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:         expectWorkspaceLabels,
			Annotations:    expectWorkspaceAnnotations,
			Extensions:     expectWorkspaceExtensions,
			Metadata:       expectWorkspaceMeta,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	workspaceTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(workspace.Metadata.Tenant),
		Name:   workspace.Metadata.Name,
	}
	stepsBuilder.GetWorkspaceV1Step("Get the created workspace", suite.Client.WorkspaceV1, workspaceTRef,
		steps.ResponseExpectsWithCondition[schema.RegionalResourceMetadata, schema.WorkspaceSpec, schema.WorkspaceStatus]{
			Labels:   expectWorkspaceLabels,
			Metadata: expectWorkspaceMeta,
			ResourceStatus: schema.WorkspaceStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Block storage setup
	block := suite.params.BlockStorage
	expectedBlockMeta := block.Metadata
	expectedBlockSpec := &block.Spec
	expectedBlockLabels := block.Labels
	expectedBlockAnnotations := block.Annotations
	expectedBlockExtensions := block.Extensions

	stepsBuilder.CreateOrUpdateBlockStorageV1Step("Create a block storage", t, suite.Client.StorageV1, block,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			Metadata:       expectedBlockMeta,
			Labels:         expectedBlockLabels,
			Annotations:    expectedBlockAnnotations,
			Extensions:     expectedBlockExtensions,
			Spec:           expectedBlockSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	blockWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(block.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(block.Metadata.Workspace),
		Name:      block.Metadata.Name,
	}
	stepsBuilder.GetBlockStorageV1Step("Get the created block storage", suite.Client.StorageV1, blockWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec, schema.BlockStorageStatus]{
			Metadata: expectedBlockMeta,
			Spec:     expectedBlockSpec,
			ResourceStatus: schema.BlockStorageStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Error scenarios — all must be rejected with 422
	stepsBuilder.CreateOrUpdateInstanceExpectViolationV1Step(
		"Create an instance with non-existent SKU — expect rejection",
		suite.Client.ComputeV1,
		suite.params.InvalidSkuInstance,
	)

	stepsBuilder.CreateOrUpdateInstanceExpectViolationV1Step(
		"Create an instance with non-existent workspace — expect rejection",
		suite.Client.ComputeV1,
		suite.params.NonExistentWorkspaceInstance,
	)

	stepsBuilder.CreateOrUpdateInstanceExpectViolationV1Step(
		"Create an instance with non-existent boot volume ref — expect rejection",
		suite.Client.ComputeV1,
		suite.params.NonExistentBootVolumeInstance,
	)

	stepsBuilder.CreateOrUpdateInstanceExpectViolationV1Step(
		"Create an instance with invalid zone — expect rejection",
		suite.Client.ComputeV1,
		suite.params.InvalidZoneInstance,
	)

	// Power state error scenarios
	// Started instance: PUT → creating → active → powerState=on
	powerInstance := suite.params.PowerInstance
	expectInstanceMeta := powerInstance.Metadata
	expectInstanceSpec := &powerInstance.Spec
	expectInstanceLabels := powerInstance.Labels
	expectInstanceAnnotations := powerInstance.Annotations
	expectInstanceExtensions := powerInstance.Extensions
	stepsBuilder.CreateOrUpdateInstanceV1Step("Create an instance", t, suite.Client.ComputeV1, powerInstance,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			Metadata:       expectInstanceMeta,
			Labels:         expectInstanceLabels,
			Annotations:    expectInstanceAnnotations,
			Extensions:     expectInstanceExtensions,
			Spec:           expectInstanceSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created instance
	instanceWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.Tenant),
		Workspace: secapi.WorkspaceID(suite.params.Workspace.Metadata.Name),
		Name:      suite.params.PowerInstance.Metadata.Name,
	}
	powerInstance = stepsBuilder.GetInstanceV1Step("Get the created instance", suite.Client.ComputeV1, instanceWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus]{
			Metadata: expectInstanceMeta,
			Spec:     expectInstanceSpec,
			ResourceStatus: schema.InstanceStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Start
	stepsBuilder.StartInstanceV1Step("Start the instance", suite.Client.ComputeV1, powerInstance)

	// Get the started instance
	powerInstance = stepsBuilder.GetInstanceV1Step("Get the started instance", suite.Client.ComputeV1, instanceWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus]{
			Metadata: expectInstanceMeta,
			Spec:     expectInstanceSpec,
			ResourceStatus: schema.InstanceStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterStartingWithoutUpdate,
			},
		},
	)

	// Start on already-on instance → 412
	stepsBuilder.InstanceOperationExpectConflictV1Step(
		"Start an already-on instance — expect 412 conflict",
		suite.Client.ComputeV1,
		powerInstance,
		func(ctx context.Context, r *schema.Instance) error {
			return suite.Client.ComputeV1.StartInstance(ctx, r)
		},
		constants.StartInstanceOperation,
	)
	// Stopped instance: PUT  → powerState=off
	stepsBuilder.StopInstanceV1Step("Stop the instance", suite.Client.ComputeV1, powerInstance)

	// Stop on already-off instance → 409
	stepsBuilder.InstanceOperationExpectConflictV1Step(
		"Stop an already-off instance — expect 409 conflict",
		suite.Client.ComputeV1,
		powerInstance,
		func(ctx context.Context, r *schema.Instance) error {
			return suite.Client.ComputeV1.StopInstance(ctx, r)
		},
		constants.StopInstanceOperation,
	)

	// Restart on already-off instance → 409
	stepsBuilder.InstanceOperationExpectConflictV1Step(
		"Restart an already-off instance — expect 409 conflict",
		suite.Client.ComputeV1,
		powerInstance,
		func(ctx context.Context, r *schema.Instance) error {
			return suite.Client.ComputeV1.RestartInstance(ctx, r)
		},
		constants.RestartInstanceOperation,
	)

	// Teardown — reverse dependency order
	stepsBuilder.DeleteInstanceV1Step("Delete started instance", t, suite.Client.ComputeV1, powerInstance)
	stepsBuilder.WatchInstanceUntilDeletedV1Step("Watch started instance deletion", t, suite.Client.ComputeV1, instanceWRef)

	stepsBuilder.DeleteBlockStorageV1Step("Delete the block storage", t, suite.Client.StorageV1, block)
	stepsBuilder.WatchBlockStorageUntilDeletedV1Step("Watch the block storage deletion", t, suite.Client.StorageV1, blockWRef)

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", t, suite.Client.WorkspaceV1, workspace)
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", t, suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *InstanceErrorV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
