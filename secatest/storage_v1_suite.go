package secatest

import (
	"context"
	"log/slog"
	"math/rand"
	"net/http"

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
	var workspaceResp *schema.Workspace

	t.WithNewStep("Create workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "CreateOrUpdateWorkspace")

		ws := &schema.Workspace{
			Metadata: &schema.RegionalResourceMetadata{
				Tenant: suite.tenant,
				Name:   workspaceName,
			},
		}
		workspaceResp, err = suite.client.WorkspaceV1.CreateOrUpdateWorkspace(ctx, ws)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, workspaceResp)

		suite.verifyStatusStep(sCtx, secalib.CreatingStatusState, *workspaceResp.Status.State)
	})

	t.WithNewStep("Get created workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "GetWorkspace")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   workspaceName,
		}
		workspaceResp, err = suite.client.WorkspaceV1.GetWorkspace(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, workspaceResp)

		suite.verifyStatusStep(sCtx, secalib.ActiveStatusState, *workspaceResp.Status.State)
	})

	// Block storage
	var blockResp *schema.BlockStorage
	var expectedBlockMeta *schema.RegionalWorkspaceResourceMetadata
	var expectedBlockSpec *schema.BlockStorageSpec

	t.WithNewStep("Create block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "CreateOrUpdateBlockStorage", workspaceName)

		bo := &schema.BlockStorage{
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
		blockResp, err = suite.client.StorageV1.CreateOrUpdateBlockStorage(ctx, bo)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockResp)

		expectedBlockMeta = secalib.NewRegionalWorkspaceResourceMetadata(blockStorageName, secalib.StorageProviderV1, blockStorageResource, secalib.ApiVersion1, secalib.BlockStorageKind,
			suite.tenant, workspaceName, suite.region)
		expectedBlockMeta.Verb = http.MethodPut
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedBlockMeta, blockResp.Metadata)

		expectedBlockSpec = &schema.BlockStorageSpec{
			SizeGB: initialStorageSize,
			SkuRef: *storageSkuRefObj,
		}
		suite.verifyBlockStorageSpecStep(sCtx, expectedBlockSpec, &blockResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.CreatingStatusState, *blockResp.Status.State)
	})

	t.WithNewStep("Get created block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "GetBlockStorage", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      blockStorageName,
		}
		blockResp, err = suite.client.StorageV1.GetBlockStorage(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockResp)

		expectedBlockMeta.Verb = http.MethodGet
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedBlockMeta, blockResp.Metadata)

		suite.verifyBlockStorageSpecStep(sCtx, expectedBlockSpec, &blockResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveStatusState, *blockResp.Status.State)
	})

	t.WithNewStep("Update block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "CreateOrUpdateBlockStorage", workspaceName)

		blockResp, err = suite.client.StorageV1.CreateOrUpdateBlockStorage(ctx, blockResp)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockResp)

		expectedBlockMeta.Verb = http.MethodPut
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedBlockMeta, blockResp.Metadata)

		expectedBlockSpec.SizeGB = updatedStorageSize
		suite.verifyBlockStorageSpecStep(sCtx, expectedBlockSpec, &blockResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.UpdatingStatusState, *blockResp.Status.State)
	})

	t.WithNewStep("Get updated block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "GetBlockStorage", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      blockStorageName,
		}
		blockResp, err = suite.client.StorageV1.GetBlockStorage(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockResp)

		expectedBlockMeta.Verb = http.MethodGet
		suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedBlockMeta, blockResp.Metadata)

		suite.verifyBlockStorageSpecStep(sCtx, expectedBlockSpec, &blockResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveStatusState, *blockResp.Status.State)
	})

	// Image
	var imageResp *schema.Image
	var expectedImageMeta *schema.RegionalResourceMetadata
	var expectedImageSpec *schema.ImageSpec

	t.WithNewStep("Create image", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "CreateOrUpdateImage", "")

		img := &schema.Image{
			Metadata: &schema.RegionalResourceMetadata{
				Tenant: suite.tenant,
				Name:   imageName,
			},
			Spec: schema.ImageSpec{
				BlockStorageRef: *blockStorageRefObj,
				CpuArchitecture: secalib.CpuArchitectureAmd64,
			},
		}
		imageResp, err = suite.client.StorageV1.CreateOrUpdateImage(ctx, img)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, imageResp)

		expectedImageMeta = secalib.NewRegionalResourceMetadata(imageName, secalib.StorageProviderV1, imageResource, secalib.ApiVersion1, secalib.ImageKind, suite.tenant, suite.region)
		expectedImageMeta.Verb = http.MethodPut
		suite.verifyRegionalResourceMetadataStep(sCtx, expectedImageMeta, imageResp.Metadata)

		expectedImageSpec = &schema.ImageSpec{
			BlockStorageRef: *blockStorageRefObj,
			CpuArchitecture: secalib.CpuArchitectureAmd64,
		}
		suite.verifyImageSpecStep(sCtx, expectedImageSpec, &imageResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.CreatingStatusState, *imageResp.Status.State)
	})

	t.WithNewStep("Get created image", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "GetImage", "")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   imageName,
		}
		imageResp, err = suite.client.StorageV1.GetImage(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, imageResp)

		expectedImageMeta.Verb = http.MethodGet
		suite.verifyRegionalResourceMetadataStep(sCtx, expectedImageMeta, imageResp.Metadata)

		suite.verifyImageSpecStep(sCtx, expectedImageSpec, &imageResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveStatusState, *imageResp.Status.State)
	})

	t.WithNewStep("Update image", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "CreateOrUpdateImage", "")

		imageResp, err = suite.client.StorageV1.CreateOrUpdateImage(ctx, imageResp)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, imageResp)

		expectedImageMeta.Verb = http.MethodPut
		suite.verifyRegionalResourceMetadataStep(sCtx, expectedImageMeta, imageResp.Metadata)

		expectedImageSpec.CpuArchitecture = secalib.CpuArchitectureArm64
		suite.verifyImageSpecStep(sCtx, expectedImageSpec, &imageResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.UpdatingStatusState, *imageResp.Status.State)
	})

	t.WithNewStep("Get updated image", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "GetImage", "")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   imageName,
		}
		imageResp, err = suite.client.StorageV1.GetImage(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, imageResp)

		expectedImageMeta.Verb = http.MethodGet
		suite.verifyRegionalResourceMetadataStep(sCtx, expectedImageMeta, imageResp.Metadata)

		suite.verifyImageSpecStep(sCtx, expectedImageSpec, &imageResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveStatusState, *imageResp.Status.State)
	})

	t.WithNewStep("Delete image", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "DeleteImage", "")

		err = suite.client.StorageV1.DeleteImage(ctx, imageResp)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get deleted image", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "GetImage", "")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   imageName,
		}
		_, err = suite.client.StorageV1.GetImage(ctx, tref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	t.WithNewStep("Delete block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "DeleteBlockStorage", workspaceName)

		err = suite.client.StorageV1.DeleteBlockStorage(ctx, blockResp)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get deleted block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "GetBlockStorage", workspaceName)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      blockStorageName,
		}
		_, err = suite.client.StorageV1.GetBlockStorage(ctx, wref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	slog.Info("Finishing " + suite.scenarioName)
}

func (suite *StorageV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}
