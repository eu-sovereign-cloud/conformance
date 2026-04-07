package workspace

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockWorkspace "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/workspace"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	labelBuilder "github.com/eu-sovereign-cloud/go-sdk/secapi/builders"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type ProviderQueriesV1TestSuite struct {
	suites.RegionalTestSuite

	params *params.WorkspaceProviderQueriesV1Params
}

func CreateProviderQueriesV1TestSuite(regionalTestSuite suites.RegionalTestSuite) *ProviderQueriesV1TestSuite {
	suite := &ProviderQueriesV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
	}
	suite.ScenarioName = constants.WorkspaceProviderQueriesV1SuiteName.String()
	return suite
}

func (suite *ProviderQueriesV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("Workspace")

	// Generate scenario data
	workspaceName := generators.GenerateWorkspaceName()
	workspaceName2 := generators.GenerateWorkspaceName()

	workspace, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}
	workspace2, err := builders.NewWorkspaceBuilder().
		Name(workspaceName2).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	workspaces := []schema.Workspace{*workspace, *workspace2}

	params := &params.WorkspaceProviderQueriesV1Params{
		Workspaces: workspaces,
	}
	suite.params = params
	err = suites.SetupMockIfEnabled(suite.TestSuite, mockWorkspace.ConfigureProviderQueriesV1, *params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *ProviderQueriesV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, sdkconsts.WorkspaceProviderV1Name, string(schema.RegionalResourceMetadataKindResourceKindWorkspace))

	stepsConfigurator := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Workspace
	workspaces := suite.params.Workspaces

	// Create workspaces
	for _, workspace := range workspaces {
		expectMeta := workspace.Metadata
		expectLabels := schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}
		stepsConfigurator.CreateOrUpdateWorkspaceV1Step("Create a workspace", suite.Client.WorkspaceV1, &workspace,
			steps.StepResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
				Labels:         expectLabels,
				Metadata:       expectMeta,
				ResourceStates: suites.CreatedResourceExpectedStates,
			},
		)

	}

	tpath := secapi.TenantPath{
		Tenant: secapi.TenantID(suite.Tenant),
	}

	// List workspaces
	stepsConfigurator.ListWorkspaceV1Step("list workspace", suite.Client.WorkspaceV1, tpath, nil)

	// List workspaces with limit
	stepsConfigurator.ListWorkspaceV1Step("list workspace", suite.Client.WorkspaceV1, tpath, secapi.NewListOptions().WithLimit(1))

	// List workspaces with label
	stepsConfigurator.ListWorkspaceV1Step("list workspace", suite.Client.WorkspaceV1, tpath,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List workspaces with label and limit
	stepsConfigurator.ListWorkspaceV1Step("list workspace", suite.Client.WorkspaceV1, tpath,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Delete all workspaces
	for _, workspace := range workspaces {
		stepsConfigurator.DeleteWorkspaceV1Step("Delete workspace 1", suite.Client.WorkspaceV1, &workspace)

		// Get the deleted workspace
		workspaceTRef := secapi.TenantReference{
			Tenant: secapi.TenantID(workspace.Metadata.Tenant),
			Name:   workspace.Metadata.Name,
		}
		stepsConfigurator.WatchWorkspaceUntilDeletedV1Step("Watch the workspace deletion", suite.Client.WorkspaceV1, workspaceTRef)
	}

	suite.FinishScenario()
}

func (suite *ProviderQueriesV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
