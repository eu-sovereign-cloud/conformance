package workspace

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockWorkspace "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/workspace"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	labelBuilder "github.com/eu-sovereign-cloud/go-sdk/secapi/builders"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type ListV1TestSuite struct {
	suites.RegionalTestSuite

	params params.WorkspaceListParamsV1
}

func (suite *ListV1TestSuite) BeforeAll(t provider.T) {
	var err error

	// Generate scenario data
	workspaceName := generators.GenerateWorkspaceName()
	workspaceName2 := generators.GenerateWorkspaceName()

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(constants.WorkspaceProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}
	workspace2, err := builders.NewWorkspaceBuilder().
		Name(workspaceName2).
		Provider(constants.WorkspaceProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	workspaces := []schema.Workspace{*workspace, *workspace2}

	params := params.WorkspaceListParamsV1{
		Workspaces: workspaces,
	}
	suite.params = params
	err = suites.SetupMockIfEnabled(&suite.TestSuite, mockWorkspace.ConfigureListScenarioV1, &suite.params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *ListV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, constants.WorkspaceProviderV1, string(schema.RegionalResourceMetadataKindResourceKindWorkspace))

	stepsBuilder := steps.NewStepsConfigurator(&suite.TestSuite, t)

	workspaces := suite.params.Workspaces

	// Create a workspace

	for _, workspace := range workspaces {

		expectMeta := workspace.Metadata
		expectLabels := schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}
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
	stepsBuilder.GetListWorkspaceV1Step("list workspace", suite.Client.WorkspaceV1, *tref, nil)
	// List workspaces with limit
	stepsBuilder.GetListWorkspaceV1Step("list workspace", suite.Client.WorkspaceV1, *tref, secapi.NewListOptions().WithLimit(1))
	// List workspaces with label
	stepsBuilder.GetListWorkspaceV1Step("list workspace", suite.Client.WorkspaceV1, *tref,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().Equals(constants.EnvLabel, constants.EnvConformanceLabel)))
	// List workspaces with label and limit
	stepsBuilder.GetListWorkspaceV1Step("list workspace", suite.Client.WorkspaceV1, *tref,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Delete all workspaces
	for _, workspace := range workspaces {
		workspaceTRef := &secapi.TenantReference{
			Tenant: secapi.TenantID(workspace.Metadata.Tenant),
			Name:   workspace.Metadata.Name,
		}
		stepsBuilder.DeleteWorkspaceV1Step("Delete workspace 1", suite.Client.WorkspaceV1, &workspace)
		stepsBuilder.GetWorkspaceWithErrorV1Step("Get deleted workspace 1", suite.Client.WorkspaceV1, *workspaceTRef, secapi.ErrResourceNotFound)
	}

	suite.FinishScenario()
}

func (suite *ListV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
