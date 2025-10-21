package secatest

import (
	"context"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// BlockStorage

func (suite *testSuite) createOrUpdateBlockStorageV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.StorageV1, role *schema.BlockStorage,
	expectedMeta *schema.RegionalWorkspaceResourceMetadata, expectedSpec *schema.BlockStorageSpec, expectedStatusState string,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "CreateOrUpdateBlockStorage")

		resp, err := api.CreateOrUpdateBlockStorage(ctx, role)
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

func (suite *testSuite) getBlockStorageV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.StorageV1, wref secapi.WorkspaceReference,
	expectedMeta *schema.RegionalWorkspaceResourceMetadata, expectedSpec *schema.BlockStorageSpec, expectedStatusState string,
) *schema.BlockStorage {
	var resp *schema.BlockStorage
	var err error

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetBlockStorage")

		resp, err = api.GetBlockStorage(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		if expectedMeta != nil {
			expectedMeta.Verb = http.MethodGet
			suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)
		}

		if expectedSpec != nil {
			suite.verifyBlockStorageSpecStep(sCtx, expectedSpec, &resp.Spec)
		}

		suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
	})
	return resp
}

func (suite *testSuite) getBlockStorageWithErrorV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.StorageV1, wref secapi.WorkspaceReference,
	expectedError error,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetBlockStorage")

		_, err := api.GetBlockStorage(ctx, wref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deleteBlockStorageV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.StorageV1, role *schema.BlockStorage) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "DeleteBlockStorage")

		err := api.DeleteBlockStorage(ctx, role)
		requireNoError(sCtx, err)
	})
}

// Image

func (suite *testSuite) createOrUpdateImageV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.StorageV1, role *schema.Image,
	expectedMeta *schema.RegionalResourceMetadata, expectedSpec *schema.ImageSpec, expectedStatusState string,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "CreateOrUpdateImage")

		resp, err := api.CreateOrUpdateImage(ctx, role)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		expectedMeta.Verb = http.MethodPut
		suite.verifyRegionalResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)

		suite.verifyImageSpecStep(sCtx, expectedSpec, &resp.Spec)

		suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
	})
}

func (suite *testSuite) getImageV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.StorageV1, tref secapi.TenantReference,
	expectedMeta *schema.RegionalResourceMetadata, expectedSpec *schema.ImageSpec, expectedStatusState string,
) *schema.Image {
	var resp *schema.Image
	var err error

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetImage")

		resp, err = api.GetImage(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		expectedMeta.Verb = http.MethodGet
		suite.verifyRegionalResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)

		suite.verifyImageSpecStep(sCtx, expectedSpec, &resp.Spec)

		suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
	})
	return resp
}

func (suite *testSuite) getImageWithErrorV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.StorageV1, tref secapi.TenantReference,
	expectedError error,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetImage")

		_, err := api.GetImage(ctx, tref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deleteImageV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.StorageV1, role *schema.Image) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "DeleteImage")

		err := api.DeleteImage(ctx, role)
		requireNoError(sCtx, err)
	})
}
