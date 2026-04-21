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

	// Workspace with name exceeding maxLength: 128 (129 chars)
	overLengthName := strings.Repeat("a", 129)
	overLengthNameWorkspace, err := builders.NewWorkspaceBuilder().
		Name(overLengthName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{"description": "Workspace with over-length name"}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build overLengthNameWorkspace: %v", err)
	}

	// Workspace with name violating kebab-case pattern (uppercase letters)
	invalidPatternName := "Invalid-Name-With-Uppercase"
	invalidPatternNameWorkspace, err := builders.NewWorkspaceBuilder().
		Name(invalidPatternName).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{"description": "Workspace with non-kebab-case name"}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build invalidPatternNameWorkspace: %v", err)
	}

	// Workspace with label value exceeding maxLength: 63 (64 chars)
	overLengthLabelWorkspace, err := builders.NewWorkspaceBuilder().
		Name(generators.GenerateWorkspaceName()).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
			"constraint-test":  strings.Repeat("x", 64),
		}).
		Annotations(schema.Annotations{"description": "Workspace with over-length label value"}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build overLengthLabelWorkspace: %v", err)
	}

	// Workspace with annotation value exceeding maxLength: 1024 (1025 chars)
	overLengthAnnotationWorkspace, err := builders.NewWorkspaceBuilder().
		Name(generators.GenerateWorkspaceName()).
		Provider(sdkconsts.WorkspaceProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).Region(suite.Region).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{
			"description":     "Workspace with over-length annotation value",
			"long-annotation": strings.Repeat("y", 1025),
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build overLengthAnnotationWorkspace: %v", err)
	}

	p := &params.WorkspaceConstraintsValidationV1Params{
		OverLengthNameWorkspace:       overLengthNameWorkspace,
		InvalidPatternNameWorkspace:   invalidPatternNameWorkspace,
		OverLengthLabelValueWorkspace: overLengthLabelWorkspace,
		OverLengthAnnotationWorkspace: overLengthAnnotationWorkspace,
	}
	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mockworkspace.ConfigureWorkspaceConstraintsViolationsV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *WorkspaceConstraintsValidationV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, sdkconsts.WorkspaceProviderV1Name, string(schema.RegionalResourceMetadataKindResourceKindWorkspace))

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// name: maxLength 128 — must be rejected
	stepsBuilder.CreateOrUpdateWorkspaceExpectViolationV1Step(
		"Create a workspace with name exceeding maxLength:128 — expect rejection",
		suite.Client.WorkspaceV1,
		suite.params.OverLengthNameWorkspace,
	)

	// name: pattern — must be rejected
	stepsBuilder.CreateOrUpdateWorkspaceExpectViolationV1Step(
		"Create a workspace with invalid name pattern (not kebab-case) — expect rejection",
		suite.Client.WorkspaceV1,
		suite.params.InvalidPatternNameWorkspace,
	)

	// labels value: maxLength 63 — must be rejected
	stepsBuilder.CreateOrUpdateWorkspaceExpectViolationV1Step(
		"Create a workspace with label value exceeding maxLength:63 — expect rejection",
		suite.Client.WorkspaceV1,
		suite.params.OverLengthLabelValueWorkspace,
	)

	// annotations value: maxLength 1024 — must be rejected
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
