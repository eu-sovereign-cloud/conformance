package workspace

import (
	"math/rand"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockworkspace "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/workspace"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// WorkspaceErrorV1TestSuite verifies that Workspace resources with invalid references
// are rejected by the API with 422 Unprocessable Entity, and that operations conflicting
// with the current state of a workspace are rejected with 409 Conflict.
//
// Scenarios tested:
//   - Create workspace with non-existent region
//   - Delete a workspace that still has an active instance — expect 409 conflict
type WorkspaceErrorV1TestSuite struct {
	suites.RegionalTestSuite

	config *WorkspaceErrorV1Config
	params *params.WorkspaceErrorV1Params
}

type WorkspaceErrorV1Config struct {
	AvailableZones []string
	InstanceSkus   []string
	StorageSkus    []string
}

func CreateWorkspaceErrorV1TestSuite(regionalTestSuite suites.RegionalTestSuite, config *WorkspaceErrorV1Config) *WorkspaceErrorV1TestSuite {
	suite := &WorkspaceErrorV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		config:            config,
	}
	suite.ScenarioName = constants.WorkspaceErrorV1SuiteName.String()
	return suite
}

func (suite *WorkspaceErrorV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite(suites.WorkspaceParentSuite)

	baseLabels := schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}

	// Non-existent region — region string does not exist in the platform
	nonExistentRegionWorkspace, err := builders.NewWorkspaceBuilder().
		Name(generators.GenerateWorkspaceName()).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region("non-existent-region").
		Labels(baseLabels).
		Annotations(schema.Annotations{"description": "Workspace with non-existent region"}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	// Select valid skus and zone for the delete-conflict scenario
	storageSkuName := suite.config.StorageSkus[rand.Intn(len(suite.config.StorageSkus))]
	instanceSkuName := suite.config.InstanceSkus[rand.Intn(len(suite.config.InstanceSkus))]
	zone := suite.config.AvailableZones[rand.Intn(len(suite.config.AvailableZones))]

	workspaceName := generators.GenerateWorkspaceName()
	blockStorageName := generators.GenerateBlockStorageName()

	storageSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, storageSkuName)
	instanceSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.ComputeProviderV1Name, suite.Tenant, instanceSkuName)
	blockStorageRefObj := generators.GenerateBlockStorageRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, workspaceName, blockStorageName)

	// Workspace that will be kept alive by an active instance
	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(baseLabels).
		Annotations(schema.Annotations{"description": "Workspace for delete conflict scenario testing"}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	blockStorage, err := builders.NewBlockStorageBuilder().
		Name(blockStorageName).
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(baseLabels).
		Annotations(schema.Annotations{"description": "BlockStorage for delete conflict scenario testing"}).
		Spec(&schema.BlockStorageSpec{
			SkuRef: *storageSkuRefObj,
			SizeGB: constants.BlockStorageInitialSize,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build BlockStorage: %v", err)
	}

	instance, err := builders.NewInstanceBuilder().
		Name(generators.GenerateInstanceName()).
		Provider(sdkconsts.ComputeProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(baseLabels).
		Annotations(schema.Annotations{"description": "Instance keeping the workspace active for delete conflict testing"}).
		Spec(&schema.InstanceSpec{
			SkuRef: *instanceSkuRefObj,
			Zone:   zone,
			BootVolume: schema.VolumeReference{
				DeviceRef: *blockStorageRefObj,
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Instance: %v", err)
	}

	p := &params.WorkspaceErrorV1Params{
		NonExistentRegionWorkspace: nonExistentRegionWorkspace,

		Workspace:    workspace,
		BlockStorage: blockStorage,
		Instance:     instance,
	}

	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mockworkspace.ConfigureWorkspaceErrorV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *WorkspaceErrorV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t, sdkconsts.WorkspaceProviderV1Name)
	suite.ConfigureResources(t, string(schema.RegionalResourceMetadataKindResourceKindWorkspace))
	suite.ConfigureDepends(t,
		string(schema.RegionalResourceMetadataKindResourceKindBlockStorage),
		string(schema.RegionalResourceMetadataKindResourceKindInstance),
	)

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Error scenarios — all must be rejected with 422
	stepsBuilder.CreateOrUpdateWorkspaceExpectViolationV1Step(
		"Create a workspace with non-existent region — expect rejection",
		suite.Client.WorkspaceV1,
		suite.params.NonExistentRegionWorkspace,
	)

	// Delete conflict scenario — workspace with an active instance must not be deletable
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

	instance := suite.params.Instance
	expectInstanceMeta := instance.Metadata
	expectInstanceSpec := &instance.Spec
	expectInstanceLabels := instance.Labels
	expectInstanceAnnotations := instance.Annotations
	expectInstanceExtensions := instance.Extensions

	stepsBuilder.CreateOrUpdateInstanceV1Step("Create an instance", t, suite.Client.ComputeV1, instance,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			Metadata:       expectInstanceMeta,
			Labels:         expectInstanceLabels,
			Annotations:    expectInstanceAnnotations,
			Extensions:     expectInstanceExtensions,
			Spec:           expectInstanceSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	instanceWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(instance.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(instance.Metadata.Workspace),
		Name:      instance.Metadata.Name,
	}
	stepsBuilder.GetInstanceV1Step("Get the created instance", suite.Client.ComputeV1, instanceWRef,
		steps.ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus]{
			Metadata: expectInstanceMeta,
			Spec:     expectInstanceSpec,
			ResourceStatus: schema.InstanceStatus{
				State:      schema.ResourceStateActive,
				Conditions: suites.GetConditionAfterCreating,
			},
		},
	)

	// Delete the workspace while the instance is still active — expect 409 conflict
	stepsBuilder.DeleteWorkspaceExpectConflictV1Step(
		"Delete a workspace with an active instance — expect 409 conflict",
		suite.Client.WorkspaceV1,
		workspace,
	)

	// Teardown — reverse dependency order
	stepsBuilder.DeleteInstanceV1Step("Delete the instance", t, suite.Client.ComputeV1, instance)
	stepsBuilder.WatchInstanceUntilDeletedV1Step("Watch the instance deletion", t, suite.Client.ComputeV1, instanceWRef)

	stepsBuilder.DeleteBlockStorageV1Step("Delete the block storage", t, suite.Client.StorageV1, block)
	stepsBuilder.WatchBlockStorageUntilDeletedV1Step("Watch the block storage deletion", t, suite.Client.StorageV1, blockWRef)

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", t, suite.Client.WorkspaceV1, workspace)
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", t, suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *WorkspaceErrorV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
