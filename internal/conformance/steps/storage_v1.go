//nolint:dupl
package steps

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/pkg/wrappers"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// BlockStorage

func (configurator *StepsConfigurator) CreateOrUpdateBlockStorageV1Step(stepName string, api secapi.StorageV1, resource *schema.BlockStorage,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	createOrUpdateWorkspaceResourceStep(configurator.t, configurator.suite,
		createOrUpdateWorkspaceResourceParams[schema.BlockStorage, schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec, schema.BlockStorageStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetStorageWorkspaceV1StepParams,
			operationName:  constants.CreateOrUpdateBlockStorageOperation,
			workspace:      resource.Metadata.Workspace,
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

func (configurator *StepsConfigurator) GetBlockStorageV1Step(stepName string, api secapi.StorageV1, wref secapi.WorkspaceReference,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.BlockStorageSpec],
) *schema.BlockStorage {
	responseExpects.Metadata.Verb = http.MethodGet
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
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

func (configurator *StepsConfigurator) ListBlockStorageV1Step(
	stepName string, api secapi.StorageV1, wref secapi.WorkspaceReference, opts *secapi.FilterOptions,
) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetStorageWorkspaceV1StepParams(sCtx, "ListBlockStorage", string(wref.Workspace))

		iter, err := api.ListBlockStorages(configurator.t.Context(), secapi.WorkspaceFilter{Tenant: wref.Tenant, Workspace: wref.Workspace, Options: opts})
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, configurator.t, *iter)
	})
}

func (configurator *StepsConfigurator) WatchBlockStorageUntilDeletedV1Step(stepName string, api secapi.StorageV1, wref secapi.WorkspaceReference) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		watchWorkspaceResourceUntilDeletedStep(configurator.t, configurator.suite,
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
	})
}

func (configurator *StepsConfigurator) DeleteBlockStorageV1Step(stepName string, api secapi.StorageV1, resource *schema.BlockStorage) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetStorageWorkspaceV1StepParams(sCtx, "DeleteBlockStorage", resource.Metadata.Workspace)

		err := api.DeleteBlockStorage(configurator.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

// Image

func (configurator *StepsConfigurator) CreateOrUpdateImageV1Step(stepName string, api secapi.StorageV1, resource *schema.Image,
	responseExpects ResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	createOrUpdateTenantResourceStep(configurator.t, configurator.suite,
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

func (configurator *StepsConfigurator) GetImageV1Step(stepName string, api secapi.StorageV1, tref secapi.TenantReference,
	responseExpects ResponseExpects[schema.RegionalResourceMetadata, schema.ImageSpec],
) *schema.Image {
	responseExpects.Metadata.Verb = http.MethodGet
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
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

func (configurator *StepsConfigurator) ListImageV1Step(
	stepName string, api secapi.StorageV1, tref secapi.TenantReference, opts *secapi.FilterOptions,
) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetStorageWorkspaceV1StepParams(sCtx, "ListImage", tref.Name)
		iter, err := api.ListImages(configurator.t.Context(), secapi.TenantFilter{Tenant: tref.Tenant, Options: opts})
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, configurator.t, *iter)
	})
}

func (configurator *StepsConfigurator) WatchImageUntilDeletedV1Step(stepName string, api secapi.StorageV1, tref secapi.TenantReference) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		watchTenantResourceUntilDeletedStep(configurator.t, configurator.suite,
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
	})
}

func (configurator *StepsConfigurator) DeleteImageV1Step(stepName string, api secapi.StorageV1, resource *schema.Image) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetStorageV1StepParams(sCtx, "DeleteImage")

		err := api.DeleteImage(configurator.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

func (configurator *StepsConfigurator) ListSkuV1Step(
	stepName string, api secapi.StorageV1, tref secapi.TenantReference, opts *secapi.FilterOptions,
) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetStorageWorkspaceV1StepParams(sCtx, "ListSku", tref.Name)
		iter, err := api.ListSkus(configurator.t.Context(), secapi.TenantFilter{Tenant: tref.Tenant, Options: opts})
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, configurator.t, *iter)
	})
}
