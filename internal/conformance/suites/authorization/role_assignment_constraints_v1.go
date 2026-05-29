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
	t.AddParentSuite(suites.AuthorizationParentSuite)

	roleName := generators.GenerateRoleName()
	roleAssignmentSub := suite.Users[rand.Intn(len(suite.Users))]

	validScope := schema.RoleAssignmentScope{
		Tenants: []string{suite.Tenant},
	}

	repeatStrings := func(value string, n int) []string {
		out := make([]string, n)
		for i := range out {
			out[i] = value
		}
		return out
	}

	repeatScopes := func(scope schema.RoleAssignmentScope, n int) []schema.RoleAssignmentScope {
		out := make([]schema.RoleAssignmentScope, n)
		for i := range out {
			out[i] = scope
		}
		return out
	}

	buildRoleAssignment := func(name string, labels schema.Labels, annotations schema.Annotations, spec *schema.RoleAssignmentSpec) *schema.RoleAssignment {
		ra, err := builders.NewRoleAssignmentBuilder().
			Name(name).
			Provider(sdkconsts.AuthorizationProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(suite.Tenant).
			Labels(labels).
			Annotations(annotations).
			Spec(spec).
			Build()
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
			&schema.RoleAssignmentSpec{
				Roles: []string{roleName},
				Subs:  []string{roleAssignmentSub},
				Scopes: []schema.RoleAssignmentScope{
					{Tenants: []string{suite.Tenant}},
				},
			},
		),
		InvalidPatternNameRoleAssignment: buildRoleAssignment(
			"Invalid-Name-With-Uppercase",
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel},
			schema.Annotations{"description": "RoleAssignment with over-length name"},
			&schema.RoleAssignmentSpec{
				Roles: []string{roleName},
				Subs:  []string{roleAssignmentSub},
				Scopes: []schema.RoleAssignmentScope{
					validScope,
				},
			},
		),
		OverLengthLabelValueRoleAssignment: buildRoleAssignment(
			generators.GenerateRoleAssignmentName(),
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel, "constraint-test": strings.Repeat("x", 64)},
			schema.Annotations{"description": "RoleAssignment with over-length label value"},
			&schema.RoleAssignmentSpec{
				Roles: []string{roleName},
				Subs:  []string{roleAssignmentSub},
				Scopes: []schema.RoleAssignmentScope{
					validScope,
				},
			},
		),
		OverLengthAnnotationRoleAssignment: buildRoleAssignment(
			generators.GenerateRoleAssignmentName(),
			schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel},
			schema.Annotations{"description": "RoleAssignment with over-length annotation value", "long-annotation": strings.Repeat("y", 1025)},
			&schema.RoleAssignmentSpec{
				Roles: []string{roleName},
				Subs:  []string{roleAssignmentSub},
				Scopes: []schema.RoleAssignmentScope{
					validScope,
				},
			},
		),
		OverLengthSubRoleAssignment: buildRoleAssignment(
			generators.GenerateRoleAssignmentName(),
			schema.Labels{},
			schema.Annotations{},
			&schema.RoleAssignmentSpec{
				Subs:   []string{strings.Repeat("a", 129)},
				Scopes: []schema.RoleAssignmentScope{validScope},
				Roles:  []string{"viewer"},
			},
		),
		OverLengthRoleNameRoleAssignment: buildRoleAssignment(
			generators.GenerateRoleAssignmentName(),
			schema.Labels{},
			schema.Annotations{},
			&schema.RoleAssignmentSpec{
				Subs:   []string{"user1@example.com"},
				Scopes: []schema.RoleAssignmentScope{validScope},
				Roles:  []string{strings.Repeat("a", 65)},
			},
		),
		OverLengthScopeTenantRoleAssignment: buildRoleAssignment(
			generators.GenerateRoleAssignmentName(),
			schema.Labels{},
			schema.Annotations{},
			&schema.RoleAssignmentSpec{
				Subs:   []string{"user1@example.com"},
				Scopes: []schema.RoleAssignmentScope{{Tenants: []string{strings.Repeat("a", 65)}}},
				Roles:  []string{"viewer"},
			},
		),
		OverLengthScopeRegionRoleAssignment: buildRoleAssignment(
			generators.GenerateRoleAssignmentName(),
			schema.Labels{},
			schema.Annotations{},
			&schema.RoleAssignmentSpec{
				Subs:   []string{"user1@example.com"},
				Scopes: []schema.RoleAssignmentScope{{Regions: []string{strings.Repeat("a", 65)}}},
				Roles:  []string{"viewer"},
			},
		),
		OverLengthScopeWorkspaceRoleAssignment: buildRoleAssignment(
			generators.GenerateRoleAssignmentName(),
			schema.Labels{},
			schema.Annotations{},
			&schema.RoleAssignmentSpec{
				Subs:   []string{"user1@example.com"},
				Scopes: []schema.RoleAssignmentScope{{Workspaces: []string{strings.Repeat("a", 65)}}},
				Roles:  []string{"viewer"},
			},
		),
		EmptyRolesRoleAssignment: buildRoleAssignment(
			generators.GenerateRoleAssignmentName(), schema.Labels{}, schema.Annotations{},
			&schema.RoleAssignmentSpec{
				Roles:  []string{},
				Subs:   []string{roleAssignmentSub},
				Scopes: []schema.RoleAssignmentScope{validScope},
			},
		),
		OverMaxItemsRolesRoleAssignment: buildRoleAssignment(
			generators.GenerateRoleAssignmentName(), schema.Labels{}, schema.Annotations{},
			&schema.RoleAssignmentSpec{
				Roles:  repeatStrings("viewer", 33),
				Subs:   []string{roleAssignmentSub},
				Scopes: []schema.RoleAssignmentScope{validScope},
			},
		),
		EmptyRoleValueRoleAssignment: buildRoleAssignment(
			generators.GenerateRoleAssignmentName(), schema.Labels{}, schema.Annotations{},
			&schema.RoleAssignmentSpec{
				Roles:  []string{""},
				Subs:   []string{roleAssignmentSub},
				Scopes: []schema.RoleAssignmentScope{validScope},
			},
		),
		EmptySubsRoleAssignment: buildRoleAssignment(
			generators.GenerateRoleAssignmentName(), schema.Labels{}, schema.Annotations{},
			&schema.RoleAssignmentSpec{
				Roles:  []string{roleName},
				Subs:   []string{},
				Scopes: []schema.RoleAssignmentScope{validScope},
			},
		),
		OverMaxItemsSubsRoleAssignment: buildRoleAssignment(
			generators.GenerateRoleAssignmentName(), schema.Labels{}, schema.Annotations{},
			&schema.RoleAssignmentSpec{
				Roles:  []string{roleName},
				Subs:   repeatStrings("user@example.com", 257),
				Scopes: []schema.RoleAssignmentScope{validScope},
			},
		),
		EmptySubValueRoleAssignment: buildRoleAssignment(
			generators.GenerateRoleAssignmentName(), schema.Labels{}, schema.Annotations{},
			&schema.RoleAssignmentSpec{
				Roles:  []string{roleName},
				Subs:   []string{""},
				Scopes: []schema.RoleAssignmentScope{validScope},
			},
		),
		EmptyScopesRoleAssignment: buildRoleAssignment(
			generators.GenerateRoleAssignmentName(), schema.Labels{}, schema.Annotations{},
			&schema.RoleAssignmentSpec{
				Roles:  []string{roleName},
				Subs:   []string{roleAssignmentSub},
				Scopes: []schema.RoleAssignmentScope{},
			},
		),
		OverMaxItemsScopesRoleAssignment: buildRoleAssignment(
			generators.GenerateRoleAssignmentName(), schema.Labels{}, schema.Annotations{},
			&schema.RoleAssignmentSpec{
				Roles:  []string{roleName},
				Subs:   []string{roleAssignmentSub},
				Scopes: repeatScopes(validScope, 257),
			},
		),
		EmptyScopeTenantValueRoleAssignment: buildRoleAssignment(
			generators.GenerateRoleAssignmentName(), schema.Labels{}, schema.Annotations{},
			&schema.RoleAssignmentSpec{
				Roles: []string{roleName},
				Subs:  []string{roleAssignmentSub},
				Scopes: []schema.RoleAssignmentScope{
					{Tenants: []string{""}},
				},
			},
		),
		OverMaxItemsScopeTenantsRoleAssignment: buildRoleAssignment(
			generators.GenerateRoleAssignmentName(), schema.Labels{}, schema.Annotations{},
			&schema.RoleAssignmentSpec{
				Roles: []string{roleName},
				Subs:  []string{roleAssignmentSub},
				Scopes: []schema.RoleAssignmentScope{
					{Tenants: repeatStrings("tenant-a", 65)},
				},
			},
		),
		EmptyScopeRegionValueRoleAssignment: buildRoleAssignment(
			generators.GenerateRoleAssignmentName(), schema.Labels{}, schema.Annotations{},
			&schema.RoleAssignmentSpec{
				Roles: []string{roleName},
				Subs:  []string{roleAssignmentSub},
				Scopes: []schema.RoleAssignmentScope{
					{Regions: []string{""}},
				},
			},
		),
		OverMaxItemsScopeRegionsRoleAssignment: buildRoleAssignment(
			generators.GenerateRoleAssignmentName(), schema.Labels{}, schema.Annotations{},
			&schema.RoleAssignmentSpec{
				Roles: []string{roleName},
				Subs:  []string{roleAssignmentSub},
				Scopes: []schema.RoleAssignmentScope{
					{Regions: repeatStrings("eu-1", 65)},
				},
			},
		),
		EmptyScopeWorkspaceValueRoleAssignment: buildRoleAssignment(
			generators.GenerateRoleAssignmentName(), schema.Labels{}, schema.Annotations{},
			&schema.RoleAssignmentSpec{
				Roles: []string{roleName},
				Subs:  []string{roleAssignmentSub},
				Scopes: []schema.RoleAssignmentScope{
					{Workspaces: []string{""}},
				},
			},
		),
		OverMaxItemsScopeWorkspacesRoleAssignment: buildRoleAssignment(
			generators.GenerateRoleAssignmentName(), schema.Labels{}, schema.Annotations{},
			&schema.RoleAssignmentSpec{
				Roles: []string{roleName},
				Subs:  []string{roleAssignmentSub},
				Scopes: []schema.RoleAssignmentScope{
					{Workspaces: repeatStrings("ws-a", 257)},
				},
			},
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
	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with sub exceeding maxLength:128 — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.OverLengthSubRoleAssignment,
	)
	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with role name exceeding maxLength:64 — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.OverLengthRoleNameRoleAssignment,
	)
	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with scope tenant exceeding maxLength:64 — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.OverLengthScopeTenantRoleAssignment,
	)
	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with scope region exceeding maxLength:64 — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.OverLengthScopeRegionRoleAssignment,
	)
	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with scope workspace exceeding maxLength:64 — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.OverLengthScopeWorkspaceRoleAssignment,
	)
	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with roles empty (minItems:1) — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.EmptyRolesRoleAssignment,
	)
	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with roles exceeding maxItems:32 — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.OverMaxItemsRolesRoleAssignment,
	)
	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with empty role value (minLength:1) — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.EmptyRoleValueRoleAssignment,
	)
	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with subs empty (minItems:1) — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.EmptySubsRoleAssignment,
	)
	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with subs exceeding maxItems:256 — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.OverMaxItemsSubsRoleAssignment,
	)
	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with empty sub value (minLength:1) — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.EmptySubValueRoleAssignment,
	)
	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with scopes empty (minItems:1) — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.EmptyScopesRoleAssignment,
	)
	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with scopes exceeding maxItems:256 — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.OverMaxItemsScopesRoleAssignment,
	)
	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with empty scope tenant value (minLength:1) — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.EmptyScopeTenantValueRoleAssignment,
	)
	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with scope tenants exceeding maxItems:64 — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.OverMaxItemsScopeTenantsRoleAssignment,
	)
	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with empty scope region value (minLength:1) — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.EmptyScopeRegionValueRoleAssignment,
	)
	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with scope regions exceeding maxItems:64 — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.OverMaxItemsScopeRegionsRoleAssignment,
	)
	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with empty scope workspace value (minLength:1) — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.EmptyScopeWorkspaceValueRoleAssignment,
	)
	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with scope workspaces exceeding maxItems:256 — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.OverMaxItemsScopeWorkspacesRoleAssignment,
	)
	suite.FinishScenario()
}

func (suite *RoleAssignmentConstraintsValidationV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
