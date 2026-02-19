//nolint:dupl
package steps

import (
	"context"
	"net/http"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (configurator *StepsConfigurator) CreateOrUpdateWorkspaceV1Step(stepName string, api secapi.WorkspaceV1, resource *schema.Workspace,
	responseExpects StepResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	configurator.logStepName(stepName)
	createOrUpdateTenantResourceStep(configurator.t, configurator.suite,
		createOrUpdateTenantResourceParams[schema.Workspace, schema.RegionalResourceMetadata, schema.WorkspaceSpec, schema.WorkspaceStatus]{
			createOrUpdateResourceParams: createOrUpdateResourceParams[schema.Workspace, schema.RegionalResourceMetadata, schema.WorkspaceSpec, schema.WorkspaceStatus]{
				resource: resource,
				createOrUpdateFunc: func(context.Context, *schema.Workspace) (
					*createOrUpdateStepFuncResponse[schema.Workspace, schema.RegionalResourceMetadata, schema.WorkspaceSpec, schema.WorkspaceStatus], error,
				) {
					if resp, err := api.CreateOrUpdateWorkspace(configurator.t.Context(), resource); err == nil {
						return newCreateOrUpdateStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), nil
					} else {
						return nil, err
					}
				},
				expectedLabels:        responseExpects.Labels,
				expectedMetadata:      responseExpects.Metadata,
				verifyMetadataFunc:    configurator.suite.VerifyRegionalResourceMetadataStep,
				expectedResourceState: responseExpects.ResourceState,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetWorkspaceV1StepParams,
			operationName:  "CreateOrUpdateWorkspace",
		},
	)
}

func (configurator *StepsConfigurator) GetWorkspaceV1Step(stepName string, api secapi.WorkspaceV1, tref secapi.TenantReference,
	responseExpects StepResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec],
) *schema.Workspace {
	responseExpects.Metadata.Verb = http.MethodGet
	configurator.logStepName(stepName)
	return getTenantResourceStep(configurator.t, configurator.suite,
		getTenantResourceParams[schema.Workspace, schema.RegionalResourceMetadata, schema.WorkspaceSpec, schema.WorkspaceStatus]{
			getResourceWithObserverParams: getResourceWithObserverParams[schema.Workspace, schema.RegionalResourceMetadata, schema.WorkspaceSpec, schema.WorkspaceStatus, secapi.TenantReference, schema.ResourceState]{
				reference: tref,
				getFunc: func(ctx context.Context, tref secapi.TenantReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (
					*getStepFuncResponse[schema.Workspace, schema.RegionalResourceMetadata, schema.WorkspaceSpec, schema.WorkspaceStatus], error,
				) {
					if resp, err := api.GetWorkspaceUntilState(ctx, tref, config); err == nil {
						return newGetStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), nil
					} else {
						return nil, err
					}
				},
				expectedLabels:        responseExpects.Labels,
				expectedMetadata:      responseExpects.Metadata,
				verifyMetadataFunc:    configurator.suite.VerifyRegionalResourceMetadataStep,
				expectedResourceState: responseExpects.ResourceState,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetWorkspaceV1StepParams,
			operationName:  "GetWorkspace",
		},
	)
}

func (configurator *StepsConfigurator) ListWorkspaceV1Step(
	stepName string,
	api secapi.WorkspaceV1,
	tref secapi.TenantReference,
	opts *secapi.ListOptions,
) {
	configurator.logStepName(stepName)
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetStorageWorkspaceV1StepParams(sCtx, "GetListWorkspace", string(tref.Tenant))
		var iter *secapi.Iterator[schema.Workspace]
		var err error
		if opts != nil {
			iter, err = api.ListWorkspacesWithFilters(configurator.t.Context(), tref.Tenant, opts)
		} else {
			iter, err = api.ListWorkspaces(configurator.t.Context(), tref.Tenant)
		}
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, configurator.t, *iter)
	})
}

func (configurator *StepsConfigurator) GetWorkspaceWithErrorV1Step(stepName string, api secapi.WorkspaceV1, tref secapi.TenantReference, expectedError error) {
	configurator.logStepName(stepName)
	getTenantResourceWithErrorStep(configurator.t,
		getTenantResourceWithErrorParams{
			getResourceWithErrorParams: getResourceWithErrorParams[secapi.TenantReference]{
				reference: tref,
				getFunc: func(ctx context.Context, tref secapi.TenantReference) error {
					_, err := api.GetWorkspace(ctx, tref)
					return err
				},
				expectedError: expectedError,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetWorkspaceV1StepParams,
			operationName:  "GetWorkspace",
		},
	)
}

func (configurator *StepsConfigurator) DeleteWorkspaceV1Step(stepName string, api secapi.WorkspaceV1, resource *schema.Workspace) {
	configurator.logStepName(stepName)
	deleteTenantResourceStep(configurator.t,
		deleteTenantResourceParams[schema.Workspace]{
			deleteResourceParams: deleteResourceParams[schema.Workspace]{
				resource: resource,
				deleteFunc: func(ctx context.Context, r *schema.Workspace) error {
					return api.DeleteWorkspace(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetWorkspaceV1StepParams,
			operationName:  "DeleteWorkspace",
		},
	)
}
