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

type ComputeLifeCycleV1TestSuite struct {
	suites.RegionalTestSuite

	config *ComputeLifeCycleV1Config
	params *params.ComputeLifeCycleV1Params
}

type ComputeLifeCycleV1Config struct {
	AvailableZones []string
	InstanceSkus   []string
	StorageSkus    []string
}

func CreateLifeCycleV1TestSuite(regionalTestSuite suites.RegionalTestSuite, config *ComputeLifeCycleV1Config) *ComputeLifeCycleV1TestSuite {
	suite := &ComputeLifeCycleV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		config:            config,
	}
	suite.ScenarioName = constants.ComputeV1LifeCycleSuiteName
	return suite
}

func (suite *ComputeLifeCycleV1TestSuite) BeforeAll(t provider.T) {
	var err error

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

//nolint:dupl
func (suite *ComputeLifeCycleV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, constants.ComputeProviderV1, string(schema.RegionalResourceMetadataKindResourceKindWorkspace), constants.WorkspaceProviderV1, constants.StorageProviderV1)

	workspace := suite.params.Workspace
	workspaceTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(suite.Tenant),
		Name:   suite.params.Workspace.Metadata.Name,
	}

	block := suite.params.BlockStorage
	blockWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.Tenant),
		Workspace: secapi.WorkspaceID(suite.params.Workspace.Metadata.Name),
		Name:      suite.params.BlockStorage.Metadata.Name,
	}

	instance := suite.params.InitialInstance
	instanceWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.Tenant),
		Workspace: secapi.WorkspaceID(suite.params.Workspace.Metadata.Name),
		Name:      suite.params.InitialInstance.Metadata.Name,
	}

	t.WithNewStep("Workspace", func(wsCtx provider.StepCtx) {
		wsSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, wsCtx)

		// Create a workspace

		expectWorkspaceMeta := workspace.Metadata
		expectWorkspaceLabels := workspace.Labels
		wsSteps.CreateOrUpdateWorkspaceV1Step("Create", suite.Client.WorkspaceV1, workspace,
			steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
				Labels:        expectWorkspaceLabels,
				Metadata:      expectWorkspaceMeta,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		// Get the created Workspace

		wsSteps.GetWorkspaceV1Step("Get", suite.Client.WorkspaceV1, workspaceTRef,
			steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
				Labels:        expectWorkspaceLabels,
				Metadata:      expectWorkspaceMeta,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	t.WithNewStep("BlockStorage", func(bsCtx provider.StepCtx) {
		bsSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, bsCtx)

		expectedBlockMeta := block.Metadata
		expectedBlockSpec := &block.Spec
		bsSteps.CreateOrUpdateBlockStorageV1Step("Create", suite.Client.StorageV1, block,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
				Metadata:      expectedBlockMeta,
				Spec:          expectedBlockSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		bsSteps.GetBlockStorageV1Step("Get", suite.Client.StorageV1, blockWRef,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
				Metadata:      expectedBlockMeta,
				Spec:          expectedBlockSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	t.WithNewStep("Instance", func(instCtx provider.StepCtx) {
		instSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, instCtx)

		// Create an instance

		expectInstanceMeta := instance.Metadata
		expectInstanceSpec := &instance.Spec
		instSteps.CreateOrUpdateInstanceV1Step("Create", suite.Client.ComputeV1, instance,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
				Metadata:      expectInstanceMeta,
				Spec:          expectInstanceSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		instance = instSteps.GetInstanceV1Step("Get", suite.Client.ComputeV1, instanceWRef,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
				Metadata:      expectInstanceMeta,
				Spec:          expectInstanceSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)

		// Update the instance
		instance.Spec.Zone = suite.params.UpdatedInstance.Spec.Zone
		expectInstanceSpec.Zone = instance.Spec.Zone
		instSteps.CreateOrUpdateInstanceV1Step("Update", suite.Client.ComputeV1, instance,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
				Metadata:      expectInstanceMeta,
				Spec:          expectInstanceSpec,
				ResourceState: schema.ResourceStateUpdating,
			},
		)

		// Get the updated instance
		instance = instSteps.GetInstanceV1Step("GetUpdated", suite.Client.ComputeV1, instanceWRef,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
				Metadata:      expectInstanceMeta,
				Spec:          expectInstanceSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)

		// Stop the instance
		instSteps.StopInstanceV1Step("Stop", suite.Client.ComputeV1, instance)

		// Get the stopped instance
		instance = instSteps.GetInstanceV1Step("GetStopped", suite.Client.ComputeV1, instanceWRef,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
				Metadata:      expectInstanceMeta,
				Spec:          expectInstanceSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)

		// Start the instance
		instSteps.StartInstanceV1Step("Start", suite.Client.ComputeV1, instance)

		// Get the started instance
		instance = instSteps.GetInstanceV1Step("GetStarted", suite.Client.ComputeV1, instanceWRef,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
				Metadata:      expectInstanceMeta,
				Spec:          expectInstanceSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)

		// Restart the instance
		instSteps.RestartInstanceV1Step("Restart", suite.Client.ComputeV1, instance)

		// Get the restarted instance
		// TODO Find an away to assert if the instance is restarted
		instance = instSteps.GetInstanceV1Step("GetAfterRestart", suite.Client.ComputeV1, instanceWRef,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
				Metadata:      expectInstanceMeta,
				Spec:          expectInstanceSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	// Resources deletion
	t.WithNewStep("Delete", func(dltCtx provider.StepCtx) {
		dltSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, dltCtx)
		dltSteps.DeleteInstanceV1Step("Delete the instance", suite.Client.ComputeV1, instance)
		dltSteps.GetInstanceWithErrorV1Step("Get the deleted instance", suite.Client.ComputeV1, instanceWRef, secapi.ErrResourceNotFound)

		dltSteps.DeleteBlockStorageV1Step("Delete the block storage", suite.Client.StorageV1, block)
		dltSteps.GetBlockStorageWithErrorV1Step("Get the deleted block storage", suite.Client.StorageV1, blockWRef, secapi.ErrResourceNotFound)

		dltSteps.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)
		dltSteps.GetWorkspaceWithErrorV1Step("Get the deleted workspace", suite.Client.WorkspaceV1, workspaceTRef, secapi.ErrResourceNotFound)
	})
	suite.FinishScenario()
}

func (suite *ComputeLifeCycleV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
