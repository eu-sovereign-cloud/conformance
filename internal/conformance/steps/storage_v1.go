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

func (configurator *StepsConfigurator) CreateOrUpdateBlockStorageV1Step(stepName string, api secapi.StorageV1, resource *schema.BlockStorage,
	responseExpects StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	configurator.logStepName(stepName)
	createOrUpdateWorkspaceResourceStep(configurator.t, configurator.suite,
		createOrUpdateWorkspaceResourceParams[schema.BlockStorage, schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec, schema.BlockStorageStatus]{
			createOrUpdateResourceParams: createOrUpdateResourceParams[schema.BlockStorage, schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec, schema.BlockStorageStatus]{
				resource: resource,
				createOrUpdateFunc: func(context.Context, *schema.BlockStorage) (
					*createOrUpdateStepFuncResponse[schema.BlockStorage, schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec, schema.BlockStorageStatus], error,
				) {
					if resp, err := api.CreateOrUpdateBlockStorage(configurator.t.Context(), resource); err == nil {
						return newCreateOrUpdateStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), nil
					} else {
						return nil, err
					}
				},
				expectedMetadata:      responseExpects.Metadata,
				verifyMetadataFunc:    configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
				expectedSpec:          responseExpects.Spec,
				verifySpecFunc:        configurator.suite.VerifyBlockStorageSpecStep,
				expectedResourceState: responseExpects.ResourceState,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetStorageWorkspaceV1StepParams,
			operationName:  "CreateOrUpdateBlockStorage",
			workspace:      resource.Metadata.Workspace,
		},
	)
}

func (configurator *StepsConfigurator) GetBlockStorageV1Step(stepName string, api secapi.StorageV1, wref secapi.WorkspaceReference,
	responseExpects StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec],
) *schema.BlockStorage {
	responseExpects.Metadata.Verb = http.MethodGet
	configurator.logStepName(stepName)
	return getWorkspaceResourceStep(configurator.t, configurator.suite,
		getWorkspaceResourceParams[schema.BlockStorage, schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec, schema.BlockStorageStatus]{
			getResourceWithObserverParams: getResourceWithObserverParams[schema.BlockStorage, schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec, schema.BlockStorageStatus, secapi.WorkspaceReference, schema.ResourceState]{
				reference: wref,
				getFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (
					*getStepFuncResponse[schema.BlockStorage, schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec, schema.BlockStorageStatus], error,
				) {
					if resp, err := api.GetBlockStorageUntilState(ctx, wref, config); err == nil {
						return newGetStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), nil
					} else {
						return nil, err
					}
				},
				expectedMetadata:      responseExpects.Metadata,
				verifyMetadataFunc:    configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
				expectedSpec:          responseExpects.Spec,
				verifySpecFunc:        configurator.suite.VerifyBlockStorageSpecStep,
				expectedResourceState: responseExpects.ResourceState,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetStorageWorkspaceV1StepParams,
			operationName:  "GetBlockStorage",
		},
	)
}

func (configurator *StepsConfigurator) ListBlockStorageV1Step(
	stepName string,
	api secapi.StorageV1,
	wref secapi.WorkspaceReference,
	opts *secapi.ListOptions,
) {
	configurator.logStepName(stepName)
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetStorageWorkspaceV1StepParams(sCtx, "GetListBlockStorage", string(wref.Workspace))

		var iter *secapi.Iterator[schema.BlockStorage]
		var err error
		if opts != nil {
			iter, err = api.ListBlockStoragesWithFilters(configurator.t.Context(), wref.Tenant, wref.Workspace, opts)
		} else {
			iter, err = api.ListBlockStorages(configurator.t.Context(), wref.Tenant, wref.Workspace)
		}
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, configurator.t, *iter)
	})
}

func (configurator *StepsConfigurator) GetBlockStorageWithErrorV1Step(stepName string, api secapi.StorageV1, wref secapi.WorkspaceReference, expectedError error) {
	configurator.logStepName(stepName)
	getWorkspaceResourceWithErrorStep(configurator.t,
		getWorkspaceResourceWithErrorParams{
			getResourceWithErrorParams: getResourceWithErrorParams[secapi.WorkspaceReference]{
				reference: wref,
				getFunc: func(ctx context.Context, wref secapi.WorkspaceReference) error {
					_, err := api.GetBlockStorage(ctx, wref)
					return err
				},
				expectedError: expectedError,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetStorageWorkspaceV1StepParams,
			operationName:  "GetBlockStorage",
		},
	)
}

func (configurator *StepsConfigurator) DeleteBlockStorageV1Step(stepName string, api secapi.StorageV1, resource *schema.BlockStorage) {
	configurator.logStepName(stepName)
	deleteWorkspaceResourceStep(configurator.t,
		deleteWorkspaceResourceParams[schema.BlockStorage]{
			deleteResourceParams: deleteResourceParams[schema.BlockStorage]{
				resource: resource,
				deleteFunc: func(ctx context.Context, r *schema.BlockStorage) error {
					return api.DeleteBlockStorage(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetStorageWorkspaceV1StepParams,
			operationName:  "DeleteBlockStorage",
			workspace:      resource.Metadata.Workspace,
		},
	)
}

// Image

func (configurator *StepsConfigurator) CreateOrUpdateImageV1Step(stepName string, api secapi.StorageV1, resource *schema.Image,
	responseExpects StepResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	configurator.logStepName(stepName)
	createOrUpdateTenantResourceStep(configurator.t, configurator.suite,
		createOrUpdateTenantResourceParams[schema.Image, schema.RegionalResourceMetadata, schema.ImageSpec, schema.ImageStatus]{
			createOrUpdateResourceParams: createOrUpdateResourceParams[schema.Image, schema.RegionalResourceMetadata, schema.ImageSpec, schema.ImageStatus]{
				resource: resource,
				createOrUpdateFunc: func(context.Context, *schema.Image) (
					*createOrUpdateStepFuncResponse[schema.Image, schema.RegionalResourceMetadata, schema.ImageSpec, schema.ImageStatus], error,
				) {
					if resp, err := api.CreateOrUpdateImage(configurator.t.Context(), resource); err == nil {
						return newCreateOrUpdateStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), nil
					} else {
						return nil, err
					}
				},
				expectedMetadata:      responseExpects.Metadata,
				verifyMetadataFunc:    configurator.suite.VerifyRegionalResourceMetadataStep,
				expectedSpec:          responseExpects.Spec,
				verifySpecFunc:        configurator.suite.VerifyImageSpecStep,
				expectedResourceState: responseExpects.ResourceState,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetStorageV1StepParams,
			operationName:  "CreateOrUpdateImage",
		},
	)
}

func (configurator *StepsConfigurator) GetImageV1Step(stepName string, api secapi.StorageV1, tref secapi.TenantReference,
	responseExpects StepResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec],
) *schema.Image {
	responseExpects.Metadata.Verb = http.MethodGet
	configurator.logStepName(stepName)
	return getTenantResourceStep(configurator.t, configurator.suite,
		getTenantResourceParams[schema.Image, schema.RegionalResourceMetadata, schema.ImageSpec, schema.ImageStatus]{
			getResourceWithObserverParams: getResourceWithObserverParams[schema.Image, schema.RegionalResourceMetadata, schema.ImageSpec, schema.ImageStatus, secapi.TenantReference, schema.ResourceState]{
				reference: tref,
				getFunc: func(ctx context.Context, tref secapi.TenantReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (
					*getStepFuncResponse[schema.Image, schema.RegionalResourceMetadata, schema.ImageSpec, schema.ImageStatus], error,
				) {
					if resp, err := api.GetImageUntilState(ctx, tref, config); err == nil {
						return newGetStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), nil
					} else {
						return nil, err
					}
				},
				expectedMetadata:      responseExpects.Metadata,
				verifyMetadataFunc:    configurator.suite.VerifyRegionalResourceMetadataStep,
				expectedSpec:          responseExpects.Spec,
				verifySpecFunc:        configurator.suite.VerifyImageSpecStep,
				expectedResourceState: responseExpects.ResourceState,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetStorageV1StepParams,
			operationName:  "GetImage",
		},
	)
}

func (configurator *StepsConfigurator) ListImageV1Step(
	stepName string,
	api secapi.StorageV1,
	tref secapi.TenantReference,
	opts *secapi.ListOptions,
) {
	configurator.logStepName(stepName)
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetStorageWorkspaceV1StepParams(sCtx, "GetListImage", tref.Name)
		var iter *secapi.Iterator[schema.Image]
		var err error
		if opts != nil {
			iter, err = api.ListImagesWithFilters(configurator.t.Context(), tref.Tenant, opts)
		} else {
			iter, err = api.ListImages(configurator.t.Context(), tref.Tenant)
		}
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, configurator.t, *iter)
	})
}

func (configurator *StepsConfigurator) GetImageWithErrorV1Step(stepName string, api secapi.StorageV1, tref secapi.TenantReference, expectedError error) {
	configurator.logStepName(stepName)
	getTenantResourceWithErrorStep(configurator.t,
		getTenantResourceWithErrorParams{
			getResourceWithErrorParams: getResourceWithErrorParams[secapi.TenantReference]{
				reference: tref,
				getFunc: func(ctx context.Context, tref secapi.TenantReference) error {
					_, err := api.GetImage(ctx, tref)
					return err
				},
				expectedError: expectedError,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetStorageV1StepParams,
			operationName:  "GetImage",
		},
	)
}

func (configurator *StepsConfigurator) DeleteImageV1Step(stepName string, api secapi.StorageV1, resource *schema.Image) {
	configurator.logStepName(stepName)
	deleteTenantResourceStep(configurator.t,
		deleteTenantResourceParams[schema.Image]{
			deleteResourceParams: deleteResourceParams[schema.Image]{
				resource: resource,
				deleteFunc: func(ctx context.Context, r *schema.Image) error {
					return api.DeleteImage(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetStorageV1StepParams,
			operationName:  "DeleteImage",
		},
	)
}

func (configurator *StepsConfigurator) ListSkuV1Step(
	stepName string,
	api secapi.StorageV1,
	tref secapi.TenantReference,
	opts *secapi.ListOptions,
) {
	configurator.logStepName(stepName)
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetStorageWorkspaceV1StepParams(sCtx, "GetListSku", tref.Name)
		var iter *secapi.Iterator[schema.StorageSku]
		var err error
		if opts != nil {
			iter, err = api.ListSkusWithFilters(configurator.t.Context(), tref.Tenant, opts)
		} else {
			iter, err = api.ListSkus(configurator.t.Context(), tref.Tenant)
		}
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, configurator.t, *iter)
	})
}
