package authorization

import (
	"math/rand"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockauthorization "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/authorization"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	labelBuilder "github.com/eu-sovereign-cloud/go-sdk/secapi/builders"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type AuthorizationListV1TestSuite struct {
	suites.GlobalTestSuite

	Users []string

	params *params.AuthorizationListV1Params
}

func NewListV1TestSuite(globalTestSuite suites.GlobalTestSuite, users []string) *AuthorizationListV1TestSuite {
	suite := &AuthorizationListV1TestSuite{
		GlobalTestSuite: globalTestSuite,
		Users:           users,
	}
	suite.ScenarioName = constants.AuthorizationListV1SuiteName.String()
	return suite
}

func (suite *AuthorizationListV1TestSuite) BeforeAll(t provider.T) {
	var err error

	// Select subs
	roleAssignmentSub1 := suite.Users[rand.Intn(len(suite.Users))]

	roleName1 := generators.GenerateRoleName()
	roleName2 := generators.GenerateRoleName()
	roleName3 := generators.GenerateRoleName()
	// Generate scenario data

	roleAssignmentName1 := generators.GenerateRoleAssignmentName()
	roleAssignmentName2 := generators.GenerateRoleAssignmentName()
	roleAssignmentName3 := generators.GenerateRoleAssignmentName()
	imageName := generators.GenerateImageName()
	imageResource := generators.GenerateImageResource(suite.Tenant, imageName)

	// Roles
	role1, err := builders.NewRoleBuilder().
		Name(roleName1).
		Provider(constants.AuthorizationProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Spec(&schema.RoleSpec{
			Permissions: []schema.Permission{
				{Provider: constants.StorageProviderV1, Resources: []string{imageResource}, Verb: []string{http.MethodGet}},
			},
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Role: %v", err)
	}

	role2, err := builders.NewRoleBuilder().
		Name(roleName2).
		Provider(constants.AuthorizationProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Spec(&schema.RoleSpec{
			Permissions: []schema.Permission{
				{Provider: constants.StorageProviderV1, Resources: []string{imageResource}, Verb: []string{http.MethodGet}},
			},
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Role: %v", err)
	}
	role3, err := builders.NewRoleBuilder().
		Name(roleName3).
		Provider(constants.AuthorizationProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Spec(&schema.RoleSpec{
			Permissions: []schema.Permission{
				{Provider: constants.StorageProviderV1, Resources: []string{imageResource}, Verb: []string{http.MethodGet}},
			},
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Role: %v", err)
	}
	roles := []schema.Role{
		*role1,
		*role2,
		*role3,
	}

	// Roles Assignment
	roleAssignment1, err := builders.NewRoleAssignmentBuilder().
		Name(roleAssignmentName1).
		Provider(constants.AuthorizationProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Spec(&schema.RoleAssignmentSpec{
			Roles: []string{roleName2},
			Subs:  []string{roleAssignmentSub1},
			Scopes: []schema.RoleAssignmentScope{
				{Tenants: &[]string{suite.Tenant}},
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build RoleAssignment: %v", err)
	}
	roleAssignment2, err := builders.NewRoleAssignmentBuilder().
		Name(roleAssignmentName2).
		Provider(constants.AuthorizationProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Spec(&schema.RoleAssignmentSpec{
			Roles: []string{roleName2},
			Subs:  []string{roleAssignmentSub1},
			Scopes: []schema.RoleAssignmentScope{
				{Tenants: &[]string{suite.Tenant}},
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build RoleAssignment: %v", err)
	}
	roleAssignment3, err := builders.NewRoleAssignmentBuilder().
		Name(roleAssignmentName3).
		Provider(constants.AuthorizationProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Spec(&schema.RoleAssignmentSpec{
			Roles: []string{roleName3},
			Subs:  []string{roleAssignmentSub1},
			Scopes: []schema.RoleAssignmentScope{
				{Tenants: &[]string{suite.Tenant}},
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build RoleAssignment: %v", err)
	}
	roleAssignments := []schema.RoleAssignment{
		*roleAssignment1,
		*roleAssignment2,
		*roleAssignment3,
	}
	params := &params.AuthorizationListV1Params{
		Roles:           roles,
		RoleAssignments: roleAssignments,
	}
	suite.params = params
	err = suites.SetupMockIfEnabled(suite.TestSuite, mockauthorization.ConfigureListScenarioV1, params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *AuthorizationListV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, constants.AuthorizationProviderV1,
		string(schema.GlobalTenantResourceMetadataKindResourceKindRole),
		string(schema.GlobalTenantResourceMetadataKindResourceKindRoleAssignment),
	)

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Role
	roles := suite.params.Roles
	// Create roles
	for _, role := range roles {
		expectRoleMeta := role.Metadata
		expectRoleSpec := role.Spec

		// Create Role
		stepsBuilder.CreateOrUpdateRoleV1Step("Create a role", suite.Client.AuthorizationV1, &role,
			steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
				Metadata:      expectRoleMeta,
				Spec:          &expectRoleSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)
	}

	roleTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.Tenant),
		Name:   suite.Tenant,
	}
	// List Roles
	stepsBuilder.GetListRoleV1Step("Get list of roles", suite.Client.AuthorizationV1, *roleTRef, nil)

	// List Roles with limit
	stepsBuilder.GetListRoleV1Step("Get list of roles with limit", suite.Client.AuthorizationV1, *roleTRef,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List Roles with Label
	stepsBuilder.GetListRoleV1Step("Get list of roles with label", suite.Client.AuthorizationV1, *roleTRef,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List Roles with Limit and label
	stepsBuilder.GetListRoleV1Step("Get list of roles with limit and label", suite.Client.AuthorizationV1, *roleTRef,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Role assignment
	roleAssignments := suite.params.RoleAssignments

	// Create role assignments
	for _, roleAssign := range roleAssignments {
		expectRoleAssignMeta := roleAssign.Metadata
		expectRoleAssignSpec := &roleAssign.Spec

		// Create a role assignment
		stepsBuilder.CreateOrUpdateRoleAssignmentV1Step("Create a role assignment", suite.Client.AuthorizationV1, &roleAssign,
			steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
				Metadata:      expectRoleAssignMeta,
				Spec:          expectRoleAssignSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)
	}
	roleAssignTRef := &secapi.TenantReference{Tenant: secapi.TenantID(suite.Tenant)}

	// List RoleAssignments
	stepsBuilder.GetListRoleAssignmentsV1("Get list of role assignments", suite.Client.AuthorizationV1, *roleAssignTRef, nil)

	// List RoleAssignments with limit
	stepsBuilder.GetListRoleAssignmentsV1("Get list of role assignments", suite.Client.AuthorizationV1, *roleAssignTRef,
		secapi.NewListOptions().WithLimit(1))

	// List RoleAssignments with Label
	stepsBuilder.GetListRoleAssignmentsV1("Get list of role assignments", suite.Client.AuthorizationV1, *roleAssignTRef,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List RoleAssignments with Limit and label
	stepsBuilder.GetListRoleAssignmentsV1("Get list of role assignments", suite.Client.AuthorizationV1, *roleAssignTRef,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Delete all role assignments
	for _, roleAssign := range roleAssignments {
		stepsBuilder.DeleteRoleAssignmentV1Step("Delete the role assignment", suite.Client.AuthorizationV1, &roleAssign)

		// Get the deleted role assignment
		roleAssignTRefSingle := &secapi.TenantReference{
			Tenant: secapi.TenantID(suite.Tenant),
			Name:   roleAssign.Metadata.Name,
		}
		stepsBuilder.GetRoleAssignmentWithErrorV1Step("Get the deleted role assignment", suite.Client.AuthorizationV1, *roleAssignTRefSingle, secapi.ErrResourceNotFound)
	}

	// Delete all roles
	for _, role := range roles {
		stepsBuilder.DeleteRoleV1Step("Delete the role", suite.Client.AuthorizationV1, &role)

		// Get the deleted role
		roleTRefSingle := &secapi.TenantReference{
			Tenant: secapi.TenantID(suite.Tenant),
			Name:   role.Metadata.Name,
		}
		stepsBuilder.GetRoleWithErrorV1Step("Get the deleted role", suite.Client.AuthorizationV1, *roleTRefSingle, secapi.ErrResourceNotFound)
	}
	suite.FinishScenario()
}

func (suite *AuthorizationListV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
