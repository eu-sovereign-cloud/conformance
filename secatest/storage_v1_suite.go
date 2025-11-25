package secatest

import (
	"context"
	"log/slog"
	"math/rand"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/conformance/secalib/builders"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/eu-sovereign-cloud/go-sdk/secapi/builders"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type StorageV1TestSuite struct {
	regionalTestSuite

	storageSkus []string
}

func (suite *StorageV1TestSuite) TestSuite(t provider.T) {
	var err error
	slog.Info("Starting " + suite.scenarioName)

	t.Title(suite.scenarioName)
	configureTags(t, secalib.StorageProviderV1,
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindBlockStorage),
		string(schema.RegionalWorkspaceResourceMetadataKindResourceKindImage),
	)

	// Select sku
	storageSkuName := suite.storageSkus[rand.Intn(len(suite.storageSkus))]

	// Generate scenario data
	workspaceName := secalib.GenerateWorkspaceName()
	workspaceResource := secalib.GenerateWorkspaceResource(suite.tenant, workspaceName)

	storageSkuRef := secalib.GenerateSkuRef(storageSkuName)
	storageSkuRefObj, err := secapi.BuildReferenceFromURN(storageSkuRef)
	if err != nil {
		t.Fatal(err)
	}

	blockStorageName := secalib.GenerateBlockStorageName()
	blockStorageResource := secalib.GenerateBlockStorageResource(suite.tenant, workspaceName, blockStorageName)
	blockStorageRef := secalib.GenerateBlockStorageRef(blockStorageName)
	blockStorageRefObj, err := secapi.BuildReferenceFromURN(blockStorageRef)
	if err != nil {
		t.Fatal(err)
	}

	imageName := secalib.GenerateImageName()
	imageResource := secalib.GenerateImageResource(suite.tenant, imageName)

	initialStorageSize := secalib.GenerateBlockStorageSize()
	updatedStorageSize := secalib.GenerateBlockStorageSize()

	// Setup mock, if configured to use
	if suite.mockEnabled {
		mockParams := &mock.StorageParamsV1{
			Params: &mock.Params{
				MockURL:   *suite.mockServerURL,
				AuthToken: suite.authToken,
				Tenant:    suite.tenant,
				Region:    suite.region,
			},
			Workspace: &mock.ResourceParams[schema.WorkspaceSpec]{
				Name: workspaceName,
				InitialLabels: schema.Labels{
					secalib.EnvLabel: secalib.EnvDevelopmentLabel,
				},
			},
			BlockStorage: &[]mock.ResourceParams[schema.BlockStorageSpec]{
				{
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
			},
			Image: &[]mock.ResourceParams[schema.ImageSpec]{
				{
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
			},
		}
		wm, err := mock.ConfigStorageLifecycleScenarioV1(suite.scenarioName, mockParams)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}
		suite.mockClient = wm
	}

	// Workspace

	// Create a workspace
	workspace := &schema.Workspace{
		Labels: schema.Labels{
			secalib.EnvLabel: secalib.EnvDevelopmentLabel,
		},
		Metadata: &schema.RegionalResourceMetadata{
			Tenant: suite.tenant,
			Name:   workspaceName,
		},
	}
	expectWorkspaceMeta, err := builders.NewRegionalResourceMetadataBuilder().
		Name(workspaceName).
		Provider(secalib.WorkspaceProviderV1).
		Resource(workspaceResource).
		ApiVersion(secalib.ApiVersion1).
		Kind(schema.RegionalResourceMetadataKindResourceKindWorkspace).
		Tenant(suite.tenant).
		Region(suite.region).
		BuildResponse()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectWorkspaceLabels := schema.Labels{secalib.EnvLabel: secalib.EnvDevelopmentLabel}

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
	expectedBlockMeta, err := builders.NewRegionalWorkspaceResourceMetadataBuilder().
		Name(blockStorageName).
		Provider(secalib.StorageProviderV1).
		Resource(blockStorageResource).
		ApiVersion(secalib.ApiVersion1).
		Kind(schema.RegionalWorkspaceResourceMetadataKindResourceKindBlockStorage).
		Tenant(suite.tenant).
		Workspace(workspaceName).
		Region(suite.region).
		BuildResponse()
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
	expectedImageMeta, err := builders.NewRegionalResourceMetadataBuilder().
		Name(imageName).
		Provider(secalib.StorageProviderV1).
		Resource(imageResource).
		ApiVersion(secalib.ApiVersion1).
		Kind(schema.RegionalResourceMetadataKindResourceKindImage).
		Tenant(suite.tenant).
		Region(suite.region).
		BuildResponse()
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

	// Delete the image
	suite.deleteImageV1Step("Delete the image", t, suite.client.StorageV1, image)

	// Get the deleted image
	suite.getImageWithErrorV1Step("Get the deleted image", t, suite.client.StorageV1, *imageTRef, secapi.ErrResourceNotFound)

	// Delete the block storage
	suite.deleteBlockStorageV1Step("Delete the block storage", t, suite.client.StorageV1, block)

	// Get the deleted block storage
	suite.getBlockStorageWithErrorV1Step("Get the deleted block storage", t, suite.client.StorageV1, *blockWRef, secapi.ErrResourceNotFound)

	// Delete the workspace
	suite.deleteWorkspaceV1Step("Delete the workspace", t, suite.client.WorkspaceV1, workspace)

	// Get the deleted workspace
	suite.getWorkspaceWithErrorV1Step("Get the deleted workspace", t, suite.client.WorkspaceV1, *workspaceTRef, secapi.ErrResourceNotFound)

	slog.Info("Finishing " + suite.scenarioName)
}

func (suite *StorageV1TestSuite) TestListSuite(t provider.T) {
	ctx := context.Background()
	var err error
	slog.Info("Starting " + suite.scenarioName)

	t.Title(suite.scenarioName)
	configureTags(t, secalib.StorageProviderV1, secalib.BlockStorageKind, secalib.ImageKind)

	// Select sku
	storageSkuName := suite.storageSkus[rand.Intn(len(suite.storageSkus))]

	// Generate scenario data
	workspaceName := secalib.GenerateWorkspaceName()

	storageSkuRef := secalib.GenerateSkuRef(storageSkuName)
	storageSkuRefObj, err := secapi.BuildReferenceFromURN(storageSkuRef)
	if err != nil {
		t.Fatal(err)
	}

	blockStorageName1 := secalib.GenerateBlockStorageName()
	blockStorageResource1 := secalib.GenerateBlockStorageResource(suite.tenant, workspaceName, blockStorageName1)
	blockStorageName2 := secalib.GenerateBlockStorageName()
	blockStorageResource2 := secalib.GenerateBlockStorageResource(suite.tenant, workspaceName, blockStorageName2)
	blockStorageName3 := secalib.GenerateBlockStorageName()
	blockStorageResource3 := secalib.GenerateBlockStorageResource(suite.tenant, workspaceName, blockStorageName3)

	blockStorageRef := secalib.GenerateBlockStorageRef(blockStorageName1)
	blockStorageRefObj, err := secapi.BuildReferenceFromURN(blockStorageRef)
	if err != nil {
		t.Fatal(err)
	}

	imageName1 := secalib.GenerateImageName()
	imageResource1 := secalib.GenerateImageResource(suite.tenant, imageName1)
	imageName2 := secalib.GenerateImageName()
	imageResource2 := secalib.GenerateImageResource(suite.tenant, imageName2)
	imageName3 := secalib.GenerateImageName()
	imageResource3 := secalib.GenerateImageResource(suite.tenant, imageName3)
	initialStorageSize := secalib.GenerateBlockStorageSize()

	// Setup mock, if configured to use
	if suite.mockEnabled {
		mockParams := &mock.StorageParamsV1{
			Params: &mock.Params{
				MockURL:   *suite.mockServerURL,
				AuthToken: suite.authToken,
				Tenant:    suite.tenant,
				Region:    suite.region,
			},
			Workspace: &mock.ResourceParams[schema.WorkspaceSpec]{
				Name: workspaceName,
				InitialLabels: schema.Labels{
					secalib.EnvLabel: secalib.EnvDevelopmentLabel,
				},
			},
			BlockStorage: &[]mock.ResourceParams[schema.BlockStorageSpec]{
				{
					Name: blockStorageName1,
					InitialLabels: map[string]string{
						secalib.EnvLabel: secalib.EnvConformance,
					},
					InitialSpec: &schema.BlockStorageSpec{
						SkuRef: *storageSkuRefObj,
						SizeGB: initialStorageSize,
					},
				},
				{
					Name: blockStorageName2,
					InitialLabels: map[string]string{
						secalib.EnvLabel: secalib.EnvConformance,
					},
					InitialSpec: &schema.BlockStorageSpec{
						SkuRef: *storageSkuRefObj,
						SizeGB: initialStorageSize,
					},
				},
				{
					Name: blockStorageName3,
					InitialLabels: map[string]string{
						secalib.EnvLabel: secalib.EnvConformance,
					},
					InitialSpec: &schema.BlockStorageSpec{
						SkuRef: *storageSkuRefObj,
						SizeGB: initialStorageSize,
					},
				},
			},
			Image: &[]mock.ResourceParams[schema.ImageSpec]{
				{
					Name: imageName1,
					InitialLabels: map[string]string{
						secalib.EnvLabel: secalib.EnvConformance,
					},
					InitialSpec: &schema.ImageSpec{
						BlockStorageRef: *blockStorageRefObj,
						CpuArchitecture: secalib.CpuArchitectureAmd64,
					},
				},
				{
					Name: imageName2,
					InitialLabels: map[string]string{
						secalib.EnvLabel: secalib.EnvConformance,
					},
					InitialSpec: &schema.ImageSpec{
						BlockStorageRef: *blockStorageRefObj,
						CpuArchitecture: secalib.CpuArchitectureAmd64,
					},
				},
				{
					Name: imageName3,
					InitialLabels: map[string]string{
						secalib.EnvLabel: secalib.EnvConformance,
					},
					InitialSpec: &schema.ImageSpec{
						BlockStorageRef: *blockStorageRefObj,
						CpuArchitecture: secalib.CpuArchitectureAmd64,
					},
				},
			},
		}
		wm, err := mock.ConfigStorageListLifecycleScenarioV1(suite.scenarioName, mockParams)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}
		suite.mockClient = wm
	}

	// Workspace

	// Create a workspace
	workspace := &schema.Workspace{
		Labels: schema.Labels{
			secalib.EnvLabel: secalib.EnvDevelopmentLabel,
		},
		Metadata: &schema.RegionalResourceMetadata{
			Tenant: suite.tenant,
			Name:   workspaceName,
		},
	}
	suite.createOrUpdateWorkspaceV1Step("Create a workspace", t, ctx, suite.client.WorkspaceV1, workspace, nil, nil, secalib.CreatingResourceState)

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
	for _, block := range *blocks {
		expectedBlockMeta := secalib.NewRegionalWorkspaceResourceMetadata(block.Metadata.Name,
			secalib.StorageProviderV1,
			block.Metadata.Resource,
			secalib.ApiVersion1,
			secalib.BlockStorageKind,
			suite.tenant,
			workspaceName,
			suite.region)
		expectedBlockSpec := &schema.BlockStorageSpec{
			SizeGB: initialStorageSize,
			SkuRef: *storageSkuRefObj,
		}
		suite.createOrUpdateBlockStorageV1Step("Create a block storage", t, ctx, suite.client.StorageV1, &block,
			expectedBlockMeta, expectedBlockSpec, secalib.CreatingResourceState)
	}
	tref := secapi.TenantReference{Tenant: secapi.TenantID(suite.tenant)}
	wref := secapi.WorkspaceReference{Workspace: secapi.WorkspaceID(workspaceName)}
	// List blockstorage
	suite.getListBlockStorageV1Step("GetList block storage", t, ctx, suite.client.StorageV1, tref, wref, nil)

	// List instances with limit
	suite.getListBlockStorageV1Step("Get List block storage with limit", t, ctx, suite.client.StorageV1, tref, wref,
		builders.NewListOptions().WithLimit(1))

	// List Instances with Label
	suite.getListBlockStorageV1Step("Get list of block storage with label", t, ctx, suite.client.StorageV1, tref, wref,
		builders.NewListOptions().WithLabels(builders.NewLabelsBuilder().
			Equals(secalib.EnvLabel, secalib.EnvConformance)))

	// List Instances with Limit and label
	suite.getListBlockStorageV1Step("Get list of block storage with limit and label", t, ctx, suite.client.StorageV1, tref, wref,
		builders.NewListOptions().WithLimit(1).WithLabels(builders.NewLabelsBuilder().
			Equals(secalib.EnvLabel, secalib.EnvConformance)))

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
				CpuArchitecture: secalib.CpuArchitectureAmd64,
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
				CpuArchitecture: secalib.CpuArchitectureAmd64,
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
				CpuArchitecture: secalib.CpuArchitectureAmd64,
			},
		},
	}

	for _, image := range *images {

		expectedImageMeta := secalib.NewRegionalResourceMetadata(image.Metadata.Name,
			secalib.StorageProviderV1,
			image.Metadata.Resource,
			secalib.ApiVersion1,
			secalib.ImageKind,
			suite.tenant,
			suite.region)
		expectedImageSpec := &schema.ImageSpec{
			BlockStorageRef: *blockStorageRefObj,
			CpuArchitecture: secalib.CpuArchitectureAmd64,
		}
		suite.createOrUpdateImageV1Step("Create an image", t, ctx, suite.client.StorageV1, &image,
			expectedImageMeta, expectedImageSpec, secalib.CreatingResourceState)
	}
	// List image
	suite.getListImageV1Step("List image", t, ctx, suite.client.StorageV1, tref, nil)

	// List image with limit
	suite.getListImageV1Step("Get list of images", t, ctx, suite.client.StorageV1, tref,
		builders.NewListOptions().WithLimit(1))

	// List image with Label
	suite.getListImageV1Step("Get list of images", t, ctx, suite.client.StorageV1, tref,
		builders.NewListOptions().WithLabels(builders.NewLabelsBuilder().
			Equals(secalib.EnvLabel, secalib.EnvConformance)))

	// List image with Limit and label
	suite.getListImageV1Step("Get list of images", t, ctx, suite.client.StorageV1, tref,
		builders.NewListOptions().WithLimit(1).WithLabels(builders.NewLabelsBuilder().
			Equals(secalib.EnvLabel, secalib.EnvConformance)))

	// Skus
	// List Skus
	suite.getListSkuV1Step("List skus", t, ctx, suite.client.StorageV1, tref, nil)

	// List Skus with limit
	suite.getListSkuV1Step("Get list of skus", t, ctx, suite.client.StorageV1, tref,
		builders.NewListOptions().WithLimit(1))

	// List Skus with Label
	suite.getListSkuV1Step("Get list of skus", t, ctx, suite.client.StorageV1, tref,
		builders.NewListOptions().WithLabels(builders.NewLabelsBuilder().
			Equals(secalib.EnvLabel, secalib.EnvConformance)))

	// List Skus with Limit and label
	suite.getListSkuV1Step("Get list of skus", t, ctx, suite.client.StorageV1, tref,
		builders.NewListOptions().WithLimit(1).WithLabels(builders.NewLabelsBuilder().
			Equals(secalib.EnvLabel, secalib.EnvConformance)))

	// Resources deletion

	// Delete all images
	imageTRef1 := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
		Name:   imageName1,
	}
	suite.deleteImageV1Step("Delete image 1", t, ctx, suite.client.StorageV1, &(*images)[0])
	suite.getImageWithErrorV1Step("Get deleted image 1", t, ctx, suite.client.StorageV1, *imageTRef1, secapi.ErrResourceNotFound)

	imageTRef2 := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
		Name:   imageName2,
	}
	suite.deleteImageV1Step("Delete image 2", t, ctx, suite.client.StorageV1, &(*images)[1])
	suite.getImageWithErrorV1Step("Get deleted image 2", t, ctx, suite.client.StorageV1, *imageTRef2, secapi.ErrResourceNotFound)

	imageTRef3 := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
		Name:   imageName3,
	}
	suite.deleteImageV1Step("Delete image 3", t, ctx, suite.client.StorageV1, &(*images)[2])
	suite.getImageWithErrorV1Step("Get deleted image 3", t, ctx, suite.client.StorageV1, *imageTRef3, secapi.ErrResourceNotFound)

	// Delete all block storages
	blockWRef1 := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      blockStorageName1,
	}
	suite.deleteBlockStorageV1Step("Delete block storage 1", t, ctx, suite.client.StorageV1, &(*blocks)[0])
	suite.getBlockStorageWithErrorV1Step("Get deleted block storage 1", t, ctx, suite.client.StorageV1, *blockWRef1, secapi.ErrResourceNotFound)

	blockWRef2 := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      blockStorageName2,
	}
	suite.deleteBlockStorageV1Step("Delete block storage 2", t, ctx, suite.client.StorageV1, &(*blocks)[1])
	suite.getBlockStorageWithErrorV1Step("Get deleted block storage 2", t, ctx, suite.client.StorageV1, *blockWRef2, secapi.ErrResourceNotFound)

	blockWRef3 := &secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      blockStorageName3,
	}
	suite.deleteBlockStorageV1Step("Delete block storage 3", t, ctx, suite.client.StorageV1, &(*blocks)[2])
	suite.getBlockStorageWithErrorV1Step("Get deleted block storage 3", t, ctx, suite.client.StorageV1, *blockWRef3, secapi.ErrResourceNotFound)

	// Delete the workspace
	workspaceTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
		Name:   workspaceName,
	}
	suite.deleteWorkspaceV1Step("Delete the workspace", t, ctx, suite.client.WorkspaceV1, workspace)
	suite.getWorkspaceWithErrorV1Step("Get the deleted workspace", t, ctx, suite.client.WorkspaceV1, *workspaceTRef, secapi.ErrResourceNotFound)

	slog.Info("Finishing " + suite.scenarioName)
}

func (suite *StorageV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}
