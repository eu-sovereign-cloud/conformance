package authorization

import (
	"math/rand"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/steps"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/authorization"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	labelBuilder "github.com/eu-sovereign-cloud/go-sdk/secapi/builders"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type AuthorizationV1ListTestSuite struct {
	suites.GlobalTestSuite

	Users []string
}

func (suite *AuthorizationV1ListTestSuite) TestListScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, constants.AuthorizationProviderV1,
		string(schema.GlobalTenantResourceMetadataKindResourceKindRole),
		string(schema.GlobalTenantResourceMetadataKindResourceKindRoleAssignment),
	)

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

	// Setup mock, if configured to use
	if suite.MockEnabled {
		mockParams := &mock.AuthorizationListParamsV1{
			BaseParams: &mock.BaseParams{
				MockURL:   *suite.MockServerURL,
				AuthToken: suite.AuthToken,
				Tenant:    suite.Tenant,
			},
			Roles: []mock.ResourceParams[schema.RoleSpec]{
				{
					Name: roleName1,

					InitialLabels: schema.Labels{
						constants.EnvLabel: constants.EnvConformanceLabel,
					},
					InitialSpec: &schema.RoleSpec{
						Permissions: []schema.Permission{
							{Provider: constants.StorageProviderV1, Resources: []string{imageResource}, Verb: []string{http.MethodGet}},
						},
					},
				},
				{
					Name: roleName2,
					InitialLabels: schema.Labels{
						constants.EnvLabel: constants.EnvConformanceLabel,
					},
					InitialSpec: &schema.RoleSpec{
						Permissions: []schema.Permission{
							{Provider: constants.StorageProviderV1, Resources: []string{imageResource}, Verb: []string{http.MethodGet}},
						},
					},
				},
				{
					Name: roleName3,
					InitialLabels: schema.Labels{
						constants.EnvLabel: constants.EnvConformanceLabel,
					},
					InitialSpec: &schema.RoleSpec{
						Permissions: []schema.Permission{
							{Provider: constants.StorageProviderV1, Resources: []string{imageResource}, Verb: []string{http.MethodGet}},
						},
					},
				},
			},
			RoleAssignments: []mock.ResourceParams[schema.RoleAssignmentSpec]{
				{
					Name: roleAssignmentName1,
					InitialLabels: schema.Labels{
						constants.EnvLabel: constants.EnvConformanceLabel,
					},
					InitialSpec: &schema.RoleAssignmentSpec{
						Roles: []string{roleName1},
						Subs:  []string{roleAssignmentSub1},
						Scopes: []schema.RoleAssignmentScope{
							{Tenants: &[]string{suite.Tenant}},
						},
					},
				},
				{
					Name: roleAssignmentName2,
					InitialLabels: schema.Labels{
						constants.EnvLabel: constants.EnvConformanceLabel,
					},
					InitialSpec: &schema.RoleAssignmentSpec{
						Roles: []string{roleName2},
						Subs:  []string{roleAssignmentSub1},
						Scopes: []schema.RoleAssignmentScope{
							{Tenants: &[]string{suite.Tenant}},
						},
					},
				},
				{
					Name: roleAssignmentName3,
					InitialLabels: schema.Labels{
						constants.EnvLabel: constants.EnvConformanceLabel,
					},
					InitialSpec: &schema.RoleAssignmentSpec{
						Roles: []string{roleName3},
						Subs:  []string{roleAssignmentSub1},
						Scopes: []schema.RoleAssignmentScope{
							{Tenants: &[]string{suite.Tenant}},
						},
					},
				},
			},
		}
		wm, err := authorization.ConfigureListScenarioV1(suite.ScenarioName, mockParams)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}
		suite.MockClient = wm
	}

	stepsBuilder := steps.NewStepsConfigurator(&suite.TestSuite, t)

	// Role
	roles := []schema.Role{
		{
			Metadata: &schema.GlobalTenantResourceMetadata{
				Tenant: suite.Tenant,
				Name:   roleName1,
			},
			Labels: map[string]string{
				constants.EnvLabel: constants.EnvConformanceLabel,
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
		},
		{
			Metadata: &schema.GlobalTenantResourceMetadata{
				Tenant: suite.Tenant,
				Name:   roleName2,
			},
			Labels: map[string]string{
				constants.EnvLabel: constants.EnvConformanceLabel,
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
		},
		{
			Metadata: &schema.GlobalTenantResourceMetadata{
				Tenant: suite.Tenant,
				Name:   roleName3,
			},
			Labels: map[string]string{
				constants.EnvLabel: constants.EnvConformanceLabel,
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
		},
	}

	// Create roles
	for _, role := range roles {

		role := &schema.Role{
			Metadata: &schema.GlobalTenantResourceMetadata{
				Tenant: suite.Tenant,
				Name:   role.Metadata.Name,
			},
			Labels: map[string]string{
				constants.EnvLabel: constants.EnvConformanceLabel,
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
			Name(role.Metadata.Name).
			Provider(constants.AuthorizationProviderV1).ApiVersion(constants.ApiVersion1).
			Tenant(suite.Tenant).
			Build()
		if err != nil {
			t.Fatalf("Failed to build Metadata: %v", err)
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

		// Create Role
		stepsBuilder.CreateOrUpdateRoleV1Step("Create a role", suite.Client.AuthorizationV1, role,
			steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
				Metadata:      expectRoleMeta,
				Spec:          expectRoleSpec,
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
	roleAssignments := []schema.RoleAssignment{
		{
			Metadata: &schema.GlobalTenantResourceMetadata{
				Tenant: suite.Tenant,
				Name:   roleAssignmentName1,
			},
			Labels: map[string]string{
				constants.EnvLabel: constants.EnvConformanceLabel,
			},
			Spec: schema.RoleAssignmentSpec{
				Roles: []string{roleName1},
				Subs:  []string{roleAssignmentSub1},
				Scopes: []schema.RoleAssignmentScope{
					{Tenants: &[]string{suite.Tenant}},
				},
			},
		},
		{
			Metadata: &schema.GlobalTenantResourceMetadata{
				Tenant: suite.Tenant,
				Name:   roleAssignmentName2,
			},
			Labels: map[string]string{
				constants.EnvLabel: constants.EnvConformanceLabel,
			},
			Spec: schema.RoleAssignmentSpec{
				Roles: []string{roleName2},
				Subs:  []string{roleAssignmentSub1},
				Scopes: []schema.RoleAssignmentScope{
					{Tenants: &[]string{suite.Tenant}},
				},
			},
		},
		{
			Metadata: &schema.GlobalTenantResourceMetadata{
				Tenant: suite.Tenant,
				Name:   roleAssignmentName3,
			},
			Labels: map[string]string{
				constants.EnvLabel: constants.EnvConformanceLabel,
			},
			Spec: schema.RoleAssignmentSpec{
				Roles: []string{roleName3},
				Subs:  []string{roleAssignmentSub1},
				Scopes: []schema.RoleAssignmentScope{
					{Tenants: &[]string{suite.Tenant}},
				},
			},
		},
	}

	// Create role assignments
	for _, roleAssign := range roleAssignments {

		expectRoleAssignMeta, err := builders.NewRoleAssignmentMetadataBuilder().
			Name(roleAssign.Metadata.Name).
			Provider(constants.AuthorizationProviderV1).ApiVersion(constants.ApiVersion1).
			Tenant(suite.Tenant).
			Build()
		if err != nil {
			t.Fatalf("Failed to build Metadata: %v", err)
		}
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

func (suite *AuthorizationV1LifeCycleTestSuite) AfterEach(t provider.T) {
	suite.ResetAllScenarios()
}
