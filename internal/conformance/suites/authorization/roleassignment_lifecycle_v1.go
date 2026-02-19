package authorization

import (
	"log/slog"
	"math/rand"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockauthorization "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/authorization"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type RoleAssignmentLifeCycleV1TestSuite struct {
	suites.GlobalTestSuite

	Users []string

	params *params.RoleAssignmentLifeCycleV1Params
}

func CreateRoleAssignmentLifeCycleV1TestSuite(globalTestSuite suites.GlobalTestSuite, users []string) *RoleAssignmentLifeCycleV1TestSuite {
	suite := &RoleAssignmentLifeCycleV1TestSuite{
		GlobalTestSuite: globalTestSuite,
		Users:           users,
	}
	suite.ScenarioName = constants.RoleAssignmentLifeCycleV1SuiteName.String()
	return suite
}

func (suite *RoleAssignmentLifeCycleV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("Authorization")

	// Select subs
	roleAssignmentSub1 := suite.Users[rand.Intn(len(suite.Users))]
	roleAssignmentSub2 := suite.Users[rand.Intn(len(suite.Users))]

	// Generate scenario data
	roleName := generators.GenerateRoleName()

	roleAssignmentName := generators.GenerateRoleAssignmentName()

	roleAssignmentInitial, err := builders.NewRoleAssignmentBuilder().
		Name(roleAssignmentName).
		Provider(constants.AuthorizationProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).
		Spec(&schema.RoleAssignmentSpec{
			Roles: []string{roleName},
			Subs:  []string{roleAssignmentSub1},
			Scopes: []schema.RoleAssignmentScope{
				{Tenants: &[]string{suite.Tenant}},
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build RoleAssignment: %v", err)
	}

	roleAssignmentUpdated, err := builders.NewRoleAssignmentBuilder().
		Name(roleAssignmentName).
		Provider(constants.AuthorizationProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).
		Spec(&schema.RoleAssignmentSpec{
			Roles: []string{roleName},
			Subs:  []string{roleAssignmentSub1, roleAssignmentSub2},
			Scopes: []schema.RoleAssignmentScope{
				{Tenants: &[]string{suite.Tenant}},
			},
		}).Build()
	if err != nil {
		t.Fatalf("Failed to build RoleAssignment: %v", err)
	}

	params := &params.RoleAssignmentLifeCycleV1Params{
		RoleAssignmentInitial: roleAssignmentInitial,
		RoleAssignmentUpdated: roleAssignmentUpdated,
	}
	suite.params = params

	err = suites.SetupMockIfEnabled(suite.TestSuite, mockauthorization.ConfigureRoleAssignmentLifecycleScenarioV1, params)
	if err != nil {
		slog.Error("Failed to setup mock", "error", err)
		t.FailNow()
	}
}

func (suite *RoleAssignmentLifeCycleV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, constants.AuthorizationProviderV1,
		string(schema.GlobalTenantResourceMetadataKindResourceKindRole),
		string(schema.GlobalTenantResourceMetadataKindResourceKindRoleAssignment),
	)

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Role assignment

	// Create a role assignment
	roleAssign := suite.params.RoleAssignmentInitial
	expectRoleAssignMeta := roleAssign.Metadata
	expectRoleAssignSpec := &roleAssign.Spec
	stepsBuilder.CreateOrUpdateRoleAssignmentV1Step("Create a role assignment", suite.Client.AuthorizationV1, roleAssign,
		steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
			Metadata:      expectRoleAssignMeta,
			Spec:          expectRoleAssignSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created role assignment
	roleAssignTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(suite.Tenant),
		Name:   roleAssign.Metadata.Name,
	}
	roleAssign = stepsBuilder.GetRoleAssignmentV1Step("Get the created role assignment", suite.Client.AuthorizationV1, roleAssignTRef,
		steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
			Metadata:      expectRoleAssignMeta,
			Spec:          expectRoleAssignSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Update the role assignment
	roleAssign.Spec = suite.params.RoleAssignmentUpdated.Spec
	expectRoleAssignSpec.Subs = roleAssign.Spec.Subs
	stepsBuilder.CreateOrUpdateRoleAssignmentV1Step("Update the role assignment", suite.Client.AuthorizationV1, roleAssign,
		steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
			Metadata:      expectRoleAssignMeta,
			Spec:          expectRoleAssignSpec,
			ResourceState: schema.ResourceStateUpdating,
		},
	)

	// Get the updated role assignment
	roleAssign = stepsBuilder.GetRoleAssignmentV1Step("Get the updated role assignment", suite.Client.AuthorizationV1, roleAssignTRef,
		steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
			Metadata:      expectRoleAssignMeta,
			Spec:          expectRoleAssignSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Resources deletion
	stepsBuilder.DeleteRoleAssignmentV1Step("Delete the role assignment", suite.Client.AuthorizationV1, roleAssign)
	stepsBuilder.GetRoleAssignmentWithErrorV1Step("Get the deleted role assignment", suite.Client.AuthorizationV1, roleAssignTRef, secapi.ErrResourceNotFound)

	suite.FinishScenario()
}

func (suite *RoleAssignmentLifeCycleV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
