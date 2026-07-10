package workspace

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockworkspace "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/workspace"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// WorkspaceErrorV1TestSuite verifies that Workspace resources with invalid references
// are rejected by the API with 422 Unprocessable Entity.
//
// Scenarios tested:
//   - Create workspace with non-existent region
type WorkspaceErrorV1TestSuite struct {
	suites.RegionalTestSuite

	params *params.WorkspaceErrorV1Params
}

func CreateWorkspaceErrorV1TestSuite(regionalTestSuite suites.RegionalTestSuite) *WorkspaceErrorV1TestSuite {
	suite := &WorkspaceErrorV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
	}
	suite.ScenarioName = constants.WorkspaceErrorV1SuiteName.String()
	return suite
}

func (suite *WorkspaceErrorV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite(suites.WorkspaceParentSuite)

	baseLabels := schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}

	// Non-existent region — region string does not exist in the platform
	nonExistentRegionWorkspace, err := builders.NewWorkspaceBuilder().
		Name(generators.GenerateWorkspaceName()).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region("non-existent-region").
		Labels(baseLabels).
		Annotations(schema.Annotations{"description": "Workspace with non-existent region"}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Workspace: %v", err)
	}

	p := &params.WorkspaceErrorV1Params{
		NonExistentRegionWorkspace: nonExistentRegionWorkspace,
	}

	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mockworkspace.ConfigureWorkspaceErrorV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *WorkspaceErrorV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t, sdkconsts.WorkspaceProviderV1Name)
	suite.ConfigureResources(t, string(schema.RegionalResourceMetadataKindResourceKindWorkspace))

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Error scenarios — all must be rejected with 422
	stepsBuilder.CreateOrUpdateWorkspaceExpectViolationV1Step(
		"Create a workspace with non-existent region — expect rejection",
		suite.Client.WorkspaceV1,
		suite.params.NonExistentRegionWorkspace,
	)

	suite.FinishScenario()
}

func (suite *WorkspaceErrorV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
