package secatest

import (
	"context"
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/secalib/builders"
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
	configureTags(t, secalib.WorkspaceProviderV1, string(schema.RegionalResourceMetadataKindResourceKindWorkspace))

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
			Workspace: &[]mock.ResourceParams[schema.WorkspaceSpec]{
				{
					Name: workspaceName,
					InitialLabels: schema.Labels{
						secalib.EnvLabel: secalib.EnvDevelopmentLabel,
					},
					UpdatedLabels: schema.Labels{
						secalib.EnvLabel: secalib.EnvConformanceLabel,
					},
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
			secalib.EnvLabel: secalib.EnvDevelopmentLabel,
		},
		Metadata: &schema.RegionalResourceMetadata{
			Tenant: suite.tenant,
			Name:   workspaceName,
		},
	}
	expectMeta, err := builders.NewRegionalResourceMetadataBuilder().
		Name(workspaceName).
		Provider(secalib.WorkspaceProviderV1).
		Resource(workspaceResource).
		ApiVersion(secalib.ApiVersion1).
		Kind(schema.RegionalResourceMetadataKindResourceKindWorkspace).
		Tenant(suite.tenant).
		Region(suite.region).
		BuildResponse()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectLabels := schema.Labels{secalib.EnvLabel: secalib.EnvDevelopmentLabel}
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
		secalib.EnvLabel: secalib.EnvProductionLabel,
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

	// Delete the workspace
	suite.deleteWorkspaceV1Step("Delete the workspace", t, suite.client.WorkspaceV1, workspace)

	// Get the deleted workspace
	suite.getWorkspaceWithErrorV1Step("Get the deleted workspace", t, suite.client.WorkspaceV1, *tref, secapi.ErrResourceNotFound)

	slog.Info("Finishing " + suite.scenarioName)
}

func (suite *WorkspaceV1TestSuite) TestSuiteList(t provider.T) {
	slog.Info("Starting " + suite.scenarioName)

	t.Title(suite.scenarioName)
	configureTags(t, secalib.WorkspaceProviderV1, string(schema.RegionalResourceMetadataKindResourceKindWorkspace))

	// Generate scenario data
	workspaceName := secalib.GenerateWorkspaceName()
	workspaceName2 := secalib.GenerateWorkspaceName()

	// Setup mock, if configured to use
	if suite.mockEnabled {
		mockParams := &mock.WorkspaceParamsV1{
			Params: &mock.Params{
				MockURL:   *suite.mockServerURL,
				AuthToken: suite.authToken,
				Tenant:    suite.tenant,
				Region:    suite.region,
			},
			Workspace: &[]mock.ResourceParams[schema.WorkspaceSpec]{
				{
					Name: workspaceName,
					InitialLabels: schema.Labels{
						secalib.EnvLabel: secalib.EnvConformanceLabel,
					},
				},
				{
					Name: workspaceName2,
					InitialLabels: schema.Labels{
						secalib.EnvLabel: secalib.EnvConformanceLabel,
					},
				},
			},
		}
		wm, err := mock.ConfigWorkspaceListLifecycleScenarioV1(suite.scenarioName, mockParams)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}
		suite.mockClient = wm
	}

	ctx := context.Background()

	// Create a workspace
	workspaces := &[]schema.Workspace{
		{
			Labels: schema.Labels{
				secalib.EnvLabel: secalib.EnvConformanceLabel,
			},
			Metadata: &schema.RegionalResourceMetadata{
				Tenant: suite.tenant,
				Name:   workspaceName,
			},
		},
		{
			Labels: schema.Labels{
				secalib.EnvLabel: secalib.EnvConformanceLabel,
			},
			Metadata: &schema.RegionalResourceMetadata{
				Tenant: suite.tenant,
				Name:   workspaceName2,
			},
		},
	}
	for _, workspace := range *workspaces {
		workspaceResource := secalib.GenerateWorkspaceResource(suite.tenant, workspace.Metadata.Name)
		expectMeta, err := builders.NewRegionalResourceMetadataBuilder().
			Name(workspace.Metadata.Name).
			Provider(secalib.WorkspaceProviderV1).
			Resource(workspaceResource).
			ApiVersion(secalib.ApiVersion1).
			Kind(schema.RegionalResourceMetadataKindResourceKindWorkspace).
			Tenant(suite.tenant).
			Region(suite.region).
			BuildResponse()
		if err != nil {
			t.Fatalf("Failed to build metadata: %v", err)
		}
		expectLabels := schema.Labels{secalib.EnvLabel: secalib.EnvDevelopmentLabel}
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
	suite.getListWorkspaceV1Step("list workspace", t, ctx, suite.client.WorkspaceV1, *tref,
		nil)
	// List workspaces with limit
	suite.getListWorkspaceV1Step("list workspace", t, ctx, suite.client.WorkspaceV1, *tref,
		secapi.NewListOptions().WithLimit(1))
	// List workspaces with label
	suite.getListWorkspaceV1Step("list workspace", t, ctx, suite.client.WorkspaceV1, *tref,
		secapi.NewListOptions().WithLabels(builders.NewLabelsBuilder().
			Equals(secalib.EnvLabel, secalib.EnvConformanceLabel)))
	// List workspaces with label and limit
	suite.getListWorkspaceV1Step("list workspace", t, ctx, suite.client.WorkspaceV1, *tref,
		secapi.NewListOptions().WithLimit(1).WithLabels(builders.NewLabelsBuilder().
			Equals(secalib.EnvLabel, secalib.EnvConformanceLabel)))
	// Resources deletion

	// Delete all workspaces
	for _, workspace := range *workspaces {
		workspaceTRef := &secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   workspace.Metadata.Name,
		}
		suite.deleteWorkspaceV1Step("Delete workspace 1", t, suite.client.WorkspaceV1, &workspace)
		suite.getWorkspaceWithErrorV1Step("Get deleted workspace 1", t, suite.client.WorkspaceV1, *workspaceTRef, secapi.ErrResourceNotFound)

	}
	slog.Info("Finishing " + suite.scenarioName)
}

func (suite *WorkspaceV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}
