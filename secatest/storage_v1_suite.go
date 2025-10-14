package secatest

import (
	"context"
	"fmt"
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
	slog.Info("Starting " + suite.scenarioName)

	t.Title(suite.scenarioName)
	configureTags(t, secalib.StorageProviderV1, secalib.BlockStorageKind, secalib.ImageKind)

	// Select sku
	storageSkuName := suite.storageSkus[rand.Intn(len(suite.storageSkus))]

	// Generate scenario data
	workspaceName := secalib.GenerateWorkspaceName()

	storageSkuRef := secalib.GenerateSkuRef(storageSkuName)

	blockStorageName := secalib.GenerateBlockStorageName()
	blockStorageResource := secalib.GenerateBlockStorageResource(suite.tenant, workspaceName, blockStorageName)
	blockStorageRef := secalib.GenerateBlockStorageRef(blockStorageName)

	imageName := secalib.GenerateImageName()
	imageResource := secalib.GenerateImageResource(suite.tenant, imageName)

	initialStorageSize := secalib.GenerateBlockStorageSize()
	updatedStorageSize := secalib.GenerateBlockStorageSize()

	// Setup mock, if configured to use
	if suite.mockEnabled {
		wm, err := mock.CreateStorageLifecycleScenarioV1(suite.scenarioName, &mock.StorageParamsV1{
			Params: &mock.Params{
				MockURL:   *suite.mockServerURL,
				AuthToken: suite.authToken,
				Tenant:    suite.tenant,
				Region:    suite.region,
			},
			Workspace: &mock.ResourceParams[secalib.WorkspaceSpecV1]{
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
			t.Fatalf("Failed to create wiremock scenario: %v", err)
		}
		suite.mockClient = wm
	}

	ctx := context.Background()
	var workResp *schema.Workspace
	var blockResp *schema.BlockStorage
	var imageResp *schema.Image
	var err error

	t.WithNewStep("Create workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "CreateOrUpdateWorkspace")

		ws := &schema.Workspace{
			Metadata: &schema.RegionalResourceMetadata{
				Tenant: suite.tenant,
				Name:   workspaceName,
			},
		}
		workResp, err = suite.client.WorkspaceV1.CreateOrUpdateWorkspace(ctx, ws)
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
			Name:   workspaceName,
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
	var expectedBlockMeta *secalib.Metadata
	var expectedBlockSpec *secalib.BlockStorageSpecV1

	t.WithNewStep("Create block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "CreateOrUpdateBlockStorage", workspaceName)

		storageSkuURN, err := secapi.BuildReferenceFromURN(storageSkuRef)
		if err != nil {
			t.Fatal(err)
		}

		bo := &schema.BlockStorage{
			Metadata: &schema.RegionalWorkspaceResourceMetadata{
				Tenant:    suite.tenant,
				Workspace: workspaceName,
				Name:      blockStorageName,
			},
			Spec: schema.BlockStorageSpec{
				SizeGB: initialStorageSize,
				SkuRef: *storageSkuURN,
			},
		}
		blockResp, err = suite.client.StorageV1.CreateOrUpdateBlockStorage(ctx, bo)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockResp)

		expectedBlockMeta = &secalib.Metadata{
			Name:       blockStorageName,
			Provider:   secalib.StorageProviderV1,
			Resource:   blockStorageResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.BlockStorageKind,
			Tenant:     suite.tenant,
			Region:     &suite.region,
		}

		verifyStorageWorkspaceMetadataStep(sCtx, expectedBlockMeta, blockResp.Metadata)

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
		verifyStorageWorkspaceMetadataStep(sCtx, expectedBlockMeta, blockResp.Metadata)

		verifyBlockStorageSpecStep(sCtx, expectedBlockSpec, blockResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*blockResp.Status.State)},
		)
	})

	t.WithNewStep("Update block storage", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "CreateOrUpdateBlockStorage", workspaceName)

		blockResp, err = suite.client.StorageV1.CreateOrUpdateBlockStorage(ctx, blockResp)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, blockResp)

		expectedBlockMeta.Verb = http.MethodPut
		verifyStorageWorkspaceMetadataStep(sCtx, expectedBlockMeta, blockResp.Metadata)

		expectedBlockSpec.SizeGB = updatedStorageSize
		verifyBlockStorageSpecStep(sCtx, expectedBlockSpec, blockResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*blockResp.Status.State)},
		)
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
		verifyStorageWorkspaceMetadataStep(sCtx, expectedBlockMeta, blockResp.Metadata)

		verifyBlockStorageSpecStep(sCtx, expectedBlockSpec, blockResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*blockResp.Status.State)},
		)
	})

	// Image
	var expectedImageMeta *secalib.Metadata
	var expectedImageSpec *secalib.ImageSpecV1

	t.WithNewStep("Create image", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "CreateOrUpdateImage", "")

		blockStorageURN, err := secapi.BuildReferenceFromURN(blockStorageRef)
		if err != nil {
			t.Fatal(err)
		}

		img := &schema.Image{
			Metadata: &schema.RegionalResourceMetadata{
				Tenant: suite.tenant,
				Name:   imageName,
			},
			Spec: schema.ImageSpec{
				BlockStorageRef: *blockStorageURN,
				CpuArchitecture: secalib.CpuArchitectureAmd64,
			},
		}
		imageResp, err = suite.client.StorageV1.CreateOrUpdateImage(ctx, img)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, imageResp)

		expectedImageMeta = &secalib.Metadata{
			Name:       imageName,
			Provider:   secalib.StorageProviderV1,
			Resource:   imageResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.ImageKind,
			Tenant:     suite.tenant,
			Region:     &suite.region,
		}
		verifyStorageRegionalMetadataStep(sCtx, expectedImageMeta, imageResp.Metadata)

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
		verifyStorageRegionalMetadataStep(sCtx, expectedImageMeta, imageResp.Metadata)

		verifyImageSpecStep(sCtx, expectedImageSpec, imageResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*imageResp.Status.State)},
		)
	})

	t.WithNewStep("Update image", func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "CreateOrUpdateImage", "")

		imageResp, err = suite.client.StorageV1.CreateOrUpdateImage(ctx, imageResp)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, imageResp)

		expectedImageMeta.Verb = http.MethodPut
		verifyStorageRegionalMetadataStep(sCtx, expectedImageMeta, imageResp.Metadata)

		expectedImageSpec.CpuArchitecture = secalib.CpuArchitectureArm64
		verifyImageSpecStep(sCtx, expectedImageSpec, imageResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*imageResp.Status.State)},
		)
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
		verifyStorageRegionalMetadataStep(sCtx, expectedImageMeta, imageResp.Metadata)

		verifyImageSpecStep(sCtx, expectedImageSpec, imageResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*imageResp.Status.State)},
		)
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

func verifyStorageRegionalMetadataStep(ctx provider.StepCtx, expected *secalib.Metadata, metadata *schema.RegionalResourceMetadata) {
	actualMetadata := &secalib.Metadata{
		Name:       metadata.Name,
		Provider:   metadata.Provider,
		Verb:       metadata.Verb,
		Resource:   metadata.Resource,
		ApiVersion: metadata.ApiVersion,
		Kind:       string(metadata.Kind),
		Tenant:     metadata.Tenant,
		Region:     &metadata.Region,
	}
	verifyRegionalMetadataStep(ctx, expected, actualMetadata)
}

func verifyStorageWorkspaceMetadataStep(ctx provider.StepCtx, expected *secalib.Metadata, metadata *schema.RegionalWorkspaceResourceMetadata) {
	actualMetadata := &secalib.Metadata{
		Name:       metadata.Name,
		Provider:   metadata.Provider,
		Verb:       metadata.Verb,
		Resource:   metadata.Resource,
		ApiVersion: metadata.ApiVersion,
		Kind:       string(metadata.Kind),
		Tenant:     metadata.Tenant,
		Workspace:  &metadata.Workspace,
		Region:     &metadata.Region,
	}
	verifyRegionalMetadataStep(ctx, expected, actualMetadata)
}

func verifyBlockStorageSpecStep(ctx provider.StepCtx, expected *secalib.BlockStorageSpecV1, actual schema.BlockStorageSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(expected.SizeGB, actual.SizeGB, "SizeGB should match expected")

		skuRef, err := asStorageReferenceURN(actual.SkuRef)
		if err != nil {
			ctx.Error(err)
		}
		stepCtx.Require().Equal(expected.SkuRef, skuRef, "SkuRef should match expected")
	})
}

func verifyImageSpecStep(ctx provider.StepCtx, expected *secalib.ImageSpecV1, actual schema.ImageSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		blockStorageRef, err := asStorageReferenceURN(actual.BlockStorageRef)
		if err != nil {
			ctx.Error(err)
		}
		stepCtx.Require().Equal(expected.BlockStorageRef, blockStorageRef, "BlockStorageRef should match expected")

		stepCtx.Require().Equal(expected.CpuArchitecture, string(actual.CpuArchitecture), "CpuArchitecture should match expected")
	})
}

func asStorageReferenceURN(ref schema.Reference) (string, error) {
	urn, err := ref.AsReferenceURN()
	if err != nil {
		return "", fmt.Errorf("error extracting URN from reference: %w", err)
	}
	return string(urn), nil
}
