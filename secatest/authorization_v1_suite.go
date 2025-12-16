package secatest

import (
	"log/slog"
	"math/rand"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	labelBuilder "github.com/eu-sovereign-cloud/go-sdk/secapi/builders"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type AuthorizationV1TestSuite struct {
	globalTestSuite

	users []string
}

func (suite *AuthorizationV1TestSuite) TestSuite(t provider.T) {
	slog.Info("Starting " + suite.scenarioName)

	t.Title(suite.scenarioName)
	configureTags(t, authorizationProviderV1,
		string(schema.GlobalTenantResourceMetadataKindResourceKindRole),
		string(schema.GlobalTenantResourceMetadataKindResourceKindRoleAssignment),
	)

	// Select subs
	roleAssignmentSub1 := suite.users[rand.Intn(len(suite.users))]
	roleAssignmentSub2 := suite.users[rand.Intn(len(suite.users))]

	// Generate scenario data
	roleName := generators.GenerateRoleName()

	roleAssignmentName := generators.GenerateRoleAssignmentName()

	imageName := generators.GenerateImageName()
	imageResource := generators.GenerateImageResource(suite.tenant, imageName)

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
						{Provider: storageProviderV1, Resources: []string{imageResource}, Verb: []string{http.MethodGet}},
					},
				},
				UpdatedSpec: &schema.RoleSpec{
					Permissions: []schema.Permission{
						{Provider: storageProviderV1, Resources: []string{imageResource}, Verb: []string{http.MethodGet, http.MethodPut}},
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
					Provider:  storageProviderV1,
					Resources: []string{imageResource},
					Verb:      []string{http.MethodGet},
				},
			},
		},
	}
	expectRoleMeta, err := builders.NewRoleMetadataBuilder().
		Name(roleName).
		Provider(authorizationProviderV1).ApiVersion(apiVersion1).
		Tenant(suite.tenant).
		Build()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectRoleSpec := &schema.RoleSpec{
		Permissions: []schema.Permission{
			{
				Provider:  storageProviderV1,
				Resources: []string{imageResource},
				Verb:      []string{http.MethodGet},
			},
		},
	}
	suite.createOrUpdateRoleV1Step("Create a role", t, suite.client.AuthorizationV1, role,
		responseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
			metadata:      expectRoleMeta,
			spec:          expectRoleSpec,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created role
	roleTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
		Name:   roleName,
	}
	role = suite.getRoleV1Step("Get the created role", t, suite.client.AuthorizationV1, *roleTRef,
		responseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
			metadata:      expectRoleMeta,
			spec:          expectRoleSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Update the role
	role.Spec.Permissions[0].Verb = []string{http.MethodGet, http.MethodPut}
	expectRoleSpec.Permissions[0] = role.Spec.Permissions[0]
	suite.createOrUpdateRoleV1Step("Update the role", t, suite.client.AuthorizationV1, role,
		responseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
			metadata:      expectRoleMeta,
			spec:          expectRoleSpec,
			resourceState: schema.ResourceStateUpdating,
		},
	)

	// Get the updated role
	role = suite.getRoleV1Step("Get the updated role", t, suite.client.AuthorizationV1, *roleTRef,
		responseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
			metadata:      expectRoleMeta,
			spec:          expectRoleSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

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
	expectRoleAssignMeta, err := builders.NewRoleAssignmentMetadataBuilder().
		Name(roleAssignmentName).
		Provider(authorizationProviderV1).ApiVersion(apiVersion1).
		Tenant(suite.tenant).
		Build()
	if err != nil {
		t.Fatalf("Failed to build metadata: %v", err)
	}
	expectRoleAssignSpec := &schema.RoleAssignmentSpec{
		Roles:  []string{roleName},
		Subs:   []string{roleAssignmentSub1},
		Scopes: []schema.RoleAssignmentScope{{Tenants: &[]string{suite.tenant}}},
	}
	suite.createOrUpdateRoleAssignmentV1Step("Create a role assignment", t, suite.client.AuthorizationV1, roleAssign,
		responseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
			metadata:      expectRoleAssignMeta,
			spec:          expectRoleAssignSpec,
			resourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created role assignment
	roleAssignTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
		Name:   roleAssignmentName,
	}
	roleAssign = suite.getRoleAssignmentV1Step("Get the created role assignment", t, suite.client.AuthorizationV1, *roleAssignTRef,
		responseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
			metadata:      expectRoleAssignMeta,
			spec:          expectRoleAssignSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Update the role assignment
	roleAssign.Spec.Subs = []string{roleAssignmentSub1, roleAssignmentSub2}
	expectRoleAssignSpec.Subs = roleAssign.Spec.Subs
	suite.createOrUpdateRoleAssignmentV1Step("Update the role assignment", t, suite.client.AuthorizationV1, roleAssign,
		responseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
			metadata:      expectRoleAssignMeta,
			spec:          expectRoleAssignSpec,
			resourceState: schema.ResourceStateUpdating,
		},
	)

	// Get the updated role assignment
	roleAssign = suite.getRoleAssignmentV1Step("Get the updated role assignment", t, suite.client.AuthorizationV1, *roleAssignTRef,
		responseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
			metadata:      expectRoleAssignMeta,
			spec:          expectRoleAssignSpec,
			resourceState: schema.ResourceStateActive,
		},
	)

	// Resources deletion

	suite.deleteRoleAssignmentV1Step("Delete the role assignment", t, suite.client.AuthorizationV1, roleAssign)
	suite.getRoleAssignmentWithErrorV1Step("Get the deleted role assignment", t, suite.client.AuthorizationV1, *roleAssignTRef, secapi.ErrResourceNotFound)

	suite.deleteRoleV1Step("Delete the role", t, suite.client.AuthorizationV1, role)
	suite.getRoleWithErrorV1Step("Get the deleted role", t, suite.client.AuthorizationV1, *roleTRef, secapi.ErrResourceNotFound)

	slog.Info("Finishing " + suite.scenarioName)
}

func (suite *AuthorizationV1TestSuite) TestSuiteListScenarios(t provider.T) {
	slog.Info("Starting " + suite.scenarioName)

	t.Title(suite.scenarioName)
	configureTags(t, authorizationProviderV1,
		string(schema.GlobalTenantResourceMetadataKindResourceKindRole),
		string(schema.GlobalTenantResourceMetadataKindResourceKindRoleAssignment),
	)

	// Select subs
	roleAssignmentSub1 := suite.users[rand.Intn(len(suite.users))]

	roleName1 := generators.GenerateRoleName()
	roleName2 := generators.GenerateRoleName()
	roleName3 := generators.GenerateRoleName()
	// Generate scenario data

	roleAssignmentName1 := generators.GenerateRoleAssignmentName()
	roleAssignmentName2 := generators.GenerateRoleAssignmentName()
	roleAssignmentName3 := generators.GenerateRoleAssignmentName()
	imageName := generators.GenerateImageName()
	imageResource := generators.GenerateImageResource(suite.tenant, imageName)

	// Setup mock, if configured to use
	if suite.mockEnabled {
		mockParams := &mock.AuthorizationListParamsV1{
			Params: &mock.Params{
				MockURL:   *suite.mockServerURL,
				AuthToken: suite.authToken,
				Tenant:    suite.tenant,
			},
			Role: &[]mock.ResourceParams[schema.RoleSpec]{
				{
					Name: roleName1,

					InitialLabels: schema.Labels{
						generators.EnvLabel: generators.EnvConformanceLabel,
					},
					InitialSpec: &schema.RoleSpec{
						Permissions: []schema.Permission{
							{Provider: storageProviderV1, Resources: []string{imageResource}, Verb: []string{http.MethodGet}},
						},
					},
				},
				{
					Name: roleName2,
					InitialLabels: schema.Labels{
						generators.EnvLabel: generators.EnvConformanceLabel,
					},
					InitialSpec: &schema.RoleSpec{
						Permissions: []schema.Permission{
							{Provider: storageProviderV1, Resources: []string{imageResource}, Verb: []string{http.MethodGet}},
						},
					},
				},
				{
					Name: roleName3,
					InitialLabels: schema.Labels{
						generators.EnvLabel: generators.EnvConformanceLabel,
					},
					InitialSpec: &schema.RoleSpec{
						Permissions: []schema.Permission{
							{Provider: storageProviderV1, Resources: []string{imageResource}, Verb: []string{http.MethodGet}},
						},
					},
				},
			},
			RoleAssignment: &[]mock.ResourceParams[schema.RoleAssignmentSpec]{
				{
					Name: roleAssignmentName1,
					InitialLabels: schema.Labels{
						generators.EnvLabel: generators.EnvConformanceLabel,
					},
					InitialSpec: &schema.RoleAssignmentSpec{
						Roles: []string{roleName1},
						Subs:  []string{roleAssignmentSub1},
						Scopes: []schema.RoleAssignmentScope{
							{Tenants: &[]string{suite.tenant}},
						},
					},
				},
				{
					Name: roleAssignmentName2,
					InitialLabels: schema.Labels{
						generators.EnvLabel: generators.EnvConformanceLabel,
					},
					InitialSpec: &schema.RoleAssignmentSpec{
						Roles: []string{roleName2},
						Subs:  []string{roleAssignmentSub1},
						Scopes: []schema.RoleAssignmentScope{
							{Tenants: &[]string{suite.tenant}},
						},
					},
				},
				{
					Name: roleAssignmentName3,
					InitialLabels: schema.Labels{
						generators.EnvLabel: generators.EnvConformanceLabel,
					},
					InitialSpec: &schema.RoleAssignmentSpec{
						Roles: []string{roleName3},
						Subs:  []string{roleAssignmentSub1},
						Scopes: []schema.RoleAssignmentScope{
							{Tenants: &[]string{suite.tenant}},
						},
					},
				},
			},
		}
		wm, err := mock.CreateAuthorizationListAndFilterScenarioV1(suite.scenarioName, mockParams)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}
		suite.mockClient = wm
	}

	// Role

	roles := []schema.Role{
		{
			Metadata: &schema.GlobalTenantResourceMetadata{
				Tenant: suite.tenant,
				Name:   roleName1,
			},
			Labels: map[string]string{
				generators.EnvLabel: generators.EnvConformanceLabel,
			},
			Spec: schema.RoleSpec{
				Permissions: []schema.Permission{
					{
						Provider:  storageProviderV1,
						Resources: []string{imageResource},
						Verb:      []string{http.MethodGet},
					},
				},
			},
		},
		{
			Metadata: &schema.GlobalTenantResourceMetadata{
				Tenant: suite.tenant,
				Name:   roleName2,
			},
			Labels: map[string]string{
				generators.EnvLabel: generators.EnvConformanceLabel,
			},
			Spec: schema.RoleSpec{
				Permissions: []schema.Permission{
					{
						Provider:  storageProviderV1,
						Resources: []string{imageResource},
						Verb:      []string{http.MethodGet},
					},
				},
			},
		},
		{
			Metadata: &schema.GlobalTenantResourceMetadata{
				Tenant: suite.tenant,
				Name:   roleName3,
			},
			Labels: map[string]string{
				generators.EnvLabel: generators.EnvConformanceLabel,
			},
			Spec: schema.RoleSpec{
				Permissions: []schema.Permission{
					{
						Provider:  storageProviderV1,
						Resources: []string{imageResource},
						Verb:      []string{http.MethodGet},
					},
				},
			},
		},
	}
	// Create a roles

	for _, role := range roles {

		role := &schema.Role{
			Metadata: &schema.GlobalTenantResourceMetadata{
				Tenant: suite.tenant,
				Name:   role.Metadata.Name,
			},
			Labels: map[string]string{
				generators.EnvLabel: generators.EnvConformanceLabel,
			},
			Spec: schema.RoleSpec{
				Permissions: []schema.Permission{
					{
						Provider:  storageProviderV1,
						Resources: []string{imageResource},
						Verb:      []string{http.MethodGet},
					},
				},
			},
		}
		expectRoleMeta, err := builders.NewRoleMetadataBuilder().
			Name(role.Metadata.Name).
			Provider(authorizationProviderV1).ApiVersion(apiVersion1).
			Tenant(suite.tenant).
			Build()
		if err != nil {
			t.Fatalf("Failed to build metadata: %v", err)
		}

		expectRoleSpec := &schema.RoleSpec{
			Permissions: []schema.Permission{
				{
					Provider:  storageProviderV1,
					Resources: []string{imageResource},
					Verb:      []string{http.MethodGet},
				},
			},
		}

		// Create Role
		suite.createOrUpdateRoleV1Step("Create a role", t, suite.client.AuthorizationV1, role,
			responseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
				metadata:      expectRoleMeta,
				spec:          expectRoleSpec,
				resourceState: schema.ResourceStateCreating,
			},
		)
	}
	roleTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
		Name:   suite.tenant,
	}
	// List Roles
	suite.getListRoleV1Step("Get list of roles", t, suite.client.AuthorizationV1, *roleTRef, nil)

	// List Roles with limit
	suite.getListRoleV1Step("Get list of roles with limit", t, suite.client.AuthorizationV1, *roleTRef,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().Equals(generators.EnvLabel, generators.EnvConformanceLabel)))

	// List Roles with Label
	suite.getListRoleV1Step("Get list of roles with label", t, suite.client.AuthorizationV1, *roleTRef,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(generators.EnvLabel, generators.EnvConformanceLabel)))

	// List Roles with Limit and label
	suite.getListRoleV1Step("Get list of roles with limit and label", t, suite.client.AuthorizationV1, *roleTRef,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(generators.EnvLabel, generators.EnvConformanceLabel)))

	// Role assignment

	roleAssignments := []schema.RoleAssignment{
		{
			Metadata: &schema.GlobalTenantResourceMetadata{
				Tenant: suite.tenant,
				Name:   roleAssignmentName1,
			},
			Labels: map[string]string{
				generators.EnvLabel: generators.EnvConformanceLabel,
			},
			Spec: schema.RoleAssignmentSpec{
				Roles: []string{roleName1},
				Subs:  []string{roleAssignmentSub1},
				Scopes: []schema.RoleAssignmentScope{
					{Tenants: &[]string{suite.tenant}},
				},
			},
		},
		{
			Metadata: &schema.GlobalTenantResourceMetadata{
				Tenant: suite.tenant,
				Name:   roleAssignmentName2,
			},
			Labels: map[string]string{
				generators.EnvLabel: generators.EnvConformanceLabel,
			},
			Spec: schema.RoleAssignmentSpec{
				Roles: []string{roleName2},
				Subs:  []string{roleAssignmentSub1},
				Scopes: []schema.RoleAssignmentScope{
					{Tenants: &[]string{suite.tenant}},
				},
			},
		},
		{
			Metadata: &schema.GlobalTenantResourceMetadata{
				Tenant: suite.tenant,
				Name:   roleAssignmentName3,
			},
			Labels: map[string]string{
				generators.EnvLabel: generators.EnvConformanceLabel,
			},
			Spec: schema.RoleAssignmentSpec{
				Roles: []string{roleName3},
				Subs:  []string{roleAssignmentSub1},
				Scopes: []schema.RoleAssignmentScope{
					{Tenants: &[]string{suite.tenant}},
				},
			},
		},
	}
	// Create a RoleAssignments

	for _, roleAssign := range roleAssignments {

		expectRoleAssignMeta, err := builders.NewRoleAssignmentMetadataBuilder().
			Name(roleAssign.Metadata.Name).
			Provider(authorizationProviderV1).ApiVersion(apiVersion1).
			Tenant(suite.tenant).
			Build()
		if err != nil {
			t.Fatalf("Failed to build metadata: %v", err)
		}
		expectRoleAssignSpec := &roleAssign.Spec
		// Create RoleAssignement
		suite.createOrUpdateRoleAssignmentV1Step("Create a role Assignment", t, suite.client.AuthorizationV1, &roleAssign,
			responseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
				metadata:      expectRoleAssignMeta,
				spec:          expectRoleAssignSpec,
				resourceState: schema.ResourceStateCreating,
			},
		)
	}
	roleAssignTRef := &secapi.TenantReference{
		Tenant: secapi.TenantID(suite.tenant),
	}
	// List RoleAssignments
	suite.getListRoleAssignmentsV1("Get list of role assignments", t, suite.client.AuthorizationV1, *roleAssignTRef, nil)
	// List RoleAssignments with limit
	suite.getListRoleAssignmentsV1("Get list of role assignments", t, suite.client.AuthorizationV1, *roleAssignTRef,
		secapi.NewListOptions().WithLimit(1))

	// List RoleAssignments with Label
	suite.getListRoleAssignmentsV1("Get list of role assignments", t, suite.client.AuthorizationV1, *roleAssignTRef,
		secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(generators.EnvLabel, generators.EnvConformanceLabel)))

	// List RoleAssignments with Limit and label
	suite.getListRoleAssignmentsV1("Get list of role assignments", t, suite.client.AuthorizationV1, *roleAssignTRef,
		secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
			Equals(generators.EnvLabel, generators.EnvConformanceLabel)))

	// Delete all role assignments
	for _, roleAssign := range roleAssignments {
		suite.deleteRoleAssignmentV1Step("Delete the role assignment", t, suite.client.AuthorizationV1, &roleAssign)

		// Get the deleted role assignment
		roleAssignTRefSingle := &secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   roleAssign.Metadata.Name,
		}
		suite.getRoleAssignmentWithErrorV1Step("Get the deleted role assignment", t, suite.client.AuthorizationV1, *roleAssignTRefSingle, secapi.ErrResourceNotFound)
	}

	// Delete all roles
	for _, role := range roles {
		suite.deleteRoleV1Step("Delete the role", t, suite.client.AuthorizationV1, &role)

		// Get the deleted role
		roleTRefSingle := &secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   role.Metadata.Name,
		}
		suite.getRoleWithErrorV1Step("Get the deleted role", t, suite.client.AuthorizationV1, *roleTRefSingle, secapi.ErrResourceNotFound)
	}

	slog.Info("Finishing " + suite.scenarioName)
}

func (suite *AuthorizationV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}
