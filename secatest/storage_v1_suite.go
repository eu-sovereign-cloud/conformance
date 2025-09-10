package secatest

import (
	"context"
	"log/slog"
	"math/rand"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/secalib"
	storage "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.storage.v1"
	workspace "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.workspace.v1"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/google/uuid"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type StorageV1TestSuite struct {
	regionalTestSuite

	storageSkus []string
}

func (suite *StorageV1TestSuite) generateLifecycleParams() *secalib.StorageLifeCycleParamsV1 {
	workspaceName := secalib.GenerateWorkspaceName()
	storageSkuName := suite.storageSkus[rand.Intn(len(suite.storageSkus))]
	blockStorageName := secalib.GenerateBlockStorageName()
	imageName := secalib.GenerateImageName()

	storageSkuRef := secalib.GenerateSkuRef(storageSkuName)
	blockStorageRef := secalib.GenerateBlockStorageRef(blockStorageName)

	// Random data
	initialStorageSize := secalib.GenerateBlockStorageSize()
	updatedStorageSize := secalib.GenerateBlockStorageSize()

	return &secalib.StorageLifeCycleParamsV1{
		Workspace: &secalib.ResourceParams[secalib.WorkspaceSpecV1]{
			Name: workspaceName,
			InitialSpec: &secalib.WorkspaceSpecV1{
				Labels: &[]secalib.Label{
					{
						Name:  secalib.EnvLabel,
						Value: secalib.EnvDevelopmentLabel,
					},
				},
			},
		},
		BlockStorage: &secalib.ResourceParams[secalib.BlockStorageSpecV1]{
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
		Image: &secalib.ResourceParams[secalib.ImageSpecV1]{
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
	}
}

func (suite *StorageV1TestSuite) TestStorageLifeCycleV1(t provider.T) {
	slog.Info("Starting Storage Lifecycle Test")

	t.Title("Storage Lifecycle Test")
	configureTags(t, secalib.StorageProviderV1, secalib.BlockStorageKind, secalib.ImageKind)

	params := suite.generateLifecycleParams()
	blockStorageResource := secalib.GenerateBlockStorageResource(suite.tenant, params.Workspace.Name, params.BlockStorage.Name)
	imageResource := secalib.GenerateImageResource(suite.tenant, params.Image.Name)

	// Setup mock, if configured to use
	if suite.isMockEnabled() {
		scenarios := mock.NewStorageScenariosV1(suite.authToken, suite.tenant, suite.region, suite.mockServerURL)

		id := uuid.New().String()
		wm, err := scenarios.ConfigureLifecycleScenario(id, params)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}
		suite.mockClient = wm
	}

	ctx := context.Background()
	var err error

	// Workspace
	var workResp *workspace.Workspace

	t.WithNewStep("Create workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "CreateOrUpdateWorkspace")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   params.Workspace.Name,
		}
		ws := &workspace.Workspace{}
		workResp, err = suite.client.WorkspaceV1.CreateOrUpdateWorkspace(ctx, tref, ws, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, workResp)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*workResp.Status.State)},
		)
	})

	t.WithNewStep("Get created workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "GetWorkspace")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   params.Workspace.Name,
		}
		workResp, err = suite.client.WorkspaceV1.GetWorkspace(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, workResp)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*workResp.Status.State)},
		)
	})

	// Block storage
	var blockResp *storage.BlockStorage
	var expectedBlockMeta *secalib.Metadata

	t.WithNewStep("Create block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "CreateOrUpdateBlockStorage", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.BlockStorage.Name,
		}
		bo := &storage.BlockStorage{
			Spec: storage.BlockStorageSpec{
				SizeGB: params.BlockStorage.InitialSpec.SizeGB,
				SkuRef: params.BlockStorage.InitialSpec.SkuRef,
			},
		}
		blockResp, err = suite.client.StorageV1.CreateOrUpdateBlockStorage(ctx, wref, bo, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockResp)

		expectedBlockMeta = &secalib.Metadata{
			Name:       params.BlockStorage.Name,
			Provider:   secalib.StorageProviderV1,
			Resource:   blockStorageResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.BlockStorageKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyStorageZonalMetadataStep(sCtx, expectedBlockMeta, blockResp.Metadata)

		verifyBlockStorageSpecStep(sCtx, params.BlockStorage.InitialSpec, blockResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*blockResp.Status.State)},
		)
	})

	t.WithNewStep("Get created block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "GetBlockStorage", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.BlockStorage.Name,
		}
		blockResp, err = suite.client.StorageV1.GetBlockStorage(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockResp)

		expectedBlockMeta.Verb = http.MethodGet
		verifyStorageZonalMetadataStep(sCtx, expectedBlockMeta, blockResp.Metadata)

		verifyBlockStorageSpecStep(sCtx, params.BlockStorage.InitialSpec, blockResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*blockResp.Status.State)},
		)
	})

	t.WithNewStep("Update block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "CreateOrUpdateBlockStorage", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.BlockStorage.Name,
		}
		blockResp, err = suite.client.StorageV1.CreateOrUpdateBlockStorage(ctx, wref, blockResp, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockResp)

		expectedBlockMeta.Verb = http.MethodPut
		verifyStorageZonalMetadataStep(sCtx, expectedBlockMeta, blockResp.Metadata)

		verifyBlockStorageSpecStep(sCtx, params.BlockStorage.UpdatedSpec, blockResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*blockResp.Status.State)},
		)
	})

	t.WithNewStep("Get updated block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "GetBlockStorage", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.BlockStorage.Name,
		}
		blockResp, err = suite.client.StorageV1.GetBlockStorage(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockResp)

		expectedBlockMeta.Verb = http.MethodGet
		verifyStorageZonalMetadataStep(sCtx, expectedBlockMeta, blockResp.Metadata)

		verifyBlockStorageSpecStep(sCtx, params.BlockStorage.UpdatedSpec, blockResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*blockResp.Status.State)},
		)
	})

	// Image
	var imageResp *storage.Image
	var expectedImageMeta *secalib.Metadata

	t.WithNewStep("Create image", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "CreateOrUpdateImage", "")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   params.Image.Name,
		}
		img := &storage.Image{
			Spec: storage.ImageSpec{
				BlockStorageRef: params.Image.InitialSpec.BlockStorageRef,
				CpuArchitecture: storage.ImageSpecCpuArchitecture(params.Image.InitialSpec.CpuArchitecture),
			},
		}
		imageResp, err = suite.client.StorageV1.CreateOrUpdateImage(ctx, tref, img, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, imageResp)

		expectedImageMeta = &secalib.Metadata{
			Name:       params.Image.Name,
			Provider:   secalib.StorageProviderV1,
			Resource:   imageResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.ImageKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyStorageRegionalMetadataStep(sCtx, expectedImageMeta, imageResp.Metadata)

		verifyImageSpecStep(sCtx, params.Image.InitialSpec, imageResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*imageResp.Status.State)},
		)
	})

	t.WithNewStep("Get created image", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "GetImage", "")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   params.Image.Name,
		}
		imageResp, err = suite.client.StorageV1.GetImage(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, imageResp)

		expectedImageMeta.Verb = http.MethodGet
		verifyStorageRegionalMetadataStep(sCtx, expectedImageMeta, imageResp.Metadata)

		verifyImageSpecStep(sCtx, params.Image.InitialSpec, imageResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*imageResp.Status.State)},
		)
	})

	t.WithNewStep("Update image", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "CreateOrUpdateImage", "")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   params.Image.Name,
		}
		imageResp, err = suite.client.StorageV1.CreateOrUpdateImage(ctx, tref, imageResp, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, imageResp)

		expectedImageMeta.Verb = http.MethodPut
		verifyStorageRegionalMetadataStep(sCtx, expectedImageMeta, imageResp.Metadata)

		verifyImageSpecStep(sCtx, params.Image.UpdatedSpec, imageResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*imageResp.Status.State)},
		)
	})

	t.WithNewStep("Get updated image", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "GetImage", "")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   params.Image.Name,
		}
		imageResp, err = suite.client.StorageV1.GetImage(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, imageResp)

		expectedImageMeta.Verb = http.MethodGet
		verifyStorageRegionalMetadataStep(sCtx, expectedImageMeta, imageResp.Metadata)

		verifyImageSpecStep(sCtx, params.Image.UpdatedSpec, imageResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*imageResp.Status.State)},
		)
	})

	t.WithNewStep("Delete image", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "DeleteImage", "")

		err = suite.client.StorageV1.DeleteImage(ctx, imageResp, nil)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get deleted image", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "GetImage", "")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   params.Image.Name,
		}
		_, err = suite.client.StorageV1.GetImage(ctx, tref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	t.WithNewStep("Delete block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "DeleteBlockStorage", params.Workspace.Name)

		err = suite.client.StorageV1.DeleteBlockStorage(ctx, blockResp, nil)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get deleted block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "GetBlockStorage", params.Workspace.Name)

		wref := secapi.WorkspaceReference{
			Tenant:    secapi.TenantID(suite.tenant),
			Workspace: secapi.WorkspaceID(params.Workspace.Name),
			Name:      params.BlockStorage.Name,
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
		stepCtx.Require().Equal(expected.SizeGB, actual.SizeGB, "SizeGB should match expected")
		stepCtx.Require().Equal(expected.SkuRef, actual.SkuRef, "SkuRef should match expected")
	})
}

func verifyImageSpecStep(ctx provider.StepCtx, expected *secalib.ImageSpecV1, actual storage.ImageSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(expected.BlockStorageRef, actual.BlockStorageRef, "BlockStorageRef should match expected")
		stepCtx.Require().Equal(expected.CpuArchitecture, string(actual.CpuArchitecture), "CpuArchitecture should match expected")
	})
}
