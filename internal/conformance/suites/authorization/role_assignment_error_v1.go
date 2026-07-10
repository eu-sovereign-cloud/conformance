package authorization

import (
	"net/http"

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

// RoleAssignmentErrorV1TestSuite verifies that RoleAssignment resources with
// invalid references are rejected by the API with 422 Unprocessable Entity.
//
// Scenarios tested:
//   - Create role assignment with non-existent role ref
//   - Create role assignment with non-existent tenant in scope
//   - Create role assignment with non-existent region in scope
//   - Create role assignment with non-existent workspace in scope
//   - Create role assignment with non-existent sub
type RoleAssignmentErrorV1TestSuite struct {
	suites.GlobalTestSuite

	Users  []string
	params *params.RoleAssignmentErrorV1Params
}

func CreateRoleAssignmentErrorV1TestSuite(globalTestSuite suites.GlobalTestSuite, users []string) *RoleAssignmentErrorV1TestSuite {
	suite := &RoleAssignmentErrorV1TestSuite{
		GlobalTestSuite: globalTestSuite,
		Users:           users,
	}
	suite.ScenarioName = constants.RoleAssignmentErrorV1SuiteName.String()
	return suite
}

func (suite *RoleAssignmentErrorV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite(suites.AuthorizationParentSuite)

	baseLabels := schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}

	// Use a valid existing role name and image resource for permission
	validRoleName := generators.GenerateRoleName()
	imageName := generators.GenerateImageName()
	imageResource := generators.GenerateImageResource(imageName)

	// A valid sub from the configured users
	validSub := suite.Users[0]

	buildRoleAssignment := func(name string, roles []string, subs []string, scopes []schema.RoleAssignmentScope) *schema.RoleAssignment {
		ra, err := builders.NewRoleAssignmentBuilder().
			Name(name).
			Provider(sdkconsts.AuthorizationProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(suite.Tenant).
			Labels(baseLabels).
			Annotations(schema.Annotations{"description": "RoleAssignment for error scenario testing"}).
			Spec(&schema.RoleAssignmentSpec{
				Roles:  roles,
				Subs:   subs,
				Scopes: scopes,
			}).Build()
		if err != nil {
			t.Fatalf("Failed to build RoleAssignment: %v", err)
		}
		return ra
	}

	// Pre-create a valid role to use in some scenarios
	validRole, err := builders.NewRoleBuilder().
		Name(validRoleName).
		Provider(sdkconsts.AuthorizationProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).
		Labels(baseLabels).
		Annotations(schema.Annotations{"description": "Role for role assignment error scenario testing"}).
		Spec(&schema.RoleSpec{
			Permissions: []schema.Permission{
				{
					Provider:  sdkconsts.StorageProviderV1Name,
					Resources: []string{imageResource},
					Verb:      []string{http.MethodGet},
				},
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build Role: %v", err)
	}

	p := &params.RoleAssignmentErrorV1Params{
		Role: validRole,

		// Non-existent role ref — role name does not exist in the tenant
		NonExistentRoleRefRoleAssignment: buildRoleAssignment(
			generators.GenerateRoleAssignmentName(),
			[]string{"non-existent-role"},
			[]string{validSub},
			[]schema.RoleAssignmentScope{
				{Tenants: []string{suite.Tenant}},
			},
		),

		// Non-existent tenant in scope — tenant does not exist
		NonExistentScopeTenantRoleAssignment: buildRoleAssignment(
			generators.GenerateRoleAssignmentName(),
			[]string{validRoleName},
			[]string{validSub},
			[]schema.RoleAssignmentScope{
				{Tenants: []string{"non-existent-tenant"}},
			},
		),

		// Non-existent region in scope — region does not exist
		NonExistentScopeRegionRoleAssignment: buildRoleAssignment(
			generators.GenerateRoleAssignmentName(),
			[]string{validRoleName},
			[]string{validSub},
			[]schema.RoleAssignmentScope{
				{
					Tenants: []string{suite.Tenant},
					Regions: []string{"non-existent-region"},
				},
			},
		),

		// Non-existent workspace in scope — workspace does not exist
		NonExistentScopeWorkspaceRoleAssignment: buildRoleAssignment(
			generators.GenerateRoleAssignmentName(),
			[]string{validRoleName},
			[]string{validSub},
			[]schema.RoleAssignmentScope{
				{
					Tenants:    []string{suite.Tenant},
					Workspaces: []string{"non-existent-workspace"},
				},
			},
		),

		// Non-existent sub — subject does not exist in identity provider
		NonExistentSubRoleAssignment: buildRoleAssignment(
			generators.GenerateRoleAssignmentName(),
			[]string{validRoleName},
			[]string{"non-existent-sub@example.com"},
			[]schema.RoleAssignmentScope{
				{Tenants: []string{suite.Tenant}},
			},
		),
	}

	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mockauthorization.ConfigureRoleAssignmentErrorV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *RoleAssignmentErrorV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t, sdkconsts.AuthorizationProviderV1Name)
	suite.ConfigureResources(t, string(schema.GlobalTenantResourceMetadataKindResourceKindRoleAssignment))
	suite.ConfigureDepends(t, string(schema.GlobalTenantResourceMetadataKindResourceKindRole))

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Error scenarios — all must be rejected with 422
	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with non-existent role ref — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.NonExistentRoleRefRoleAssignment,
	)

	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with non-existent tenant in scope — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.NonExistentScopeTenantRoleAssignment,
	)

	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with non-existent region in scope — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.NonExistentScopeRegionRoleAssignment,
	)

	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with non-existent workspace in scope — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.NonExistentScopeWorkspaceRoleAssignment,
	)

	stepsBuilder.CreateOrUpdateRoleAssignmentExpectViolationV1Step(
		"Create a role assignment with non-existent sub — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.NonExistentSubRoleAssignment,
	)

	suite.FinishScenario()
}

func (suite *RoleAssignmentErrorV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
