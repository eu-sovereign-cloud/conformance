package storage

import (
	"math/rand"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockstorage "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/storage"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type ProviderLifeCycleV1TestSuite struct {
	suites.RegionalTestSuite

	StorageSkus []string

	params *params.StorageProviderLifeCycleV1Params
}

func CreateProviderLifeCycleV1TestSuite(regionalTestSuite suites.RegionalTestSuite, storageSkus []string) *ProviderLifeCycleV1TestSuite {
	suite := &ProviderLifeCycleV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		StorageSkus:       storageSkus,
	}
	suite.ScenarioName = constants.StorageProviderLifeCycleV1SuiteName.String()
	return suite
}

func (suite *ProviderLifeCycleV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("Storage")

	// Select sku
	storageSkuName := suite.StorageSkus[rand.Intn(len(suite.StorageSkus))]

	// Generate scenario data
	workspaceName := generators.GenerateWorkspaceName()

	storageSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, storageSkuName)

	blockStorageName := generators.GenerateBlockStorageName()
	blockStorageRefObj := generators.GenerateBlockStorageRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, workspaceName, blockStorageName)

	imageName := generators.GenerateImageName()

	initialStorageSize := constants.BlockStorageInitialSize
	updatedStorageSize := constants.BlockStorageUpdatedSize

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvDevelopmentLabel,
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	blockStorageInitial, err := builders.NewBlockStorageBuilder().
		Name(blockStorageName).
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.BlockStorageSpec{
			SkuRef: *storageSkuRefObj,
			SizeGB: initialStorageSize,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build BlockStorage: %v", err)
	}

	blockStorageUpdated, err := builders.NewBlockStorageBuilder().
		Name(blockStorageName).
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.BlockStorageSpec{
			SkuRef: *storageSkuRefObj,
			SizeGB: updatedStorageSize,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build BlockStorage: %v", err)
	}

	imageInitial, err := builders.NewImageBuilder().
		Name(imageName).
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Spec(&schema.ImageSpec{
			BlockStorageRef: *blockStorageRefObj,
			CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Image: %v", err)
	}

	imageUpdated, err := builders.NewImageBuilder().
		Name(imageName).
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Spec(&schema.ImageSpec{
			BlockStorageRef: *blockStorageRefObj,
			CpuArchitecture: schema.ImageSpecCpuArchitectureArm64,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Image: %v", err)
	}

	params := &params.StorageProviderLifeCycleV1Params{
		Workspace:           workspace,
		BlockStorageInitial: blockStorageInitial,
		BlockStorageUpdated: blockStorageUpdated,
		ImageInitial:        imageInitial,
		ImageUpdated:        imageUpdated,
	}
	suite.params = params
	err = suites.SetupMockIfEnabled(suite.TestSuite, mockstorage.ConfigureProviderLifecycleScenarioV1, *params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *ProviderLifeCycleV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, sdkconsts.StorageProviderV1Name,
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindBlockStorage),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindImage),
	)

	stepsConfigurator := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Workspace

	// Create a workspace
	workspace := suite.params.Workspace
	expectWorkspaceMeta := workspace.Metadata
	expectWorkspaceLabels := workspace.Labels

	stepsConfigurator.CreateOrUpdateWorkspaceV1Step("Create a workspace", t, suite.Client.WorkspaceV1, workspace,
		steps.StepResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:         expectWorkspaceLabels,
			Metadata:       expectWorkspaceMeta,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created Workspace
	workspaceTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(workspace.Metadata.Tenant),
		Name:   workspace.Metadata.Name,
	}
	stepsConfigurator.GetWorkspaceV1Step("Get the created workspace", suite.Client.WorkspaceV1, workspaceTRef,
		steps.StepResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:         expectWorkspaceLabels,
			Metadata:       expectWorkspaceMeta,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Block storage

	// Create a block storage
	block := suite.params.BlockStorageInitial
	expectedBlockMeta := block.Metadata
	expectedBlockSpec := &block.Spec
	stepsConfigurator.CreateOrUpdateBlockStorageV1Step("Create a block storage", t, suite.Client.StorageV1, block,
		steps.StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			Metadata:       expectedBlockMeta,
			Spec:           expectedBlockSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created block storage
	blockWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(block.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(block.Metadata.Workspace),
		Name:      block.Metadata.Name,
	}
	stepsConfigurator.GetBlockStorageV1Step("Get the created block storage", suite.Client.StorageV1, blockWRef,
		steps.StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			Metadata:       expectedBlockMeta,
			Spec:           expectedBlockSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Update the block storage
	block = suite.params.BlockStorageUpdated
	expectedBlockSpec.SizeGB = block.Spec.SizeGB
	stepsConfigurator.CreateOrUpdateBlockStorageV1Step("Update the block storage", t, suite.Client.StorageV1, block,
		steps.StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			Metadata:       expectedBlockMeta,
			Spec:           expectedBlockSpec,
			ResourceStates: suites.UpdatedResourceExpectedStates,
		},
	)

	// Get the updated block storage
	block = stepsConfigurator.GetBlockStorageV1Step("Get the updated block storage", suite.Client.StorageV1, blockWRef,
		steps.StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			Metadata:       expectedBlockMeta,
			Spec:           expectedBlockSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Image

	// Create an image
	image := suite.params.ImageInitial
	expectedImageMeta := image.Metadata
	expectedImageSpec := &image.Spec
	stepsConfigurator.CreateOrUpdateImageV1Step("Create an image", t, suite.Client.StorageV1, image,
		steps.StepResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
			Metadata:       expectedImageMeta,
			Spec:           expectedImageSpec,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Get the created image
	imageTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(image.Metadata.Tenant),
		Name:   image.Metadata.Name,
	}
	stepsConfigurator.GetImageV1Step("Get the created image", suite.Client.StorageV1, imageTRef,
		steps.StepResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
			Metadata:       expectedImageMeta,
			Spec:           expectedImageSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Update the image
	image = suite.params.ImageUpdated
	expectedImageSpec.CpuArchitecture = image.Spec.CpuArchitecture
	stepsConfigurator.CreateOrUpdateImageV1Step("Update the image", t, suite.Client.StorageV1, image,
		steps.StepResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
			Metadata:       expectedImageMeta,
			Spec:           expectedImageSpec,
			ResourceStates: suites.UpdatedResourceExpectedStates,
		},
	)

	// Get the updated image
	image = stepsConfigurator.GetImageV1Step("Get the updated image", suite.Client.StorageV1, imageTRef,
		steps.StepResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
			Metadata:       expectedImageMeta,
			Spec:           expectedImageSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Resources deletion

	stepsConfigurator.DeleteImageV1Step("Delete the image", t, suite.Client.StorageV1, image)
	stepsConfigurator.WatchImageUntilDeletedV1Step("Watch the image deletion", t, suite.Client.StorageV1, imageTRef)

	stepsConfigurator.DeleteBlockStorageV1Step("Delete the block storage", t, suite.Client.StorageV1, block)
	stepsConfigurator.WatchBlockStorageUntilDeletedV1Step("Watch the block storage deletion", t, suite.Client.StorageV1, blockWRef)

	stepsConfigurator.DeleteWorkspaceV1Step("Delete the workspace", t, suite.Client.WorkspaceV1, workspace)
	stepsConfigurator.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", t, suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *ProviderLifeCycleV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
