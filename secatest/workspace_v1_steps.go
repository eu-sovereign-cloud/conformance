package secatest

import (
	"context"
	"net/http"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (suite *testSuite) createOrUpdateWorkspaceV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.WorkspaceV1,
	resource *schema.Workspace,
	expectedMeta *schema.RegionalResourceMetadata,
	expectedLabels schema.Labels,
	expectedState schema.ResourceState,
) {
	expectedMeta.Verb = http.MethodPut
	createOrUpdateResourceStep(
		t,
		ctx,
		suite,
		stepName,
		suite.setWorkspaceV1StepParams,
		"CreateOrUpdateWorkspace",
		resource,
		func(context.Context, *schema.Workspace) (*stepFuncResponse[schema.Workspace, schema.RegionalResourceMetadata, schema.WorkspaceSpec], error) {
			resp, err := api.CreateOrUpdateWorkspace(ctx, resource)
			return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
		},
		expectedLabels,
		expectedMeta,
		suite.verifyRegionalResourceMetadataStep,
		nil,
		nil,
		expectedState,
	)
}

func (suite *testSuite) getWorkspaceV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.WorkspaceV1,
	tref secapi.TenantReference,
	expectedMeta *schema.RegionalResourceMetadata,
	expectedLabels schema.Labels,
	expectedState schema.ResourceState,
) *schema.Workspace {
	expectedMeta.Verb = http.MethodGet
	return getTenantResourceStep(
		t,
		ctx,
		suite,
		stepName,
		suite.setWorkspaceV1StepParams,
		"GetWorkspace",
		tref,
		func(context.Context, secapi.TenantReference) (*stepFuncResponse[schema.Workspace, schema.RegionalResourceMetadata, schema.WorkspaceSpec], error) {
			resp, err := api.GetWorkspace(ctx, tref)
			return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
		},
		expectedLabels,
		expectedMeta,
		suite.verifyRegionalResourceMetadataStep,
		nil,
		nil,
		expectedState,
	)
}

func (suite *testSuite) getWorkspaceWithErrorV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.WorkspaceV1,
	tref secapi.TenantReference,
	expectedError error,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "GetWorkspace")

		_, err := api.GetWorkspace(ctx, tref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deleteWorkspaceV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.WorkspaceV1, resource *schema.Workspace) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "DeleteWorkspace")

		err := api.DeleteWorkspace(ctx, resource)
		requireNoError(sCtx, err)
	})
}
