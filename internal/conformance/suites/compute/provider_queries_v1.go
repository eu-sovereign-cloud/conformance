package compute

import (
	"math/rand"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockCompute "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/compute"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	labelBuilder "github.com/eu-sovereign-cloud/go-sdk/secapi/builders"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type ProviderQueriesV1TestSuite struct {
	suites.RegionalTestSuite

	config *ProviderQueriesV1Config
	params *params.ComputeProviderQueriesV1Params
}

type ProviderQueriesV1Config struct {
	AvailableZones []string
	InstanceSkus   []string
	StorageSkus    []string
}

func CreateProviderQueriesV1TestSuite(regionalTestSuite suites.RegionalTestSuite, config *ProviderQueriesV1Config) *ProviderQueriesV1TestSuite {
	suite := &ProviderQueriesV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		config:            config,
	}
	suite.ScenarioName = constants.ComputeProviderQueriesV1SuiteName.String()
	return suite
}

func (suite *ProviderQueriesV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("Compute")

	// Select skus
	instanceSkuName := suite.config.InstanceSkus[rand.Intn(len(suite.config.InstanceSkus))]
	storageSkuName := suite.config.StorageSkus[rand.Intn(len(suite.config.StorageSkus))]

	// Select zones
	zone := suite.config.AvailableZones[rand.Intn(len(suite.config.AvailableZones))]

	// Generate scenario data
	workspaceName := generators.GenerateWorkspaceName()

	instanceSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.ComputeProviderV1Name, suite.Tenant, instanceSkuName)

	instanceName1 := generators.GenerateInstanceName()
	instanceName2 := generators.GenerateInstanceName()
	instanceName3 := generators.GenerateInstanceName()

	storageSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, storageSkuName)

	blockStorageName := generators.GenerateBlockStorageName()
	blockStorageRefObj := generators.GenerateBlockStorageRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, workspaceName, blockStorageName)
	blockStorageSize := constants.BlockStorageInitialSize

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	blockStorage, err := builders.NewBlockStorageBuilder().
		Name(blockStorageName).
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.BlockStorageSpec{
			SkuRef: *storageSkuRefObj,
			SizeGB: blockStorageSize,
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build BlockStorage: %v", err)
	}

	instance1, err := builders.NewInstanceBuilder().
		Name(instanceName1).
		Provider(sdkconsts.ComputeProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.InstanceSpec{
			SkuRef: *instanceSkuRefObj,
			Zone:   zone,
			BootVolume: schema.VolumeReference{
				DeviceRef: *blockStorageRefObj,
			},
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Instance: %v", err)
	}

	instance2, err := builders.NewInstanceBuilder().
		Name(instanceName2).
		Provider(sdkconsts.ComputeProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.InstanceSpec{
			SkuRef: *instanceSkuRefObj,
			Zone:   zone,
			BootVolume: schema.VolumeReference{
				DeviceRef: *blockStorageRefObj,
			},
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Instance: %v", err)
	}

	instance3, err := builders.NewInstanceBuilder().
		Name(instanceName3).
		Provider(sdkconsts.ComputeProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.InstanceSpec{
			SkuRef: *instanceSkuRefObj,
			Zone:   zone,
			BootVolume: schema.VolumeReference{
				DeviceRef: *blockStorageRefObj,
			},
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Instance: %v", err)
	}

	instances := []schema.Instance{*instance1, *instance2, *instance3}

	params := &params.ComputeProviderQueriesV1Params{
		Workspace:    workspace,
		BlockStorage: blockStorage,
		Instances:    instances,
	}
	suite.params = params
	err = suites.SetupMockIfEnabled(suite.TestSuite, mockCompute.ConfigureProviderQueriesV1, *params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *ProviderQueriesV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, sdkconsts.ComputeProviderV1Name, string(schema.RegionalResourceMetadataKindResourceKindWorkspace))

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Workspace
	workspace := suite.params.Workspace

	// Create a workspace
	expectWorkspaceMeta := workspace.Metadata
	expectWorkspaceLabels := workspace.Labels
	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", suite.Client.WorkspaceV1, workspace,
		steps.StepResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:         expectWorkspaceLabels,
			Metadata:       expectWorkspaceMeta,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Block storage
	block := suite.params.BlockStorage

	// Create a block storage
	expectedBlockMeta := block.Metadata
	expectedBlockSpec := &block.Spec
	stepsBuilder.CreateOrUpdateBlockStorageV1Step("Create a block storage", suite.Client.StorageV1, block,
		steps.StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			Metadata:       expectedBlockMeta,
			Spec:           expectedBlockSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Instance
	instances := suite.params.Instances

	// Create instances
	for _, instance := range instances {
		expectInstanceMeta := instance.Metadata
		expectInstanceSpec := &instance.Spec
		stepsBuilder.CreateOrUpdateInstanceV1Step("Create an instance", suite.Client.ComputeV1, &instance,
			steps.StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
				Metadata:       expectInstanceMeta,
				Spec:           expectInstanceSpec,
				ResourceStates: suites.CreatedResourceExpectedStates,
			},
		)
	}

	// List instances
	wpath := secapi.WorkspacePath{
		Tenant:    secapi.TenantID(workspace.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(workspace.Metadata.Name),
	}
	stepsBuilder.ListInstanceV1Step("List instances", suite.Client.ComputeV1, wpath, nil)

	// List instances with limit
	stepsBuilder.ListInstanceV1Step("List instances", suite.Client.ComputeV1, wpath,
		secapi.NewListOptions().WithLimit(1))

	// List Instances with label
	stepsBuilder.ListInstanceV1Step("List instances", suite.Client.ComputeV1, wpath,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List Instances with limit and label
	stepsBuilder.ListInstanceV1Step("List instances", suite.Client.ComputeV1, wpath,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Skus

	// List skus
	stepsBuilder.ListSkusV1Step("List skus", suite.Client.ComputeV1, secapi.TenantPath{Tenant: secapi.TenantID(workspace.Metadata.Tenant)}, nil)

	// List skus with limit
	stepsBuilder.ListSkusV1Step("List skus", suite.Client.ComputeV1, secapi.TenantPath{Tenant: secapi.TenantID(workspace.Metadata.Tenant)},
		secapi.NewListOptions().WithLimit(1))

	// List skus with label
	stepsBuilder.ListSkusV1Step("List skus with label", suite.Client.ComputeV1, secapi.TenantPath{Tenant: secapi.TenantID(workspace.Metadata.Tenant)},
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.TierLabel, constants.TierSkuD2XSLabel)))

	// List skus with limit and label
	stepsBuilder.ListSkusV1Step("List skus with limit and label", suite.Client.ComputeV1, secapi.TenantPath{Tenant: secapi.TenantID(workspace.Metadata.Tenant)},
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.TierLabel, constants.TierSkuD2XSLabel)))

	// Delete all instances
	for _, instance := range instances {
		stepsBuilder.DeleteInstanceV1Step("Delete the instance", suite.Client.ComputeV1, &instance)

		// Get the deleted instance
		instanceWRef := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(instance.Metadata.Tenant),
			Workspace: secapi.WorkspaceID(instance.Metadata.Workspace),
			Name:      instance.Metadata.Name,
		}
		stepsBuilder.WatchInstanceUntilDeletedV1Step("Watch the instance deletion", suite.Client.ComputeV1, instanceWRef)
	}

	// Delete the block storage
	stepsBuilder.DeleteBlockStorageV1Step("Delete the block storage", suite.Client.StorageV1, block)

	// Get the deleted block storage
	blockWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(block.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(block.Metadata.Workspace),
		Name:      block.Metadata.Name,
	}
	stepsBuilder.WatchBlockStorageUntilDeletedV1Step("Watch the block storage deletion", suite.Client.StorageV1, blockWRef)

	// Delete the workspace
	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)

	// Get the deleted workspace
	workspaceTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(workspace.Metadata.Tenant),
		Name:   workspace.Metadata.Name,
	}
	stepsBuilder.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *ProviderQueriesV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
