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

type WorkspaceListV1TestSuite struct {
	suites.RegionalTestSuite

	params params.WorkspaceListV1Params
}

func CreateListV1TestSuite(regionalTestSuite suites.RegionalTestSuite) *WorkspaceListV1TestSuite {
	suite := &WorkspaceListV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
	}
	suite.ScenarioName = constants.WorkspaceV1ListSuiteName
	return suite
}

func (suite *WorkspaceListV1TestSuite) BeforeAll(t provider.T) {
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

	params := params.WorkspaceListV1Params{
		Workspaces: workspaces,
	}
	suite.params = params
	err = suites.SetupMockIfEnabled(suite.TestSuite, mockWorkspace.ConfigureListScenarioV1, &suite.params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *WorkspaceListV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, constants.WorkspaceProviderV1, string(schema.RegionalResourceMetadataKindResourceKindWorkspace))

	// Workspace list scenario grouped
	workspaces := suite.params.Workspaces
	tref := &secapi.TenantReference{Tenant: secapi.TenantID(suite.Tenant)}

	t.WithNewStep("Workspace", func(wsCtx provider.StepCtx) {
		wsSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, wsCtx)

		// Create workspaces
		for _, workspace := range workspaces {
			expectMeta := workspace.Metadata
			expectLabels := schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}
			wsSteps.CreateOrUpdateWorkspaceV1Step("Create", suite.Client.WorkspaceV1, &workspace,
				steps.ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec]{
					Labels:        expectLabels,
					Metadata:      expectMeta,
					ResourceState: schema.ResourceStateCreating,
				},
			)
		}

		// List workspaces with different filters
		wsSteps.GetListWorkspaceV1Step("ListAll", suite.Client.WorkspaceV1, *tref, nil)
		wsSteps.GetListWorkspaceV1Step("ListWithLimit", suite.Client.WorkspaceV1, *tref, secapi.NewListOptions().WithLimit(1))
		wsSteps.GetListWorkspaceV1Step("ListWithLabel", suite.Client.WorkspaceV1, *tref,
			secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().Equals(constants.EnvLabel, constants.EnvConformanceLabel)))
		wsSteps.GetListWorkspaceV1Step("ListWithLabelAndLimit", suite.Client.WorkspaceV1, *tref,
			secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().Equals(constants.EnvLabel, constants.EnvConformanceLabel)))
	})

	t.WithNewStep("Delete", func(delCtx provider.StepCtx) {
		delSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, delCtx)

		for _, workspace := range workspaces {
			ws := workspace
			delSteps.DeleteWorkspaceV1Step("Delete", suite.Client.WorkspaceV1, &ws)

			workspaceTRef := secapi.TenantReference{
				Tenant: secapi.TenantID(ws.Metadata.Tenant),
				Name:   ws.Metadata.Name,
			}
			delSteps.GetWorkspaceWithErrorV1Step("GetDeleted", suite.Client.WorkspaceV1, workspaceTRef, secapi.ErrResourceNotFound)
		}
	})

	suite.FinishScenario()
}

func (suite *WorkspaceListV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
