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

type WorkspaceLifeCycleV1TestSuite struct {
	suites.RegionalTestSuite
	params *params.WorkspaceLifeCycleV1Params
}

func CreateLifeCycleV1TestSuite(regionalTestSuite suites.RegionalTestSuite) *WorkspaceLifeCycleV1TestSuite {
	suite := &WorkspaceLifeCycleV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
	}
	suite.ScenarioName = constants.WorkspaceV1LifeCycleSuiteName
	return suite
}

func (suite *WorkspaceLifeCycleV1TestSuite) BeforeAll(t provider.T) {
	var err error

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

//nolint:dupl
func (suite *WorkspaceLifeCycleV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, constants.WorkspaceProviderV1, string(schema.RegionalResourceMetadataKindResourceKindWorkspace))

	// Workspace lifecycle grouped
	workspace := suite.params.WorkspaceInitial
	expectMeta := workspace.Metadata
	expectLabels := workspace.Labels
	tref := secapi.TenantReference{
		Tenant: secapi.TenantID(workspace.Metadata.Tenant),
		Name:   workspace.Metadata.Name,
	}

	t.WithNewStep("Workspace", func(wsCtx provider.StepCtx) {
		wsSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, wsCtx)

		wsSteps.CreateOrUpdateWorkspaceV1Step("Create", suite.Client.WorkspaceV1, workspace,
			steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
				Labels:        expectLabels,
				Metadata:      expectMeta,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		workspace = wsSteps.GetWorkspaceV1Step("Get", suite.Client.WorkspaceV1, tref,
			steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
				Labels:        expectLabels,
				Metadata:      expectMeta,
				ResourceState: schema.ResourceStateActive,
			},
		)

		workspace.Labels = suite.params.WorkspaceUpdated.Labels
		expectLabels = workspace.Labels
		wsSteps.CreateOrUpdateWorkspaceV1Step("Update", suite.Client.WorkspaceV1, workspace,
			steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
				Labels:        expectLabels,
				Metadata:      expectMeta,
				ResourceState: schema.ResourceStateUpdating,
			},
		)

		workspace = wsSteps.GetWorkspaceV1Step("GetUpdated", suite.Client.WorkspaceV1, tref,
			steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
				Labels:        expectLabels,
				Metadata:      expectMeta,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	t.WithNewStep("Delete", func(delCtx provider.StepCtx) {
		// reuse root configurator for delete checks
		delSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, delCtx)
		delSteps.DeleteWorkspaceV1Step("Delete the workspace", suite.Client.WorkspaceV1, workspace)
		delSteps.GetWorkspaceWithErrorV1Step("Get the deleted workspace", suite.Client.WorkspaceV1, tref, secapi.ErrResourceNotFound)
	})

	suite.FinishScenario()
}

func (suite *WorkspaceLifeCycleV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
