package secatest

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	storage "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.storage.v1"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

const (
	// Test constants
	storageSkuMinVolumeSize10 = 10
	storageSkuIops500         = 500
	blockStorageSizeGB100     = 100
	blockStorageSizeGB200     = 200
)

type StorageV1TestSuite struct {
	regionalTestSuite
}

type verifyStorageSkuLabelsStepParams struct {
	provider string
	tier     string
}

type verifyStorageSkuSpecStepParams struct {
	iops          int
	storageType   string
	minVolumeSize int
}

type verifyBlockStorageSpecStepParams struct {
	sizeGB int
	skuRef string
}

type verifyImageSpecStepParams struct {
	blockStorageRef string
	cpuArchitecture string
}

func (suite *StorageV1TestSuite) TestStorageV1(t provider.T) {
	t.Title("Storage Lifecycle Test")
	configureTags(t, storageV1Provider, storageSkuKind, blockStorageKind, imageKind, instanceSkuKind)

	// TODO Export to configuration
	skuProvider := "seca"
	skuTier := "gold"

	// TODO Create before the scenario
	workspaceName := suite.generateWorkspaceName()

	// Generate scenario data
	storageSkuName := suite.generateSkuName()
	storageSkuRef := suite.generateSkuRef(storageSkuName)

	blockStorageName := suite.generateBlockStorageName()
	blockStorageResource := suite.generateBlockStorageResource(workspaceName, blockStorageName)
	blockStorageRef := suite.generateBlockStorageRef(blockStorageName)

	imageName := suite.generateImageName()
	imageResource := suite.generateImageResource(imageName)

	// Setup mock, if configured to use
	if suite.isMockEnabled() {
		wm, err := mock.CreateStorageLifecycleScenarioV1("Storage Lifecycle",
			mock.StorageParamsV1{
				Params: mock.Params{
					MockURL:   suite.mockServerURL,
					AuthToken: suite.authToken,
					Tenant:    suite.tenant,
					Workspace: workspaceName,
					Region:    suite.region,
				},
				Sku: mock.StorageSkuParamsV1{
					Name:          storageSkuName,
					Provider:      skuProvider,
					Tier:          skuTier,
					Iops:          storageSkuIops500,
					StorageType:   storageTypeLocalEphemeral,
					MinVolumeSize: storageSkuMinVolumeSize10,
				},
				BlockStorage: mock.BlockStorageParamsV1{
					Name:          blockStorageName,
					SkuRef:        storageSkuRef,
					SizeGBInitial: blockStorageSizeGB100,
					SizeGBUpdated: blockStorageSizeGB200,
				},
				Image: mock.ImageParamsV1{
					Name:                   imageName,
					BlockStorageRef:        blockStorageRef,
					CpuArchitectureInitial: cpuArchitectureAmd64,
					CpuArchitectureUpdated: cpuArchitectureArm64,
				},
			})
		if err != nil {
			slog.Error("Failed to create storage scenario", "error", err)
			return
		}
		suite.mockClient = wm
	}

	ctx := context.Background()
	var skuResp *storage.StorageSku
	var blockResp *storage.BlockStorage
	var imageResp *storage.Image
	var err error

	// Step 1
	t.WithNewStep("Get storage sku", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetStorageSku",
			tenantStepParameter, suite.tenant,
		)

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   storageSkuName,
		}
		skuResp, err = suite.client.StorageV1.GetSku(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, skuResp)

		verifySkuMetadataStep(sCtx, verifySkuMetadataStepParams{name: storageSkuName}, verifySkuMetadataStepParams{name: skuResp.Metadata.Name})

		verifyLabelsParams := verifyStorageSkuLabelsStepParams{
			provider: skuProvider,
			tier:     skuTier,
		}
		verifyStorageSkuLabelsStep(sCtx, verifyLabelsParams, *skuResp.Labels)

		verifySpecParams := verifyStorageSkuSpecStepParams{
			iops:          storageSkuIops500,
			storageType:   storageTypeLocalEphemeral,
			minVolumeSize: storageSkuMinVolumeSize10,
		}
		verifyStorageSkuSpecStep(sCtx, verifySpecParams, skuResp.Spec)
	})

	expectedBlockMetadata := verifyRegionalMetadataStepParams{
		name:       blockStorageName,
		provider:   storageV1Provider,
		resource:   blockStorageResource,
		apiVersion: version1,
		kind:       blockStorageKind,
		tenant:     suite.tenant,
		region:     suite.region,
	}

	// Step 2
	t.WithNewStep("Create block storage", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateBlockStorage",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		bo := &storage.BlockStorage{
			Metadata: &storage.ZonalResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: &workspaceName,
				Name:      blockStorageName,
			},
			Spec: storage.BlockStorageSpec{
				SizeGB: blockStorageSizeGB100,
				SkuRef: storageSkuRef,
			},
		}
		blockResp, err = suite.client.StorageV1.CreateOrUpdateBlockStorage(ctx, bo)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockResp)

		expectedBlockMetadata.verb = http.MethodPut
		verifyStorageZonalMetadataStep(sCtx, expectedBlockMetadata, blockResp.Metadata)

		verifySpecParams := verifyBlockStorageSpecStepParams{
			sizeGB: blockStorageSizeGB100,
			skuRef: storageSkuRef,
		}
		verifyBlockStorageSpecStep(sCtx, verifySpecParams, blockResp.Spec)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: creatingStatusState,
			actualState:   string(*blockResp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
	})

	// Step 3
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

		expectedBlockMetadata.verb = http.MethodGet
		verifyStorageZonalMetadataStep(sCtx, expectedBlockMetadata, blockResp.Metadata)

		verifySpecParams := verifyBlockStorageSpecStepParams{
			sizeGB: blockStorageSizeGB100,
			skuRef: storageSkuRef,
		}
		verifyBlockStorageSpecStep(sCtx, verifySpecParams, blockResp.Spec)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: activeStatusState,
			actualState:   string(*blockResp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
	})

	// Step 4
	t.WithNewStep("Update block storage", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateBlockStorage",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		blockResp.Spec.SizeGB = blockStorageSizeGB200
		blockResp, err = suite.client.StorageV1.CreateOrUpdateBlockStorage(ctx, blockResp)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockResp)

		expectedBlockMetadata.verb = http.MethodPut
		verifyStorageZonalMetadataStep(sCtx, expectedBlockMetadata, blockResp.Metadata)

		verifySpecParams := verifyBlockStorageSpecStepParams{
			sizeGB: blockStorageSizeGB200,
			skuRef: storageSkuRef,
		}
		verifyBlockStorageSpecStep(sCtx, verifySpecParams, blockResp.Spec)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: updatingStatusState,
			actualState:   string(*blockResp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
	})

	// Step 5
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

		expectedBlockMetadata.verb = http.MethodGet
		verifyStorageZonalMetadataStep(sCtx, expectedBlockMetadata, blockResp.Metadata)

		verifySpecParams := verifyBlockStorageSpecStepParams{
			sizeGB: blockStorageSizeGB200,
			skuRef: storageSkuRef,
		}
		verifyBlockStorageSpecStep(sCtx, verifySpecParams, blockResp.Spec)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: activeStatusState,
			actualState:   string(*blockResp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
	})

	expectedImageMetadata := verifyRegionalMetadataStepParams{
		name:       imageName,
		provider:   storageV1Provider,
		resource:   imageResource,
		apiVersion: version1,
		kind:       imageKind,
		tenant:     suite.tenant,
		region:     suite.region,
	}

	// Step 6
	t.WithNewStep("Create image", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateImage",
			tenantStepParameter, suite.tenant,
		)

		img := &storage.Image{
			Metadata: &storage.RegionalResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: &workspaceName,
				Name:      imageName,
			},
			Spec: storage.ImageSpec{
				BlockStorageRef: blockStorageRef,
				CpuArchitecture: cpuArchitectureAmd64,
			},
		}
		imageResp, err = suite.client.StorageV1.CreateOrUpdateImage(ctx, img)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, imageResp)

		expectedImageMetadata.verb = http.MethodPut
		verifyStorageRegionalMetadataStep(sCtx, expectedImageMetadata, imageResp.Metadata)

		verifySpecParams := verifyImageSpecStepParams{
			blockStorageRef: blockStorageRef,
			cpuArchitecture: cpuArchitectureAmd64,
		}
		verifyImageSpecStep(sCtx, verifySpecParams, imageResp.Spec)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: creatingStatusState,
			actualState:   string(*imageResp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
	})

	// Step 7
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

		expectedImageMetadata.verb = http.MethodGet
		verifyStorageRegionalMetadataStep(sCtx, expectedImageMetadata, imageResp.Metadata)

		verifySpecParams := verifyImageSpecStepParams{
			blockStorageRef: blockStorageRef,
			cpuArchitecture: cpuArchitectureAmd64,
		}
		verifyImageSpecStep(sCtx, verifySpecParams, imageResp.Spec)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: activeStatusState,
			actualState:   string(*imageResp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
	})

	// Step 8
	t.WithNewStep("Update image", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateImage",
			tenantStepParameter, suite.tenant,
		)

		imageResp.Spec.CpuArchitecture = cpuArchitectureArm64
		imageResp, err = suite.client.StorageV1.CreateOrUpdateImage(ctx, imageResp)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, imageResp)

		expectedImageMetadata.verb = http.MethodPut
		verifyStorageRegionalMetadataStep(sCtx, expectedImageMetadata, imageResp.Metadata)

		verifySpecParams := verifyImageSpecStepParams{
			blockStorageRef: blockStorageRef,
			cpuArchitecture: cpuArchitectureArm64,
		}
		verifyImageSpecStep(sCtx, verifySpecParams, imageResp.Spec)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: updatingStatusState,
			actualState:   string(*imageResp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
	})

	// Step 9
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

		expectedImageMetadata.verb = http.MethodGet
		verifyStorageRegionalMetadataStep(sCtx, expectedImageMetadata, imageResp.Metadata)

		verifySpecParams := verifyImageSpecStepParams{
			blockStorageRef: blockStorageRef,
			cpuArchitecture: cpuArchitectureArm64,
		}
		verifyImageSpecStep(sCtx, verifySpecParams, imageResp.Spec)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: activeStatusState,
			actualState:   string(*imageResp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
	})

	// Step 10
	t.WithNewStep("Delete image", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "DeleteImage",
			tenantStepParameter, suite.tenant,
		)

		err = suite.client.StorageV1.DeleteImage(ctx, imageResp)
		requireNoError(sCtx, err)
	})

	// Step 11
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
		imageResp, err = suite.client.StorageV1.GetImage(ctx, tref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	// Step 12
	t.WithNewStep("Delete block storage", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "DeleteBlockStorage",
			tenantStepParameter, suite.tenant,
		)

		err = suite.client.StorageV1.DeleteBlockStorage(ctx, blockResp)
		requireNoError(sCtx, err)
	})

	// Step 13
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
		blockResp, err = suite.client.StorageV1.GetBlockStorage(ctx, wref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})
}

func (suite *StorageV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}

func verifyStorageZonalMetadataStep(ctx provider.StepCtx, expected verifyRegionalMetadataStepParams, metadata *storage.ZonalResourceMetadata) {
	actualMetadata := verifyRegionalMetadataStepParams{
		name:       metadata.Name,
		provider:   metadata.Provider,
		verb:       metadata.Verb,
		resource:   metadata.Resource,
		apiVersion: metadata.ApiVersion,
		kind:       string(metadata.Kind),
		tenant:     metadata.Tenant,
		workspace:  *metadata.Workspace,
		region:     metadata.Region,
	}
	verifyRegionalMetadataStep(ctx, expected, actualMetadata)
}

func verifyStorageRegionalMetadataStep(ctx provider.StepCtx, expected verifyRegionalMetadataStepParams, metadata *storage.RegionalResourceMetadata) {
	actualMetadata := verifyRegionalMetadataStepParams{
		name:       metadata.Name,
		provider:   metadata.Provider,
		verb:       metadata.Verb,
		resource:   metadata.Resource,
		apiVersion: metadata.ApiVersion,
		kind:       string(metadata.Kind),
		tenant:     metadata.Tenant,
		workspace:  *metadata.Workspace,
		region:     metadata.Region,
	}
	verifyRegionalMetadataStep(ctx, expected, actualMetadata)
}

func verifyStorageSkuLabelsStep(ctx provider.StepCtx, expected verifyStorageSkuLabelsStepParams, labels map[string]string) {
	ctx.WithNewStep("Verify labels", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters(
			"expected_provider", expected.provider,
			"actual_provider", labels[providerLabel],

			"expected_tier", expected.tier,
			"actual_tier", labels[tierLabel],
		)
		stepCtx.Require().Equal(expected.provider, labels[providerLabel], "Provider should match expected")
		stepCtx.Require().Equal(expected.tier, labels[tierLabel], "Tier should match expected")
	})
}

func verifyStorageSkuSpecStep(ctx provider.StepCtx, expected verifyStorageSkuSpecStepParams, actual *storage.StorageSkuSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters(
			"expected_iops", expected.iops,
			"actual_iops", actual.Iops,

			"expected_type", expected.storageType,
			"actual_type", actual.Type,

			"expected_minVolumeSize", expected.minVolumeSize,
			"actual_minVolumeSize", actual.MinVolumeSize,
		)
		stepCtx.Require().Equal(expected.iops, actual.Iops, "Iops should match expected")
		stepCtx.Require().Equal(expected.storageType, string(actual.Type), "Type should match expected")
		stepCtx.Require().Equal(expected.minVolumeSize, actual.MinVolumeSize, "MinVolumeSize should match expected")
	})
}

func verifyBlockStorageSpecStep(ctx provider.StepCtx, expected verifyBlockStorageSpecStepParams, actual storage.BlockStorageSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters(
			"expected_sizeGB", expected.sizeGB,
			"actual_sizeGB", actual.SizeGB,

			"expected_skuRef", expected.skuRef,
			"actual_skuRef", actual.SkuRef,
		)
		stepCtx.Require().Equal(expected.sizeGB, actual.SizeGB, "SizeGB should match expected")
		stepCtx.Require().Equal(expected.skuRef, actual.SkuRef, "SkuRef should match expected")
	})
}

func verifyImageSpecStep(ctx provider.StepCtx, expected verifyImageSpecStepParams, actual storage.ImageSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.WithNewParameters(
			"expected_blockStorageRef", expected.blockStorageRef,
			"actual_blockStorageRef", actual.BlockStorageRef,

			"expected_cpuArchitecture", expected.cpuArchitecture,
			"actual_cpuArchitecture", actual.CpuArchitecture,
		)
		stepCtx.Require().Equal(expected.blockStorageRef, actual.BlockStorageRef, "BlockStorageRef should match expected")
		stepCtx.Require().Equal(expected.cpuArchitecture, string(actual.CpuArchitecture), "CpuArchitecture should match expected")
	})
}
