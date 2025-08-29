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
	slog.Info("Starting Workspace Lifecycle Test")

	t.Title("Workspace Lifecycle Test")
	configureTags(t, secalib.WorkspaceProviderV1, secalib.WorkspaceKind)

	// Generate scenario data
	workspaceName := secalib.GenerateWorkspaceName()
	resource := secalib.GenerateWorkspaceResource(suite.tenant, workspaceName)

	// Setup mock, if configured to use
	if suite.isMockEnabled() {
		wm, err := mock.CreateWorkspaceLifecycleScenarioV1("Workspace Lifecycle",
			mock.WorkspaceParamsV1{
				Params: &mock.Params{
					MockURL:   suite.mockServerURL,
					AuthToken: suite.authToken,
					Tenant:    suite.tenant,
					Region:    suite.region,
				},
				Workspace: &mock.ResourceParams[secalib.WorkspaceSpecV1]{
					Name:        workspaceName,
					InitialSpec: &secalib.WorkspaceSpecV1{},
				},
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

	var expectedMetadata *secalib.Metadata

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

		expectedMetadata = &secalib.Metadata{
			Name:       workspaceName,
			Provider:   secalib.WorkspaceProviderV1,
			Resource:   resource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.WorkspaceKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyWorkspaceMetadataStep(sCtx, expectedMetadata, resp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*resp.Status.State)},
		)
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

		expectedMetadata.Verb = http.MethodGet
		verifyWorkspaceMetadataStep(sCtx, expectedMetadata, resp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*resp.Status.State)},
		)
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

		expectedMetadata.Verb = http.MethodPut
		verifyWorkspaceMetadataStep(sCtx, expectedMetadata, resp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*resp.Status.State)},
		)
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

		expectedMetadata.Verb = http.MethodGet
		verifyWorkspaceMetadataStep(sCtx, expectedMetadata, resp.Metadata)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*resp.Status.State)},
		)
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

	slog.Info("Finishing Workspace Lifecycle Test")
}

func (suite *WorkspaceV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}

func verifyWorkspaceMetadataStep(ctx provider.StepCtx, expected *secalib.Metadata, metadata *workspace.RegionalResourceMetadata) {
	actualMetadata := &secalib.Metadata{
		Name:       metadata.Name,
		Provider:   metadata.Provider,
		Verb:       metadata.Verb,
		Resource:   metadata.Resource,
		ApiVersion: metadata.ApiVersion,
		Kind:       string(metadata.Kind),
		Tenant:     metadata.Tenant,
		Region:     metadata.Region,
	}
	verifyRegionalMetadataStep(ctx, expected, actualMetadata)
}
