package secatest

import (
	"context"
	"log/slog"
	"math/rand"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/secalib"
	storage "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.storage.v1"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type StorageV1TestSuite struct {
	regionalTestSuite

	storageSkus []string
}

func (suite *StorageV1TestSuite) TestStorageV1(t provider.T) {
	slog.Info("Starting Storage Lifecycle Test")

	t.Title("Storage Lifecycle Test")
	configureTags(t, secalib.StorageProviderV1, secalib.BlockStorageKind, secalib.ImageKind)

	// Select sku
	storageSkuName := suite.storageSkus[rand.Intn(len(suite.storageSkus))]

	// TODO Dynamically create before the scenario
	workspaceName := secalib.GenerateWorkspaceName()

	// Generate scenario data
	storageSkuRef := secalib.GenerateSkuRef(storageSkuName)

	blockStorageName := secalib.GenerateBlockStorageName()
	blockStorageResource := secalib.GenerateBlockStorageResource(suite.tenant, workspaceName, blockStorageName)
	blockStorageRef := secalib.GenerateBlockStorageRef(blockStorageName)

	imageName := secalib.GenerateImageName()
	imageResource := secalib.GenerateImageResource(suite.tenant, imageName)

	initialStorageSize := secalib.GenerateBlockStorageSize()
	updatedStorageSize := secalib.GenerateBlockStorageSize()

	// Setup mock, if configured to use
	if suite.isMockEnabled() {
		wm, err := mock.CreateStorageLifecycleScenarioV1("Storage Lifecycle",
			mock.StorageParamsV1{
				Params: &mock.Params{
					MockURL:   suite.mockServerURL,
					AuthToken: suite.authToken,
					Tenant:    suite.tenant,
					Workspace: workspaceName,
					Region:    suite.region,
				},
				BlockStorage: &mock.ResourceParams[secalib.BlockStorageSpecV1]{
					Name: blockStorageName,
					InitialSpec: &secalib.BlockStorageSpecV1{
						SkuRef: storageSkuRef,
						SizeGB: initialStorageSize,
					},
					UpdatedSpec: &secalib.BlockStorageSpecV1{
						SkuRef: storageSkuRef,
						SizeGB: updatedStorageSize,
					},
				},
				Image: &mock.ResourceParams[secalib.ImageSpecV1]{
					Name: imageName,
					InitialSpec: &secalib.ImageSpecV1{
						BlockStorageRef: blockStorageRef,
						CpuArchitecture: secalib.CpuArchitectureAmd64,
					},
					UpdatedSpec: &secalib.ImageSpecV1{
						BlockStorageRef: blockStorageRef,
						CpuArchitecture: secalib.CpuArchitectureArm64,
					},
				},
			})
		if err != nil {
			slog.Error("Failed to create storage scenario", "error", err)
			return
		}
		suite.mockClient = wm
	}

	ctx := context.Background()
	var blockResp *storage.BlockStorage
	var imageResp *storage.Image
	var err error

	var expectedBlockMetadata *secalib.Metadata
	var expectedBlockSpec *secalib.BlockStorageSpecV1

	// Step 1
	t.WithNewStep("Create block storage", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateBlockStorage",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      blockStorageName,
		}
		bo := &storage.BlockStorage{
			Spec: storage.BlockStorageSpec{
				SizeGB: initialStorageSize,
				SkuRef: storageSkuRef,
			},
		}
		blockResp, err = suite.client.StorageV1.CreateOrUpdateBlockStorage(ctx, wref, bo, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockResp)

		expectedBlockMetadata = &secalib.Metadata{
			Name:       blockStorageName,
			Provider:   secalib.StorageProviderV1,
			Resource:   blockStorageResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.BlockStorageKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}

		verifyStorageZonalMetadataStep(sCtx, expectedBlockMetadata, blockResp.Metadata)

		expectedBlockSpec = &secalib.BlockStorageSpecV1{
			SizeGB: initialStorageSize,
			SkuRef: storageSkuRef,
		}
		verifyBlockStorageSpecStep(sCtx, expectedBlockSpec, blockResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*blockResp.Status.State)},
		)
	})

	// Step 2
	t.WithNewStep("Get created block storage", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetBlockStorage",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      blockStorageName,
		}
		blockResp, err = suite.client.StorageV1.GetBlockStorage(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockResp)

		expectedBlockMetadata.Verb = http.MethodGet
		verifyStorageZonalMetadataStep(sCtx, expectedBlockMetadata, blockResp.Metadata)

		verifyBlockStorageSpecStep(sCtx, expectedBlockSpec, blockResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*blockResp.Status.State)},
		)
	})

	// Step 3
	t.WithNewStep("Update block storage", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateBlockStorage",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      blockStorageName,
		}
		blockResp, err = suite.client.StorageV1.CreateOrUpdateBlockStorage(ctx, wref, blockResp, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockResp)

		expectedBlockMetadata.Verb = http.MethodPut
		verifyStorageZonalMetadataStep(sCtx, expectedBlockMetadata, blockResp.Metadata)

		expectedBlockSpec.SizeGB = updatedStorageSize
		verifyBlockStorageSpecStep(sCtx, expectedBlockSpec, blockResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*blockResp.Status.State)},
		)
	})

	// Step 4
	t.WithNewStep("Get updated block storage", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetBlockStorage",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      blockStorageName,
		}
		blockResp, err = suite.client.StorageV1.GetBlockStorage(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockResp)

		expectedBlockMetadata.Verb = http.MethodGet
		verifyStorageZonalMetadataStep(sCtx, expectedBlockMetadata, blockResp.Metadata)

		verifyBlockStorageSpecStep(sCtx, expectedBlockSpec, blockResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*blockResp.Status.State)},
		)
	})

	var expectedImageMetadata *secalib.Metadata
	var expectedImageSpec *secalib.ImageSpecV1

	// Step 5
	t.WithNewStep("Create image", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateImage",
			tenantStepParameter, suite.tenant,
		)

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   imageName,
		}
		img := &storage.Image{
			Spec: storage.ImageSpec{
				BlockStorageRef: blockStorageRef,
				CpuArchitecture: secalib.CpuArchitectureAmd64,
			},
		}
		imageResp, err = suite.client.StorageV1.CreateOrUpdateImage(ctx, tref, img, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, imageResp)

		expectedImageMetadata = &secalib.Metadata{
			Name:       imageName,
			Provider:   secalib.StorageProviderV1,
			Resource:   imageResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.ImageKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyStorageRegionalMetadataStep(sCtx, expectedImageMetadata, imageResp.Metadata)

		expectedImageSpec = &secalib.ImageSpecV1{
			BlockStorageRef: blockStorageRef,
			CpuArchitecture: secalib.CpuArchitectureAmd64,
		}
		verifyImageSpecStep(sCtx, expectedImageSpec, imageResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*imageResp.Status.State)},
		)
	})

	// Step 6
	t.WithNewStep("Get created image", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetImage",
			tenantStepParameter, suite.tenant,
		)

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   imageName,
		}
		imageResp, err = suite.client.StorageV1.GetImage(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, imageResp)

		expectedImageMetadata.Verb = http.MethodGet
		verifyStorageRegionalMetadataStep(sCtx, expectedImageMetadata, imageResp.Metadata)

		verifyImageSpecStep(sCtx, expectedImageSpec, imageResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*imageResp.Status.State)},
		)
	})

	// Step 7
	t.WithNewStep("Update image", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateImage",
			tenantStepParameter, suite.tenant,
		)

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   imageName,
		}
		imageResp, err = suite.client.StorageV1.CreateOrUpdateImage(ctx, tref, imageResp, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, imageResp)

		expectedImageMetadata.Verb = http.MethodPut
		verifyStorageRegionalMetadataStep(sCtx, expectedImageMetadata, imageResp.Metadata)

		expectedImageSpec.CpuArchitecture = secalib.CpuArchitectureArm64
		verifyImageSpecStep(sCtx, expectedImageSpec, imageResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*imageResp.Status.State)},
		)
	})

	// Step 8
	t.WithNewStep("Get updated image", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetImage",
			tenantStepParameter, suite.tenant,
		)

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   imageName,
		}
		imageResp, err = suite.client.StorageV1.GetImage(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, imageResp)

		expectedImageMetadata.Verb = http.MethodGet
		verifyStorageRegionalMetadataStep(sCtx, expectedImageMetadata, imageResp.Metadata)

		verifyImageSpecStep(sCtx, expectedImageSpec, imageResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*imageResp.Status.State)},
		)
	})

	// Step 9
	t.WithNewStep("Delete image", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "DeleteImage",
			tenantStepParameter, suite.tenant,
		)

		err = suite.client.StorageV1.DeleteImage(ctx, imageResp, nil)
		requireNoError(sCtx, err)
	})

	// Step 10
	t.WithNewStep("Get deleted image", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetImage",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   imageName,
		}
		_, err = suite.client.StorageV1.GetImage(ctx, tref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	// Step 11
	t.WithNewStep("Delete block storage", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "DeleteBlockStorage",
			tenantStepParameter, suite.tenant,
		)

		err = suite.client.StorageV1.DeleteBlockStorage(ctx, blockResp, nil)
		requireNoError(sCtx, err)
	})

	// Step 12
	t.WithNewStep("Get deleted block storage", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetBlockStorage",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(workspaceName),
			Name:      blockStorageName,
		}
		_, err = suite.client.StorageV1.GetBlockStorage(ctx, wref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	slog.Info("Finishing Storage Lifecycle Test")
}

func (suite *StorageV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}

func verifyStorageZonalMetadataStep(ctx provider.StepCtx, expected *secalib.Metadata, metadata *storage.ZonalResourceMetadata) {
	actualMetadata := &secalib.Metadata{
		Name:       metadata.Name,
		Provider:   metadata.Provider,
		Verb:       metadata.Verb,
		Resource:   metadata.Resource,
		ApiVersion: metadata.ApiVersion,
		Kind:       string(metadata.Kind),
		Tenant:     metadata.Tenant,
		Workspace:  *metadata.Workspace,
		Region:     metadata.Region,
	}
	verifyRegionalMetadataStep(ctx, expected, actualMetadata)
}

func verifyStorageRegionalMetadataStep(ctx provider.StepCtx, expected *secalib.Metadata, metadata *storage.RegionalResourceMetadata) {
	actualMetadata := &secalib.Metadata{
		Name:       metadata.Name,
		Provider:   metadata.Provider,
		Verb:       metadata.Verb,
		Resource:   metadata.Resource,
		ApiVersion: metadata.ApiVersion,
		Kind:       string(metadata.Kind),
		Tenant:     metadata.Tenant,
		Workspace:  *metadata.Workspace,
		Region:     metadata.Region,
	}
	verifyRegionalMetadataStep(ctx, expected, actualMetadata)
}

func verifyBlockStorageSpecStep(ctx provider.StepCtx, expected *secalib.BlockStorageSpecV1, actual storage.BlockStorageSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters(
			"expected_sizeGB", expected.SizeGB,
			"actual_sizeGB", actual.SizeGB,

			"expected_skuRef", expected.SkuRef,
			"actual_skuRef", actual.SkuRef,
		)
		stepCtx.Require().Equal(expected.SizeGB, actual.SizeGB, "SizeGB should match expected")
		stepCtx.Require().Equal(expected.SkuRef, actual.SkuRef, "SkuRef should match expected")
	})
}

func verifyImageSpecStep(ctx provider.StepCtx, expected *secalib.ImageSpecV1, actual storage.ImageSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters(
			"expected_blockStorageRef", expected.BlockStorageRef,
			"actual_blockStorageRef", actual.BlockStorageRef,

			"expected_cpuArchitecture", expected.CpuArchitecture,
			"actual_cpuArchitecture", actual.CpuArchitecture,
		)
		stepCtx.Require().Equal(expected.BlockStorageRef, actual.BlockStorageRef, "BlockStorageRef should match expected")
		stepCtx.Require().Equal(expected.CpuArchitecture, string(actual.CpuArchitecture), "CpuArchitecture should match expected")
	})
}
