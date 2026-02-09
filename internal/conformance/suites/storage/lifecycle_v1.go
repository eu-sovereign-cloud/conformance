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
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type StorageLifeCycleV1TestSuite struct {
	suites.RegionalTestSuite

	StorageSkus []string

	params *params.StorageLifeCycleV1Params
}

func CreateLifeCycleV1TestSuite(regionalTestSuite suites.RegionalTestSuite, storageSkus []string) *StorageLifeCycleV1TestSuite {
	suite := &StorageLifeCycleV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		StorageSkus:       storageSkus,
	}
	suite.ScenarioName = constants.StorageV1LifeCycleSuiteName
	return suite
}

func (suite *StorageLifeCycleV1TestSuite) BeforeAll(t provider.T) {
	var err error

	// Select sku
	storageSkuName := suite.StorageSkus[rand.Intn(len(suite.StorageSkus))]

	// Generate scenario data
	workspaceName := generators.GenerateWorkspaceName()

	storageSkuRefObj, err := generators.GenerateSkuRefObject(storageSkuName)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	blockStorageName := generators.GenerateBlockStorageName()
	blockStorageRefObj, err := generators.GenerateBlockStorageRefObject(blockStorageName)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	imageName := generators.GenerateImageName()

	initialStorageSize := generators.GenerateBlockStorageSize()
	updatedStorageSize := generators.GenerateBlockStorageSize()

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

	blockStorageInitial, err := builders.NewBlockStorageBuilder().
		Name(blockStorageName).
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
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
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
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
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
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
		Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Spec(&schema.ImageSpec{
			BlockStorageRef: *blockStorageRefObj,
			CpuArchitecture: schema.ImageSpecCpuArchitectureArm64,
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Image: %v", err)
	}

	params := &params.StorageLifeCycleV1Params{
		Workspace:           workspace,
		BlockStorageInitial: blockStorageInitial,
		BlockStorageUpdated: blockStorageUpdated,
		ImageInitial:        imageInitial,
		ImageUpdated:        imageUpdated,
	}

	suite.params = params

	err = suites.SetupMockIfEnabled(suite.TestSuite, mockstorage.ConfigureLifecycleScenarioV1, params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *StorageLifeCycleV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, constants.StorageProviderV1,
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindBlockStorage),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindImage),
	)

	// Workspace
	workspace := suite.params.Workspace
	expectWorkspaceMeta := workspace.Metadata
	expectWorkspaceLabels := workspace.Labels
	workspaceTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(workspace.Metadata.Tenant),
		Name:   workspace.Metadata.Name,
	}

	t.WithNewStep("Workspace", func(wsCtx provider.StepCtx) {
		wsSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, wsCtx)

		wsSteps.CreateOrUpdateWorkspaceV1Step("Create", suite.Client.WorkspaceV1, workspace,
			steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
				Labels:        expectWorkspaceLabels,
				Metadata:      expectWorkspaceMeta,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		wsSteps.GetWorkspaceV1Step("Get", suite.Client.WorkspaceV1, workspaceTRef,
			steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
				Labels:        expectWorkspaceLabels,
				Metadata:      expectWorkspaceMeta,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	// Block storage
	block := suite.params.BlockStorageInitial
	expectedBlockMeta := block.Metadata
	expectedBlockSpec := &block.Spec
	blockWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(block.Metadata.Tenant),
		Workspace: secapi.WorkspaceID(block.Metadata.Workspace),
		Name:      block.Metadata.Name,
	}

	t.WithNewStep("BlockStorage", func(bsCtx provider.StepCtx) {
		bsSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, bsCtx)

		bsSteps.CreateOrUpdateBlockStorageV1Step("Create", suite.Client.StorageV1, block,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
				Metadata:      expectedBlockMeta,
				Spec:          expectedBlockSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		block = bsSteps.GetBlockStorageV1Step("Get", suite.Client.StorageV1, blockWRef,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
				Metadata:      expectedBlockMeta,
				Spec:          expectedBlockSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)

		block.Spec = suite.params.BlockStorageUpdated.Spec
		expectedBlockSpec.SizeGB = block.Spec.SizeGB
		bsSteps.CreateOrUpdateBlockStorageV1Step("Update", suite.Client.StorageV1, block,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
				Metadata:      expectedBlockMeta,
				Spec:          expectedBlockSpec,
				ResourceState: schema.ResourceStateUpdating,
			},
		)

		block = bsSteps.GetBlockStorageV1Step("GetUpdated", suite.Client.StorageV1, blockWRef,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
				Metadata:      expectedBlockMeta,
				Spec:          expectedBlockSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	// Image
	image := suite.params.ImageInitial
	expectedImageMeta := image.Metadata
	expectedImageSpec := &image.Spec
	imageTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(image.Metadata.Tenant),
		Name:   image.Metadata.Name,
	}

	t.WithNewStep("Image", func(imgCtx provider.StepCtx) {
		imgSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, imgCtx)

		imgSteps.CreateOrUpdateImageV1Step("Create", suite.Client.StorageV1, image,
			steps.ResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
				Metadata:      expectedImageMeta,
				Spec:          expectedImageSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		image = imgSteps.GetImageV1Step("Get", suite.Client.StorageV1, imageTRef,
			steps.ResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
				Metadata:      expectedImageMeta,
				Spec:          expectedImageSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)

		image.Spec = suite.params.ImageUpdated.Spec
		expectedImageSpec.CpuArchitecture = image.Spec.CpuArchitecture
		imgSteps.CreateOrUpdateImageV1Step("Update", suite.Client.StorageV1, image,
			steps.ResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
				Metadata:      expectedImageMeta,
				Spec:          expectedImageSpec,
				ResourceState: schema.ResourceStateUpdating,
			},
		)

		image = imgSteps.GetImageV1Step("GetUpdated", suite.Client.StorageV1, imageTRef,
			steps.ResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
				Metadata:      expectedImageMeta,
				Spec:          expectedImageSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	// Resources deletion
	t.WithNewStep("Delete", func(delCtx provider.StepCtx) {
		delSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, delCtx)
		delSteps.DeleteImageV1Step("Delete the image", suite.Client.StorageV1, image)
		delSteps.GetImageWithErrorV1Step("Get the deleted image", suite.Client.StorageV1, imageTRef, secapi.ErrResourceNotFound)

		delSteps.DeleteBlockStorageV1Step("Delete the block storage", suite.Client.StorageV1, block)
		delSteps.GetBlockStorageWithErrorV1Step("Get the deleted block storage", suite.Client.StorageV1, blockWRef, secapi.ErrResourceNotFound)

		delSteps.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)
		delSteps.GetWorkspaceWithErrorV1Step("Get the deleted workspace", suite.Client.WorkspaceV1, workspaceTRef, secapi.ErrResourceNotFound)
	})

	suite.FinishScenario()
}

func (suite *StorageLifeCycleV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
