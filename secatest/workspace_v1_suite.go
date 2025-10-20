package secatest

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type WorkspaceV1TestSuite struct {
	regionalTestSuite
}

func (suite *WorkspaceV1TestSuite) TestSuite(t provider.T) {
	ctx := context.Background()
	var err error
	slog.Info("Starting " + suite.scenarioName)

	t.Title(suite.scenarioName)
	configureTags(t, secalib.WorkspaceProviderV1, secalib.WorkspaceKind)

	// Generate scenario data
	workspaceName := secalib.GenerateWorkspaceName()
	workspaceResource := secalib.GenerateWorkspaceResource(suite.tenant, workspaceName)

	// Setup mock, if configured to use
	if suite.mockEnabled {
		wm, err := mock.CreateWorkspaceLifecycleScenarioV1(suite.scenarioName, &mock.WorkspaceParamsV1{
			Params: &mock.Params{
				MockURL:   *suite.mockServerURL,
				AuthToken: suite.authToken,
				Tenant:    suite.tenant,
				Region:    suite.region,
			},
			Workspace: &mock.ResourceParams[schema.WorkspaceSpec]{
				Name: workspaceName,
				InitialLabels: schema.Labels{
					secalib.EnvLabel: secalib.EnvDevelopmentLabel,
				},
				UpdatedLabels: schema.Labels{
					secalib.EnvLabel: secalib.EnvProductionLabel,
				},
			},
		})
		if err != nil {
			t.Fatalf("Failed to create wiremock scenario: %v", err)
		}
		suite.mockClient = wm
	}

	var resp *schema.Workspace
	var expectedMeta *schema.RegionalResourceMetadata
	var expectedLabels schema.Labels

	t.WithNewStep("Create workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "CreateOrUpdateWorkspace")

		// TODO Create a builder to help to create this object to use on the client
		ws := &schema.Workspace{
			Metadata: &schema.RegionalResourceMetadata{
				Tenant: suite.tenant,
				Name:   workspaceName,
			},
		}
		resp, err = suite.client.WorkspaceV1.CreateOrUpdateWorkspace(ctx, ws)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		expectedMeta = &schema.RegionalResourceMetadata{
			Name:       workspaceName,
			Provider:   secalib.WorkspaceProviderV1,
			Resource:   workspaceResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.WorkspaceKind,
			Tenant:     suite.tenant,
			Region:     suite.region,
		}
		verifyRegionalResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)

		verifyStatusStep(sCtx, secalib.CreatingStatusState, *resp.Status.State)
	})

	t.WithNewStep("Get created workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "GetWorkspace")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   workspaceName,
		}
		resp, err = suite.client.WorkspaceV1.GetWorkspace(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		expectedMeta.Verb = http.MethodGet
		verifyRegionalResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)

		expectedLabels = schema.Labels{
			secalib.EnvLabel: secalib.EnvDevelopmentLabel,
		}
		verifyLabelStep(sCtx, expectedLabels, resp.Labels)

		verifyStatusStep(sCtx, secalib.ActiveStatusState, *resp.Status.State)
	})

	t.WithNewStep("Update workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "CreateOrUpdateWorkspace")

		resp, err = suite.client.WorkspaceV1.CreateOrUpdateWorkspace(ctx, resp)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		expectedMeta.Verb = http.MethodPut
		verifyRegionalResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)

		verifyStatusStep(sCtx, secalib.UpdatingStatusState, *resp.Status.State)
	})

	t.WithNewStep("Get updated workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "GetWorkspace")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   workspaceName,
		}
		resp, err = suite.client.WorkspaceV1.GetWorkspace(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		expectedMeta.Verb = http.MethodGet
		verifyRegionalResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)

		expectedLabels = schema.Labels{
			secalib.EnvLabel: secalib.EnvDevelopmentLabel,
		}
		verifyLabelStep(sCtx, expectedLabels, resp.Labels)

		verifyStatusStep(sCtx, secalib.ActiveStatusState, *resp.Status.State)
	})

	t.WithNewStep("Delete workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "DeleteWorkspace")

		err = suite.client.WorkspaceV1.DeleteWorkspace(ctx, resp)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get deleted workspace", func(sCtx provider.StepCtx) {
		suite.setWorkspaceV1StepParams(sCtx, "GetWorkspace")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   workspaceName,
		}
		_, err = suite.client.WorkspaceV1.GetWorkspace(ctx, tref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	slog.Info("Finishing " + suite.scenarioName)
}

func (suite *WorkspaceV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}
