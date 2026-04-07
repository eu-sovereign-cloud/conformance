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
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	labelBuilder "github.com/eu-sovereign-cloud/go-sdk/secapi/builders"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type ProviderQueriesV1TestSuite struct {
	suites.GlobalTestSuite

	Users []string

	params *params.AuthorizationProviderQueriesV1Params
}

func CreateProviderQueriesV1TestSuite(globalTestSuite suites.GlobalTestSuite, users []string) *ProviderQueriesV1TestSuite {
	suite := &ProviderQueriesV1TestSuite{
		GlobalTestSuite: globalTestSuite,
		Users:           users,
	}
	suite.ScenarioName = constants.AuthorizationProviderQueriesV1SuiteName.String()
	return suite
}

func (suite *ProviderQueriesV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite(suites.AuthorizationParentSuite)

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
	imageResource := generators.GenerateImageResource(imageName)

	// Roles
	role1, err := builders.NewRoleBuilder().
		Name(roleName1).
		Provider(sdkconsts.AuthorizationProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Spec(&schema.RoleSpec{
			Permissions: []schema.Permission{
				{Provider: sdkconsts.StorageProviderV1Name, Resources: []string{imageResource}, Verb: []string{http.MethodGet}},
			},
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Role: %v", err)
	}

	role2, err := builders.NewRoleBuilder().
		Name(roleName2).
		Provider(sdkconsts.AuthorizationProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Spec(&schema.RoleSpec{
			Permissions: []schema.Permission{
				{Provider: sdkconsts.StorageProviderV1Name, Resources: []string{imageResource}, Verb: []string{http.MethodGet}},
			},
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Role: %v", err)
	}
	role3, err := builders.NewRoleBuilder().
		Name(roleName3).
		Provider(sdkconsts.AuthorizationProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Spec(&schema.RoleSpec{
			Permissions: []schema.Permission{
				{Provider: sdkconsts.StorageProviderV1Name, Resources: []string{imageResource}, Verb: []string{http.MethodGet}},
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
		Provider(sdkconsts.AuthorizationProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Spec(&schema.RoleAssignmentSpec{
			Roles: []string{roleName2},
			Subs:  []string{roleAssignmentSub1},
			Scopes: []schema.RoleAssignmentScope{
				{Tenants: []string{suite.Tenant}},
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build RoleAssignment: %v", err)
	}
	roleAssignment2, err := builders.NewRoleAssignmentBuilder().
		Name(roleAssignmentName2).
		Provider(sdkconsts.AuthorizationProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Spec(&schema.RoleAssignmentSpec{
			Roles: []string{roleName2},
			Subs:  []string{roleAssignmentSub1},
			Scopes: []schema.RoleAssignmentScope{
				{Tenants: []string{suite.Tenant}},
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build RoleAssignment: %v", err)
	}
	roleAssignment3, err := builders.NewRoleAssignmentBuilder().
		Name(roleAssignmentName3).
		Provider(sdkconsts.AuthorizationProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Tenant(suite.Tenant).
		Labels(schema.Labels{constants.EnvLabel: constants.EnvConformanceLabel}).
		Spec(&schema.RoleAssignmentSpec{
			Roles: []string{roleName3},
			Subs:  []string{roleAssignmentSub1},
			Scopes: []schema.RoleAssignmentScope{
				{Tenants: []string{suite.Tenant}},
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
	params := &params.AuthorizationProviderQueriesV1Params{
		Roles:           roles,
		RoleAssignments: roleAssignments,
	}
	suite.params = params
	err = suites.SetupMockIfEnabled(suite.TestSuite, mockauthorization.ConfigureProviderQueriesV1, *params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

func (suite *ProviderQueriesV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, sdkconsts.AuthorizationProviderV1Name,
		string(schema.GlobalTenantResourceMetadataKindResourceKindRole),
		string(schema.GlobalTenantResourceMetadataKindResourceKindRoleAssignment),
	)

	stepsConfigurator := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Role
	roles := suite.params.Roles

	// Create roles
	steps.BulkCreateRolesStepsV1(stepsConfigurator, suite.GlobalTestSuite, "Create roles", roles)

	tpath := secapi.TenantPath{
		Tenant: secapi.TenantID(suite.Tenant),
	}

	// List Roles
	stepsConfigurator.ListRoleV1Step("List roles", suite.Client.AuthorizationV1, tpath, nil)

	// List Roles with limit
	stepsConfigurator.ListRoleV1Step("List roles with limit", suite.Client.AuthorizationV1, tpath,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List Roles with Label
	stepsConfigurator.ListRoleV1Step("List roles with label", suite.Client.AuthorizationV1, tpath,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List Roles with Limit and label
	stepsConfigurator.ListRoleV1Step("List roles with limit and label", suite.Client.AuthorizationV1, tpath,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Role assignment
	roleAssignments := suite.params.RoleAssignments

	// Create role assignments
	steps.BulkCreateRoleAssignmentsStepsV1(stepsConfigurator, suite.GlobalTestSuite, "Create role assignments", roleAssignments)

	tpath = secapi.TenantPath{
		Tenant: secapi.TenantID(suite.Tenant),
	}

	// List RoleAssignments
	stepsConfigurator.ListRoleAssignmentsV1("List role assignments", suite.Client.AuthorizationV1, tpath, nil)

	// List RoleAssignments with limit
	stepsConfigurator.ListRoleAssignmentsV1("List role assignments with limit", suite.Client.AuthorizationV1, tpath,
		secapi.NewListOptions().WithLimit(1))

	// List RoleAssignments with Label
	stepsConfigurator.ListRoleAssignmentsV1("List role assignments", suite.Client.AuthorizationV1, tpath,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// List RoleAssignments with Limit and label
	stepsConfigurator.ListRoleAssignmentsV1("List role assignments with limit and label", suite.Client.AuthorizationV1, tpath,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(constants.EnvLabel, constants.EnvConformanceLabel)))

	// Delete all role assignments
	steps.BulkDeleteRoleAssignmentsStepsV1(stepsConfigurator, suite.GlobalTestSuite, "Delete all role assignments", roleAssignments)

	// Delete all roles
	steps.BulkDeleteRolesStepsV1(stepsConfigurator, suite.GlobalTestSuite, "Delete all roles", roles)

	suite.FinishScenario()
}

func (suite *ProviderQueriesV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
