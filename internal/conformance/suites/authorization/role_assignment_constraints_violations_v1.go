package authorization

import (
	"math/rand"
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

// RoleAssignmentConstraintsViolationsV1TestSuite verifies that RoleAssignment resources
// violating field constraints are rejected by the API with 422 Unprocessable Entity.
//
// Constraints tested:
//   - name: maxLength 128 (NameMetadata)
//   - name: pattern ^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$ (NameMetadata)
//   - labels values: maxLength 64 (UserResourceMetadata)
//   - annotations values: maxLength 1024 (UserResourceMetadata)
type RoleAssignmentConstraintsViolationsV1TestSuite struct {
	suites.GlobalTestSuite
	Users  []string
	params *params.RoleAssignmentConstraintsViolationsV1Params
}

func CreateRoleAssignmentConstraintsViolationsV1TestSuite(globalTestSuite suites.GlobalTestSuite, users []string) *RoleAssignmentConstraintsViolationsV1TestSuite {
	suite := &RoleAssignmentConstraintsViolationsV1TestSuite{
		GlobalTestSuite: globalTestSuite,
		Users:           users,
	}
	suite.ScenarioName = constants.RoleAssignmentConstraintsV1SuiteName.String()
	return suite
}

func (suite *RoleAssignmentConstraintsViolationsV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("Constraints")

	roleName := generators.GenerateRoleName()
	roleAssignmentSub := suite.Users[rand.Intn(len(suite.Users))]

	// RoleAssignment with name exceeding maxLength: 128 (129 chars)
	overLengthName := strings.Repeat("a", 129)
	overLengthNameRoleAssignment, err := builders.NewRoleAssignmentBuilder().
		Name(overLengthName).
		Provider(sdkconsts.AuthorizationProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{"description": "RoleAssignment with over-length name"}).
		Spec(&schema.RoleAssignmentSpec{
			Roles: []string{roleName},
			Subs:  []string{roleAssignmentSub},
			Scopes: []schema.RoleAssignmentScope{
				{Tenants: []string{suite.Tenant}},
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build overLengthNameRoleAssignment: %v", err)
	}

	// RoleAssignment with name violating kebab-case pattern (uppercase letters)
	invalidPatternName := "Invalid-Name-With-Uppercase"
	invalidPatternNameRoleAssignment, err := builders.NewRoleAssignmentBuilder().
		Name(invalidPatternName).
		Provider(sdkconsts.AuthorizationProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{"description": "RoleAssignment with non-kebab-case name"}).
		Spec(&schema.RoleAssignmentSpec{
			Roles: []string{roleName},
			Subs:  []string{roleAssignmentSub},
			Scopes: []schema.RoleAssignmentScope{
				{Tenants: []string{suite.Tenant}},
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build invalidPatternNameRoleAssignment: %v", err)
	}

	// RoleAssignment with label value exceeding maxLength: 63 (64 chars)
	overLengthLabelRoleAssignment, err := builders.NewRoleAssignmentBuilder().
		Name(generators.GenerateRoleAssignmentName()).
		Provider(sdkconsts.AuthorizationProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).
		Labels(schema.Labels{
			constants.EnvLabel: constants.EnvConformanceLabel,
			"constraint-test":  strings.Repeat("x", 64),
		}).
		Annotations(schema.Annotations{"description": "RoleAssignment with over-length label value"}).
		Spec(&schema.RoleAssignmentSpec{
			Roles: []string{roleName},
			Subs:  []string{roleAssignmentSub},
			Scopes: []schema.RoleAssignmentScope{
				{Tenants: []string{suite.Tenant}},
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build overLengthLabelRoleAssignment: %v", err)
	}

	// RoleAssignment with annotation value exceeding maxLength: 1024 (1025 chars)
	overLengthAnnotationRoleAssignment, err := builders.NewRoleAssignmentBuilder().
		Name(generators.GenerateRoleAssignmentName()).
		Provider(sdkconsts.AuthorizationProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Annotations(schema.Annotations{
			"description":     "RoleAssignment with over-length annotation value",
			"long-annotation": strings.Repeat("y", 1025),
		}).
		Spec(&schema.RoleAssignmentSpec{
			Roles: []string{roleName},
			Subs:  []string{roleAssignmentSub},
			Scopes: []schema.RoleAssignmentScope{
				{Tenants: []string{suite.Tenant}},
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build overLengthAnnotationRoleAssignment: %v", err)
	}

	p := &params.RoleAssignmentConstraintsViolationsV1Params{
		OverLengthNameRoleAssignment:       overLengthNameRoleAssignment,
		InvalidPatternNameRoleAssignment:   invalidPatternNameRoleAssignment,
		OverLengthLabelValueRoleAssignment: overLengthLabelRoleAssignment,
		OverLengthAnnotationRoleAssignment: overLengthAnnotationRoleAssignment,
	}
	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mockauthorization.ConfigureRoleAssignmentConstraintsViolationsV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *RoleAssignmentConstraintsViolationsV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, sdkconsts.AuthorizationProviderV1Name, string(schema.GlobalTenantResourceMetadataKindResourceKindRoleAssignment))

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// name: maxLength 128 — must be rejected
	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with name exceeding maxLength:128 — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.OverLengthNameRoleAssignment,
	)

	// name: pattern — must be rejected
	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with invalid name pattern (not kebab-case) — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.InvalidPatternNameRoleAssignment,
	)

	// labels value: maxLength 64 — must be rejected
	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with label value exceeding maxLength:64 — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.OverLengthLabelValueRoleAssignment,
	)

	// annotations value: maxLength 1024 — must be rejected
	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with annotation value exceeding maxLength:1024 — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.OverLengthAnnotationRoleAssignment,
	)

	suite.FinishScenario()
}

func (suite *RoleAssignmentConstraintsViolationsV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
