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

	// Create a workspace
	workspace := &schema.Workspace{
		Labels: schema.Labels{
			secalib.EnvLabel: secalib.EnvDevelopmentLabel,
		},
		Metadata: &schema.RegionalResourceMetadata{
			Tenant: suite.tenant,
			Name:   workspaceName,
		},
	}
	expectMeta := secalib.NewRegionalResourceMetadata(workspaceName, secalib.WorkspaceProviderV1, workspaceResource, secalib.ApiVersion1, secalib.WorkspaceKind,
		suite.tenant, suite.region)
	expectLabels := schema.Labels{secalib.EnvLabel: secalib.EnvDevelopmentLabel}
	suite.createOrUpdateWorkspaceV1Step("Create a workspace", t, ctx, suite.client.WorkspaceV1, workspace, expectMeta, expectLabels, secalib.CreatingResourceState)

	// Get the created Workspace
	tref := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
		Name:   workspaceName,
	}
	workspace = suite.getWorkspaceV1Step("Get the created workspace", t, ctx, suite.client.WorkspaceV1, *tref, expectMeta, expectLabels, secalib.ActiveResourceState)

	// Update the workspace
	workspace.Labels = schema.Labels{
		secalib.EnvLabel: secalib.EnvProductionLabel,
	}
	expectLabels = schema.Labels{secalib.EnvLabel: secalib.EnvProductionLabel}
	suite.createOrUpdateWorkspaceV1Step("Update the workspace", t, ctx, suite.client.WorkspaceV1, workspace, expectMeta, expectLabels, secalib.UpdatingResourceState)

	// Get the updated workspace
	workspace = suite.getWorkspaceV1Step("Get the updated workspace", t, ctx, suite.client.WorkspaceV1, *tref, expectMeta, expectLabels, secalib.ActiveResourceState)

	// Delete the workspace
	suite.deleteWorkspaceV1Step("Delete the workspace", t, ctx, suite.client.WorkspaceV1, workspace)

	// Get the deleted workspace
	suite.getWorkspaceWithErrorV1Step("Get the deleted workspace", t, ctx, suite.client.WorkspaceV1, *tref, secapi.ErrResourceNotFound)

	slog.Info("Finishing " + suite.scenarioName)
}

func (suite *WorkspaceV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}
