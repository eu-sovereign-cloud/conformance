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
	t.AddParentSuite(suites.AuthorizationParentSuite)

	imageName := generators.GenerateImageName()
	imageResource := generators.GenerateImageResource(imageName)

	basePermission := schema.Permission{
		Provider:  sdkconsts.StorageProviderV1Name,
		Resources: []string{imageResource},
		Verb:      []string{http.MethodGet},
	}
	basePermissions := []schema.Permission{basePermission}

	repeatStrings := func(v string, n int) []string {
		out := make([]string, n)
		for i := range out {
			out[i] = v
		}
		return out
	}
	repeatPermissions := func(p schema.Permission, n int) []schema.Permission {
		out := make([]schema.Permission, n)
		for i := range out {
			out[i] = p
		}
		return out
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

	buildRoleWithPermission := func(name string, permissions []schema.Permission) *schema.Role {
		role, err := builders.NewRoleBuilder().
			Name(name).
			Provider(sdkconsts.AuthorizationProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
			Tenant(suite.Tenant).
			Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
			Spec(&schema.RoleSpec{Permissions: permissions}).
			Build()
		if err != nil {
			t.Fatalf("Failed to build Role: %v", err)
		}
		return role
	}

	buildValidRoleAndMutate := func(name string, mutate func(*schema.Role)) *schema.Role {
		role := buildRoleWithPermission(name, []schema.Permission{basePermission})
		mutate(role)
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
		OverLengthPermissionProviderRole: buildRoleWithPermission(
			generators.GenerateRoleName(),
			[]schema.Permission{{
				Provider:  strings.Repeat("a", 65),
				Resources: []string{"images/*"},
				Verb:      []string{"get"},
			}},
		),
		OverLengthPermissionResourceRole: buildRoleWithPermission(
			generators.GenerateRoleName(),
			[]schema.Permission{{
				Provider:  "seca.storage/v1",
				Resources: []string{strings.Repeat("a", 257)},
				Verb:      []string{"get"},
			}},
		),
		OverLengthPermissionVerbRole: buildRoleWithPermission(
			generators.GenerateRoleName(),
			[]schema.Permission{{
				Provider:  "seca.storage/v1",
				Resources: []string{"images/*"},
				Verb:      []string{strings.Repeat("a", 8)},
			}},
		),
		EmptyPermissionsRole: buildValidRoleAndMutate(
			generators.GenerateRoleName(),
			func(r *schema.Role) { r.Spec.Permissions = []schema.Permission{} },
		),
		EmptyPermissionProviderRole: buildValidRoleAndMutate(
			generators.GenerateRoleName(),
			func(r *schema.Role) { r.Spec.Permissions[0].Provider = "" },
		),
		EmptyPermissionResourcesRole: buildValidRoleAndMutate(
			generators.GenerateRoleName(),
			func(r *schema.Role) { r.Spec.Permissions[0].Resources = []string{} },
		),
		EmptyPermissionResourceValueRole: buildValidRoleAndMutate(
			generators.GenerateRoleName(),
			func(r *schema.Role) { r.Spec.Permissions[0].Resources = []string{""} },
		),
		EmptyPermissionVerbsRole: buildValidRoleAndMutate(
			generators.GenerateRoleName(),
			func(r *schema.Role) { r.Spec.Permissions[0].Verb = []string{} },
		),
		EmptyPermissionVerbValueRole: buildValidRoleAndMutate(
			generators.GenerateRoleName(),
			func(r *schema.Role) { r.Spec.Permissions[0].Verb = []string{""} },
		),
		OverMaxItemsPermissionsRole: buildRoleWithPermission(
			generators.GenerateRoleName(),
			repeatPermissions(basePermission, 257),
		),
		OverMaxItemsPermissionResourcesRole: buildRoleWithPermission(
			generators.GenerateRoleName(),
			[]schema.Permission{{
				Provider:  "seca.storage/v1",
				Resources: repeatStrings("images/*", 257),
				Verb:      []string{"get"},
			}},
		),
		OverMaxItemsPermissionVerbsRole: buildRoleWithPermission(
			generators.GenerateRoleName(),
			[]schema.Permission{{
				Provider:  "seca.storage/v1",
				Resources: []string{"images/*"},
				Verb:      repeatStrings("get", 17),
			}},
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
	stepsBuilder.CreateOrUpdateRoleExpectViolationV1Step(
		"Create a role with permission provider exceeding maxLength:64 — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.OverLengthPermissionProviderRole,
	)
	stepsBuilder.CreateOrUpdateRoleExpectViolationV1Step(
		"Create a role with permission resource exceeding maxLength:256 — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.OverLengthPermissionResourceRole,
	)
	stepsBuilder.CreateOrUpdateRoleExpectViolationV1Step(
		"Create a role with permission verb exceeding maxLength:7 — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.OverLengthPermissionVerbRole,
	)
	stepsBuilder.CreateOrUpdateRoleExpectViolationV1Step(
		"Create a role with permissions empty (minItems:1) — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.EmptyPermissionsRole,
	)
	stepsBuilder.CreateOrUpdateRoleExpectViolationV1Step(
		"Create a role with permissions exceeding maxItems:256 — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.OverMaxItemsPermissionsRole,
	)
	stepsBuilder.CreateOrUpdateRoleExpectViolationV1Step(
		"Create a role with empty permission provider (minLength:1) — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.EmptyPermissionProviderRole,
	)
	stepsBuilder.CreateOrUpdateRoleExpectViolationV1Step(
		"Create a role with permission resources empty (minItems:1) — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.EmptyPermissionResourcesRole,
	)
	stepsBuilder.CreateOrUpdateRoleExpectViolationV1Step(
		"Create a role with permission resources exceeding maxItems:256 — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.OverMaxItemsPermissionResourcesRole,
	)
	stepsBuilder.CreateOrUpdateRoleExpectViolationV1Step(
		"Create a role with empty permission resource value (minLength:1) — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.EmptyPermissionResourceValueRole,
	)
	stepsBuilder.CreateOrUpdateRoleExpectViolationV1Step(
		"Create a role with permission verbs empty (minItems:1) — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.EmptyPermissionVerbsRole,
	)
	stepsBuilder.CreateOrUpdateRoleExpectViolationV1Step(
		"Create a role with permission verbs exceeding maxItems:16 — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.OverMaxItemsPermissionVerbsRole,
	)
	stepsBuilder.CreateOrUpdateRoleExpectViolationV1Step(
		"Create a role with empty permission verb value (minLength:1) — expect rejection",
		suite.Client.AuthorizationV1,
		suite.params.EmptyPermissionVerbValueRole,
	)

	suite.FinishScenario()
}

func (suite *RoleConstraintsValidationV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
