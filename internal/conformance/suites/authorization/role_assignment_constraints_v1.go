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

// RoleAssignmentConstraintsValidationV1TestSuite verifies that RoleAssignment resources
// violating field constraints are rejected by the API with 422 Unprocessable Entity.
//
// Constraints tested:
//   - name: maxLength 128 (NameMetadata)
//   - name: pattern ^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$ (NameMetadata)
//   - labels values: maxLength 63 (UserResourceMetadata)
//   - annotations values: maxLength 1024 (UserResourceMetadata)
type RoleAssignmentConstraintsValidationV1TestSuite struct {
	suites.GlobalTestSuite
	Users  []string
	params *params.RoleAssignmentConstraintsValidationV1Params
}

func CreateRoleAssignmentConstraintsValidationV1TestSuite(globalTestSuite suites.GlobalTestSuite, users []string) *RoleAssignmentConstraintsValidationV1TestSuite {
	suite := &RoleAssignmentConstraintsValidationV1TestSuite{
		GlobalTestSuite: globalTestSuite,
		Users:           users,
	}
	suite.ScenarioName = constants.RoleAssignmentConstraintsValidationV1SuiteName.String()
	return suite
}

func (suite *RoleAssignmentConstraintsValidationV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("Constraints")

	roleName := generators.GenerateRoleName()
	roleAssignmentSub := suite.Users[rand.Intn(len(suite.Users))]

	buildRoleAssignment := func(name string, labels schema.Labels, annotations schema.Annotations) *schema.RoleAssignment {
		ra, err := builders.NewRoleAssignmentBuilder().
			Name(name).
			Provider(sdkconsts.AuthorizationProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(suite.Tenant).
			Labels(labels).
			Annotations(annotations).
			Spec(&schema.RoleAssignmentSpec{
				Roles: []string{roleName},
				Subs:  []string{roleAssignmentSub},
				Scopes: []schema.RoleAssignmentScope{
					{Tenants: []string{suite.Tenant}},
				},
			}).Build()
		if err != nil {
			t.Fatalf("Failed to build RoleAssignment: %v", err)
		}
		return ra
	}

	p := &params.RoleAssignmentConstraintsValidationV1Params{
		OverLengthNameRoleAssignment: buildRoleAssignment(
			strings.Repeat("a", 129),
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel},
			schema.Annotations{"description": "RoleAssignment with over-length name"},
		),
		InvalidPatternNameRoleAssignment: buildRoleAssignment(
			"Invalid-Name-With-Uppercase",
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel},
			schema.Annotations{"description": "RoleAssignment with non-kebab-case name"},
		),
		OverLengthLabelValueRoleAssignment: buildRoleAssignment(
			generators.GenerateRoleAssignmentName(),
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel, "constraint-test": strings.Repeat("x", 64)},
			schema.Annotations{"description": "RoleAssignment with over-length label value"},
		),
		OverLengthAnnotationRoleAssignment: buildRoleAssignment(
			generators.GenerateRoleAssignmentName(),
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel},
			schema.Annotations{"description": "RoleAssignment with over-length annotation value", "long-annotation": strings.Repeat("y", 1025)},
		),
	}
	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mockauthorization.ConfigureRoleAssignmentConstraintsValidationV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *RoleAssignmentConstraintsValidationV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t, sdkconsts.AuthorizationProviderV1Name)
	suite.ConfigureResources(t, string(schema.GlobalTenantResourceMetadataKindResourceKindRoleAssignment))

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with name exceeding maxLength:128 — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.OverLengthNameRoleAssignment,
	)
	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with invalid name pattern (not kebab-case) — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.InvalidPatternNameRoleAssignment,
	)
	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with label value exceeding maxLength:63 — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.OverLengthLabelValueRoleAssignment,
	)
	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with annotation value exceeding maxLength:1024 — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.OverLengthAnnotationRoleAssignment,
	)

	suite.FinishScenario()
}

func (suite *RoleAssignmentConstraintsValidationV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
