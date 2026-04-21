package authorization

import (
	"net/http"
	"strings"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockauthorization "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/authorization"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// RoleConstraintsValidationV1TestSuite verifies that Role resources violating field constraints
// are rejected by the API with an error (422 Unprocessable Entity).
//
// Constraints tested:
//   - name: maxLength 128 (NameMetadata)
//   - name: pattern ^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$ (NameMetadata)
//   - labels values: maxLength 63 (UserResourceMetadata)
//   - annotations values: maxLength 1024 (UserResourceMetadata)
type RoleConstraintsValidationV1TestSuite struct {
	suites.GlobalTestSuite

	params *params.RoleConstraintsValidationV1Params
}

func CreateRoleConstraintsValidationV1TestSuite(globalTestSuite suites.GlobalTestSuite) *RoleConstraintsValidationV1TestSuite {
	suite := &RoleConstraintsValidationV1TestSuite{
		GlobalTestSuite: globalTestSuite,
	}
	suite.ScenarioName = constants.RoleConstraintsValidationV1SuiteName.String()
	return suite
}

func (suite *RoleConstraintsValidationV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("Constraints")

	imageName := generators.GenerateImageName()
	imageResource := generators.GenerateImageResource(imageName)
	basePermissions := []schema.Permission{
		{Provider: sdkconsts.StorageProviderV1Name, Resources: []string{imageResource}, Verb: []string{http.MethodGet}},
	}

	// Role with name exceeding maxLength: 128 (129 chars)
	overLengthName := strings.Repeat("a", 129)
	overLengthNameRole, err := builders.NewRoleBuilder().
		Name(overLengthName).
		Provider(sdkconsts.AuthorizationProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{"description": "Role with over-length name"}).
		Spec(&schema.RoleSpec{Permissions: basePermissions}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build overLengthNameRole: %v", err)
	}

	// Role with name violating kebab-case pattern (uppercase letters)
	invalidPatternName := "Invalid-Name-With-Uppercase"
	invalidPatternNameRole, err := builders.NewRoleBuilder().
		Name(invalidPatternName).
		Provider(sdkconsts.AuthorizationProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{"description": "Role with non-kebab-case name"}).
		Spec(&schema.RoleSpec{Permissions: basePermissions}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build invalidPatternNameRole: %v", err)
	}

	// Role with label value exceeding maxLength: 63 (64 chars)
	overLengthLabelRole, err := builders.NewRoleBuilder().
		Name(generators.GenerateRoleName()).
		Provider(sdkconsts.AuthorizationProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
			"constraint-test":  strings.Repeat("x", 64),
		}).
		Annotations(schema.Annotations{"description": "Role with over-length label value"}).
		Spec(&schema.RoleSpec{Permissions: basePermissions}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build overLengthLabelRole: %v", err)
	}

	// Role with annotation value exceeding maxLength: 1024 (1025 chars)
	overLengthAnnotationRole, err := builders.NewRoleBuilder().
		Name(generators.GenerateRoleName()).
		Provider(sdkconsts.AuthorizationProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{
			"description":     "Role with over-length annotation value",
			"long-annotation": strings.Repeat("y", 1025),
		}).
		Spec(&schema.RoleSpec{Permissions: basePermissions}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build overLengthAnnotationRole: %v", err)
	}

	p := &params.RoleConstraintsValidationV1Params{
		OverLengthNameRole:       overLengthNameRole,
		InvalidPatternNameRole:   invalidPatternNameRole,
		OverLengthLabelValueRole: overLengthLabelRole,
		OverLengthAnnotationRole: overLengthAnnotationRole,
	}
	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mockauthorization.ConfigureRoleConstraintsValidationV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *RoleConstraintsValidationV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, sdkconsts.AuthorizationProviderV1Name, string(schema.GlobalTenantResourceMetadataKindResourceKindRole))

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// name: maxLength 128 — must be rejected
	stepsBuilder.CreateOrUpdateRoleExpectViolationV1Step(
		"Create a role with name exceeding maxLength:128 — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.OverLengthNameRole,
	)

	// name: pattern — must be rejected
	stepsBuilder.CreateOrUpdateRoleExpectViolationV1Step(
		"Create a role with invalid name pattern (not kebab-case) — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.InvalidPatternNameRole,
	)

	// labels value: maxLength 63 — must be rejected
	stepsBuilder.CreateOrUpdateRoleExpectViolationV1Step(
		"Create a role with label value exceeding maxLength:63 — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.OverLengthLabelValueRole,
	)

	// annotations value: maxLength 1024 — must be rejected
	stepsBuilder.CreateOrUpdateRoleExpectViolationV1Step(
		"Create a role with annotation value exceeding maxLength:1024 — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.OverLengthAnnotationRole,
	)

	suite.FinishScenario()
}

func (suite *RoleConstraintsValidationV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
