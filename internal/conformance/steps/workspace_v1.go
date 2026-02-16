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
	responseExpects ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	configurator.logStepName(stepName)
	createOrUpdateTenantResourceStep(configurator.t, configurator.suite,
		createOrUpdateTenantResourceParams[schema.Workspace, schema.RegionalResourceMetadata, schema.WorkspaceSpec, schema.WorkspaceStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetWorkspaceV1StepParams,
			operationName:  "CreateOrUpdateWorkspace",
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.Workspace) (
				*stepFuncResponse[schema.Workspace, schema.RegionalResourceMetadata, schema.WorkspaceSpec, schema.WorkspaceStatus], error,
			) {
				resp, err := api.CreateOrUpdateWorkspace(configurator.t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), err
			},
			expectedLabels:        responseExpects.Labels,
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    configurator.suite.VerifyRegionalResourceMetadataStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (configurator *StepsConfigurator) GetWorkspaceV1Step(stepName string, api secapi.WorkspaceV1, tref secapi.TenantReference,
	responseExpects ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec],
) *schema.Workspace {
	responseExpects.Metadata.Verb = http.MethodGet
	configurator.logStepName(stepName)
	return getTenantResourceStep(configurator.t, configurator.suite,
		getTenantResourceParams[schema.Workspace, schema.RegionalResourceMetadata, schema.WorkspaceSpec, schema.WorkspaceStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetWorkspaceV1StepParams,
			operationName:  "GetWorkspace",
			tref:           tref,
			getFunc: func(ctx context.Context, tref secapi.TenantReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (
				*stepFuncResponse[schema.Workspace, schema.RegionalResourceMetadata, schema.WorkspaceSpec, schema.WorkspaceStatus], error,
			) {
				resp, err := api.GetWorkspaceUntilState(ctx, tref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), err
			},
			expectedLabels:        responseExpects.Labels,
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    configurator.suite.VerifyRegionalResourceMetadataStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (configurator *StepsConfigurator) GetListWorkspaceV1Step(
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
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetWorkspaceV1StepParams,
			operationName:  "GetWorkspace",
			tref:           tref,
			getFunc: func(ctx context.Context, tref secapi.TenantReference) error {
				_, err := api.GetWorkspace(ctx, tref)
				return err
			},
			expectedError: expectedError,
		},
	)
}

func (configurator *StepsConfigurator) DeleteWorkspaceV1Step(stepName string, api secapi.WorkspaceV1, resource *schema.Workspace) {
	configurator.logStepName(stepName)
	deleteTenantResourceStep(configurator.t,
		deleteTenantResourceParams[schema.Workspace]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetWorkspaceV1StepParams,
			operationName:  "DeleteWorkspace",
			resource:       resource,
			deleteFunc: func(ctx context.Context, r *schema.Workspace) error {
				return api.DeleteWorkspace(ctx, r)
			},
		},
	)
}
