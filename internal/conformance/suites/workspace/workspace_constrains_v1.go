package workspace

import (
	"strings"

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

// WorkspaceConstraintsValidationV1TestSuite verifies that Workspace resources violating
// field constraints are rejected by the API with 422 Unprocessable Entity.
//
// Constraints tested:
//   - name: maxLength 128 (NameMetadata)
//   - name: pattern ^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$ (NameMetadata)
//   - labels values: maxLength 63 (UserResourceMetadata)
//   - annotations values: maxLength 1024 (UserResourceMetadata)
type WorkspaceConstraintsValidationV1TestSuite struct {
	suites.RegionalTestSuite

	params *params.WorkspaceConstraintsValidationV1Params
}

func CreateWorkspaceConstraintsValidationV1TestSuite(regionalTestSuite suites.RegionalTestSuite) *WorkspaceConstraintsValidationV1TestSuite {
	suite := &WorkspaceConstraintsValidationV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
	}
	suite.ScenarioName = constants.WorkspaceConstraintsValidationV1SuiteName.String()
	return suite
}

func (suite *WorkspaceConstraintsValidationV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("Constraints")

	buildWorkspace := func(name string, labels schema.Labels, annotations schema.Annotations) *schema.Workspace {
		ws, err := builders.NewWorkspaceBuilder().
			Name(name).
			Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(suite.Tenant).Region(suite.Region).
			Labels(labels).Annotations(annotations).
			Build()
		if err != nil {
			t.Fatalf("Failed to build Workspace: %v", err)
		}
		return ws
	}

	p := &params.WorkspaceConstraintsValidationV1Params{
		OverLengthNameWorkspace: buildWorkspace(
			strings.Repeat("a", 129),
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel},
			schema.Annotations{"description": "Workspace with over-length name"},
		),
		InvalidPatternNameWorkspace: buildWorkspace(
			"Invalid-Name-With-Uppercase",
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel},
			schema.Annotations{"description": "Workspace with non-kebab-case name"},
		),
		OverLengthLabelValueWorkspace: buildWorkspace(
			generators.GenerateWorkspaceName(),
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel, "constraint-test": strings.Repeat("x", 64)},
			schema.Annotations{"description": "Workspace with over-length label value"},
		),
		OverLengthAnnotationWorkspace: buildWorkspace(
			generators.GenerateWorkspaceName(),
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel},
			schema.Annotations{"description": "Workspace with over-length annotation value", "long-annotation": strings.Repeat("y", 1025)},
		),
	}
	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mockworkspace.ConfigureWorkspaceConstraintsViolationsV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *WorkspaceConstraintsValidationV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t, sdkconsts.WorkspaceProviderV1Name)
	suite.ConfigureResources(t, string(schema.RegionalResourceMetadataKindResourceKindWorkspace))

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	stepsBuilder.CreateOrUpdateWorkspaceExpectViolationV1Step(
		"Create a workspace with name exceeding maxLength:128 — expect rejection",
		suite.Client.WorkspaceV1,
		suite.params.OverLengthNameWorkspace,
	)
	stepsBuilder.CreateOrUpdateWorkspaceExpectViolationV1Step(
		"Create a workspace with invalid name pattern (not kebab-case) — expect rejection",
		suite.Client.WorkspaceV1,
		suite.params.InvalidPatternNameWorkspace,
	)
	stepsBuilder.CreateOrUpdateWorkspaceExpectViolationV1Step(
		"Create a workspace with label value exceeding maxLength:63 — expect rejection",
		suite.Client.WorkspaceV1,
		suite.params.OverLengthLabelValueWorkspace,
	)
	stepsBuilder.CreateOrUpdateWorkspaceExpectViolationV1Step(
		"Create a workspace with annotation value exceeding maxLength:1024 — expect rejection",
		suite.Client.WorkspaceV1,
		suite.params.OverLengthAnnotationWorkspace,
	)

	suite.FinishScenario()
}

func (suite *WorkspaceConstraintsValidationV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
