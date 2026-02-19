package storage

import (
	"math/rand"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockStorage "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/storage"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	labelBuilder "github.com/eu-sovereign-cloud/go-sdk/secapi/builders"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type ProviderQueriesV1TestSuite struct {
	suites.RegionalTestSuite

	StorageSkus []string

	params *params.StorageProviderQueriesV1Params
}

func CreateProviderQueriesV1TestSuite(regionalTestSuite suites.RegionalTestSuite, storageSkus []string) *ProviderQueriesV1TestSuite {
	suite := &ProviderQueriesV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		StorageSkus:       storageSkus,
	}
	suite.ScenarioName = constants.StorageProviderQueriesV1SuiteName.String()
	return suite
}

func (suite *ProviderQueriesV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("Storage")

	// Select sku
	storageSkuName := suite.StorageSkus[rand.Intn(len(suite.StorageSkus))]

	// Generate scenario data
	workspaceName := generators.GenerateWorkspaceName()

	storageSkuRefObj, err := generators.GenerateSkuRefObject(storageSkuName)
	if err != nil {
		t.Fatal(err)
	}

	blockStorageName1 := generators.GenerateBlockStorageName()
	blockStorageName2 := generators.GenerateBlockStorageName()
	blockStorageName3 := generators.GenerateBlockStorageName()

	blockStorageRefObj, err := generators.GenerateBlockStorageRefObject(blockStorageName1)
	if err != nil {
		t.Fatal(err)
	}

	imageName1 := generators.GenerateImageName()
	imageName2 := generators.GenerateImageName()
	imageName3 := generators.GenerateImageName()
	initialStorageSize := generators.GenerateBlockStorageSize()

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(constants.WorkspaceProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	blockStorage1, err := builders.NewBlockStorageBuilder().
		Name(blockStorageName1).
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).
		Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.BlockStorageSpec{
			SkuRef: *storageSkuRefObj,
			SizeGB: initialStorageSize,
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build BlockStorage: %v", err)
	}

	blockStorage2, err := builders.NewBlockStorageBuilder().
		Name(blockStorageName2).
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).
		Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.BlockStorageSpec{
			SkuRef: *storageSkuRefObj,
			SizeGB: initialStorageSize,
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build BlockStorage: %v", err)
	}

	blockStorage3, err := builders.NewBlockStorageBuilder().
		Name(blockStorageName3).
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Workspace(workspaceName).
		Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.BlockStorageSpec{
			SkuRef: *storageSkuRefObj,
			SizeGB: initialStorageSize,
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build BlockStorage: %v", err)
	}

	blockStorages := []schema.BlockStorage{*blockStorage1, *blockStorage2, *blockStorage3}

	image1, err := builders.NewImageBuilder().
		Name(imageName1).
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.ImageSpec{
			BlockStorageRef: *blockStorageRefObj,
			CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Image: %v", err)
	}

	image2, err := builders.NewImageBuilder().
		Name(imageName2).
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).
		Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.ImageSpec{
			BlockStorageRef: *blockStorageRefObj,
			CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Image: %v", err)
	}

	image3, err := builders.NewImageBuilder().
		Name(imageName3).
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).
		Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		}).
		Spec(&schema.ImageSpec{
			BlockStorageRef: *blockStorageRefObj,
			CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Image: %v", err)
	}
	images := []schema.Image{*image1, *image2, *image3}

	params := &params.StorageProviderQueriesV1Params{
		Workspace:     workspace,
		BlockStorages: blockStorages,
		Images:        images,
	}

	suite.params = params
	err = suites.SetupMockIfEnabledV2(suite.TestSuite, mockStorage.ConfigureProviderQueriesV1, params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *ProviderQueriesV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, constants.StorageProviderV1,
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindBlockStorage),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindImage),
	)

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Workspace
	workspace := suite.params.Workspace

	// Create a workspace
	expectWorkspaceMeta := workspace.Metadata
	expectWorkspaceLabels := workspace.Labels
	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:        expectWorkspaceLabels,
			Metadata:      expectWorkspaceMeta,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Block storage
	blocks := suite.params.BlockStorages

	// Create the block storages
	for _, block := range blocks {
		expectedBlockMeta := block.Metadata
		expectedBlockSpec := &block.Spec
		stepsBuilder.CreateOrUpdateBlockStorageV1Step("Create a block storage", suite.Client.StorageV1, &block,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
				Metadata:      expectedBlockMeta,
				Spec:          expectedBlockSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)
	}

	// List block storages
	wref := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(workspace.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(workspace.Metadata.Name),
	}
	stepsBuilder.GetListBlockStorageV1Step("GetList block storage", suite.Client.StorageV1, wref, nil)

	// List block storages with limit
	stepsBuilder.GetListBlockStorageV1Step("Get List block storage with limit", suite.Client.StorageV1, wref,
		secapi.NewListOptions().WithLimit(1))

	// List block storages with label
	stepsBuilder.GetListBlockStorageV1Step("Get list of block storage with label", suite.Client.StorageV1, wref,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvDevelopmentLabel)))

	// List block storages with limit and label
	stepsBuilder.GetListBlockStorageV1Step("Get list of block storage with limit and label", suite.Client.StorageV1, wref,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvDevelopmentLabel)))

	// Image
	images := suite.params.Images

	// Create images
	for _, image := range images {
		expectedImageMeta := image.Metadata
		expectedImageSpec := &image.Spec
		stepsBuilder.CreateOrUpdateImageV1Step("Create an image", suite.Client.StorageV1, &image,
			steps.ResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
				Metadata:      expectedImageMeta,
				Spec:          expectedImageSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)
	}

	// List images
	tref := secapi.TenantReference{
		Name:   workspace.Metadata.Tenant,
		Tenant: secapi.TenantID(workspace.Metadata.Tenant),
	}
	stepsBuilder.GetListImageV1Step("List image", suite.Client.StorageV1, tref, nil)

	// List images with limit
	stepsBuilder.GetListImageV1Step("Get list of images", suite.Client.StorageV1, tref,
		secapi.NewListOptions().WithLimit(1))

	// List images with label
	stepsBuilder.GetListImageV1Step("Get list of images", suite.Client.StorageV1, tref,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List images with limit and label
	stepsBuilder.GetListImageV1Step("Get list of images", suite.Client.StorageV1, tref,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Skus

	// List Skus
	stepsBuilder.GetListSkuV1Step("List skus", suite.Client.StorageV1, tref, nil)

	// List Skus with limit
	stepsBuilder.GetListSkuV1Step("Get list of skus", suite.Client.StorageV1, tref,
		secapi.NewListOptions().WithLimit(1))

	// Delete all images
	for _, image := range images {
		stepsBuilder.DeleteImageV1Step("Delete image", suite.Client.StorageV1, &image)

		// Get the deleted image
		imageTRef := &secapi.TenantReference{
			Tenant: secapi.TenantID(workspace.Metadata.Tenant),
			Name:   image.Metadata.Name,
		}
		stepsBuilder.GetImageWithErrorV1Step("Get deleted image ", suite.Client.StorageV1, *imageTRef, secapi.ErrResourceNotFound)
	}

	// Delete all block storages
	for _, block := range blocks {
		stepsBuilder.DeleteBlockStorageV1Step("Delete block storage 1", suite.Client.StorageV1, &block)

		// Get the deleted block storage
		blockWRef := &secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(block.Metadata.Tenant),
			Workspace: secapi.WorkspaceID(block.Metadata.Workspace),
			Name:      block.Metadata.Name,
		}
		stepsBuilder.GetBlockStorageWithErrorV1Step("Get deleted block storage 1", suite.Client.StorageV1, *blockWRef, secapi.ErrResourceNotFound)
	}

	// Delete the workspace
	workspaceTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(workspace.Metadata.Tenant),
		Name:   workspace.Metadata.Name,
	}
	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)
	stepsBuilder.GetWorkspaceWithErrorV1Step("Get the deleted workspace", suite.Client.WorkspaceV1, *workspaceTRef, secapi.ErrResourceNotFound)

	suite.FinishScenario()
}

func (suite *ProviderQueriesV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
