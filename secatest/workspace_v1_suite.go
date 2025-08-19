package secatest

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"math/rand"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	workspace "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.workspace.v1"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type WorkspaceV1TestSuite struct {
	regionalTestSuite
}

func (suite *WorkspaceV1TestSuite) TestWorkspaceV1(t provider.T) {
	t.Title("Workspace Lifecycle Test")
	t.Tags(
		"provider:"+workspaceV1Provider,
		"resource:"+workspaceKind,
	)

	// Generate scenario data
	workspaceName := fmt.Sprintf("workspace-%d", rand.Intn(math.MaxInt32))
	resource := fmt.Sprintf(workspaceResource, suite.tenant, workspaceName)

	// Setup mock, if configured to use
	if suite.mockEnabled == "true" {
		wm, err := mock.CreateWorkspaceScenarioV1("Workspace Lifecycle",
			mock.WorkspaceParamsV1{
				Params: mock.Params{
					MockURL:   suite.mockServerURL,
					AuthToken: suite.authToken,
					Tenant:    suite.tenant,
					Region:    suite.region,
				},
				Name: workspaceName,
			})
		if err != nil {
			slog.Error("Failed to create workspace scenario", "error", err)
			return
		}
		suite.mockClient = wm
	}

	expectedMetadataParams := verifyRegionalMetadataStepParams{
		name:       workspaceName,
		provider:   workspaceV1Provider,
		resource:   resource,
		apiVersion: version1,
		kind:       workspaceKind,
		tenant:     suite.tenant,
		region:     suite.region,
	}

	ctx := context.Background()
	var resp *workspace.Workspace
	var err error

	// Step 1
	t.WithNewStep("Create Workspace", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateWorkspace",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		ws := &workspace.Workspace{
			Metadata: &workspace.RegionalResourceMetadata{
				Tenant: suite.tenant,
				Name:   workspaceName,
			},
		}
		resp, err = suite.client.WorkspaceV1.CreateOrUpdateWorkspace(ctx, ws)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		expectedMetadataParams.verb = http.MethodPut
		actualMetadataParams := verifyRegionalMetadataStepParams{
			name:       resp.Metadata.Name,
			provider:   resp.Metadata.Provider,
			verb:       resp.Metadata.Verb,
			resource:   resp.Metadata.Resource,
			apiVersion: resp.Metadata.ApiVersion,
			kind:       string(resp.Metadata.Kind),
			tenant:     resp.Metadata.Tenant,
			region:     resp.Metadata.Region,
		}
		verifyRegionalMetadataStep(sCtx, expectedMetadataParams, actualMetadataParams)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: creatingStatusState,
			actualState:   string(*resp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
	})

	// Step 2
	t.WithNewStep("Get Created Workspace", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetWorkspace",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   workspaceName,
		}
		resp, err = suite.client.WorkspaceV1.GetWorkspace(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		expectedMetadataParams.verb = http.MethodGet
		actualMetadataParams := verifyRegionalMetadataStepParams{
			name:       resp.Metadata.Name,
			provider:   resp.Metadata.Provider,
			verb:       resp.Metadata.Verb,
			resource:   resp.Metadata.Resource,
			apiVersion: resp.Metadata.ApiVersion,
			kind:       string(resp.Metadata.Kind),
			tenant:     resp.Metadata.Tenant,
			region:     resp.Metadata.Region,
		}
		verifyRegionalMetadataStep(sCtx, expectedMetadataParams, actualMetadataParams)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: activeStatusState,
			actualState:   string(*resp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
	})

	// Step 3
	t.WithNewStep("Update Workspace", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateWorkspace",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		ws := &workspace.Workspace{
			Metadata: &workspace.RegionalResourceMetadata{
				Tenant: suite.tenant,
				Name:   workspaceName,
			},
		}
		resp, err = suite.client.WorkspaceV1.CreateOrUpdateWorkspace(ctx, ws)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		expectedMetadataParams.verb = http.MethodPut
		actualMetadataParams := verifyRegionalMetadataStepParams{
			name:       resp.Metadata.Name,
			provider:   resp.Metadata.Provider,
			verb:       resp.Metadata.Verb,
			resource:   resp.Metadata.Resource,
			apiVersion: resp.Metadata.ApiVersion,
			kind:       string(resp.Metadata.Kind),
			tenant:     resp.Metadata.Tenant,
			region:     resp.Metadata.Region,
		}
		verifyRegionalMetadataStep(sCtx, expectedMetadataParams, actualMetadataParams)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: updatingStatusState,
			actualState:   string(*resp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
	})

	// Step 4
	t.WithNewStep("Get Updated Workspace", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetWorkspace",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   workspaceName,
		}
		resp, err = suite.client.WorkspaceV1.GetWorkspace(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		expectedMetadataParams.verb = http.MethodGet
		actualMetadataParams := verifyRegionalMetadataStepParams{
			name:       resp.Metadata.Name,
			provider:   resp.Metadata.Provider,
			verb:       resp.Metadata.Verb,
			resource:   resp.Metadata.Resource,
			apiVersion: resp.Metadata.ApiVersion,
			kind:       string(resp.Metadata.Kind),
			tenant:     resp.Metadata.Tenant,
			region:     resp.Metadata.Region,
		}
		verifyRegionalMetadataStep(sCtx, expectedMetadataParams, actualMetadataParams)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: activeStatusState,
			actualState:   string(*resp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
	})

	// Step 5
	t.WithNewStep("Delete Workspace", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "DeleteWorkspace",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		err = suite.client.WorkspaceV1.DeleteWorkspace(ctx, resp)
		requireNoError(sCtx, err)
	})

	// Step 6
	t.WithNewStep("Re-Delete Workspace", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "DeleteWorkspace",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		err = suite.client.WorkspaceV1.DeleteWorkspace(ctx, resp)
		requireNoError(sCtx, err)
	})

	// Step 7
	t.WithNewStep("Get Deleted Workspace", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetWorkspace",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   workspaceName,
		}
		resp, err = suite.client.WorkspaceV1.GetWorkspace(ctx, tref)
		requireNoError(sCtx, err)
		requireNilResponse(sCtx, resp)
	})
}

func (suite *WorkspaceV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}
