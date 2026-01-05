package workspace

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/workspace"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	labelBuilder "github.com/eu-sovereign-cloud/go-sdk/secapi/builders"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type WorkspaceV1ListTestSuite struct {
	suites.RegionalTestSuite
}

func (suite *WorkspaceV1ListTestSuite) TestListScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, conformance.WorkspaceProviderV1, string(schema.RegionalResourceMetadataKindResourceKindWorkspace))

	// Generate scenario data
	workspaceName := generators.GenerateWorkspaceName()
	workspaceName2 := generators.GenerateWorkspaceName()

	// Setup mock, if configured to use
	if suite.MockEnabled {
		mockParams := &mock.WorkspaceListParamsV1{
			BaseParams: &mock.BaseParams{
				MockURL:   *suite.MockServerURL,
				AuthToken: suite.AuthToken,
				Tenant:    suite.Tenant,
				Region:    suite.Region,
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
		wm, err := workspace.ConfigureListScenarioV1(suite.ScenarioName, mockParams)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}
		suite.MockClient = wm
	}

	stepsBuilder := steps.NewBuilder(&suite.TestSuite, t)

	// Create a workspace
	workspaces := &[]schema.Workspace{
		{
			Labels: schema.Labels{
				generators.EnvLabel: generators.EnvConformanceLabel,
			},
			Metadata: &schema.RegionalResourceMetadata{
				Tenant: suite.Tenant,
				Name:   workspaceName,
			},
		},
		{
			Labels: schema.Labels{
				generators.EnvLabel: generators.EnvConformanceLabel,
			},
			Metadata: &schema.RegionalResourceMetadata{
				Tenant: suite.Tenant,
				Name:   workspaceName2,
			},
		},
	}
	for _, workspace := range *workspaces {
		expectMeta, err := builders.NewWorkspaceMetadataBuilder().
			Name(workspace.Metadata.Name).
			Provider(conformance.WorkspaceProviderV1).ApiVersion(conformance.ApiVersion1).
			Tenant(suite.Tenant).Region(suite.Region).
			Build()
		if err != nil {
			t.Fatalf("Failed to build Metadata: %v", err)
		}
		expectLabels := schema.Labels{generators.EnvLabel: generators.EnvConformanceLabel}
		stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", suite.Client.WorkspaceV1, &workspace,
			steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
				Labels:        expectLabels,
				Metadata:      expectMeta,
				ResourceState: schema.ResourceStateCreating,
			},
		)

	}
	tref := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.Tenant),
	}

	// List workspaces
	stepsBuilder.GetListWorkspaceV1Step("list workspace", suite.Client.WorkspaceV1, *tref,
		nil)
	// List workspaces with limit
	stepsBuilder.GetListWorkspaceV1Step("list workspace", suite.Client.WorkspaceV1, *tref,
		secapi.NewListOptions().WithLimit(1))
	// List workspaces with label
	stepsBuilder.GetListWorkspaceV1Step("list workspace", suite.Client.WorkspaceV1, *tref,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(generators.EnvLabel, generators.EnvConformanceLabel)))
	// List workspaces with label and limit
	stepsBuilder.GetListWorkspaceV1Step("list workspace", suite.Client.WorkspaceV1, *tref,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(generators.EnvLabel, generators.EnvConformanceLabel)))

	// Delete all workspaces
	for _, workspace := range *workspaces {
		workspaceTRef := &secapi.TenantReference{
			Tenant: secapi.TenantID(suite.Tenant),
			Name:   workspace.Metadata.Name,
		}
		stepsBuilder.DeleteWorkspaceV1Step("Delete workspace 1", suite.Client.WorkspaceV1, &workspace)
		stepsBuilder.GetWorkspaceWithErrorV1Step("Get deleted workspace 1", suite.Client.WorkspaceV1, *workspaceTRef, secapi.ErrResourceNotFound)

	}
	suite.FinishScenario()
}

func (suite *WorkspaceV1LifeCycleTestSuite) AfterEach(t provider.T) {
	suite.ResetAllScenarios()
}
