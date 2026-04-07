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
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
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

	storageSkuRefObj := generators.GenerateSkuRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, storageSkuName)

	blockStorageName1 := generators.GenerateBlockStorageName()
	blockStorageName2 := generators.GenerateBlockStorageName()
	blockStorageName3 := generators.GenerateBlockStorageName()
	blockStorageRefObj := generators.GenerateBlockStorageRefObject(sdkconsts.StorageProviderV1Name, suite.Tenant, workspaceName, blockStorageName1)

	imageName1 := generators.GenerateImageName()
	imageName2 := generators.GenerateImageName()
	imageName3 := generators.GenerateImageName()

	initialStorageSize := constants.BlockStorageInitialSize

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

	blockStorage1, err := builders.NewBlockStorageBuilder().
		Name(blockStorageName1).
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
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
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
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
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
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
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
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
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
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
		Provider(sdkconsts.StorageProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
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
	err = suites.SetupMockIfEnabled(suite.TestSuite, mockStorage.ConfigureProviderQueriesV1, *params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *ProviderQueriesV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, sdkconsts.StorageProviderV1Name,
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindBlockStorage),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindImage),
	)

	stepsConfigurator := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Workspace
	workspace := suite.params.Workspace

	// Create a workspace
	expectWorkspaceMeta := workspace.Metadata
	expectWorkspaceLabels := workspace.Labels
	stepsConfigurator.CreateOrUpdateWorkspaceV1Step("Create a workspace", suite.Client.WorkspaceV1, workspace,
		steps.StepResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:         expectWorkspaceLabels,
			Metadata:       expectWorkspaceMeta,
			ResourceStates: suites.CreatedResourceExpectedStates,
		},
	)

	// Block storage
	blocks := suite.params.BlockStorages

	// Create the block storages
	for _, block := range blocks {
		expectedBlockMeta := block.Metadata
		expectedBlockSpec := &block.Spec
		stepsConfigurator.CreateOrUpdateBlockStorageV1Step("Create a block storage", suite.Client.StorageV1, &block,
			steps.StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
				Metadata:       expectedBlockMeta,
				Spec:           expectedBlockSpec,
				ResourceStates: suites.CreatedResourceExpectedStates,
			},
		)
	}

	wpath := secapi.WorkspacePath{
		Tenant:    secapi.TenantID(workspace.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(workspace.Metadata.Name),
	}

	// List block storages
	stepsConfigurator.ListBlockStorageV1Step("List block storage", suite.Client.StorageV1, wpath, nil)

	// List block storages with limit
	stepsConfigurator.ListBlockStorageV1Step("List block storage with limit", suite.Client.StorageV1, wpath,
		secapi.NewListOptions().WithLimit(1))

	// List block storages with label
	stepsConfigurator.ListBlockStorageV1Step("List block storage with label", suite.Client.StorageV1, wpath,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List block storages with limit and label
	stepsConfigurator.ListBlockStorageV1Step("List block storage with limit and label", suite.Client.StorageV1, wpath,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Image
	images := suite.params.Images

	// Create images
	for _, image := range images {
		expectedImageMeta := image.Metadata
		expectedImageSpec := &image.Spec
		stepsConfigurator.CreateOrUpdateImageV1Step("Create an image", suite.Client.StorageV1, &image,
			steps.StepResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
				Metadata:       expectedImageMeta,
				Spec:           expectedImageSpec,
				ResourceStates: suites.CreatedResourceExpectedStates,
			},
		)
	}

	tpath := secapi.TenantPath{
		Tenant: secapi.TenantID(workspace.Metadata.Tenant),
	}

	// List images
	stepsConfigurator.ListImageV1Step("List image", suite.Client.StorageV1, tpath, nil)

	// List images with limit
	stepsConfigurator.ListImageV1Step("List images", suite.Client.StorageV1, tpath,
		secapi.NewListOptions().WithLimit(1))

	// List images with label
	stepsConfigurator.ListImageV1Step("List images", suite.Client.StorageV1, tpath,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List images with limit and label
	stepsConfigurator.ListImageV1Step("List images", suite.Client.StorageV1, tpath,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Skus

	// List Skus
	stepsConfigurator.ListSkuV1Step("List skus", suite.Client.StorageV1, tpath, nil)

	// List Skus with limit
	stepsConfigurator.ListSkuV1Step("List skus", suite.Client.StorageV1, tpath,
		secapi.NewListOptions().WithLimit(1))

	// Delete all images
	for _, image := range images {
		stepsConfigurator.DeleteImageV1Step("Delete image", suite.Client.StorageV1, &image)

		// Get the deleted image
		imageTRef := secapi.TenantReference{
			Tenant: secapi.TenantID(workspace.Metadata.Tenant),
			Name:   image.Metadata.Name,
		}
		stepsConfigurator.WatchImageUntilDeletedV1Step("Watch the image deletion", suite.Client.StorageV1, imageTRef)
	}

	// Delete all block storages
	for _, block := range blocks {
		stepsConfigurator.DeleteBlockStorageV1Step("Delete block storage 1", suite.Client.StorageV1, &block)

		// Get the deleted block storage
		blockWRef := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(block.Metadata.Tenant),
			Workspace: secapi.WorkspaceID(block.Metadata.Workspace),
			Name:      block.Metadata.Name,
		}
		stepsConfigurator.WatchBlockStorageUntilDeletedV1Step("Watch the block storage deletion", suite.Client.StorageV1, blockWRef)
	}

	// Delete the workspace
	workspaceTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(workspace.Metadata.Tenant),
		Name:   workspace.Metadata.Name,
	}
	stepsConfigurator.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)
	stepsConfigurator.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", suite.Client.WorkspaceV1, workspaceTRef)

	suite.FinishScenario()
}

func (suite *ProviderQueriesV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
