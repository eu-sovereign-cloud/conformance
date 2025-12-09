package secatest

import (
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/secalib/builders"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/secalib/generators"
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
	configureTags(t, workspaceProviderV1, string(schema.RegionalResourceMetadataKindResourceKindWorkspace))

	// Generate scenario data
	workspaceName := generators.GenerateWorkspaceName()
	workspaceResource := generators.GenerateWorkspaceResource(suite.tenant, workspaceName)

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
					envLabel: envDevelopmentLabel,
				},
				UpdatedLabels: schema.Labels{
					envLabel: envProductionLabel,
				},
			},
		}
		wm, err := mock.ConfigWorkspaceLifecycleScenarioV1(suite.scenarioName, mockParams)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}
		suite.mockClient = wm
	}

	// Create a workspace
	workspace := &schema.Workspace{
		Labels: schema.Labels{
			envLabel: envDevelopmentLabel,
		},
		Metadata: &schema.RegionalResourceMetadata{
			Tenant: suite.tenant,
			Name:   workspaceName,
		},
	}
	expectMeta, err := builders.NewRegionalResourceMetadataBuilder().
		Name(workspaceName).Resource(workspaceResource).
		Provider(workspaceProviderV1).ApiVersion(apiVersion1).
		Kind(schema.RegionalResourceMetadataKindResourceKindWorkspace).
		Tenant(suite.tenant).Region(suite.region).
		BuildResponse()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectLabels := schema.Labels{envLabel: envDevelopmentLabel}
	suite.createOrUpdateWorkspaceV1Step("Create a workspace", t, suite.client.WorkspaceV1, workspace,
		responseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			labels:        expectLabels,
			metadata:      expectMeta,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created Workspace
	tref := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
		Name:   workspaceName,
	}
	workspace = suite.getWorkspaceV1Step("Get the created workspace", t, suite.client.WorkspaceV1, *tref,
		responseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			labels:        expectLabels,
			metadata:      expectMeta,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Update the workspace labels
	workspace.Labels = schema.Labels{
		envLabel: envProductionLabel,
	}
	expectLabels = workspace.Labels
	suite.createOrUpdateWorkspaceV1Step("Update the workspace", t, suite.client.WorkspaceV1, workspace,
		responseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			labels:        expectLabels,
			metadata:      expectMeta,
			resourceState: schema.ResourceStateUpdating,
		},
	)

	// Get the updated workspace
	workspace = suite.getWorkspaceV1Step("Get the updated workspace", t, suite.client.WorkspaceV1, *tref,
		responseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			labels:        expectLabels,
			metadata:      expectMeta,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Resources deletion

	suite.deleteWorkspaceV1Step("Delete the workspace", t, suite.client.WorkspaceV1, workspace)
	suite.getWorkspaceWithErrorV1Step("Get the deleted workspace", t, suite.client.WorkspaceV1, *tref, secapi.ErrResourceNotFound)

	slog.Info("Finishing " + suite.scenarioName)
}

func (suite *WorkspaceV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}
