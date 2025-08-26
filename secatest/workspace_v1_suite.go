package secatest

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/secalib"
	workspace "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.workspace.v1"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type WorkspaceV1TestSuite struct {
	regionalTestSuite
}

func (suite *WorkspaceV1TestSuite) TestWorkspaceV1(t provider.T) {
	t.Title("Workspace Lifecycle Test")
	configureTags(t, secalib.WorkspaceProviderV1, secalib.WorkspaceKind)

	// Generate scenario data
	workspaceName := secalib.GenerateWorkspaceName()
	resource := secalib.GenerateWorkspaceResource(suite.tenant, workspaceName)

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

	ctx := context.Background()
	var resp *workspace.Workspace
	var err error

	var expectedMetadata verifyRegionalMetadataStepParams

	// Step 1
	t.WithNewStep("Create workspace", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateWorkspace",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   workspaceName,
		}
		ws := &workspace.Workspace{}
		resp, err = suite.client.WorkspaceV1.CreateOrUpdateWorkspace(ctx, tref, ws, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		expectedMetadata = verifyRegionalMetadataStepParams{
			name:       workspaceName,
			provider:   secalib.WorkspaceProviderV1,
			resource:   resource,
			verb:       http.MethodPut,
			apiVersion: secalib.ApiVersion1,
			kind:       secalib.WorkspaceKind,
			tenant:     suite.tenant,
			region:     suite.region,
		}
		verifyWorkspaceMetadataStep(sCtx, expectedMetadata, resp.Metadata)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: secalib.CreatingStatusState,
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
			expectedState: secalib.ActiveStatusState,
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

		// TODO Find a attribute to update

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   workspaceName,
		}
		resp, err = suite.client.WorkspaceV1.CreateOrUpdateWorkspace(ctx, tref, resp, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		expectedMetadata.verb = http.MethodPut
		verifyWorkspaceMetadataStep(sCtx, expectedMetadata, resp.Metadata)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: secalib.UpdatingStatusState,
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
			expectedState: secalib.ActiveStatusState,
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

		err = suite.client.WorkspaceV1.DeleteWorkspace(ctx, resp, nil)
		requireNoError(sCtx, err)
	})

	// Step 6
	t.WithNewStep("Re-delete workspace", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "DeleteWorkspace",
			tenantStepParameter, suite.tenant,
			workspaceStepParameter, workspaceName,
		)

		err = suite.client.WorkspaceV1.DeleteWorkspace(ctx, resp, nil)
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
		_, err = suite.client.WorkspaceV1.GetWorkspace(ctx, tref)
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
