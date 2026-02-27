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

type ImageLifeCycleV1TestSuite struct {
	suites.RegionalTestSuite

	StorageSkus []string

	params *params.ImageLifeCycleV1Params
}

func CreateImageLifeCycleV1TestSuite(regionalTestSuite suites.RegionalTestSuite, storageSkus []string) *ImageLifeCycleV1TestSuite {
	suite := &ImageLifeCycleV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		StorageSkus:       storageSkus,
	}
	suite.ScenarioName = constants.ImageLifeCycleV1SuiteName.String()
	return suite
}

func (suite *ImageLifeCycleV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("Storage")

	// Select sku
	storageSkuName := suite.StorageSkus[rand.Intn(len(suite.StorageSkus))]

	// Generate scenario data
	workspaceName := generators.GenerateWorkspaceName()

	storageSkuRefObj := generators.GenerateSkuRefObject(storageSkuName)

	blockStorageName := generators.GenerateBlockStorageName()
	blockStorageRefObj := generators.GenerateBlockStorageRefObject(blockStorageName)

	imageName := generators.GenerateImageName()

	storageSize := generators.GenerateBlockStorageSize()

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

	blockStorage, err := builders.NewBlockStorageBuilder().
		Name(blockStorageName).
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
		Spec(&schema.BlockStorageSpec{
			SkuRef: *storageSkuRefObj,
			SizeGB: storageSize,
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

	params := &params.ImageLifeCycleV1Params{
		Workspace:    workspace,
		BlockStorage: blockStorage,
		ImageInitial: imageInitial,
		ImageUpdated: imageUpdated,
	}

	suite.params = params

	err = suites.SetupMockIfEnabled(suite.TestSuite, mockstorage.ConfigureImageLifecycleScenarioV1, params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *ImageLifeCycleV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, sdkconsts.StorageProviderV1Name,
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindImage),
	)

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Workspace

	// Create a workspace
	workspace := suite.params.Workspace
	expectWorkspaceMeta := workspace.Metadata
	expectWorkspaceLabels := workspace.Labels

	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
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
	stepsBuilder.GetWorkspaceV1Step("Get the created workspace", suite.Client.WorkspaceV1, workspaceTRef,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:         expectWorkspaceLabels,
			Metadata:       expectWorkspaceMeta,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Block storage

	// Create a block storage
	block := suite.params.BlockStorage
	expectedBlockMeta := block.Metadata
	expectedBlockSpec := &block.Spec
	stepsBuilder.CreateOrUpdateBlockStorageV1Step("Create a block storage", suite.Client.StorageV1, block,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
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
	block = stepsBuilder.GetBlockStorageV1Step("Get the created block storage", suite.Client.StorageV1, blockWRef,
		steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
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
	stepsBuilder.CreateOrUpdateImageV1Step("Create an image", suite.Client.StorageV1, image,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
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
	image = stepsBuilder.GetImageV1Step("Get the created image", suite.Client.StorageV1, imageTRef,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
			Metadata:       expectedImageMeta,
			Spec:           expectedImageSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Update the image
	image.Spec = suite.params.ImageUpdated.Spec
	expectedImageSpec.CpuArchitecture = image.Spec.CpuArchitecture
	stepsBuilder.CreateOrUpdateImageV1Step("Update the image", suite.Client.StorageV1, image,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
			Metadata:       expectedImageMeta,
			Spec:           expectedImageSpec,
			ResourceStates: suites.UpdatedResourceExpectedStates,
		},
	)

	// Get the updated image
	image = stepsBuilder.GetImageV1Step("Get the updated image", suite.Client.StorageV1, imageTRef,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
			Metadata:       expectedImageMeta,
			Spec:           expectedImageSpec,
			ResourceStates: []schema.ResourceState{schema.ResourceStateActive},
		},
	)

	// Resources deletion

	stepsBuilder.DeleteImageV1Step("Delete the image", suite.Client.StorageV1, image)
	stepsBuilder.GetImageWithErrorV1Step("Get the deleted image", suite.Client.StorageV1, imageTRef, secapi.ErrResourceNotFound)

	stepsBuilder.DeleteBlockStorageV1Step("Delete the block storage", suite.Client.StorageV1, block)
	stepsBuilder.GetBlockStorageWithErrorV1Step("Get the deleted block storage", suite.Client.StorageV1, blockWRef, secapi.ErrResourceNotFound)

	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)
	stepsBuilder.GetWorkspaceWithErrorV1Step("Get the deleted workspace", suite.Client.WorkspaceV1, workspaceTRef, secapi.ErrResourceNotFound)

	suite.FinishScenario()
}

func (suite *ImageLifeCycleV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
