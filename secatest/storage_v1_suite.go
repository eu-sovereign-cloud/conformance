package secatest

import (
	"math/rand"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	labelBuilder "github.com/eu-sovereign-cloud/go-sdk/secapi/builders"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type StorageV1LifeCycleTestSuite struct {
	regionalTestSuite

	storageSkus []string
}
type StorageV1ListTestSuite struct {
	regionalTestSuite

	storageSkus []string
}

func (suite *StorageV1LifeCycleTestSuite) TestLifeCycleScenario(t provider.T) {
	suite.startScenario(t)
	configureTags(t, storageProviderV1,
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindBlockStorage),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindImage),
	)

	var err error

	// Select sku
	storageSkuName := suite.storageSkus[rand.Intn(len(suite.storageSkus))]

	// Generate scenario data
	workspaceName := generators.GenerateWorkspaceName()

	storageSkuRef := generators.GenerateSkuRef(storageSkuName)
	storageSkuRefObj, err := secapi.BuildReferenceFromURN(storageSkuRef)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	blockStorageName := generators.GenerateBlockStorageName()
	blockStorageRef := generators.GenerateBlockStorageRef(blockStorageName)
	blockStorageRefObj, err := secapi.BuildReferenceFromURN(blockStorageRef)
	if err != nil {
		t.Fatalf("Failed to build URN: %v", err)
	}

	imageName := generators.GenerateImageName()

	initialStorageSize := generators.GenerateBlockStorageSize()
	updatedStorageSize := generators.GenerateBlockStorageSize()

	// Setup mock, if configured to use
	if suite.mockEnabled {
		mockParams := &mock.StorageLifeCycleParamsV1{
			BaseParams: &mock.BaseParams{
				MockURL:   *suite.mockServerURL,
				AuthToken: suite.authToken,
				Tenant:    suite.tenant,
				Region:    suite.region,
			},
			Workspace: &mock.ResourceParams[schema.WorkspaceSpec]{
				Name: workspaceName,
				InitialLabels: schema.Labels{
					envLabel: envDevelopmentLabel,
				},
			},
			BlockStorage: &mock.ResourceParams[schema.BlockStorageSpec]{
				Name: blockStorageName,
				InitialSpec: &schema.BlockStorageSpec{
					SkuRef: *storageSkuRefObj,
					SizeGB: initialStorageSize,
				},
				UpdatedSpec: &schema.BlockStorageSpec{
					SkuRef: *storageSkuRefObj,
					SizeGB: updatedStorageSize,
				},
			},
			Image: &mock.ResourceParams[schema.ImageSpec]{
				Name: imageName,
				InitialSpec: &schema.ImageSpec{
					BlockStorageRef: *blockStorageRefObj,
					CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
				},
				UpdatedSpec: &schema.ImageSpec{
					BlockStorageRef: *blockStorageRefObj,
					CpuArchitecture: schema.ImageSpecCpuArchitectureArm64,
				},
			},
		}
		wm, err := mock.ConfigureStorageLifecycleScenarioV1(suite.scenarioName, mockParams)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}
		suite.mockClient = wm
	}

	// Workspace

	// Create a workspace
	workspace := &schema.Workspace{
		Labels: schema.Labels{
			envLabel: envDevelopmentLabel,
		},
		Metadata: &schema.RegionalResourceMetadata{
			Tenant: suite.tenant,
			Name:   workspaceName,
		},
	}
	expectWorkspaceMeta, err := builders.NewWorkspaceMetadataBuilder().
		Name(workspaceName).
		Provider(workspaceProviderV1).ApiVersion(apiVersion1).
		Tenant(suite.tenant).Region(suite.region).
		Build()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectWorkspaceLabels := schema.Labels{envLabel: envDevelopmentLabel}

	suite.createOrUpdateWorkspaceV1Step("Create a workspace", t, suite.client.WorkspaceV1, workspace,
		responseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			labels:        expectWorkspaceLabels,
			metadata:      expectWorkspaceMeta,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created Workspace
	workspaceTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
		Name:   workspaceName,
	}
	suite.getWorkspaceV1Step("Get the created workspace", t, suite.client.WorkspaceV1, *workspaceTRef,
		responseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			labels:        expectWorkspaceLabels,
			metadata:      expectWorkspaceMeta,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Block storage

	// Create a block storage
	block := &schema.BlockStorage{
		Metadata: &schema.RegionalWorkspaceResourceMetadata{
			Tenant:    suite.tenant,
			Workspace: workspaceName,
			Name:      blockStorageName,
		},
		Spec: schema.BlockStorageSpec{
			SizeGB: initialStorageSize,
			SkuRef: *storageSkuRefObj,
		},
	}
	expectedBlockMeta, err := builders.NewBlockStorageMetadataBuilder().
		Name(blockStorageName).
		Provider(storageProviderV1).ApiVersion(apiVersion1).
		Tenant(suite.tenant).Workspace(workspaceName).Region(suite.region).
		Build()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectedBlockSpec := &schema.BlockStorageSpec{
		SizeGB: initialStorageSize,
		SkuRef: *storageSkuRefObj,
	}
	suite.createOrUpdateBlockStorageV1Step("Create a block storage", t, suite.client.StorageV1, block,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			metadata:      expectedBlockMeta,
			spec:          expectedBlockSpec,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created block storage
	blockWRef := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      blockStorageName,
	}
	block = suite.getBlockStorageV1Step("Get the created block storage", t, suite.client.StorageV1, *blockWRef,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			metadata:      expectedBlockMeta,
			spec:          expectedBlockSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Update the block storage
	block.Spec.SizeGB = updatedStorageSize
	expectedBlockSpec.SizeGB = block.Spec.SizeGB
	suite.createOrUpdateBlockStorageV1Step("Update the block storage", t, suite.client.StorageV1, block,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			metadata:      expectedBlockMeta,
			spec:          expectedBlockSpec,
			resourceState: schema.ResourceStateUpdating,
		},
	)

	// Get the updated block storage
	block = suite.getBlockStorageV1Step("Get the updated block storage", t, suite.client.StorageV1, *blockWRef,
		responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			metadata:      expectedBlockMeta,
			spec:          expectedBlockSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Image

	// Create an image
	image := &schema.Image{
		Metadata: &schema.RegionalResourceMetadata{
			Tenant: suite.tenant,
			Name:   imageName,
		},
		Spec: schema.ImageSpec{
			BlockStorageRef: *blockStorageRefObj,
			CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
		},
	}
	expectedImageMeta, err := builders.NewImageMetadataBuilder().
		Name(imageName).
		Provider(storageProviderV1).ApiVersion(apiVersion1).
		Tenant(suite.tenant).Region(suite.region).
		Build()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectedImageSpec := &schema.ImageSpec{
		BlockStorageRef: *blockStorageRefObj,
		CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
	}
	suite.createOrUpdateImageV1Step("Create an image", t, suite.client.StorageV1, image,
		responseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
			metadata:      expectedImageMeta,
			spec:          expectedImageSpec,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created image
	imageTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
		Name:   imageName,
	}
	image = suite.getImageV1Step("Get the created image", t, suite.client.StorageV1, *imageTRef,
		responseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
			metadata:      expectedImageMeta,
			spec:          expectedImageSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Update the image
	image.Spec.CpuArchitecture = schema.ImageSpecCpuArchitectureArm64
	expectedImageSpec.CpuArchitecture = image.Spec.CpuArchitecture
	suite.createOrUpdateImageV1Step("Update the image", t, suite.client.StorageV1, image,
		responseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
			metadata:      expectedImageMeta,
			spec:          expectedImageSpec,
			resourceState: schema.ResourceStateUpdating,
		},
	)

	// Get the updated image
	image = suite.getImageV1Step("Get the updated image", t, suite.client.StorageV1, *imageTRef,
		responseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
			metadata:      expectedImageMeta,
			spec:          expectedImageSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Resources deletion

	suite.deleteImageV1Step("Delete the image", t, suite.client.StorageV1, image)
	suite.getImageWithErrorV1Step("Get the deleted image", t, suite.client.StorageV1, *imageTRef, secapi.ErrResourceNotFound)

	suite.deleteBlockStorageV1Step("Delete the block storage", t, suite.client.StorageV1, block)
	suite.getBlockStorageWithErrorV1Step("Get the deleted block storage", t, suite.client.StorageV1, *blockWRef, secapi.ErrResourceNotFound)

	suite.deleteWorkspaceV1Step("Delete the workspace", t, suite.client.WorkspaceV1, workspace)
	suite.getWorkspaceWithErrorV1Step("Get the deleted workspace", t, suite.client.WorkspaceV1, *workspaceTRef, secapi.ErrResourceNotFound)

	suite.finishScenario()
}

func (suite *StorageV1ListTestSuite) TestListScenario(t provider.T) {
	suite.startScenario(t)
	configureTags(t, storageProviderV1,
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindBlockStorage),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindImage),
	)

	// Select sku
	storageSkuName := suite.storageSkus[rand.Intn(len(suite.storageSkus))]

	// Generate scenario data
	workspaceName := generators.GenerateWorkspaceName()
	storageSkuRef := generators.GenerateSkuRef(storageSkuName)
	storageSkuRefObj, err := secapi.BuildReferenceFromURN(storageSkuRef)
	if err != nil {
		t.Fatal(err)
	}

	blockStorageName1 := generators.GenerateBlockStorageName()
	blockStorageResource1 := generators.GenerateBlockStorageResource(suite.tenant, workspaceName, blockStorageName1)
	blockStorageName2 := generators.GenerateBlockStorageName()
	blockStorageResource2 := generators.GenerateBlockStorageResource(suite.tenant, workspaceName, blockStorageName2)
	blockStorageName3 := generators.GenerateBlockStorageName()
	blockStorageResource3 := generators.GenerateBlockStorageResource(suite.tenant, workspaceName, blockStorageName3)

	blockStorageRef := generators.GenerateBlockStorageRef(blockStorageName1)
	blockStorageRefObj, err := secapi.BuildReferenceFromURN(blockStorageRef)
	if err != nil {
		t.Fatal(err)
	}

	imageName1 := generators.GenerateImageName()
	imageResource1 := generators.GenerateImageResource(suite.tenant, imageName1)
	imageName2 := generators.GenerateImageName()
	imageResource2 := generators.GenerateImageResource(suite.tenant, imageName2)
	imageName3 := generators.GenerateImageName()
	imageResource3 := generators.GenerateImageResource(suite.tenant, imageName3)
	initialStorageSize := generators.GenerateBlockStorageSize()

	// Setup mock, if configured to use
	if suite.mockEnabled {
		mockParams := &mock.StorageListParamsV1{
			BaseParams: &mock.BaseParams{
				MockURL:   *suite.mockServerURL,
				AuthToken: suite.authToken,
				Tenant:    suite.tenant,
				Region:    suite.region,
			},
			Workspace: &mock.ResourceParams[schema.WorkspaceSpec]{
				Name: workspaceName,
				InitialLabels: schema.Labels{
					generators.EnvLabel: generators.EnvConformanceLabel,
				},
			},
			BlockStorages: []mock.ResourceParams[schema.BlockStorageSpec]{
				{
					Name: blockStorageName1,
					InitialLabels: map[string]string{
						generators.EnvLabel: generators.EnvConformanceLabel,
					},
					InitialSpec: &schema.BlockStorageSpec{
						SkuRef: *storageSkuRefObj,
						SizeGB: initialStorageSize,
					},
				},
				{
					Name: blockStorageName2,
					InitialLabels: map[string]string{
						generators.EnvLabel: generators.EnvConformanceLabel,
					},
					InitialSpec: &schema.BlockStorageSpec{
						SkuRef: *storageSkuRefObj,
						SizeGB: initialStorageSize,
					},
				},
				{
					Name: blockStorageName3,
					InitialLabels: map[string]string{
						generators.EnvLabel: generators.EnvConformanceLabel,
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
						generators.EnvLabel: generators.EnvConformanceLabel,
					},
					InitialSpec: &schema.ImageSpec{
						BlockStorageRef: *blockStorageRefObj,
						CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
					},
				},
				{
					Name: imageName2,
					InitialLabels: map[string]string{
						generators.EnvLabel: generators.EnvConformanceLabel,
					},
					InitialSpec: &schema.ImageSpec{
						BlockStorageRef: *blockStorageRefObj,
						CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
					},
				},
				{
					Name: imageName3,
					InitialLabels: map[string]string{
						generators.EnvLabel: generators.EnvConformanceLabel,
					},
					InitialSpec: &schema.ImageSpec{
						BlockStorageRef: *blockStorageRefObj,
						CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
					},
				},
			},
		}
		wm, err := mock.ConfigureStorageListScenarioV1(suite.scenarioName, mockParams)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}
		suite.mockClient = wm
	}

	// Workspace

	// Create a workspace
	workspace := &schema.Workspace{
		Labels: schema.Labels{
			generators.EnvLabel: generators.EnvConformanceLabel,
		},
		Metadata: &schema.RegionalResourceMetadata{
			Tenant: suite.tenant,
			Name:   workspaceName,
		},
	}
	expectWorkspaceMeta, err := builders.NewWorkspaceMetadataBuilder().
		Name(workspaceName).
		Provider(workspaceProviderV1).ApiVersion(apiVersion1).
		Tenant(suite.tenant).Region(suite.region).
		Build()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectWorkspaceLabels := schema.Labels{generators.EnvLabel: generators.EnvConformanceLabel}

	suite.createOrUpdateWorkspaceV1Step("Create a workspace", t, suite.client.WorkspaceV1, workspace,
		responseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			labels:        expectWorkspaceLabels,
			metadata:      expectWorkspaceMeta,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Block storage

	// Create a block storage
	blocks := &[]schema.BlockStorage{
		{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.tenant,
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
				Tenant:    suite.tenant,
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
				Tenant:    suite.tenant,
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
			Provider(storageProviderV1).ApiVersion(apiVersion1).
			Tenant(suite.tenant).Workspace(workspaceName).Region(suite.region).
			Build()
		if err != nil {
			t.Fatalf("Failed to build metadata: %v", err)
		}
		expectedBlockSpec := &schema.BlockStorageSpec{
			SizeGB: initialStorageSize,
			SkuRef: *storageSkuRefObj,
		}
		suite.createOrUpdateBlockStorageV1Step("Create a block storage", t, suite.client.StorageV1, &block,
			responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
				metadata:      expectedBlockMeta,
				spec:          expectedBlockSpec,
				resourceState: schema.ResourceStateCreating,
			},
		)
	}

	wref := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
	}

	// List block storages
	suite.getListBlockStorageV1Step("GetList block storage", t, suite.client.StorageV1, wref, nil)

	// List block storages with limit
	suite.getListBlockStorageV1Step("Get List block storage with limit", t, suite.client.StorageV1, wref,
		secapi.NewListOptions().WithLimit(1))

	// List block storages with label
	suite.getListBlockStorageV1Step("Get list of block storage with label", t, suite.client.StorageV1, wref,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(generators.EnvLabel, generators.EnvConformanceLabel)))

	// List block storages with limit and label
	suite.getListBlockStorageV1Step("Get list of block storage with limit and label", t, suite.client.StorageV1, wref,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(generators.EnvLabel, generators.EnvConformanceLabel)))

	// Image

	// Create an image
	images := &[]schema.Image{
		{
			Metadata: &schema.RegionalResourceMetadata{
				Tenant:   suite.tenant,
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
				Tenant:   suite.tenant,
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
				Tenant:   suite.tenant,
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
			Provider(storageProviderV1).ApiVersion(apiVersion1).
			Tenant(suite.tenant).Region(suite.region).
			Build()
		if err != nil {
			t.Fatalf("Failed to build metadata: %v", err)
		}
		expectedImageSpec := &schema.ImageSpec{
			BlockStorageRef: *blockStorageRefObj,
			CpuArchitecture: schema.ImageSpecCpuArchitectureAmd64,
		}
		suite.createOrUpdateImageV1Step("Create an image", t, suite.client.StorageV1, &image,
			responseExpects[schema.RegionalResourceMetadata, schema.ImageSpec]{
				metadata:      expectedImageMeta,
				spec:          expectedImageSpec,
				resourceState: schema.ResourceStateCreating,
			},
		)
	}
	tref := secapi.TenantReference{
		Name:   suite.tenant,
		Tenant: secapi.TenantID(suite.tenant),
	}
	// List image
	suite.getListImageV1Step("List image", t, suite.client.StorageV1, tref, nil)

	// List image with limit
	suite.getListImageV1Step("Get list of images", t, suite.client.StorageV1, tref,
		secapi.NewListOptions().WithLimit(1))

	// List image with Label
	suite.getListImageV1Step("Get list of images", t, suite.client.StorageV1, tref,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(generators.EnvLabel, generators.EnvConformanceLabel)))

	// List image with Limit and label
	suite.getListImageV1Step("Get list of images", t, suite.client.StorageV1, tref,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(generators.EnvLabel, generators.EnvConformanceLabel)))

	// Skus
	// List Skus
	suite.getListSkuV1Step("List skus", t, suite.client.StorageV1, tref, nil)

	// List Skus with limit
	suite.getListSkuV1Step("Get list of skus", t, suite.client.StorageV1, tref,
		secapi.NewListOptions().WithLimit(1))

	// Delete all images
	for _, image := range *images {
		imageTRef := &secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   image.Metadata.Name,
		}

		suite.deleteImageV1Step("Delete image", t, suite.client.StorageV1, &image)
		suite.getImageWithErrorV1Step("Get deleted image ", t, suite.client.StorageV1, *imageTRef, secapi.ErrResourceNotFound)
	}
	// Delete all block storages
	for _, block := range *blocks {
		blockWRef := &secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      block.Metadata.Name,
		}
		suite.deleteBlockStorageV1Step("Delete block storage 1", t, suite.client.StorageV1, &block)
		suite.getBlockStorageWithErrorV1Step("Get deleted block storage 1", t, suite.client.StorageV1, *blockWRef, secapi.ErrResourceNotFound)
	}

	// Delete the workspace
	workspaceTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
		Name:   workspaceName,
	}
	suite.deleteWorkspaceV1Step("Delete the workspace", t, suite.client.WorkspaceV1, workspace)
	suite.getWorkspaceWithErrorV1Step("Get the deleted workspace", t, suite.client.WorkspaceV1, *workspaceTRef, secapi.ErrResourceNotFound)

	suite.finishScenario()
}

func (suite *StorageV1LifeCycleTestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}
