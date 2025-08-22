package secatest

import (
	"context"
	"log/slog"
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
	configureTags(t, workspaceV1Provider, workspaceKind)

	// Generate scenario data
	workspaceName := suite.generateWorkspaceName()
	resource := suite.generateWorkspaceResource(workspaceName)

	// Setup mock, if configured to use
	if suite.isMockEnabled() {
		wm, err := mock.CreateWorkspaceLifecycleScenarioV1("Workspace Lifecycle",
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

	expectedMetadata := verifyRegionalMetadataStepParams{
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
	t.WithNewStep("Create workspace", func(sCtx provider.StepCtx) {
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

		expectedMetadata.verb = http.MethodPut
		verifyWorkspaceMetadataStep(sCtx, expectedMetadata, resp.Metadata)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: creatingStatusState,
			actualState:   string(*resp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
	})

	// Step 2
	t.WithNewStep("Get created workspace", func(sCtx provider.StepCtx) {
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

		expectedMetadata.verb = http.MethodGet
		verifyWorkspaceMetadataStep(sCtx, expectedMetadata, resp.Metadata)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: activeStatusState,
			actualState:   string(*resp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
	})

	// Step 3
	t.WithNewStep("Update workspace", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateWorkspace",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		resp, err = suite.client.WorkspaceV1.CreateOrUpdateWorkspace(ctx, resp)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		expectedMetadata.verb = http.MethodPut
		verifyWorkspaceMetadataStep(sCtx, expectedMetadata, resp.Metadata)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: updatingStatusState,
			actualState:   string(*resp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
	})

	// Step 4
	t.WithNewStep("Get updated workspace", func(sCtx provider.StepCtx) {
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

		expectedMetadata.verb = http.MethodGet
		verifyWorkspaceMetadataStep(sCtx, expectedMetadata, resp.Metadata)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: activeStatusState,
			actualState:   string(*resp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
	})

	// Step 5
	t.WithNewStep("Delete workspace", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "DeleteWorkspace",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		err = suite.client.WorkspaceV1.DeleteWorkspace(ctx, resp)
		requireNoError(sCtx, err)
	})

	// Step 6
	t.WithNewStep("Re-delete workspace", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "DeleteWorkspace",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		err = suite.client.WorkspaceV1.DeleteWorkspace(ctx, resp)
		requireNoError(sCtx, err)
	})

	// Step 7
	t.WithNewStep("Get deleted workspace", func(sCtx provider.StepCtx) {
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
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})
}

func (suite *WorkspaceV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}

func verifyWorkspaceMetadataStep(ctx provider.StepCtx, expected verifyRegionalMetadataStepParams, metadata *workspace.RegionalResourceMetadata) {
	actualMetadata := verifyRegionalMetadataStepParams{
		name:       metadata.Name,
		provider:   metadata.Provider,
		verb:       metadata.Verb,
		resource:   metadata.Resource,
		apiVersion: metadata.ApiVersion,
		kind:       string(metadata.Kind),
		tenant:     metadata.Tenant,
		region:     metadata.Region,
	}
	verifyRegionalMetadataStep(ctx, expected, actualMetadata)
}
