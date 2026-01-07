package storage

import (
	"math/rand"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/storage"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	labelBuilder "github.com/eu-sovereign-cloud/go-sdk/secapi/builders"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type StorageV1ListTestSuite struct {
	suites.RegionalTestSuite

	StorageSkus []string
}

func (suite *StorageV1ListTestSuite) TestListScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, constants.StorageProviderV1,
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindBlockStorage),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindImage),
	)

	// Select sku
	storageSkuName := suite.StorageSkus[rand.Intn(len(suite.StorageSkus))]

	// Generate scenario data
	workspaceName := generators.GenerateWorkspaceName()
	storageSkuRef := generators.GenerateSkuRef(storageSkuName)
	storageSkuRefObj, err := secapi.BuildReferenceFromURN(storageSkuRef)
	if err != nil {
		t.Fatal(err)
	}

	blockStorageName1 := generators.GenerateBlockStorageName()
	blockStorageResource1 := generators.GenerateBlockStorageResource(suite.Tenant, workspaceName, blockStorageName1)
	blockStorageName2 := generators.GenerateBlockStorageName()
	blockStorageResource2 := generators.GenerateBlockStorageResource(suite.Tenant, workspaceName, blockStorageName2)
	blockStorageName3 := generators.GenerateBlockStorageName()
	blockStorageResource3 := generators.GenerateBlockStorageResource(suite.Tenant, workspaceName, blockStorageName3)

	blockStorageRef := generators.GenerateBlockStorageRef(blockStorageName1)
	blockStorageRefObj, err := secapi.BuildReferenceFromURN(blockStorageRef)
	if err != nil {
		t.Fatal(err)
	}

	imageName1 := generators.GenerateImageName()
	imageResource1 := generators.GenerateImageResource(suite.Tenant, imageName1)
	imageName2 := generators.GenerateImageName()
	imageResource2 := generators.GenerateImageResource(suite.Tenant, imageName2)
	imageName3 := generators.GenerateImageName()
	imageResource3 := generators.GenerateImageResource(suite.Tenant, imageName3)
	initialStorageSize := generators.GenerateBlockStorageSize()

	// Setup mock, if configured to use
	if suite.MockEnabled {
		mockParams := &mock.StorageListParamsV1{
			BaseParams: &mock.BaseParams{
				MockURL:   *suite.MockServerURL,
				AuthToken: suite.AuthToken,
				Tenant:    suite.Tenant,
				Region:    suite.Region,
			},
			Workspace: &mock.ResourceParams[schema.WorkspaceSpec]{
				Name: workspaceName,
				InitialLabels: schema.Labels{
					constants.EnvLabel: constants.EnvConformanceLabel,
				},
			},
			BlockStorages: []mock.ResourceParams[schema.BlockStorageSpec]{
				{
					Name: blockStorageName1,
					InitialLabels: map[string]string{
						constants.EnvLabel: constants.EnvConformanceLabel,
					},
					InitialSpec: &schema.BlockStorageSpec{
						SkuRef: *storageSkuRefObj,
						SizeGB: initialStorageSize,
					},
				},
				{
					Name: blockStorageName2,
					InitialLabels: map[string]string{
						constants.EnvLabel: constants.EnvConformanceLabel,
					},
					InitialSpec: &schema.BlockStorageSpec{
						SkuRef: *storageSkuRefObj,
						SizeGB: initialStorageSize,
					},
				},
				{
					Name: blockStorageName3,
					InitialLabels: map[string]string{
						constants.EnvLabel: constants.EnvConformanceLabel,
					},
					InitialSpec: &schema.BlockStorageSpec{
						SkuRef: *storageSkuRefObj,
						SizeGB: initialStorageSize,
					},
				},
			},
			Images: []mock.ResourceParams[schema.ImageSpec]{
				{
					Name: imageName1,
					InitialLabels: map[string]string{
						constants.EnvLabel: constants.EnvConformanceLabel,
					},
					InitialSpec: &schema.ImageSpec{
						BlockStorageRef: *blockStorageRefObj,
						CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
					},
				},
				{
					Name: imageName2,
					InitialLabels: map[string]string{
						constants.EnvLabel: constants.EnvConformanceLabel,
					},
					InitialSpec: &schema.ImageSpec{
						BlockStorageRef: *blockStorageRefObj,
						CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
					},
				},
				{
					Name: imageName3,
					InitialLabels: map[string]string{
						constants.EnvLabel: constants.EnvConformanceLabel,
					},
					InitialSpec: &schema.ImageSpec{
						BlockStorageRef: *blockStorageRefObj,
						CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
					},
				},
			},
		}
		wm, err := storage.ConfigureListScenarioV1(suite.ScenarioName, mockParams)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}
		suite.MockClient = wm
	}

	stepsBuilder := steps.NewStepsConfigurator(&suite.TestSuite, t)

	// Workspace

	// Create a workspace
	workspace := &schema.Workspace{
		Labels: schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
		},
		Metadata: &schema.RegionalResourceMetadata{
			Tenant: suite.Tenant,
			Name:   workspaceName,
		},
	}
	expectWorkspaceMeta, err := builders.NewWorkspaceMetadataBuilder().
		Name(workspaceName).
		Provider(constants.WorkspaceProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Metadata: %v", err)
	}
	expectWorkspaceLabels := schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}

	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:        expectWorkspaceLabels,
			Metadata:      expectWorkspaceMeta,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Block storage

	// Create a block storage
	blocks := &[]schema.BlockStorage{
		{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.Tenant,
				Workspace: workspaceName,
				Name:      blockStorageName1,
				Resource:  blockStorageResource1,
			},
			Spec: schema.BlockStorageSpec{
				SizeGB: initialStorageSize,
				SkuRef: *storageSkuRefObj,
			},
		},
		{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.Tenant,
				Workspace: workspaceName,
				Name:      blockStorageName2,
				Resource:  blockStorageResource2,
			},
			Spec: schema.BlockStorageSpec{
				SizeGB: initialStorageSize,
				SkuRef: *storageSkuRefObj,
			},
		},
		{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.Tenant,
				Workspace: workspaceName,
				Name:      blockStorageName3,
				Resource:  blockStorageResource3,
			},
			Spec: schema.BlockStorageSpec{
				SizeGB: initialStorageSize,
				SkuRef: *storageSkuRefObj,
			},
		},
	}

	// Create the block storages
	for _, block := range *blocks {
		expectedBlockMeta, err := builders.NewBlockStorageMetadataBuilder().
			Name(block.Metadata.Name).
			Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
			Tenant(suite.Tenant).Workspace(workspaceName).Region(suite.Region).
			Build()
		if err != nil {
			t.Fatalf("Failed to build Metadata: %v", err)
		}
		expectedBlockSpec := &schema.BlockStorageSpec{
			SizeGB: initialStorageSize,
			SkuRef: *storageSkuRefObj,
		}
		stepsBuilder.CreateOrUpdateBlockStorageV1Step("Create a block storage", suite.Client.StorageV1, &block,
			steps.ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
				Metadata:      expectedBlockMeta,
				Spec:          expectedBlockSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)
	}

	wref := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.Tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
	}

	// List block storages
	stepsBuilder.GetListBlockStorageV1Step("GetList block storage", t, suite.Client.StorageV1, wref, nil)

	// List block storages with limit
	stepsBuilder.GetListBlockStorageV1Step("Get List block storage with limit", t, suite.Client.StorageV1, wref,
		secapi.NewListOptions().WithLimit(1))

	// List block storages with label
	stepsBuilder.GetListBlockStorageV1Step("Get list of block storage with label", t, suite.Client.StorageV1, wref,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvDevelopmentLabel)))

	// List block storages with limit and label
	stepsBuilder.GetListBlockStorageV1Step("Get list of block storage with limit and label", t, suite.Client.StorageV1, wref,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvDevelopmentLabel)))

	// Image

	// Create an image
	images := &[]schema.Image{
		{
			Metadata: &schema.RegionalResourceMetadata{
				Tenant:   suite.Tenant,
				Name:     imageName1,
				Resource: imageResource1,
			},
			Spec: schema.ImageSpec{
				BlockStorageRef: *blockStorageRefObj,
				CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
			},
		},
		{
			Metadata: &schema.RegionalResourceMetadata{
				Tenant:   suite.Tenant,
				Name:     imageName2,
				Resource: imageResource2,
			},
			Spec: schema.ImageSpec{
				BlockStorageRef: *blockStorageRefObj,
				CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
			},
		},
		{
			Metadata: &schema.RegionalResourceMetadata{
				Tenant:   suite.Tenant,
				Name:     imageName3,
				Resource: imageResource3,
			},
			Spec: schema.ImageSpec{
				BlockStorageRef: *blockStorageRefObj,
				CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
			},
		},
	}

	for _, image := range *images {

		expectedImageMeta, err := builders.NewImageMetadataBuilder().
			Name(image.Metadata.Name).
			Provider(constants.StorageProviderV1).ApiVersion(constants.ApiVersion1).
			Tenant(suite.Tenant).Region(suite.Region).
			Build()
		if err != nil {
			t.Fatalf("Failed to build Metadata: %v", err)
		}
		expectedImageSpec := &schema.ImageSpec{
			BlockStorageRef: *blockStorageRefObj,
			CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
		}
		stepsBuilder.CreateOrUpdateImageV1Step("Create an image", suite.Client.StorageV1, &image,
			steps.ResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
				Metadata:      expectedImageMeta,
				Spec:          expectedImageSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)
	}
	tref := secapi.TenantReference{
		Name:   suite.Tenant,
		Tenant: secapi.TenantID(suite.Tenant),
	}
	// List image
	stepsBuilder.GetListImageV1Step("List image", suite.Client.StorageV1, tref, nil)

	// List image with limit
	stepsBuilder.GetListImageV1Step("Get list of images", suite.Client.StorageV1, tref,
		secapi.NewListOptions().WithLimit(1))

	// List image with Label
	stepsBuilder.GetListImageV1Step("Get list of images", suite.Client.StorageV1, tref,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List image with Limit and label
	stepsBuilder.GetListImageV1Step("Get list of images", suite.Client.StorageV1, tref,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Skus
	// List Skus
	stepsBuilder.GetListSkuV1Step("List skus", t, suite.Client.StorageV1, tref, nil)

	// List Skus with limit
	stepsBuilder.GetListSkuV1Step("Get list of skus", t, suite.Client.StorageV1, tref,
		secapi.NewListOptions().WithLimit(1))

	// Delete all images
	for _, image := range *images {
		imageTRef := &secapi.TenantReference{
			Tenant: secapi.TenantID(suite.Tenant),
			Name:   image.Metadata.Name,
		}

		stepsBuilder.DeleteImageV1Step("Delete image", suite.Client.StorageV1, &image)
		stepsBuilder.GetImageWithErrorV1Step("Get deleted image ", suite.Client.StorageV1, *imageTRef, secapi.ErrResourceNotFound)
	}
	// Delete all block storages
	for _, block := range *blocks {
		blockWRef := &secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.Tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      block.Metadata.Name,
		}
		stepsBuilder.DeleteBlockStorageV1Step("Delete block storage 1", suite.Client.StorageV1, &block)
		stepsBuilder.GetBlockStorageWithErrorV1Step("Get deleted block storage 1", suite.Client.StorageV1, *blockWRef, secapi.ErrResourceNotFound)
	}

	// Delete the workspace
	workspaceTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.Tenant),
		Name:   workspaceName,
	}
	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)
	stepsBuilder.GetWorkspaceWithErrorV1Step("Get the deleted workspace", suite.Client.WorkspaceV1, *workspaceTRef, secapi.ErrResourceNotFound)

	suite.FinishScenario()
}

func (suite *StorageV1LifeCycleTestSuite) AfterEach(t provider.T) {
	suite.ResetAllScenarios()
}
