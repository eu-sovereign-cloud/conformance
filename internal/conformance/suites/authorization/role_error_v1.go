package authorization

import (
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

// RoleErrorV1TestSuite verifies that Role resources with invalid references or
// invalid spec values are rejected by the API with 422 Unprocessable Entity.
//
// Scenarios tested:
//   - Create role with invalid provider in permission (non-existent provider)
//   - Create role with empty permissions list
//   - Create role with invalid HTTP verb in permission
type RoleErrorV1TestSuite struct {
	suites.GlobalTestSuite

	params *params.RoleErrorV1Params
}

func CreateRoleErrorV1TestSuite(globalTestSuite suites.GlobalTestSuite) *RoleErrorV1TestSuite {
	suite := &RoleErrorV1TestSuite{
		GlobalTestSuite: globalTestSuite,
	}
	suite.ScenarioName = constants.RoleErrorV1SuiteName.String()
	return suite
}

func (suite *RoleErrorV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite(suites.AuthorizationParentSuite)

	baseLabels := schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}

	imageName := generators.GenerateImageName()
	imageResource := generators.GenerateImageResource(imageName)

	buildRole := func(name string, permissions []schema.Permission) *schema.Role {
		role, err := builders.NewRoleBuilder().
			Name(name).
			Provider(sdkconsts.AuthorizationProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(suite.Tenant).
			Labels(baseLabels).
			Annotations(schema.Annotations{"description": "Role for error scenario testing"}).
			Spec(&schema.RoleSpec{Permissions: permissions}).
			Build()
		if err != nil {
			t.Fatalf("Failed to build Role: %v", err)
		}
		return role
	}

	p := &params.RoleErrorV1Params{
		// Non-existent provider in permission — provider name does not exist
		InvalidProviderPermissionRole: buildRole(
			generators.GenerateRoleName(),
			[]schema.Permission{
				{
					Provider:  "non-existent-provider",
					Resources: []string{imageResource},
					Verb:      []string{"GET"},
				},
			},
		),

		// Empty permissions list — must have at least 1 permission
		EmptyPermissionsRole: buildRole(
			generators.GenerateRoleName(),
			[]schema.Permission{},
		),

		// Invalid HTTP verb — "CONNECT" is not a valid SECA verb
		InvalidVerbRole: buildRole(
			generators.GenerateRoleName(),
			[]schema.Permission{
				{
					Provider:  sdkconsts.StorageProviderV1Name,
					Resources: []string{imageResource},
					Verb:      []string{"CONNECT"},
				},
			},
		),
	}

	suite.params = p
	if err := suites.SetupMockIfEnabled(suite.TestSuite, mockauthorization.ConfigureRoleErrorV1, *p); err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *RoleErrorV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t, sdkconsts.AuthorizationProviderV1Name)
	suite.ConfigureResources(t, string(schema.GlobalTenantResourceMetadataKindResourceKindRole))

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Error scenarios — all must be rejected with 422
	stepsBuilder.CreateOrUpdateRoleExpectViolationV1Step(
		"Create a role with non-existent provider in permission — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.InvalidProviderPermissionRole,
	)

	stepsBuilder.CreateOrUpdateRoleExpectViolationV1Step(
		"Create a role with empty permissions list — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.EmptyPermissionsRole,
	)

	stepsBuilder.CreateOrUpdateRoleExpectViolationV1Step(
		"Create a role with invalid HTTP verb (CONNECT not allowed) — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.InvalidVerbRole,
	)

	suite.FinishScenario()
}

func (suite *RoleErrorV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
