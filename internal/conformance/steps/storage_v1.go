//nolint:dupl
package steps

import (
	"context"
	"net/http"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// BlockStorage

func (builder *Builder) CreateOrUpdateBlockStorageV1Step(stepName string, api *secapi.StorageV1, resource *schema.BlockStorage,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	createOrUpdateWorkspaceResourceStep(builder.t, builder.suite,
		createOrUpdateWorkspaceResourceParams[schema.BlockStorage, schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			stepName:       stepName,
			stepParamsFunc: builder.suite.SetStorageWorkspaceV1StepParams,
			operationName:  "CreateOrUpdateBlockStorage",
			workspace:      resource.Metadata.Workspace,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.BlockStorage) (*stepFuncResponse[schema.BlockStorage, schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec], error) {
				resp, err := api.CreateOrUpdateBlockStorage(builder.t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    builder.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        builder.suite.VerifyBlockStorageSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (builder *Builder) GetBlockStorageV1Step(stepName string, api *secapi.StorageV1, wref secapi.WorkspaceReference,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec],
) *schema.BlockStorage {
	responseExpects.Metadata.Verb = http.MethodGet
	return getWorkspaceResourceStep(builder.t, builder.suite,
		getWorkspaceResourceParams[schema.BlockStorage, schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec]{
			stepName:       stepName,
			stepParamsFunc: builder.suite.SetStorageWorkspaceV1StepParams,
			operationName:  "GetBlockStorage",
			wref:           wref,
			getFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[schema.BlockStorage, schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec], error) {
				resp, err := api.GetBlockStorageUntilState(ctx, wref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    builder.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        builder.suite.VerifyBlockStorageSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (builder *Builder) GetListBlockStorageV1Step(
	stepName string,
	t provider.T,
	api *secapi.StorageV1,
	wref secapi.WorkspaceReference,
	opts *secapi.ListOptions,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetStorageWorkspaceV1StepParams(sCtx, "GetListBlockStorage", string(wref.Workspace))

		var iter *secapi.Iterator[schema.BlockStorage]
		var err error
		if opts != nil {
			iter, err = api.ListBlockStoragesWithFilters(t.Context(), wref.Tenant, wref.Workspace, opts)
		} else {
			iter, err = api.ListBlockStorages(t.Context(), wref.Tenant, wref.Workspace)
		}
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, t, *iter)
	})
}

func (builder *Builder) GetBlockStorageWithErrorV1Step(stepName string, api *secapi.StorageV1, wref secapi.WorkspaceReference, expectedError error) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetStorageWorkspaceV1StepParams(sCtx, "GetBlockStorage", string(wref.Workspace))

		_, err := api.GetBlockStorage(builder.t.Context(), wref)
		requireError(sCtx, err, expectedError)
	})
}

func (builder *Builder) DeleteBlockStorageV1Step(stepName string, api *secapi.StorageV1, resource *schema.BlockStorage) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetStorageWorkspaceV1StepParams(sCtx, "DeleteBlockStorage", resource.Metadata.Workspace)

		err := api.DeleteBlockStorage(builder.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

// Image

func (builder *Builder) CreateOrUpdateImageV1Step(stepName string, api *secapi.StorageV1, resource *schema.Image,
	responseExpects ResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	createOrUpdateTenantResourceStep(builder.t, builder.suite,
		createOrUpdateTenantResourceParams[schema.Image, schema.RegionalResourceMetadata, schema.ImageSpec]{
			stepName:       stepName,
			stepParamsFunc: builder.suite.SetStorageV1StepParams,
			operationName:  "CreateOrUpdateImage",
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.Image) (*stepFuncResponse[schema.Image, schema.RegionalResourceMetadata, schema.ImageSpec], error) {
				resp, err := api.CreateOrUpdateImage(builder.t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    builder.suite.VerifyRegionalResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        builder.suite.VerifyImageSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (builder *Builder) GetImageV1Step(stepName string, api *secapi.StorageV1, tref secapi.TenantReference,
	responseExpects ResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec],
) *schema.Image {
	responseExpects.Metadata.Verb = http.MethodGet
	return getTenantResourceStep(builder.t, builder.suite,
		getTenantResourceParams[schema.Image, schema.RegionalResourceMetadata, schema.ImageSpec]{
			stepName:       stepName,
			stepParamsFunc: builder.suite.SetStorageV1StepParams,
			operationName:  "GetImage",
			tref:           tref,
			getFunc: func(ctx context.Context, tref secapi.TenantReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[schema.Image, schema.RegionalResourceMetadata, schema.ImageSpec], error) {
				resp, err := api.GetImageUntilState(ctx, tref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    builder.suite.VerifyRegionalResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        builder.suite.VerifyImageSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (builder *Builder) GetListImageV1Step(
	stepName string,
	api *secapi.StorageV1,
	tref secapi.TenantReference,
	opts *secapi.ListOptions,
) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetStorageWorkspaceV1StepParams(sCtx, "GetListImage", tref.Name)
		var iter *secapi.Iterator[schema.Image]
		var err error
		if opts != nil {
			iter, err = api.ListImagesWithFilters(builder.t.Context(), tref.Tenant, opts)
		} else {
			iter, err = api.ListImages(builder.t.Context(), tref.Tenant)
		}
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, builder.t, *iter)
	})
}

func (builder *Builder) GetImageWithErrorV1Step(stepName string, api *secapi.StorageV1, tref secapi.TenantReference, expectedError error) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetStorageV1StepParams(sCtx, "GetImage")

		_, err := api.GetImage(builder.t.Context(), tref)
		requireError(sCtx, err, expectedError)
	})
}

func (builder *Builder) DeleteImageV1Step(stepName string, api *secapi.StorageV1, resource *schema.Image) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetStorageV1StepParams(sCtx, "DeleteImage")

		err := api.DeleteImage(builder.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

func (builder *Builder) GetListSkuV1Step(
	stepName string,
	t provider.T,
	api *secapi.StorageV1,
	tref secapi.TenantReference,
	opts *secapi.ListOptions,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetStorageWorkspaceV1StepParams(sCtx, "GetListSku", tref.Name)
		var iter *secapi.Iterator[schema.StorageSku]
		var err error
		if opts != nil {
			iter, err = api.ListSkusWithFilters(t.Context(), tref.Tenant, opts)
		} else {
			iter, err = api.ListSkus(t.Context(), tref.Tenant)
		}
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, t, *iter)
	})
}
