//nolint:dupl
package steps

import (
	"context"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/pkg/wrappers"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
)

// Sku

func (configurator *StepsConfigurator) CreateOrUpdateBlockStorageV1Step(stepName string, stepCreator StepCreator, api secapi.StorageV1, resource *schema.BlockStorage,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	createOrUpdateWorkspaceResourceStep(configurator.t.Context(), configurator.suite, stepCreator,
		createOrUpdateWorkspaceResourceParams[schema.BlockStorage, schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec, schema.BlockStorageStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetStorageWorkspaceV1StepParams,
			operationName:  constants.CreateOrUpdateBlockStorageOperation,
			workspace:      secapi.WorkspaceID(resource.Metadata.Workspace),
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.BlockStorage) (
				wrappers.ResourceWrapper[schema.BlockStorage, schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec, schema.BlockStorageStatus], error,
			) {
				resp, err := api.CreateOrUpdateBlockStorage(configurator.t.Context(), resource)
				return wrappers.NewBlockStorageWrapper(resp), err
			},
			expectedMetadata:       responseExpects.Metadata,
			verifyMetadataFunc:     configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:           responseExpects.Spec,
			verifySpecFunc:         configurator.suite.VerifyBlockStorageSpecStep,
			expectedResourceStates: responseExpects.ResourceStates,
		},
	)
}

func (configurator *StepsConfigurator) ListSkuV1Step(stepName string, api secapi.StorageV1, tpath secapi.TenantPath, opts *secapi.ListOptions) {
	listTenantResourcesStep(configurator.t, configurator.suite,
		listTenantResourcesParams[schema.StorageSku, schema.SkuResourceMetadata]{
			listResourcesParams: listResourcesParams[schema.StorageSku, schema.SkuResourceMetadata, secapi.TenantPath]{
				path: tpath, listOptions: opts,
				listFunc: func(ctx context.Context, path secapi.TenantPath, options *secapi.ListOptions) (*secapi.Iterator[schema.StorageSku], error) {
					return api.ListSkusWithOptions(ctx, path, options)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetStorageV1StepParams,
			operationName:  constants.ListSkusOperation,
		},
	)
}

// BlockStorage

func (configurator *StepsConfigurator) ListBlockStorageV1Step(stepName string, api secapi.StorageV1, wpath secapi.WorkspacePath, opts *secapi.ListOptions) {
	listWorkspaceResourcesStep(configurator.t, configurator.suite,
		listWorkspaceResourcesParams[schema.BlockStorage, schema.RegionalWorkspaceResourceMetadata]{
			listResourcesParams: listResourcesParams[schema.BlockStorage, schema.RegionalWorkspaceResourceMetadata, secapi.WorkspacePath]{
				path: wpath, listOptions: opts,
				listFunc: func(ctx context.Context, path secapi.WorkspacePath, options *secapi.ListOptions) (*secapi.Iterator[schema.BlockStorage], error) {
					return api.ListBlockStoragesWithOptions(ctx, path, options)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetStorageWorkspaceV1StepParams,
			operationName:  constants.ListBlockStorageOperation,
			workspace:      wpath.Workspace,
		},
	)
}

func (configurator *StepsConfigurator) GetBlockStorageV1Step(stepName string, api secapi.StorageV1, wref secapi.WorkspaceReference,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec],
) *schema.BlockStorage {
	responseExpects.Metadata.Verb = http.MethodGet
	return getWorkspaceResourceStep(configurator.t, configurator.suite,
		getWorkspaceResourceParams[schema.BlockStorage, schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec, schema.BlockStorageStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetStorageWorkspaceV1StepParams,
			operationName:  constants.GetBlockStorageOperation,
			wref:           wref,
			getValueFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverUntilValueConfig[schema.ResourceState]) (
				wrappers.ResourceWrapper[schema.BlockStorage, schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec, schema.BlockStorageStatus], error,
			) {
				resp, err := api.GetBlockStorageUntilState(ctx, wref, config)
				return wrappers.NewBlockStorageWrapper(resp), err
			},
			expectedMetadata:       responseExpects.Metadata,
			verifyMetadataFunc:     configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:           responseExpects.Spec,
			verifySpecFunc:         configurator.suite.VerifyBlockStorageSpecStep,
			expectedResourceStates: responseExpects.ResourceStates,
		},
	)
}

func (configurator *StepsConfigurator) WatchBlockStorageUntilDeletedV1Step(stepName string, stepCreator StepCreator, api secapi.StorageV1, wref secapi.WorkspaceReference) {
	watchWorkspaceResourceUntilDeletedStep(configurator.t.Context(), configurator.suite, stepCreator,
		watchWorkspaceResourceUntilDeletedParams{
			watchResourceUntilDeletedParams: watchResourceUntilDeletedParams[secapi.WorkspaceReference]{
				reference: wref,
				getErrorFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig) error {
					return api.WatchBlockStorageUntilDeleted(configurator.t.Context(), wref, config)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetStorageWorkspaceV1StepParams,
			operationName:  constants.GetBlockStorageOperation,
		},
	)
}

func (configurator *StepsConfigurator) DeleteBlockStorageV1Step(stepName string, stepCreator StepCreator, api secapi.StorageV1, resource *schema.BlockStorage) {
	deleteWorkspaceResourceStep(configurator.t.Context(), configurator.suite, stepCreator,
		deleteWorkspaceResourceParams[schema.BlockStorage]{
			deleteResourceParams: deleteResourceParams[schema.BlockStorage]{
				resource: resource,
				deleteFunc: func(ctx context.Context, r *schema.BlockStorage) error {
					return api.DeleteBlockStorage(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetStorageWorkspaceV1StepParams,
			operationName:  constants.DeleteBlockStorageOperation,
			workspace:      secapi.WorkspaceID(resource.Metadata.Workspace),
		},
	)
}

// Image

func (configurator *StepsConfigurator) ListImageV1Step(stepName string, api secapi.StorageV1, tpath secapi.TenantPath, opts *secapi.ListOptions) {
	listTenantResourcesStep(configurator.t, configurator.suite,
		listTenantResourcesParams[schema.Image, schema.GlobalTenantResourceMetadata]{
			listResourcesParams: listResourcesParams[schema.Image, schema.GlobalTenantResourceMetadata, secapi.TenantPath]{
				path: tpath, listOptions: opts,
				listFunc: func(ctx context.Context, path secapi.TenantPath, options *secapi.ListOptions) (*secapi.Iterator[schema.Image], error) {
					return api.ListImagesWithOptions(ctx, path, options)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetStorageV1StepParams,
			operationName:  constants.ListImagesOperation,
		},
	)
}

func (configurator *StepsConfigurator) GetImageV1Step(stepName string, api secapi.StorageV1, tref secapi.TenantReference,
	responseExpects ResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec],
) *schema.Image {
	responseExpects.Metadata.Verb = http.MethodGet
	return getTenantResourceStep(configurator.t, configurator.suite,
		getTenantResourceParams[schema.Image, schema.RegionalResourceMetadata, schema.ImageSpec, schema.ImageStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetStorageV1StepParams,
			operationName:  constants.GetImageOperation,
			tref:           tref,
			getValueFunc: func(ctx context.Context, tref secapi.TenantReference, config secapi.ResourceObserverUntilValueConfig[schema.ResourceState]) (
				wrappers.ResourceWrapper[schema.Image, schema.RegionalResourceMetadata, schema.ImageSpec, schema.ImageStatus], error,
			) {
				resp, err := api.GetImageUntilState(ctx, tref, config)
				return wrappers.NewImageWrapper(resp), err
			},
			expectedMetadata:       responseExpects.Metadata,
			verifyMetadataFunc:     configurator.suite.VerifyRegionalResourceMetadataStep,
			expectedSpec:           responseExpects.Spec,
			verifySpecFunc:         configurator.suite.VerifyImageSpecStep,
			expectedResourceStates: responseExpects.ResourceStates,
		},
	)
}

func (configurator *StepsConfigurator) WatchImageUntilDeletedV1Step(stepName string, stepCreator StepCreator, api secapi.StorageV1, tref secapi.TenantReference) {
	watchTenantResourceUntilDeletedStep(configurator.t.Context(), configurator.suite, stepCreator,
		watchTenantResourceUntilDeletedParams{
			watchResourceUntilDeletedParams: watchResourceUntilDeletedParams[secapi.TenantReference]{
				reference: tref,
				getErrorFunc: func(ctx context.Context, tref secapi.TenantReference, config secapi.ResourceObserverConfig) error {
					return api.WatchImageUntilDeleted(configurator.t.Context(), tref, config)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetStorageV1StepParams,
			operationName:  constants.GetImageOperation,
		},
	)
}

func (configurator *StepsConfigurator) CreateOrUpdateImageV1Step(stepName string, stepCreator StepCreator, api secapi.StorageV1, resource *schema.Image,
	responseExpects ResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	createOrUpdateTenantResourceStep(configurator.t.Context(), configurator.suite, stepCreator,
		createOrUpdateTenantResourceParams[schema.Image, schema.RegionalResourceMetadata, schema.ImageSpec, schema.ImageStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetStorageV1StepParams,
			operationName:  constants.CreateOrUpdateImageOperation,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.Image) (
				wrappers.ResourceWrapper[schema.Image, schema.RegionalResourceMetadata, schema.ImageSpec, schema.ImageStatus], error,
			) {
				resp, err := api.CreateOrUpdateImage(configurator.t.Context(), resource)
				return wrappers.NewImageWrapper(resp), err
			},
			expectedMetadata:       responseExpects.Metadata,
			verifyMetadataFunc:     configurator.suite.VerifyRegionalResourceMetadataStep,
			expectedSpec:           responseExpects.Spec,
			verifySpecFunc:         configurator.suite.VerifyImageSpecStep,
			expectedResourceStates: responseExpects.ResourceStates,
		},
	)
}

func (configurator *StepsConfigurator) DeleteImageV1Step(stepName string, stepCreator StepCreator, api secapi.StorageV1, resource *schema.Image) {
	deleteTenantResourceStep(configurator.t.Context(), configurator.suite, stepCreator,
		deleteTenantResourceParams[schema.Image]{
			deleteResourceParams: deleteResourceParams[schema.Image]{
				resource: resource,
				deleteFunc: func(ctx context.Context, r *schema.Image) error {
					return api.DeleteImage(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetStorageV1StepParams,
			operationName:  constants.DeleteImageOperation,
		},
	)
}
