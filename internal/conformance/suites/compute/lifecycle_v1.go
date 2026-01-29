package compute

import (
	"math/rand"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockcompute "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/compute"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type LifeCycleV1TestSuite struct {
	suites.RegionalTestSuite

	config *LifeCycleV1Config
	params *params.ComputeLifeCycleV1Params
}

type LifeCycleV1Config struct {
	AvailableZones []string
	InstanceSkus   []string
	StorageSkus    []string
}

func NewLifeCycleV1TestSuite(regionalTestSuite suites.RegionalTestSuite, config *LifeCycleV1Config) *LifeCycleV1TestSuite {
	suite := &LifeCycleV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		config:            config,
	}
	suite.ScenarioName = constants.ComputeLifeCycleV1SuiteName.String()
	return suite
}

func (suite *LifeCycleV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("compute")

	// Select skus
	instanceSkuName := suite.config.InstanceSkus[rand.Intn(len(suite.config.InstanceSkus))]
	storageSkuName := suite.config.StorageSkus[rand.Intn(len(suite.config.StorageSkus))]

	// Select zones
	initialInstanceZone := suite.config.AvailableZones[rand.Intn(len(suite.config.AvailableZones))]
	updatedInstanceZone := suite.config.AvailableZones[rand.Intn(len(suite.config.AvailableZones))]

	// Generate scenario data
	workspaceName := generators.GenerateWorkspaceName()

	blockStorageName := generators.GenerateBlockStorageName()
	blockStorageSize := generators.GenerateBlockStorageSize()

	storageSkuRefObj, err := generators.GenerateSkuRefObject(storageSkuName)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	blockStorageRefObj, err := generators.GenerateBlockStorageRefObject(blockStorageName)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	instanceName := generators.GenerateInstanceName()

	instanceSkuRefObj, err := generators.GenerateSkuRefObject(instanceSkuName)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(constants.WorkspaceProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	blockStorage, err := builders.NewBlockStorageBuilder().
		Name(blockStorageName).
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.BlockStorageSpec{
			SkuRef: *storageSkuRefObj,
			SizeGB: blockStorageSize,
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build BlockStorage: %v", err)
	}

	initialInstance, err := builders.NewInstanceBuilder().
		Name(instanceName).
		Provider(constants.ComputeProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.InstanceSpec{
			SkuRef: *instanceSkuRefObj,
			Zone:   initialInstanceZone,
			BootVolume: schema.VolumeReference{
				DeviceRef: *blockStorageRefObj,
			},
		},
		).Build()
	if err != nil {
		t.Fatalf("Failed to build Instance: %v", err)
	}

	updatedInstance, err := builders.NewInstanceBuilder().
		Name(instanceName).
		Provider(constants.ComputeProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.InstanceSpec{
			SkuRef: *instanceSkuRefObj,
			Zone:   updatedInstanceZone,
			BootVolume: schema.VolumeReference{
				DeviceRef: *blockStorageRefObj,
			},
		},
		).Build()
	if err != nil {
		t.Fatalf("Failed to build Instance: %v", err)
	}

	params := &params.ComputeLifeCycleV1Params{
		Workspace:       workspace,
		BlockStorage:    blockStorage,
		InitialInstance: initialInstance,
		UpdatedInstance: updatedInstance,
	}
	suite.params = params

	err = suites.SetupMockIfEnabled(suite.TestSuite, mockcompute.ConfigureLifecycleScenarioV1, params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *LifeCycleV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, constants.ComputeProviderV1, string(schema.RegionalResourceMetadataKindResourceKindWorkspace))

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Workspace

	// Create a workspace
	workspace := suite.params.Workspace
	expectWorkspaceMeta := workspace.Metadata
	expectWorkspaceLabels := workspace.Labels
	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:        expectWorkspaceLabels,
			Metadata:      expectWorkspaceMeta,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created Workspace
	workspaceTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(suite.Tenant),
		Name:   suite.params.Workspace.Metadata.Name,
	}
	stepsBuilder.GetWorkspaceV1Step("Get the created workspace", suite.Client.WorkspaceV1, workspaceTRef,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:        expectWorkspaceLabels,
			Metadata:      expectWorkspaceMeta,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Block storage

	// Create a block storage
	block := suite.params.BlockStorage
	expectedBlockMeta := block.Metadata
	expectedBlockSpec := &block.Spec
	stepsBuilder.CreateOrUpdateBlockStorageV1Step("Create a block storage", suite.Client.StorageV1, block,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			Metadata:      expectedBlockMeta,
			Spec:          expectedBlockSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created block storage
	blockWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.Tenant),
		Workspace: secapi.WorkspaceID(suite.params.Workspace.Metadata.Name),
		Name:      suite.params.BlockStorage.Metadata.Name,
	}
	stepsBuilder.GetBlockStorageV1Step("Get the created block storage", suite.Client.StorageV1, blockWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			Metadata:      expectedBlockMeta,
			Spec:          expectedBlockSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Instance

	// Create an instance
	instance := suite.params.InitialInstance
	expectInstanceMeta := instance.Metadata
	expectInstanceSpec := &instance.Spec
	stepsBuilder.CreateOrUpdateInstanceV1Step("Create an instance", suite.Client.ComputeV1, instance,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			Metadata:      expectInstanceMeta,
			Spec:          expectInstanceSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created instance
	instanceWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.Tenant),
		Workspace: secapi.WorkspaceID(suite.params.Workspace.Metadata.Name),
		Name:      suite.params.InitialInstance.Metadata.Name,
	}
	instance = stepsBuilder.GetInstanceV1Step("Get the created instance", suite.Client.ComputeV1, instanceWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			Metadata:      expectInstanceMeta,
			Spec:          expectInstanceSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Update the instance
	instance.Spec.Zone = suite.params.UpdatedInstance.Spec.Zone
	expectInstanceSpec.Zone = instance.Spec.Zone
	stepsBuilder.CreateOrUpdateInstanceV1Step("Update the instance", suite.Client.ComputeV1, instance,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			Metadata:      expectInstanceMeta,
			Spec:          expectInstanceSpec,
			ResourceState: schema.ResourceStateUpdating,
		},
	)

	// Get the updated instance
	instance = stepsBuilder.GetInstanceV1Step("Get the updated instance", suite.Client.ComputeV1, instanceWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			Metadata:      expectInstanceMeta,
			Spec:          expectInstanceSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Stop the instance
	stepsBuilder.StopInstanceV1Step("Stop the instance", suite.Client.ComputeV1, instance)

	// Get the stoped instance
	instance = stepsBuilder.GetInstanceV1Step("Get the updated instance", suite.Client.ComputeV1, instanceWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			Metadata:      expectInstanceMeta,
			Spec:          expectInstanceSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Start the instance
	stepsBuilder.StartInstanceV1Step("Start the instance", suite.Client.ComputeV1, instance)

	// Get the started instance
	instance = stepsBuilder.GetInstanceV1Step("Get the started instance", suite.Client.ComputeV1, instanceWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			Metadata:      expectInstanceMeta,
			Spec:          expectInstanceSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Restart the instance
	stepsBuilder.RestartInstanceV1Step("Restart the instance", suite.Client.ComputeV1, instance)

	// Get the restarted instance
	// TODO Find an away to assert if the instance is restarted
	instance = stepsBuilder.GetInstanceV1Step("Get the updated instance", suite.Client.ComputeV1, instanceWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			Metadata:      expectInstanceMeta,
			Spec:          expectInstanceSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Resources deletion

	stepsBuilder.DeleteInstanceV1Step("Delete the instance", suite.Client.ComputeV1, instance)
	stepsBuilder.GetInstanceWithErrorV1Step("Get the deleted instance", suite.Client.ComputeV1, instanceWRef, secapi.ErrResourceNotFound)

	stepsBuilder.DeleteBlockStorageV1Step("Delete the block storage", suite.Client.StorageV1, block)
	stepsBuilder.GetBlockStorageWithErrorV1Step("Get the deleted block storage", suite.Client.StorageV1, blockWRef, secapi.ErrResourceNotFound)

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)
	stepsBuilder.GetWorkspaceWithErrorV1Step("Get the deleted workspace", suite.Client.WorkspaceV1, workspaceTRef, secapi.ErrResourceNotFound)

	suite.FinishScenario()
}

func (suite *LifeCycleV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
