package secatest

import (
	"context"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/secalib"
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
	expectedStatusState string,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "CreateOrUpdateWorkspace")

		resp, err := api.CreateOrUpdateWorkspace(ctx, resource)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		if expectedMeta != nil {
			expectedMeta.Verb = http.MethodPut
			suite.verifyRegionalResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)
		}

		if expectedLabels != nil {
			suite.verifyLabelsStep(sCtx, expectedLabels, resp.Labels)
		}

		suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
	})
}

func (suite *testSuite) getWorkspaceV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.WorkspaceV1,
	tref secapi.TenantReference,
	expectedMeta *schema.RegionalResourceMetadata,
	expectedLabels schema.Labels,
	expectedStatusState string,
) *schema.Workspace {
	var resp *schema.Workspace

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "GetWorkspace")
		retry := newStepRetry(
			suite.baseDelay,
			suite.baseInterval,
			suite.maxAttempts,
			func() schema.ResourceState {
				var err error
				resp, err = api.GetWorkspace(ctx, tref)
				requireNoError(sCtx, err)
				requireNotNilResponse(sCtx, resp)

				suite.requireNotNilStatus(sCtx, resp.Status)
				return *resp.Status.State
			},
			func() {
				expectedMeta.Verb = http.MethodGet
				suite.verifyRegionalResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)

				suite.verifyLabelsStep(sCtx, expectedLabels, resp.Labels)

				suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
			},
		)
		retry.run(sCtx, "GetWorkspace", expectedStatusState)
	})
	return resp
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
