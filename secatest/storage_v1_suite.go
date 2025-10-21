package secatest

import (
	"context"
	"log/slog"
	"math/rand"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type StorageV1TestSuite struct {
	regionalTestSuite

	storageSkus []string
}

func (suite *StorageV1TestSuite) TestSuite(t provider.T) {
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
					CpuArchitecture: secalib.CpuArchitectureAmd64,
				},
				UpdatedSpec: &schema.ImageSpec{
					BlockStorageRef: *blockStorageRefObj,
					CpuArchitecture: secalib.CpuArchitectureArm64,
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
	suite.createOrUpdateWorkspaceV1Step("Create a workspace", t, ctx, suite.client.WorkspaceV1, workspace, nil, nil, secalib.CreatingResourceState)

	// Get the created Workspace
	workspaceTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
		Name:   workspaceName,
	}
	workspace = suite.getWorkspaceV1Step("Get the created workspace", t, ctx, suite.client.WorkspaceV1, *workspaceTRef, nil, nil, secalib.ActiveResourceState)

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
	expectedBlockMeta := secalib.NewRegionalWorkspaceResourceMetadata(blockStorageName,
		secalib.StorageProviderV1,
		blockStorageResource,
		secalib.ApiVersion1,
		secalib.BlockStorageKind,
		suite.tenant,
		workspaceName,
		suite.region)
	expectedBlockSpec := block.Spec
	suite.createOrUpdateBlockStorageV1Step("Create a block storage", t, ctx, suite.client.StorageV1, block,
		expectedBlockMeta, &expectedBlockSpec, secalib.CreatingResourceState)

	// Get the created block storage
	blockWRef := secapi.WorkspaceReference{
		Tenant:    secapi.TenantID(suite.tenant),
		Workspace: secapi.WorkspaceID(workspaceName),
		Name:      blockStorageName,
	}
	block = suite.getBlockStorageV1Step("Get the created block storage", t, ctx, suite.client.StorageV1, blockWRef,
		expectedBlockMeta, &expectedBlockSpec, secalib.ActiveResourceState)

	// Update the block storage
	block.Spec.SizeGB = updatedStorageSize
	expectedBlockSpec = block.Spec
	suite.createOrUpdateBlockStorageV1Step("Update the block storage", t, ctx, suite.client.StorageV1, block,
		expectedBlockMeta, &expectedBlockSpec, secalib.UpdatingResourceState)

	// Get the updated block storage
	block = suite.getBlockStorageV1Step("Get the updated block storage", t, ctx, suite.client.StorageV1, blockWRef,
		expectedBlockMeta, &expectedBlockSpec, secalib.ActiveResourceState)

	// Image

	// Create an image
	image := &schema.Image{
		Metadata: &schema.RegionalResourceMetadata{
			Tenant: suite.tenant,
			Name:   imageName,
		},
		Spec: schema.ImageSpec{
			BlockStorageRef: *blockStorageRefObj,
			CpuArchitecture: secalib.CpuArchitectureAmd64,
		},
	}
	expectedImageMeta := secalib.NewRegionalResourceMetadata(imageName,
		secalib.StorageProviderV1,
		imageResource,
		secalib.ApiVersion1,
		secalib.ImageKind,
		suite.tenant,
		suite.region)
	expectedImageSpec := image.Spec
	suite.createOrUpdateImageV1Step("Create an image", t, ctx, suite.client.StorageV1, image,
		expectedImageMeta, &expectedImageSpec, secalib.CreatingResourceState)

	// Get the created image
	imageTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
		Name:   imageName,
	}
	image = suite.getImageV1Step("Get the created image", t, ctx, suite.client.StorageV1, imageTRef,
		expectedImageMeta, &expectedImageSpec, secalib.ActiveResourceState)

	// Update the image
	image.Spec.CpuArchitecture = secalib.CpuArchitectureArm64
	expectedImageSpec = image.Spec
	suite.createOrUpdateImageV1Step("Update the image", t, ctx, suite.client.StorageV1, image,
		expectedImageMeta, &expectedImageSpec, secalib.UpdatingResourceState)

	// Get the updated image
	image = suite.getImageV1Step("Get the updated image", t, ctx, suite.client.StorageV1, imageTRef,
		expectedImageMeta, &expectedImageSpec, secalib.ActiveResourceState)

	// Delete the image
	suite.deleteImageV1Step("Delete the image", t, ctx, suite.client.StorageV1, image)

	// Get the deleted image
	suite.getImageWithErrorV1Step("Get the deleted image", t, ctx, suite.client.StorageV1, imageTRef, secapi.ErrResourceNotFound)

	// Delete the block storage
	suite.deleteBlockStorageV1Step("Delete the block storage", t, ctx, suite.client.StorageV1, block)

	// Get the deleted block storage
	suite.getBlockStorageWithErrorV1Step("Get the deleted block storage", t, ctx, suite.client.StorageV1, blockWRef, secapi.ErrResourceNotFound)

	slog.Info("Finishing " + suite.scenarioName)
}

func (suite *StorageV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}
