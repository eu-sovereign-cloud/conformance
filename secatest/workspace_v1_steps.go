package secatest

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/eu-sovereign-cloud/go-sdk/secapi/builders"

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

		resp, err = api.GetWorkspace(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		if expectedMeta != nil {
			expectedMeta.Verb = http.MethodGet
			suite.verifyRegionalResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)
		}

		if expectedLabels != nil {
			suite.verifyLabelsStep(sCtx, expectedLabels, resp.Labels)
		}

		suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
	})
	return resp
}

func (suite *testSuite) getListWorkspaceV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.WorkspaceV1,
	tref secapi.TenantReference,
	opts *builders.ListOptions,
) []*schema.Workspace {
	var respNext []*schema.Workspace
	var respAll []*schema.Workspace
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setStorageWorkspaceV1StepParams(sCtx, "GetListWorkspace", string(tref.Tenant))
		var iter *secapi.Iterator[schema.Workspace]
		var err error
		if opts != nil {
			iter, err = api.ListWorkspacesWithFilters(ctx, tref.Tenant, opts)
		} else {
			iter, err = api.ListWorkspaces(ctx, tref.Tenant)
		}
		requireNoError(sCtx, err)
		for {
			item, err := iter.Next(context.Background())
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				break
			}
			respNext = append(respNext, item)
		}
		requireNotNilResponse(sCtx, respNext)
		requireLenResponse(sCtx, len(respNext))
		/*
			respAll, err = iter.All(ctx)
			requireNoError(sCtx, err)
			requireNotNilResponse(sCtx, respAll)
			requireLenResponse(sCtx, len(respAll))

			compareIteratorsResponse(sCtx, len(respNext), len(respAll))
		*/
	})
	return respAll
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
