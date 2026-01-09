package authorization

import (
	"math/rand"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	mockauthorization "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/authorization"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type LifeCycleV1TestSuite struct {
	suites.GlobalTestSuite

	Users []string
}

func (suite *LifeCycleV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, constants.AuthorizationProviderV1,
		string(schema.GlobalTenantResourceMetadataKindResourceKindRole),
		string(schema.GlobalTenantResourceMetadataKindResourceKindRoleAssignment),
	)

	// Select subs
	roleAssignmentSub1 := suite.Users[rand.Intn(len(suite.Users))]
	roleAssignmentSub2 := suite.Users[rand.Intn(len(suite.Users))]

	// Generate scenario data
	roleName := generators.GenerateRoleName()

	roleAssignmentName := generators.GenerateRoleAssignmentName()

	imageName := generators.GenerateImageName()
	imageResource := generators.GenerateImageResource(suite.Tenant, imageName)

	// Setup mock, if configured to use
	if suite.MockEnabled {
		mockParams := &params.AuthorizationLifeCycleParamsV1{
			BaseParams: &params.BaseParams{
				Tenant: suite.Tenant,
				MockParams: &mock.MockParams{
					ServerURL: *suite.MockServerURL,
					AuthToken: suite.AuthToken,
				},
			},
			Role: &params.ResourceParams[schema.RoleSpec]{
				Name: roleName,
				InitialSpec: &schema.RoleSpec{
					Permissions: []schema.Permission{
						{Provider: constants.StorageProviderV1, Resources: []string{imageResource}, Verb: []string{http.MethodGet}},
					},
				},
				UpdatedSpec: &schema.RoleSpec{
					Permissions: []schema.Permission{
						{Provider: constants.StorageProviderV1, Resources: []string{imageResource}, Verb: []string{http.MethodGet, http.MethodPut}},
					},
				},
			},
			RoleAssignment: &params.ResourceParams[schema.RoleAssignmentSpec]{
				Name: roleAssignmentName,
				InitialSpec: &schema.RoleAssignmentSpec{
					Roles: []string{roleName},
					Subs:  []string{roleAssignmentSub1},
					Scopes: []schema.RoleAssignmentScope{
						{Tenants: &[]string{suite.Tenant}},
					},
				},
				UpdatedSpec: &schema.RoleAssignmentSpec{
					Roles: []string{roleName},
					Subs:  []string{roleAssignmentSub1, roleAssignmentSub2},
					Scopes: []schema.RoleAssignmentScope{
						{Tenants: &[]string{suite.Tenant}},
					},
				},
			},
		}
		wm, err := mockauthorization.ConfigureLifecycleScenarioV1(suite.ScenarioName, mockParams)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}
		suite.MockClient = wm
	}

	stepsBuilder := steps.NewStepsConfigurator(&suite.TestSuite, t)

	// Role

	// Create a role
	role := &schema.Role{
		Metadata: &schema.GlobalTenantResourceMetadata{
			Tenant: suite.Tenant,
			Name:   roleName,
		},
		Spec: schema.RoleSpec{
			Permissions: []schema.Permission{
				{
					Provider:  constants.StorageProviderV1,
					Resources: []string{imageResource},
					Verb:      []string{http.MethodGet},
				},
			},
		},
	}
	expectRoleMeta, err := builders.NewRoleMetadataBuilder().
		Name(roleName).
		Provider(constants.AuthorizationProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).
		Build()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectRoleSpec := &schema.RoleSpec{
		Permissions: []schema.Permission{
			{
				Provider:  constants.StorageProviderV1,
				Resources: []string{imageResource},
				Verb:      []string{http.MethodGet},
			},
		},
	}
	stepsBuilder.CreateOrUpdateRoleV1Step("Create a role", suite.Client.AuthorizationV1, role,
		steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
			Metadata:      expectRoleMeta,
			Spec:          expectRoleSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created role
	roleTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.Tenant),
		Name:   roleName,
	}
	role = stepsBuilder.GetRoleV1Step("Get the created role", suite.Client.AuthorizationV1, *roleTRef,
		steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
			Metadata:      expectRoleMeta,
			Spec:          expectRoleSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Update the role
	role.Spec.Permissions[0].Verb = []string{http.MethodGet, http.MethodPut}
	expectRoleSpec.Permissions[0] = role.Spec.Permissions[0]
	stepsBuilder.CreateOrUpdateRoleV1Step("Update the role", suite.Client.AuthorizationV1, role,
		steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
			Metadata:      expectRoleMeta,
			Spec:          expectRoleSpec,
			ResourceState: schema.ResourceStateUpdating,
		},
	)

	// Get the updated role
	role = stepsBuilder.GetRoleV1Step("Get the updated role", suite.Client.AuthorizationV1, *roleTRef,
		steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
			Metadata:      expectRoleMeta,
			Spec:          expectRoleSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Role assignment

	// Create a role assignment
	roleAssign := &schema.RoleAssignment{
		Metadata: &schema.GlobalTenantResourceMetadata{
			Tenant: suite.Tenant,
			Name:   roleAssignmentName,
		},
		Spec: schema.RoleAssignmentSpec{
			Roles:  []string{roleName},
			Subs:   []string{roleAssignmentSub1},
			Scopes: []schema.RoleAssignmentScope{{Tenants: &[]string{suite.Tenant}}},
		},
	}
	expectRoleAssignMeta, err := builders.NewRoleAssignmentMetadataBuilder().
		Name(roleAssignmentName).
		Provider(constants.AuthorizationProviderV1).ApiVersion(constants.ApiVersion1).
		Tenant(suite.Tenant).
		Build()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectRoleAssignSpec := &schema.RoleAssignmentSpec{
		Roles:  []string{roleName},
		Subs:   []string{roleAssignmentSub1},
		Scopes: []schema.RoleAssignmentScope{{Tenants: &[]string{suite.Tenant}}},
	}
	stepsBuilder.CreateOrUpdateRoleAssignmentV1Step("Create a role assignment", suite.Client.AuthorizationV1, roleAssign,
		steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
			Metadata:      expectRoleAssignMeta,
			Spec:          expectRoleAssignSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created role assignment
	roleAssignTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.Tenant),
		Name:   roleAssignmentName,
	}
	roleAssign = stepsBuilder.GetRoleAssignmentV1Step("Get the created role assignment", suite.Client.AuthorizationV1, *roleAssignTRef,
		steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
			Metadata:      expectRoleAssignMeta,
			Spec:          expectRoleAssignSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Update the role assignment
	roleAssign.Spec.Subs = []string{roleAssignmentSub1, roleAssignmentSub2}
	expectRoleAssignSpec.Subs = roleAssign.Spec.Subs
	stepsBuilder.CreateOrUpdateRoleAssignmentV1Step("Update the role assignment", suite.Client.AuthorizationV1, roleAssign,
		steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
			Metadata:      expectRoleAssignMeta,
			Spec:          expectRoleAssignSpec,
			ResourceState: schema.ResourceStateUpdating,
		},
	)

	// Get the updated role assignment
	roleAssign = stepsBuilder.GetRoleAssignmentV1Step("Get the updated role assignment", suite.Client.AuthorizationV1, *roleAssignTRef,
		steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
			Metadata:      expectRoleAssignMeta,
			Spec:          expectRoleAssignSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Resources deletion

	stepsBuilder.DeleteRoleAssignmentV1Step("Delete the role assignment", suite.Client.AuthorizationV1, roleAssign)
	stepsBuilder.GetRoleAssignmentWithErrorV1Step("Get the deleted role assignment", suite.Client.AuthorizationV1, *roleAssignTRef, secapi.ErrResourceNotFound)

	stepsBuilder.DeleteRoleV1Step("Delete the role", suite.Client.AuthorizationV1, role)
	stepsBuilder.GetRoleWithErrorV1Step("Get the deleted role", suite.Client.AuthorizationV1, *roleTRef, secapi.ErrResourceNotFound)

	suite.FinishScenario()
}

func (suite *LifeCycleV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
