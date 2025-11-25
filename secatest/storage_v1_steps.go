package secatest

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/eu-sovereign-cloud/go-sdk/secapi/builders"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// BlockStorage

func (suite *testSuite) createOrUpdateBlockStorageV1Step(stepName string, t provider.T, api *secapi.StorageV1, resource *schema.BlockStorage,
	responseExpects responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec],
) {
	responseExpects.metadata.Verb = http.MethodPut
	createOrUpdateWorkspaceResourceStep(t, suite,
		createOrUpdateWorkspaceResourceParams[schema.BlockStorage, schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			stepName:       stepName,
			stepParamsFunc: suite.setStorageWorkspaceV1StepParams,
			operationName:  "CreateOrUpdateBlockStorage",
			workspace:      resource.Metadata.Workspace,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.BlockStorage) (*stepFuncResponse[schema.BlockStorage, schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec], error) {
				resp, err := api.CreateOrUpdateBlockStorage(t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.metadata,
			verifyMetadataFunc:    suite.verifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.spec,
			verifySpecFunc:        suite.verifyBlockStorageSpecStep,
			expectedResourceState: responseExpects.resourceState,
		},
	)
}

func (suite *testSuite) getBlockStorageV1Step(stepName string, t provider.T, api *secapi.StorageV1, wref secapi.WorkspaceReference,
	responseExpects responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec],
) *schema.BlockStorage {
	responseExpects.metadata.Verb = http.MethodGet
	return getWorkspaceResourceStep(t, suite,
		getWorkspaceResourceParams[schema.BlockStorage, schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			stepName:       stepName,
			stepParamsFunc: suite.setStorageWorkspaceV1StepParams,
			operationName:  "GetBlockStorage",
			wref:           wref,
			getFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[schema.BlockStorage, schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec], error) {
				resp, err := api.GetBlockStorageUntilState(ctx, wref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.metadata,
			verifyMetadataFunc:    suite.verifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.spec,
			verifySpecFunc:        suite.verifyBlockStorageSpecStep,
			expectedResourceState: responseExpects.resourceState,
		},
	)
}

func (suite *testSuite) getListBlockStorageV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.StorageV1,
	tref secapi.TenantReference,
	wref secapi.WorkspaceReference,
	opts *builders.ListOptions,
) []*schema.BlockStorage {
	var respNext []*schema.BlockStorage
	var respAll []*schema.BlockStorage

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageWorkspaceV1StepParams(sCtx, "GetListBlockStorage", string(wref.Workspace))

		var iter *secapi.Iterator[schema.BlockStorage]
		var err error
		if opts != nil {
			iter, err = api.ListBlockStoragesWithFilters(ctx, tref.Tenant, wref.Workspace, opts)
		} else {
			iter, err = api.ListBlockStorages(ctx, tref.Tenant, wref.Workspace)
		}
		requireNoError(sCtx, err)
		for {
			item, err := iter.Next(context.Background())
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				break
			}
			respNext = append(respNext, item)
		}
		requireNotNilResponse(sCtx, respNext)
		requireLenResponse(sCtx, len(respNext))
		/*
			respAll, err = iter.All(ctx)
			requireNoError(sCtx, err)
			requireNotNilResponse(sCtx, respAll)
			requireLenResponse(sCtx, len(respAll))

			compareIteratorsResponse(sCtx, len(respNext), len(respAll))
		*/
	})
	return respAll
}

func (suite *testSuite) getBlockStorageWithErrorV1Step(stepName string, t provider.T, api *secapi.StorageV1, wref secapi.WorkspaceReference, expectedError error) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageWorkspaceV1StepParams(sCtx, "GetBlockStorage", string(wref.Workspace))

		_, err := api.GetBlockStorage(t.Context(), wref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deleteBlockStorageV1Step(stepName string, t provider.T, api *secapi.StorageV1, resource *schema.BlockStorage) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageWorkspaceV1StepParams(sCtx, "DeleteBlockStorage", resource.Metadata.Workspace)

		err := api.DeleteBlockStorage(t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

// Image

func (suite *testSuite) createOrUpdateImageV1Step(stepName string, t provider.T, api *secapi.StorageV1, resource *schema.Image,
	responseExpects responseExpects[schema.RegionalResourceMetadata, schema.ImageSpec],
) {
	responseExpects.metadata.Verb = http.MethodPut
	createOrUpdateTenantResourceStep(t, suite,
		createOrUpdateTenantResourceParams[schema.Image, schema.RegionalResourceMetadata, schema.ImageSpec]{
			stepName:       stepName,
			stepParamsFunc: suite.setStorageV1StepParams,
			operationName:  "CreateOrUpdateImage",
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.Image) (*stepFuncResponse[schema.Image, schema.RegionalResourceMetadata, schema.ImageSpec], error) {
				resp, err := api.CreateOrUpdateImage(t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.metadata,
			verifyMetadataFunc:    suite.verifyRegionalResourceMetadataStep,
			expectedSpec:          responseExpects.spec,
			verifySpecFunc:        suite.verifyImageSpecStep,
			expectedResourceState: responseExpects.resourceState,
		},
	)
}

func (suite *testSuite) getImageV1Step(stepName string, t provider.T, api *secapi.StorageV1, tref secapi.TenantReference,
	responseExpects responseExpects[schema.RegionalResourceMetadata, schema.ImageSpec],
) *schema.Image {
	responseExpects.metadata.Verb = http.MethodGet
	return getTenantResourceStep(t, suite,
		getTenantResourceParams[schema.Image, schema.RegionalResourceMetadata, schema.ImageSpec]{
			stepName:       stepName,
			stepParamsFunc: suite.setStorageV1StepParams,
			operationName:  "GetImage",
			tref:           tref,
			getFunc: func(ctx context.Context, tref secapi.TenantReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[schema.Image, schema.RegionalResourceMetadata, schema.ImageSpec], error) {
				resp, err := api.GetImageUntilState(ctx, tref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.metadata,
			verifyMetadataFunc:    suite.verifyRegionalResourceMetadataStep,
			expectedSpec:          responseExpects.spec,
			verifySpecFunc:        suite.verifyImageSpecStep,
			expectedResourceState: responseExpects.resourceState,
		},
	)
}

func (suite *testSuite) getListImageV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.StorageV1,
	tref secapi.TenantReference,
	opts *builders.ListOptions,
) []*schema.Image {
	var respNext []*schema.Image
	var respAll []*schema.Image
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageWorkspaceV1StepParams(sCtx, "GetListImage", string(tref.Name))
		var iter *secapi.Iterator[schema.Image]
		var err error
		if opts != nil {
			iter, err = api.ListImagesWithFilters(ctx, tref.Tenant, opts)
		} else {
			iter, err = api.ListImages(ctx, tref.Tenant)
		}
		requireNoError(sCtx, err)
		for {
			item, err := iter.Next(context.Background())
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				break
			}
			respNext = append(respNext, item)
		}
		requireNotNilResponse(sCtx, respNext)
		requireLenResponse(sCtx, len(respNext))
		/*
			respAll, err = iter.All(ctx)
			requireNoError(sCtx, err)
			requireNotNilResponse(sCtx, respAll)
			requireLenResponse(sCtx, len(respAll))

			compareIteratorsResponse(sCtx, len(respNext), len(respAll))
		*/
	})
	return respAll
}

func (suite *testSuite) getImageWithErrorV1Step(stepName string, t provider.T, api *secapi.StorageV1, tref secapi.TenantReference, expectedError error) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "GetImage")

		_, err := api.GetImage(t.Context(), tref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deleteImageV1Step(stepName string, t provider.T, api *secapi.StorageV1, resource *schema.Image) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageV1StepParams(sCtx, "DeleteImage")

		err := api.DeleteImage(t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

func (suite *testSuite) getListSkuV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.StorageV1,
	tref secapi.TenantReference,
	opts *builders.ListOptions,
) []*schema.StorageSku {
	var respNext []*schema.StorageSku
	var respAll []*schema.StorageSku
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageWorkspaceV1StepParams(sCtx, "GetListSku", string(tref.Name))
		var iter *secapi.Iterator[schema.StorageSku]
		var err error
		if opts != nil {
			iter, err = api.ListSkusWithFilters(ctx, tref.Tenant, opts)
		} else {
			iter, err = api.ListSkus(ctx, tref.Tenant)
		}
		requireNoError(sCtx, err)
		for {
			item, err := iter.Next(context.Background())
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				break
			}
			respNext = append(respNext, item)
		}
		requireNotNilResponse(sCtx, respNext)
		requireLenResponse(sCtx, len(respNext))
		/*
			respAll, err = iter.All(ctx)
			requireNoError(sCtx, err)
			requireNotNilResponse(sCtx, respAll)
			requireLenResponse(sCtx, len(respAll))

			compareIteratorsResponse(sCtx, len(respNext), len(respAll))
		*/
	})
	return respAll
}
