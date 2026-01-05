//nolint:dupl
package steps

import (
	"context"
	"net/http"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (builder *Builder) CreateOrUpdateWorkspaceV1Step(stepName string, api *secapi.WorkspaceV1, resource *schema.Workspace,
	responseExpects ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	createOrUpdateTenantResourceStep(builder.t, builder.suite,
		createOrUpdateTenantResourceParams[schema.Workspace, schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			stepName:       stepName,
			stepParamsFunc: builder.suite.SetWorkspaceV1StepParams,
			operationName:  "CreateOrUpdateWorkspace",
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.Workspace) (*stepFuncResponse[schema.Workspace, schema.RegionalResourceMetadata, schema.WorkspaceSpec], error) {
				resp, err := api.CreateOrUpdateWorkspace(builder.t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedLabels:        responseExpects.Labels,
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    builder.suite.VerifyRegionalResourceMetadataStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (builder *Builder) GetWorkspaceV1Step(stepName string, api *secapi.WorkspaceV1, tref secapi.TenantReference,
	responseExpects ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec],
) *schema.Workspace {
	responseExpects.Metadata.Verb = http.MethodGet
	return getTenantResourceStep(builder.t, builder.suite,
		getTenantResourceParams[schema.Workspace, schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			stepName:       stepName,
			stepParamsFunc: builder.suite.SetWorkspaceV1StepParams,
			operationName:  "GetWorkspace",
			tref:           tref,
			getFunc: func(ctx context.Context, tref secapi.TenantReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[schema.Workspace, schema.RegionalResourceMetadata, schema.WorkspaceSpec], error) {
				resp, err := api.GetWorkspaceUntilState(ctx, tref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedLabels:        responseExpects.Labels,
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    builder.suite.VerifyRegionalResourceMetadataStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (builder *Builder) GetListWorkspaceV1Step(
	stepName string,
	api *secapi.WorkspaceV1,
	tref secapi.TenantReference,
	opts *secapi.ListOptions,
) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetStorageWorkspaceV1StepParams(sCtx, "GetListWorkspace", string(tref.Tenant))
		var iter *secapi.Iterator[schema.Workspace]
		var err error
		if opts != nil {
			iter, err = api.ListWorkspacesWithFilters(builder.t.Context(), tref.Tenant, opts)
		} else {
			iter, err = api.ListWorkspaces(builder.t.Context(), tref.Tenant)
		}
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, builder.t, *iter)
	})
}

func (builder *Builder) GetWorkspaceWithErrorV1Step(stepName string, api *secapi.WorkspaceV1, tref secapi.TenantReference, expectedError error) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetWorkspaceV1StepParams(sCtx, "GetWorkspace")

		_, err := api.GetWorkspace(builder.t.Context(), tref)
		requireError(sCtx, err, expectedError)
	})
}

func (builder *Builder) DeleteWorkspaceV1Step(stepName string, api *secapi.WorkspaceV1, resource *schema.Workspace) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetWorkspaceV1StepParams(sCtx, "DeleteWorkspace")

		err := api.DeleteWorkspace(builder.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}
