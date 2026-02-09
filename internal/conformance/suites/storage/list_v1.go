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

type StorageListV1TestSuite struct {
	suites.RegionalTestSuite

	StorageSkus []string

	params *params.StorageListV1Params
}

func CreateListV1TestSuite(regionalTestSuite suites.RegionalTestSuite, storageSkus []string) *StorageListV1TestSuite {
	suite := &StorageListV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		StorageSkus:       storageSkus,
	}
	suite.ScenarioName = constants.StorageV1ListSuiteName
	return suite
}

func (suite *StorageListV1TestSuite) BeforeAll(t provider.T) {
	var err error

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

	params := &params.StorageListV1Params{
		Workspace:     workspace,
		BlockStorages: blockStorages,
		Images:        images,
	}

	suite.params = params
	err = suites.SetupMockIfEnabled(suite.TestSuite, mockStorage.ConfigureListScenarioV1, params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *StorageListV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, constants.StorageProviderV1,
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindBlockStorage),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindImage),
	)

	workspace := suite.params.Workspace
	blocks := suite.params.BlockStorages
	images := suite.params.Images

	t.WithNewStep("Workspace", func(wsCtx provider.StepCtx) {
		wsSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, wsCtx)

		expectWorkspaceMeta := workspace.Metadata
		expectWorkspaceLabels := workspace.Labels
		wsSteps.CreateOrUpdateWorkspaceV1Step("Create", suite.Client.WorkspaceV1, workspace,
			steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
				Labels:        expectWorkspaceLabels,
				Metadata:      expectWorkspaceMeta,
				ResourceState: schema.ResourceStateCreating,
			},
		)
	})

	t.WithNewStep("BlockStorage", func(bsCtx provider.StepCtx) {
		bsSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, bsCtx)

		for _, block := range blocks {
			b := block
			expectedBlockMeta := b.Metadata
			expectedBlockSpec := &b.Spec
			bsSteps.CreateOrUpdateBlockStorageV1Step("Create", suite.Client.StorageV1, &b,
				steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
					Metadata:      expectedBlockMeta,
					Spec:          expectedBlockSpec,
					ResourceState: schema.ResourceStateCreating,
				},
			)
		}

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(workspace.Metadata.Tenant),
			Workspace: secapi.WorkspaceID(workspace.Metadata.Name),
		}

		bsSteps.GetListBlockStorageV1Step("ListAll", suite.Client.StorageV1, wref, nil)
		bsSteps.GetListBlockStorageV1Step("ListWithLimit", suite.Client.StorageV1, wref,
			secapi.NewListOptions().WithLimit(1))
		bsSteps.GetListBlockStorageV1Step("ListWithLabel", suite.Client.StorageV1, wref,
			secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
				Equals(constants.EnvLabel, constants.EnvDevelopmentLabel)))
		bsSteps.GetListBlockStorageV1Step("ListWithLabelAndLimit", suite.Client.StorageV1, wref,
			secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
				Equals(constants.EnvLabel, constants.EnvDevelopmentLabel)))
	})

	t.WithNewStep("Image", func(imgCtx provider.StepCtx) {
		imgSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, imgCtx)

		for _, image := range images {
			im := image
			expectedImageMeta := im.Metadata
			expectedImageSpec := &im.Spec
			imgSteps.CreateOrUpdateImageV1Step("Create", suite.Client.StorageV1, &im,
				steps.ResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
					Metadata:      expectedImageMeta,
					Spec:          expectedImageSpec,
					ResourceState: schema.ResourceStateCreating,
				},
			)
		}

		tref := secapi.TenantReference{
			Name:   workspace.Metadata.Tenant,
			Tenant: secapi.TenantID(workspace.Metadata.Tenant),
		}
		imgSteps.GetListImageV1Step("ListAll", suite.Client.StorageV1, tref, nil)
		imgSteps.GetListImageV1Step("ListWithLimit", suite.Client.StorageV1, tref,
			secapi.NewListOptions().WithLimit(1))
		imgSteps.GetListImageV1Step("ListWithLabel", suite.Client.StorageV1, tref,
			secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
				Equals(constants.EnvLabel, constants.EnvConformanceLabel)))
		imgSteps.GetListImageV1Step("ListWithLabelAndLimit", suite.Client.StorageV1, tref,
			secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
				Equals(constants.EnvLabel, constants.EnvConformanceLabel)))
	})

	t.WithNewStep("Skus", func(skuCtx provider.StepCtx) {
		skuSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, skuCtx)

		tref := secapi.TenantReference{
			Name:   workspace.Metadata.Tenant,
			Tenant: secapi.TenantID(workspace.Metadata.Tenant),
		}
		skuSteps.GetListSkuV1Step("ListAll", suite.Client.StorageV1, tref, nil)
		skuSteps.GetListSkuV1Step("ListWithLimit", suite.Client.StorageV1, tref,
			secapi.NewListOptions().WithLimit(1))
	})

	t.WithNewStep("Deletes", func(delCtx provider.StepCtx) {
		delSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, delCtx)

		for _, image := range images {
			im := image
			delSteps.DeleteImageV1Step("Image", suite.Client.StorageV1, &im)

			imageTRef := secapi.TenantReference{
				Tenant: secapi.TenantID(workspace.Metadata.Tenant),
				Name:   im.Metadata.Name,
			}
			delSteps.GetImageWithErrorV1Step("GetDeletedImage", suite.Client.StorageV1, imageTRef, secapi.ErrResourceNotFound)
		}

		for _, block := range blocks {
			b := block
			delSteps.DeleteBlockStorageV1Step("BlockStorage", suite.Client.StorageV1, &b)

			blockWRef := secapi.WorkspaceReference{
				Tenant:    secapi.TenantID(b.Metadata.Tenant),
				Workspace: secapi.WorkspaceID(b.Metadata.Workspace),
				Name:      b.Metadata.Name,
			}
			delSteps.GetBlockStorageWithErrorV1Step("GetDeletedBlockStorage", suite.Client.StorageV1, blockWRef, secapi.ErrResourceNotFound)
		}

		workspaceTRef := secapi.TenantReference{
			Tenant: secapi.TenantID(workspace.Metadata.Tenant),
			Name:   workspace.Metadata.Name,
		}
		delSteps.DeleteWorkspaceV1Step("Workspace", suite.Client.WorkspaceV1, workspace)
		delSteps.GetWorkspaceWithErrorV1Step("GetDeletedWorkspace", suite.Client.WorkspaceV1, workspaceTRef, secapi.ErrResourceNotFound)
	})

	suite.FinishScenario()
}

func (suite *StorageListV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
