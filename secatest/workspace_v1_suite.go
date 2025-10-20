package secatest

import (
	"context"
	"log/slog"

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
	slog.Info("Starting " + suite.scenarioName)

	t.Title(suite.scenarioName)
	configureTags(t, secalib.WorkspaceProviderV1, secalib.WorkspaceKind)

	// Generate scenario data
	workspaceName := secalib.GenerateWorkspaceName()
	workspaceResource := secalib.GenerateWorkspaceResource(suite.tenant, workspaceName)

	// Setup mock, if configured to use
	if suite.mockEnabled {
		mockParams := &mock.WorkspaceParamsV1{
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
		}
		wm, err := mock.ConfigWorkspaceLifecycleScenarioV1(suite.scenarioName, mockParams)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}
		suite.mockClient = wm
	}

	ctx := context.Background()
	var tref *secapi.TenantReference
	var resp *schema.Workspace
	var expectedMeta *schema.RegionalResourceMetadata
	var expectedLabels schema.Labels

	// Create Workspace
	ws := &schema.Workspace{
		Labels: schema.Labels{
			secalib.EnvLabel: secalib.EnvDevelopmentLabel,
		},
		Metadata: &schema.RegionalResourceMetadata{
			Tenant: suite.tenant,
			Name:   workspaceName,
		},
	}
	expectedMeta = secalib.NewRegionalResourceMetadata(workspaceName, secalib.WorkspaceProviderV1, workspaceResource, secalib.ApiVersion1, secalib.WorkspaceKind,
		suite.tenant, suite.region)
	expectedLabels = schema.Labels{secalib.EnvLabel: secalib.EnvDevelopmentLabel}
	CreateOrUpdateWorkspaceV1Step("Create a Workspace", t, ctx, &suite.testSuite, suite.client.WorkspaceV1, ws, expectedMeta, expectedLabels, secalib.CreatingStatusState)

	// Get the created Workspace
	tref = &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
		Name:   workspaceName,
	}
	resp = GetWorkspaceV1Step("Get the created Workspace", t, ctx, &suite.testSuite, suite.client.WorkspaceV1, *tref, expectedMeta, expectedLabels, secalib.ActiveStatusState)

	// Update the Workspace
	resp.Labels = schema.Labels{
		secalib.EnvLabel: secalib.EnvProductionLabel,
	}
	expectedLabels = schema.Labels{secalib.EnvLabel: secalib.EnvProductionLabel}
	CreateOrUpdateWorkspaceV1Step("Update the Workspace", t, ctx, &suite.testSuite, suite.client.WorkspaceV1, resp, expectedMeta, expectedLabels, secalib.UpdatingStatusState)

	// Get the updated Workspace	
	resp = GetWorkspaceV1Step("Get the updated Workspace", t, ctx, &suite.testSuite, suite.client.WorkspaceV1, *tref, expectedMeta, expectedLabels, secalib.ActiveStatusState)

	// Delete the Workspace
	DeleteWorkspaceV1Step("Delete the Workspace", t, ctx, &suite.testSuite, suite.client.WorkspaceV1, resp)

	// Get the deleted Workspace
	GetWorkspaceWithErrorV1Step("Get the deleted Workspace", t, ctx, &suite.testSuite, suite.client.WorkspaceV1, *tref, secapi.ErrResourceNotFound)

	slog.Info("Finishing " + suite.scenarioName)
}

func (suite *WorkspaceV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}
