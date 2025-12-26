//nolint:dupl
package secatest

import (
	"context"
	"net/http"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (suite *testSuite) createOrUpdateWorkspaceV1Step(stepName string, t provider.T, api *secapi.WorkspaceV1, resource *schema.Workspace,
	responseExpects responseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec],
) {
	responseExpects.metadata.Verb = http.MethodPut
	createOrUpdateTenantResourceStep(t, suite,
		createOrUpdateTenantResourceParams[schema.Workspace, schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			stepName:       stepName,
			stepParamsFunc: suite.setWorkspaceV1StepParams,
			operationName:  "CreateOrUpdateWorkspace",
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.Workspace) (*stepFuncResponse[schema.Workspace, schema.RegionalResourceMetadata, schema.WorkspaceSpec], error) {
				resp, err := api.CreateOrUpdateWorkspace(t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedLabels:        responseExpects.labels,
			expectedMetadata:      responseExpects.metadata,
			verifyMetadataFunc:    suite.verifyRegionalResourceMetadataStep,
			expectedResourceState: responseExpects.resourceState,
		},
	)
}

func (suite *testSuite) getWorkspaceV1Step(stepName string, t provider.T, api *secapi.WorkspaceV1, tref secapi.TenantReference,
	responseExpects responseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec],
) *schema.Workspace {
	responseExpects.metadata.Verb = http.MethodGet
	return getTenantResourceStep(t, suite,
		getTenantResourceParams[schema.Workspace, schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			stepName:       stepName,
			stepParamsFunc: suite.setWorkspaceV1StepParams,
			operationName:  "GetWorkspace",
			tref:           tref,
			getFunc: func(ctx context.Context, tref secapi.TenantReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[schema.Workspace, schema.RegionalResourceMetadata, schema.WorkspaceSpec], error) {
				resp, err := api.GetWorkspaceUntilState(ctx, tref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedLabels:        responseExpects.labels,
			expectedMetadata:      responseExpects.metadata,
			verifyMetadataFunc:    suite.verifyRegionalResourceMetadataStep,
			expectedResourceState: responseExpects.resourceState,
		},
	)
}

func (suite *testSuite) getListWorkspaceV1Step(
	stepName string,
	t provider.T,
	api *secapi.WorkspaceV1,
	tref secapi.TenantReference,
	opts *secapi.ListOptions,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageWorkspaceV1StepParams(sCtx, "GetListWorkspace", string(tref.Tenant))
		var iter *secapi.Iterator[schema.Workspace]
		var err error
		if opts != nil {
			iter, err = api.ListWorkspacesWithFilters(t.Context(), tref.Tenant, opts)
		} else {
			iter, err = api.ListWorkspaces(t.Context(), tref.Tenant)
		}
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, t, *iter)
	})
}

func (suite *testSuite) getWorkspaceWithErrorV1Step(stepName string, t provider.T, api *secapi.WorkspaceV1, tref secapi.TenantReference, expectedError error) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "GetWorkspace")

		_, err := api.GetWorkspace(t.Context(), tref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deleteWorkspaceV1Step(stepName string, t provider.T, api *secapi.WorkspaceV1, resource *schema.Workspace) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "DeleteWorkspace")

		err := api.DeleteWorkspace(t.Context(), resource)
		requireNoError(sCtx, err)
	})
}
