package secatest

import (
	"context"
	"net/http"
	"time"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// BlockStorage

func (suite *testSuite) createOrUpdateBlockStorageV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.StorageV1,
	resource *schema.BlockStorage,
	expectedMeta *schema.RegionalWorkspaceResourceMetadata,
	expectedSpec *schema.BlockStorageSpec,
	expectedStatusState string,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageWorkspaceV1StepParams(sCtx, "CreateOrUpdateBlockStorage", resource.Metadata.Workspace)

		resp, err := api.CreateOrUpdateBlockStorage(ctx, resource)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		if expectedMeta != nil {
			expectedMeta.Verb = http.MethodPut
			suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)
		}

		if expectedSpec != nil {
			suite.verifyBlockStorageSpecStep(sCtx, expectedSpec, &resp.Spec)
		}

		suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)

	})
}

func (suite *testSuite) getBlockStorageV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.StorageV1,
	wref secapi.WorkspaceReference,
	expectedMeta *schema.RegionalWorkspaceResourceMetadata,
	expectedSpec *schema.BlockStorageSpec,
	expectedStatusState string,
) *schema.BlockStorage {
	var resp *schema.BlockStorage
	var err error

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageWorkspaceV1StepParams(sCtx, "GetBlockStorage", string(wref.Workspace))
		time.Sleep(time.Duration(suite.initialDelay) * time.Second)
		for attempt := 1; attempt <= suite.maxAttempts; attempt++ {
			resp, err = api.GetBlockStorage(ctx, wref)
			requireNoError(sCtx, err)
			requireNotNilResponse(sCtx, resp)
			if resp.Status.State != nil && *resp.Status.State == secalib.ActiveResourceState {
				if expectedMeta != nil {
					expectedMeta.Verb = http.MethodPut
					suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)
				}

				if expectedSpec != nil {
					suite.verifyBlockStorageSpecStep(sCtx, expectedSpec, &resp.Spec)
				}

				suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
				return
			} else {
				time.Sleep(time.Duration(suite.baseInterval) * time.Second)
			}

		}
	})
	return resp

}

func (suite *testSuite) getBlockStorageWithErrorV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.StorageV1,
	wref secapi.WorkspaceReference,
	expectedError error,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageWorkspaceV1StepParams(sCtx, "GetBlockStorage", string(wref.Workspace))

		_, err := api.GetBlockStorage(ctx, wref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deleteBlockStorageV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.StorageV1, resource *schema.BlockStorage) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageWorkspaceV1StepParams(sCtx, "DeleteBlockStorage", resource.Metadata.Workspace)

		err := api.DeleteBlockStorage(ctx, resource)
		requireNoError(sCtx, err)
	})
}

// Image

func (suite *testSuite) createOrUpdateImageV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.StorageV1,
	resource *schema.Image,
	expectedMeta *schema.RegionalResourceMetadata,
	expectedSpec *schema.ImageSpec,
	expectedStatusState string,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageTenantV1StepParams(sCtx, "CreateOrUpdateImage")

		resp, err := api.CreateOrUpdateImage(ctx, resource)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		expectedMeta.Verb = http.MethodPut
		suite.verifyRegionalResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)

		suite.verifyImageSpecStep(sCtx, expectedSpec, &resp.Spec)

		suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
	})
}

func (suite *testSuite) getImageV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.StorageV1,
	tref secapi.TenantReference,
	expectedMeta *schema.RegionalResourceMetadata, expectedSpec *schema.ImageSpec, expectedStatusState string,
) *schema.Image {
	var resp *schema.Image
	var err error

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageTenantV1StepParams(sCtx, "GetImage")
		time.Sleep(time.Duration(suite.initialDelay) * time.Second)
		for attempt := 1; attempt <= suite.maxAttempts; attempt++ {
			resp, err = api.GetImage(ctx, tref)
			requireNoError(sCtx, err)
			requireNotNilResponse(sCtx, resp)
			if resp.Status.State != nil && *resp.Status.State == secalib.ActiveResourceState {

				expectedMeta.Verb = http.MethodGet
				suite.verifyRegionalResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)

				suite.verifyImageSpecStep(sCtx, expectedSpec, &resp.Spec)

				suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
				return
			}
		}
	})
	return resp
}

func (suite *testSuite) getImageWithErrorV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.StorageV1,
	tref secapi.TenantReference,
	expectedError error,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageTenantV1StepParams(sCtx, "GetImage")

		_, err := api.GetImage(ctx, tref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deleteImageV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.StorageV1, resource *schema.Image) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageTenantV1StepParams(sCtx, "DeleteImage")

		err := api.DeleteImage(ctx, resource)
		requireNoError(sCtx, err)
	})
}
