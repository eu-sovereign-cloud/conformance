package secatest

import (
	"context"
	"net/http"

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
	expectedState schema.ResourceState,
) {
	expectedMeta.Verb = http.MethodPut
	createOrUpdateWorkspaceResourceStep(
		t,
		ctx,
		suite,
		stepName,
		suite.setStorageWorkspaceV1StepParams,
		"CreateOrUpdateBlockStorage",
		resource.Metadata.Workspace,
		resource,
		func(context.Context, *schema.BlockStorage) (*stepFuncResponse[schema.BlockStorage, schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec], error) {
			resp, err := api.CreateOrUpdateBlockStorage(ctx, resource)
			return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
		},
		nil,
		expectedMeta,
		suite.verifyRegionalWorkspaceResourceMetadataStep,
		expectedSpec,
		suite.verifyBlockStorageSpecStep,
		expectedState,
	)
}

func (suite *testSuite) getBlockStorageV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.StorageV1,
	wref secapi.WorkspaceReference,
	expectedMeta *schema.RegionalWorkspaceResourceMetadata,
	expectedSpec *schema.BlockStorageSpec,
	expectedState schema.ResourceState,
) *schema.BlockStorage {
	var resp *schema.BlockStorage

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageWorkspaceV1StepParams(sCtx, "GetBlockStorage", string(wref.Workspace))
		retry := newStepResourceStateRetry(
			suite.baseDelay,
			suite.baseInterval,
			suite.maxAttempts,
			func() (schema.ResourceState, error) {
				var err error
				resp, err = api.GetBlockStorage(ctx, wref)
				requireNoError(sCtx, err)
				requireNotNilResponse(sCtx, resp)

				suite.requireNotNilStatus(sCtx, resp.Status)
				return *resp.Status.State, nil
			},
			func() {
				expectedMeta.Verb = http.MethodGet
				suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)

				suite.verifyBlockStorageSpecStep(sCtx, expectedSpec, &resp.Spec)

				suite.verifyStatusStep(sCtx, expectedState, *resp.Status.State)
			},
		)
		retry.run(sCtx, "GetBlockStorage", expectedState)
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
	expectedState schema.ResourceState,
) {
	expectedMeta.Verb = http.MethodPut
	createOrUpdateResourceStep(
		t,
		ctx,
		suite,
		stepName,
		suite.setStorageV1StepParams,
		"CreateOrUpdateImage",
		resource,
		func(context.Context, *schema.Image) (*stepFuncResponse[schema.Image, schema.RegionalResourceMetadata, schema.ImageSpec], error) {
			resp, err := api.CreateOrUpdateImage(ctx, resource)
			return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
		},
		nil,
		expectedMeta,
		suite.verifyRegionalResourceMetadataStep,
		expectedSpec,
		suite.verifyImageSpecStep,
		expectedState,
	)
}

func (suite *testSuite) getImageV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.StorageV1,
	tref secapi.TenantReference,
	expectedMeta *schema.RegionalResourceMetadata,
	expectedSpec *schema.ImageSpec,
	expectedState schema.ResourceState,
) *schema.Image {
	var resp *schema.Image

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "GetImage")
		retry := newStepResourceStateRetry(
			suite.baseDelay,
			suite.baseInterval,
			suite.maxAttempts,
			func() (schema.ResourceState, error) {
				var err error
				resp, err = api.GetImage(ctx, tref)
				requireNoError(sCtx, err)
				requireNotNilResponse(sCtx, resp)

				suite.requireNotNilStatus(sCtx, resp.Status)
				return *resp.Status.State, nil
			},
			func() {
				expectedMeta.Verb = http.MethodGet
				suite.verifyRegionalResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)

				suite.verifyImageSpecStep(sCtx, expectedSpec, &resp.Spec)

				suite.verifyStatusStep(sCtx, expectedState, *resp.Status.State)
			},
		)
		retry.run(sCtx, "GetImage", expectedState)
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
		suite.setStorageV1StepParams(sCtx, "GetImage")

		_, err := api.GetImage(ctx, tref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deleteImageV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.StorageV1, resource *schema.Image) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "DeleteImage")

		err := api.DeleteImage(ctx, resource)
		requireNoError(sCtx, err)
	})
}
