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

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type LifeCycleV1TestSuite struct {
	suites.RegionalTestSuite
	params *params.WorkspaceLifeCycleV1Params
}

func NewLifeCycleV1TestSuite(regionalTestSuite suites.RegionalTestSuite) *LifeCycleV1TestSuite {
	suite := &LifeCycleV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
	}
	suite.ScenarioName = constants.WorkspaceLifeCycleV1SuiteName.String()
	return suite
}

func (suite *LifeCycleV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("workspace")

	// Generate scenario data
	workspaceName := generators.GenerateWorkspaceName()

	workspaceInitial, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(constants.WorkspaceProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvDevelopmentLabel}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	workspaceUpdated, err := builders.NewWorkspaceBuilder().
		Name(workspaceName).
		Provider(constants.WorkspaceProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	params := &params.WorkspaceLifeCycleV1Params{
		WorkspaceInitial: workspaceInitial,
		WorkspaceUpdated: workspaceUpdated,
	}
	suite.params = params
	err = suites.SetupMockIfEnabled(suite.TestSuite, mockWorkspace.ConfigureLifecycleScenarioV1, params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *LifeCycleV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, constants.WorkspaceProviderV1, string(schema.RegionalResourceMetadataKindResourceKindWorkspace))

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Create a workspace
	workspace := suite.params.WorkspaceInitial
	expectMeta := workspace.Metadata
	expectLabels := workspace.Labels
	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Create a workspace", suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:        expectLabels,
			Metadata:      expectMeta,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created Workspace
	tref := secapi.TenantReference{
		Tenant: secapi.TenantID(workspace.Metadata.Tenant),
		Name:   workspace.Metadata.Name,
	}
	workspace = stepsBuilder.GetWorkspaceV1Step("Get the created workspace", suite.Client.WorkspaceV1, tref,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:        expectLabels,
			Metadata:      expectMeta,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Update the workspace labels
	workspace.Labels = suite.params.WorkspaceUpdated.Labels
	expectLabels = workspace.Labels
	stepsBuilder.CreateOrUpdateWorkspaceV1Step("Update the workspace", suite.Client.WorkspaceV1, workspace,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:        expectLabels,
			Metadata:      expectMeta,
			ResourceState: schema.ResourceStateUpdating,
		},
	)

	// Get the updated workspace
	workspace = stepsBuilder.GetWorkspaceV1Step("Get the updated workspace", suite.Client.WorkspaceV1, tref,
		steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
			Labels:        expectLabels,
			Metadata:      expectMeta,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Resources deletion
	stepsBuilder.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)
	stepsBuilder.GetWorkspaceWithErrorV1Step("Get the deleted workspace", suite.Client.WorkspaceV1, tref, secapi.ErrResourceNotFound)

	suite.FinishScenario()
}

func (suite *LifeCycleV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
