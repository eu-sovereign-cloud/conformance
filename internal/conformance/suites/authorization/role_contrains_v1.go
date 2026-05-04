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

	buildRole := func(name string, labels schema.Labels, annotations schema.Annotations) *schema.Role {
		role, err := builders.NewRoleBuilder().
			Name(name).
			Provider(sdkconsts.AuthorizationProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(suite.Tenant).
			Labels(labels).
			Annotations(annotations).
			Spec(&schema.RoleSpec{Permissions: basePermissions}).
			Build()
		if err != nil {
			t.Fatalf("Failed to build Role: %v", err)
		}
		return role
	}

	p := &params.RoleConstraintsValidationV1Params{
		OverLengthNameRole: buildRole(
			strings.Repeat("a", 129),
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel},
			schema.Annotations{"description": "Role with over-length name"},
		),
		InvalidPatternNameRole: buildRole(
			"Invalid-Name-With-Uppercase",
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel},
			schema.Annotations{"description": "Role with non-kebab-case name"},
		),
		OverLengthLabelValueRole: buildRole(
			generators.GenerateRoleName(),
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel, "constraint-test": strings.Repeat("x", 64)},
			schema.Annotations{"description": "Role with over-length label value"},
		),
		OverLengthAnnotationRole: buildRole(
			generators.GenerateRoleName(),
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel},
			schema.Annotations{"description": "Role with over-length annotation value", "long-annotation": strings.Repeat("y", 1025)},
		),
	}
	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mockauthorization.ConfigureRoleConstraintsValidationV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *RoleConstraintsValidationV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t, sdkconsts.AuthorizationProviderV1Name)
	suite.ConfigureResources(t, string(schema.GlobalTenantResourceMetadataKindResourceKindRole))

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	stepsBuilder.CreateOrUpdateRoleExpectViolationV1Step(
		"Create a role with name exceeding maxLength:128 — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.OverLengthNameRole,
	)
	stepsBuilder.CreateOrUpdateRoleExpectViolationV1Step(
		"Create a role with invalid name pattern (not kebab-case) — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.InvalidPatternNameRole,
	)
	stepsBuilder.CreateOrUpdateRoleExpectViolationV1Step(
		"Create a role with label value exceeding maxLength:63 — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.OverLengthLabelValueRole,
	)
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
