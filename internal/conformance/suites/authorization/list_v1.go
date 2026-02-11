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

func CreateListV1TestSuite(globalTestSuite suites.GlobalTestSuite, users []string) *AuthorizationListV1TestSuite {
	suite := &AuthorizationListV1TestSuite{
		GlobalTestSuite: globalTestSuite,
		Users:           users,
	}
	suite.ScenarioName = constants.AuthorizationV1ListSuiteName
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

//nolint:dupl
func (suite *AuthorizationListV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, constants.AuthorizationProviderV1,
		string(schema.GlobalTenantResourceMetadataKindResourceKindRole),
		string(schema.GlobalTenantResourceMetadataKindResourceKindRoleAssignment),
	)

	roles := suite.params.Roles
	roleAssignments := suite.params.RoleAssignments

	t.WithNewStep("Role", func(roleCtx provider.StepCtx) {
		roleSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, roleCtx)

		for _, role := range roles {
			r := role
			expectRoleMeta := r.Metadata
			expectRoleSpec := r.Spec

			roleSteps.CreateOrUpdateRoleV1Step("Create", suite.Client.AuthorizationV1, &r,
				steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
					Metadata:      expectRoleMeta,
					Spec:          &expectRoleSpec,
					ResourceState: schema.ResourceStateCreating,
				},
			)
		}

		roleTRef := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.Tenant),
			Name:   suite.Tenant,
		}
		roleSteps.GetListRoleV1Step("ListAll", suite.Client.AuthorizationV1, roleTRef, nil)
		roleSteps.GetListRoleV1Step("ListWithLimit", suite.Client.AuthorizationV1, roleTRef,
			secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().Equals(constants.EnvLabel, constants.EnvConformanceLabel)))
		roleSteps.GetListRoleV1Step("ListWithLabel", suite.Client.AuthorizationV1, roleTRef,
			secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
				Equals(constants.EnvLabel, constants.EnvConformanceLabel)))
		roleSteps.GetListRoleV1Step("ListWithLimitAndLabel", suite.Client.AuthorizationV1, roleTRef,
			secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
				Equals(constants.EnvLabel, constants.EnvConformanceLabel)))
	})

	t.WithNewStep("RoleAssignment", func(raCtx provider.StepCtx) {
		raSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, raCtx)

		for _, roleAssign := range roleAssignments {
			ra := roleAssign
			expectRoleAssignMeta := ra.Metadata
			expectRoleAssignSpec := &ra.Spec

			raSteps.CreateOrUpdateRoleAssignmentV1Step("Create", suite.Client.AuthorizationV1, &ra,
				steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
					Metadata:      expectRoleAssignMeta,
					Spec:          expectRoleAssignSpec,
					ResourceState: schema.ResourceStateCreating,
				},
			)
		}

		roleAssignTRef := secapi.TenantReference{Tenant: secapi.TenantID(suite.Tenant)}
		raSteps.GetListRoleAssignmentsV1("ListAll", suite.Client.AuthorizationV1, roleAssignTRef, nil)
		raSteps.GetListRoleAssignmentsV1("ListWithLimit", suite.Client.AuthorizationV1, roleAssignTRef,
			secapi.NewListOptions().WithLimit(1))
		raSteps.GetListRoleAssignmentsV1("ListWithLabel", suite.Client.AuthorizationV1, roleAssignTRef,
			secapi.NewListOptions().WithLabels(labelBuilder.NewLabelsBuilder().
				Equals(constants.EnvLabel, constants.EnvConformanceLabel)))
		raSteps.GetListRoleAssignmentsV1("ListWithLimitAndLabel", suite.Client.AuthorizationV1, roleAssignTRef,
			secapi.NewListOptions().WithLimit(1).WithLabels(labelBuilder.NewLabelsBuilder().
				Equals(constants.EnvLabel, constants.EnvConformanceLabel)))
	})

	t.WithNewStep("Delete", func(delCtx provider.StepCtx) {
		delSteps := steps.NewStepsConfiguratorWithCtx(suite.TestSuite, t, delCtx)

		for _, roleAssign := range roleAssignments {
			ra := roleAssign
			delSteps.DeleteRoleAssignmentV1Step("RoleAssignment", suite.Client.AuthorizationV1, &ra)

			roleAssignTRefSingle := secapi.TenantReference{
				Tenant: secapi.TenantID(suite.Tenant),
				Name:   ra.Metadata.Name,
			}
			delSteps.GetRoleAssignmentWithErrorV1Step("GetDeletedRoleAssignment", suite.Client.AuthorizationV1, roleAssignTRefSingle, secapi.ErrResourceNotFound)
		}

		for _, role := range roles {
			r := role
			delSteps.DeleteRoleV1Step("Role", suite.Client.AuthorizationV1, &r)

			roleTRefSingle := secapi.TenantReference{
				Tenant: secapi.TenantID(suite.Tenant),
				Name:   r.Metadata.Name,
			}
			delSteps.GetRoleWithErrorV1Step("GetDeletedRole", suite.Client.AuthorizationV1, roleTRefSingle, secapi.ErrResourceNotFound)
		}
	})
	suite.FinishScenario()
}

func (suite *AuthorizationListV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
