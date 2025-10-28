package secatest

import (
	"context"
	"net/http"
	"time"

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
	var err error

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "GetWorkspace")
		time.Sleep(time.Duration(suite.baseDelay) * time.Second)
		for attempt := 1; attempt <= suite.maxAttempts; attempt++ {
			resp, err = api.GetWorkspace(ctx, tref)
			requireNoError(sCtx, err)
			requireNotNilResponse(sCtx, resp)
			if resp.Status.State != nil && *resp.Status.State == *secalib.SetResourceState(expectedStatusState) {

				if expectedMeta != nil {
					expectedMeta.Verb = http.MethodGet
					suite.verifyRegionalResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)
				}

				if expectedLabels != nil {
					suite.verifyLabelsStep(sCtx, expectedLabels, resp.Labels)
				}

				suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
				return
			} else {
				time.Sleep(time.Duration(suite.baseInterval) * time.Second)
			}
			suite.verifyMaxAttempts(sCtx, attempt, "GetWorkspace", expectedStatusState)
		}
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
