package secatest

import (
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	labelBuilder "github.com/eu-sovereign-cloud/go-sdk/secapi/builders"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type WorkspaceV1TestSuite struct {
	regionalTestSuite
}

func (suite *WorkspaceV1TestSuite) TestLifeCycleScenario(t provider.T) {
	suite.startScenario(t)
	configureTags(t, workspaceProviderV1, string(schema.RegionalResourceMetadataKindResourceKindWorkspace))

	// Generate scenario data
	workspaceName := generators.GenerateWorkspaceName()

	// Setup mock, if configured to use
	if suite.mockEnabled {
		mockParams := &mock.WorkspaceLifeCycleParamsV1{
			BaseParams: &mock.BaseParams{
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
		wm, err := mock.ConfigureWorkspaceLifecycleScenarioV1(suite.scenarioName, mockParams)
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
	expectMeta, err := builders.NewWorkspaceMetadataBuilder().
		Name(workspaceName).
		Provider(workspaceProviderV1).ApiVersion(apiVersion1).
		Tenant(suite.tenant).Region(suite.region).
		Build()
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

	suite.finishScenario()
}

func (suite *WorkspaceV1TestSuite) TestListScenario(t provider.T) {
	suite.startScenario(t)
	configureTags(t, workspaceProviderV1, string(schema.RegionalResourceMetadataKindResourceKindWorkspace))

	// Generate scenario data
	workspaceName := generators.GenerateWorkspaceName()
	workspaceName2 := generators.GenerateWorkspaceName()

	// Setup mock, if configured to use
	if suite.mockEnabled {
		mockParams := &mock.WorkspaceListParamsV1{
			BaseParams: &mock.BaseParams{
				MockURL:   *suite.mockServerURL,
				AuthToken: suite.authToken,
				Tenant:    suite.tenant,
				Region:    suite.region,
			},
			Workspaces: []mock.ResourceParams[schema.WorkspaceSpec]{
				{
					Name: workspaceName,
					InitialLabels: schema.Labels{
						generators.EnvLabel: generators.EnvConformanceLabel,
					},
				},
				{
					Name: workspaceName2,
					InitialLabels: schema.Labels{
						generators.EnvLabel: generators.EnvConformanceLabel,
					},
				},
			},
		}
		wm, err := mock.ConfigureWorkspaceListScenarioV1(suite.scenarioName, mockParams)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}
		suite.mockClient = wm
	}

	// Create a workspace
	workspaces := &[]schema.Workspace{
		{
			Labels: schema.Labels{
				generators.EnvLabel: generators.EnvConformanceLabel,
			},
			Metadata: &schema.RegionalResourceMetadata{
				Tenant: suite.tenant,
				Name:   workspaceName,
			},
		},
		{
			Labels: schema.Labels{
				generators.EnvLabel: generators.EnvConformanceLabel,
			},
			Metadata: &schema.RegionalResourceMetadata{
				Tenant: suite.tenant,
				Name:   workspaceName2,
			},
		},
	}
	for _, workspace := range *workspaces {
		expectMeta, err := builders.NewWorkspaceMetadataBuilder().
			Name(workspace.Metadata.Name).
			Provider(workspaceProviderV1).ApiVersion(apiVersion1).
			Tenant(suite.tenant).Region(suite.region).
			Build()
		if err != nil {
			t.Fatalf("Failed to build metadata: %v", err)
		}
		expectLabels := schema.Labels{generators.EnvLabel: generators.EnvConformanceLabel}
		suite.createOrUpdateWorkspaceV1Step("Create a workspace", t, suite.client.WorkspaceV1, &workspace,
			responseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
				labels:        expectLabels,
				metadata:      expectMeta,
				resourceState: schema.ResourceStateCreating,
			},
		)

	}
	tref := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
	}

	// List workspaces
	suite.getListWorkspaceV1Step("list workspace", t, suite.client.WorkspaceV1, *tref,
		nil)
	// List workspaces with limit
	suite.getListWorkspaceV1Step("list workspace", t, suite.client.WorkspaceV1, *tref,
		secapi.NewListOptions().WithLimit(1))
	// List workspaces with label
	suite.getListWorkspaceV1Step("list workspace", t, suite.client.WorkspaceV1, *tref,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(generators.EnvLabel, generators.EnvConformanceLabel)))
	// List workspaces with label and limit
	suite.getListWorkspaceV1Step("list workspace", t, suite.client.WorkspaceV1, *tref,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(generators.EnvLabel, generators.EnvConformanceLabel)))

	// Delete all workspaces
	for _, workspace := range *workspaces {
		workspaceTRef := &secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   workspace.Metadata.Name,
		}
		suite.deleteWorkspaceV1Step("Delete workspace 1", t, suite.client.WorkspaceV1, &workspace)
		suite.getWorkspaceWithErrorV1Step("Get deleted workspace 1", t, suite.client.WorkspaceV1, *workspaceTRef, secapi.ErrResourceNotFound)

	}
	suite.finishScenario()
}

func (suite *WorkspaceV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}
