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

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type AuthorizationLifeCycleV1TestSuite struct {
	suites.GlobalTestSuite

	Users []string

	params *params.AuthorizationLifeCycleV1Params
}

func CreateLifeCycleV1TestSuite(globalTestSuite suites.GlobalTestSuite, users []string) *AuthorizationLifeCycleV1TestSuite {
	suite := &AuthorizationLifeCycleV1TestSuite{
		GlobalTestSuite: globalTestSuite,
		Users:           users,
	}
	suite.ScenarioName = constants.AuthorizationV1LifeCycleSuiteName
	return suite
}

func (suite *AuthorizationLifeCycleV1TestSuite) BeforeAll(t provider.T) {
	var err error

	// Select subs
	roleAssignmentSub1 := suite.Users[rand.Intn(len(suite.Users))]
	roleAssignmentSub2 := suite.Users[rand.Intn(len(suite.Users))]

	// Generate scenario data
	roleName := generators.GenerateRoleName()

	roleAssignmentName := generators.GenerateRoleAssignmentName()

	imageName := generators.GenerateImageName()
	imageResource := generators.GenerateImageResource(suite.Tenant, imageName)

	roleInitial, err := builders.NewRoleBuilder().
		Name(roleName).
		Provider(constants.AuthorizationProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).
		Spec(&schema.RoleSpec{
			Permissions: []schema.Permission{
				{Provider: constants.StorageProviderV1, Resources: []string{imageResource}, Verb: []string{http.MethodGet}},
			},
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Role: %v", err)
	}

	roleUpdated, err := builders.NewRoleBuilder().
		Name(roleName).
		Provider(constants.AuthorizationProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).
		Spec(&schema.RoleSpec{
			Permissions: []schema.Permission{
				{Provider: constants.StorageProviderV1, Resources: []string{imageResource}, Verb: []string{http.MethodGet, http.MethodPut}},
			},
		}).
		Build()
	if err != nil {
		t.Fatalf("Failed to build Role: %v", err)
	}

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

	params := &params.AuthorizationLifeCycleV1Params{
		RoleInitial:           roleInitial,
		RoleUpdated:           roleUpdated,
		RoleAssignmentInitial: roleAssignmentInitial,
		RoleAssignmentUpdated: roleAssignmentUpdated,
	}
	suite.params = params
	err = suites.SetupMockIfEnabled(suite.TestSuite, mockauthorization.ConfigureLifecycleScenarioV1, params)
	if err != nil {
		t.Fatalf("Failed to setup mock: %v", err)
	}
}

//nolint:dupl
func (suite *AuthorizationLifeCycleV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, constants.AuthorizationProviderV1,
		string(schema.GlobalTenantResourceMetadataKindResourceKindRole),
		string(schema.GlobalTenantResourceMetadataKindResourceKindRoleAssignment),
	)

	// Role
	role := suite.params.RoleInitial
	expectRoleMeta := role.Metadata
	expectRoleSpec := &role.Spec
	roleTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(suite.Tenant),
		Name:   role.Metadata.Name,
	}

	t.WithNewStep("Role", func(roleCtx provider.StepCtx) {
		roleSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, roleCtx)

		roleSteps.CreateOrUpdateRoleV1Step("Create", suite.Client.AuthorizationV1, role,
			steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
				Metadata:      expectRoleMeta,
				Spec:          expectRoleSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		role = roleSteps.GetRoleV1Step("Get", suite.Client.AuthorizationV1, roleTRef,
			steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
				Metadata:      expectRoleMeta,
				Spec:          expectRoleSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)

		role.Spec = suite.params.RoleUpdated.Spec
		expectRoleSpec = &role.Spec
		roleSteps.CreateOrUpdateRoleV1Step("Update", suite.Client.AuthorizationV1, role,
			steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
				Metadata:      expectRoleMeta,
				Spec:          expectRoleSpec,
				ResourceState: schema.ResourceStateUpdating,
			},
		)

		role = roleSteps.GetRoleV1Step("GetUpdated", suite.Client.AuthorizationV1, roleTRef,
			steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
				Metadata:      expectRoleMeta,
				Spec:          expectRoleSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	// Role assignment
	roleAssign := suite.params.RoleAssignmentInitial
	expectRoleAssignMeta := roleAssign.Metadata
	expectRoleAssignSpec := &roleAssign.Spec
	roleAssignTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(suite.Tenant),
		Name:   roleAssign.Metadata.Name,
	}

	t.WithNewStep("RoleAssignment", func(raCtx provider.StepCtx) {
		raSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, raCtx)

		raSteps.CreateOrUpdateRoleAssignmentV1Step("Create", suite.Client.AuthorizationV1, roleAssign,
			steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
				Metadata:      expectRoleAssignMeta,
				Spec:          expectRoleAssignSpec,
				ResourceState: schema.ResourceStateCreating,
			},
		)

		roleAssign = raSteps.GetRoleAssignmentV1Step("Get", suite.Client.AuthorizationV1, roleAssignTRef,
			steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
				Metadata:      expectRoleAssignMeta,
				Spec:          expectRoleAssignSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)

		roleAssign.Spec = suite.params.RoleAssignmentUpdated.Spec
		expectRoleAssignSpec.Subs = roleAssign.Spec.Subs
		raSteps.CreateOrUpdateRoleAssignmentV1Step("Update", suite.Client.AuthorizationV1, roleAssign,
			steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
				Metadata:      expectRoleAssignMeta,
				Spec:          expectRoleAssignSpec,
				ResourceState: schema.ResourceStateUpdating,
			},
		)

		roleAssign = raSteps.GetRoleAssignmentV1Step("GetUpdated", suite.Client.AuthorizationV1, roleAssignTRef,
			steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
				Metadata:      expectRoleAssignMeta,
				Spec:          expectRoleAssignSpec,
				ResourceState: schema.ResourceStateActive,
			},
		)
	})

	t.WithNewStep("Delete", func(delctx provider.StepCtx) {
		delStep := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, delctx)

		delStep.DeleteRoleAssignmentV1Step("Delete the role assignment", suite.Client.AuthorizationV1, roleAssign)
		delStep.GetRoleAssignmentWithErrorV1Step("Get the deleted role assignment", suite.Client.AuthorizationV1, roleAssignTRef, secapi.ErrResourceNotFound)

		delStep.DeleteRoleV1Step("Delete the role", suite.Client.AuthorizationV1, role)
		delStep.GetRoleWithErrorV1Step("Get the deleted role", suite.Client.AuthorizationV1, roleTRef, secapi.ErrResourceNotFound)
	})

	suite.FinishScenario()
}

func (suite *AuthorizationLifeCycleV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
