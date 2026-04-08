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

func (configurator *StepsConfigurator) CreateOrUpdateWorkspaceV1Step(stepName string, stepCreator StepCreator, api secapi.WorkspaceV1, resource *schema.Workspace,
	responseExpects ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	createOrUpdateTenantResourceStep(configurator.t.Context(), configurator.suite, stepCreator,
		createOrUpdateTenantResourceParams[schema.Workspace, schema.RegionalResourceMetadata, schema.WorkspaceSpec, schema.WorkspaceStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetWorkspaceV1StepParams,
			operationName:  constants.CreateOrUpdateWorkspaceOperation,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.Workspace) (
				wrappers.ResourceWrapper[schema.Workspace, schema.RegionalResourceMetadata, schema.WorkspaceSpec, schema.WorkspaceStatus], error,
			) {
				resp, err := api.CreateOrUpdateWorkspace(configurator.t.Context(), resource)
				return wrappers.NewWorkspaceWrapper(resp), err
			},
			expectedLabels:         responseExpects.Labels,
			expectedMetadata:       responseExpects.Metadata,
			verifyMetadataFunc:     configurator.suite.VerifyRegionalResourceMetadataStep,
			expectedResourceStates: responseExpects.ResourceStates,
		},
	)
}

func (configurator *StepsConfigurator) ListWorkspaceV1Step(stepName string, api secapi.WorkspaceV1, tpath secapi.TenantPath, opts *secapi.ListOptions) {
	listTenantResourcesStep(configurator.t, configurator.suite,
		listTenantResourcesParams[schema.Workspace, schema.GlobalTenantResourceMetadata]{
			listResourcesParams: listResourcesParams[schema.Workspace, schema.GlobalTenantResourceMetadata, secapi.TenantPath]{
				path: tpath, listOptions: opts,
				listFunc: func(ctx context.Context, path secapi.TenantPath, options *secapi.ListOptions) (*secapi.Iterator[schema.Workspace], error) {
					return api.ListWorkspacesWithOptions(ctx, path, options)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetWorkspaceV1StepParams,
			operationName:  constants.ListWorkspacesOperation,
		},
	)
}

func (configurator *StepsConfigurator) GetWorkspaceV1Step(stepName string, api secapi.WorkspaceV1, tref secapi.TenantReference,
	responseExpects ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec],
) *schema.Workspace {
	responseExpects.Metadata.Verb = http.MethodGet
	return getTenantResourceStep(configurator.t, configurator.suite,
		getTenantResourceParams[schema.Workspace, schema.RegionalResourceMetadata, schema.WorkspaceSpec, schema.WorkspaceStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetWorkspaceV1StepParams,
			operationName:  constants.GetWorkspaceOperation,
			tref:           tref,
			getValueFunc: func(ctx context.Context, tref secapi.TenantReference, config secapi.ResourceObserverUntilValueConfig[schema.ResourceState]) (
				wrappers.ResourceWrapper[schema.Workspace, schema.RegionalResourceMetadata, schema.WorkspaceSpec, schema.WorkspaceStatus], error,
			) {
				resp, err := api.GetWorkspaceUntilState(ctx, tref, config)
				return wrappers.NewWorkspaceWrapper(resp), err
			},
			expectedLabels:         responseExpects.Labels,
			expectedMetadata:       responseExpects.Metadata,
			verifyMetadataFunc:     configurator.suite.VerifyRegionalResourceMetadataStep,
			expectedResourceStates: responseExpects.ResourceStates,
		},
	)
}

func (configurator *StepsConfigurator) WatchWorkspaceUntilDeletedV1Step(stepName string, stepCreator StepCreator, api secapi.WorkspaceV1, tref secapi.TenantReference) {
	watchTenantResourceUntilDeletedStep(configurator.t.Context(), configurator.suite, stepCreator,
		watchTenantResourceUntilDeletedParams{
			watchResourceUntilDeletedParams: watchResourceUntilDeletedParams[secapi.TenantReference]{
				reference: tref,
				getErrorFunc: func(ctx context.Context, tref secapi.TenantReference, config secapi.ResourceObserverConfig) error {
					return api.WatchWorkspaceUntilDeleted(configurator.t.Context(), tref, config)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetWorkspaceV1StepParams,
			operationName:  constants.GetWorkspaceOperation,
		},
	)
}

func (configurator *StepsConfigurator) DeleteWorkspaceV1Step(stepName string, stepCreator StepCreator, api secapi.WorkspaceV1, resource *schema.Workspace) {
	deleteTenantResourceStep(configurator.t.Context(), configurator.suite, stepCreator,
		deleteTenantResourceParams[schema.Workspace]{
			deleteResourceParams: deleteResourceParams[schema.Workspace]{
				resource: resource,
				deleteFunc: func(ctx context.Context, r *schema.Workspace) error {
					return api.DeleteWorkspace(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetWorkspaceV1StepParams,
			operationName:  constants.DeleteWorkspaceOperation,
		},
	)
}
