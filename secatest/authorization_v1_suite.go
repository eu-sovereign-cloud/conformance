package secatest

import (
	"context"
	"log/slog"
	"math/rand"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type AuthorizationV1TestSuite struct {
	globalTestSuite

	users []string
}

func (suite *AuthorizationV1TestSuite) TestSuite(t provider.T) {
	slog.Info("Starting " + suite.scenarioName)

	t.Title(suite.scenarioName)
	configureTags(t, secalib.AuthorizationProviderV1, secalib.RoleKind, secalib.RoleAssignmentKind)

	// Select subs
	roleAssignmentSub1 := suite.users[rand.Intn(len(suite.users))]
	roleAssignmentSub2 := suite.users[rand.Intn(len(suite.users))]

	// Generate scenario data
	roleName := secalib.GenerateRoleName()
	roleResource := secalib.GenerateRoleResource(suite.tenant, roleName)

	roleAssignmentName := secalib.GenerateRoleAssignmentName()
	roleAssignmentResource := secalib.GenerateRoleAssignmentResource(suite.tenant, roleAssignmentName)

	imageName := secalib.GenerateImageName()
	imageResource := secalib.GenerateImageResource(suite.tenant, imageName)

	// Setup mock, if configured to use
	if suite.mockEnabled {
		mockParams := &mock.AuthorizationParamsV1{
			Params: &mock.Params{
				MockURL:   *suite.mockServerURL,
				AuthToken: suite.authToken,
				Tenant:    suite.tenant,
			},
			Role: &mock.ResourceParams[schema.RoleSpec]{
				Name: roleName,
				InitialSpec: &schema.RoleSpec{
					Permissions: []schema.Permission{
						{Provider: secalib.StorageProviderV1, Resources: []string{imageResource}, Verb: []string{http.MethodGet}},
					},
				},
				UpdatedSpec: &schema.RoleSpec{
					Permissions: []schema.Permission{
						{Provider: secalib.StorageProviderV1, Resources: []string{imageResource}, Verb: []string{http.MethodGet, http.MethodPut}},
					},
				},
			},
			RoleAssignment: &mock.ResourceParams[schema.RoleAssignmentSpec]{
				Name: roleAssignmentName,
				InitialSpec: &schema.RoleAssignmentSpec{
					Roles: []string{roleName},
					Subs:  []string{roleAssignmentSub1},
					Scopes: []schema.RoleAssignmentScope{
						{Tenants: &[]string{suite.tenant}},
					},
				},
				UpdatedSpec: &schema.RoleAssignmentSpec{
					Roles: []string{roleName},
					Subs:  []string{roleAssignmentSub1, roleAssignmentSub2},
					Scopes: []schema.RoleAssignmentScope{
						{Tenants: &[]string{suite.tenant}},
					},
				},
			},
		}
		wm, err := mock.CreateAuthorizationLifecycleScenarioV1(suite.scenarioName, mockParams)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}
		suite.mockClient = wm
	}

	ctx := context.Background()

	// Role

	// Create a role
	role := &schema.Role{
		Metadata: &schema.GlobalTenantResourceMetadata{
			Tenant: suite.tenant,
			Name:   roleName,
		},
		Spec: schema.RoleSpec{
			Permissions: []schema.Permission{
				{
					Provider:  secalib.StorageProviderV1,
					Resources: []string{imageResource},
					Verb:      []string{http.MethodGet},
				},
			},
		},
	}
	expectRoleMeta := secalib.NewGlobalTenantResourceMetadata(roleName,
		secalib.AuthorizationProviderV1,
		roleResource,
		secalib.ApiVersion1,
		secalib.RoleKind,
		suite.tenant)
	expectRoleSpec := role.Spec
	suite.createOrUpdateRoleV1Step("Create a role", t, ctx, suite.client.AuthorizationV1, role,
		expectRoleMeta, &expectRoleSpec, secalib.CreatingResourceState)

	// Get the created role
	roleTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
		Name:   roleName,
	}
	role = suite.getRoleV1Step("Get the created role", t, ctx, suite.client.AuthorizationV1, *roleTRef,
		expectRoleMeta, &expectRoleSpec, secalib.ActiveResourceState)

	// Update the role
	role.Spec.Permissions[0].Verb = []string{http.MethodGet, http.MethodPut}
	expectRoleSpec = role.Spec
	suite.createOrUpdateRoleV1Step("Update the role", t, ctx, suite.client.AuthorizationV1, role,
		expectRoleMeta, &expectRoleSpec, secalib.UpdatingResourceState)

	// Get the updated role
	role = suite.getRoleV1Step("Get the updated role", t, ctx, suite.client.AuthorizationV1, *roleTRef,
		expectRoleMeta, &expectRoleSpec, secalib.ActiveResourceState)

	// Role assignment

	// Create a role assignment
	roleAssign := &schema.RoleAssignment{
		Metadata: &schema.GlobalTenantResourceMetadata{
			Tenant: suite.tenant,
			Name:   roleAssignmentName,
		},
		Spec: schema.RoleAssignmentSpec{
			Roles:  []string{roleName},
			Subs:   []string{roleAssignmentSub1},
			Scopes: []schema.RoleAssignmentScope{{Tenants: &[]string{suite.tenant}}},
		},
	}
	expectRoleAssignMeta := secalib.NewGlobalTenantResourceMetadata(roleAssignmentName,
		secalib.AuthorizationProviderV1,
		roleAssignmentResource,
		secalib.ApiVersion1,
		secalib.RoleAssignmentKind,
		suite.tenant)
	expectRoleAssignSpec := roleAssign.Spec
	suite.createOrUpdateRoleAssignmentV1Step("Create a role assignment", t, ctx, suite.client.AuthorizationV1, roleAssign,
		expectRoleAssignMeta, &expectRoleAssignSpec, secalib.CreatingResourceState)

	// Get the created role assignment
	roleAssignTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
		Name:   roleAssignmentName,
	}
	roleAssign = suite.getRoleAssignmentV1Step("Get the created role assignment", t, ctx, suite.client.AuthorizationV1, *roleAssignTRef,
		expectRoleAssignMeta, &expectRoleAssignSpec, secalib.ActiveResourceState)

	// Update the role assignment
	roleAssign.Spec.Subs = []string{roleAssignmentSub1, roleAssignmentSub2}
	expectRoleAssignSpec = roleAssign.Spec
	suite.createOrUpdateRoleAssignmentV1Step("Update the role assignment", t, ctx, suite.client.AuthorizationV1, roleAssign,
		expectRoleAssignMeta, &expectRoleAssignSpec, secalib.UpdatingResourceState)

	// Get the updated role assignment
	roleAssign = suite.getRoleAssignmentV1Step("Get the updated role assignment", t, ctx, suite.client.AuthorizationV1, *roleAssignTRef,
		expectRoleAssignMeta, &expectRoleAssignSpec, secalib.ActiveResourceState)

	// Resources deletion

	// Delete the role assignment
	suite.deleteRoleAssignmentV1Step("Delete the role assignment", t, ctx, suite.client.AuthorizationV1, roleAssign)

	// Get the deleted role assignment
	suite.getRoleAssignmentWithErrorV1Step("Get the deleted role assignment", t, ctx, suite.client.AuthorizationV1, *roleAssignTRef, secapi.ErrResourceNotFound)

	// Delete the role
	suite.deleteRoleV1Step("Delete the role", t, ctx, suite.client.AuthorizationV1, role)

	// Get the deleted role
	suite.getRoleWithErrorV1Step("Get the deleted role", t, ctx, suite.client.AuthorizationV1, *roleTRef, secapi.ErrResourceNotFound)

	slog.Info("Finishing " + suite.scenarioName)
}

func (suite *AuthorizationV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}
