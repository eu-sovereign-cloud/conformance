package authorization

import (
	"log/slog"
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

type RoleLifeCycleV1TestSuite struct {
	suites.GlobalTestSuite

	params *params.RoleLifeCycleV1Params
}

func CreateRoleLifeCycleV1TestSuite(globalTestSuite suites.GlobalTestSuite) *RoleLifeCycleV1TestSuite {
	suite := &RoleLifeCycleV1TestSuite{
		GlobalTestSuite: globalTestSuite,
	}
	suite.ScenarioName = constants.RoleLifeCycleV1SuiteName.String()
	return suite
}

func (suite *RoleLifeCycleV1TestSuite) BeforeAll(t provider.T) {
	t.AddParentSuite("Authorization")

	// Generate scenario data
	roleName := generators.GenerateRoleName()

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

	params := &params.RoleLifeCycleV1Params{
		RoleInitial: roleInitial,
		RoleUpdated: roleUpdated,
	}
	suite.params = params

	err = suites.SetupMockIfEnabled(suite.TestSuite, mockauthorization.ConfigureRoleLifecycleScenarioV1, params)
	if err != nil {
		slog.Error("Failed to setup mock", "error", err)
		t.FailNow()
	}
}

func (suite *RoleLifeCycleV1TestSuite) TestScenario(t provider.T) {
	suite.StartScenario(t)
	suite.ConfigureTags(t, constants.AuthorizationProviderV1,
		string(schema.GlobalTenantResourceMetadataKindResourceKindRole),
		string(schema.GlobalTenantResourceMetadataKindResourceKindRoleAssignment),
	)

	stepsBuilder := steps.NewStepsConfigurator(suite.TestSuite, t)

	// Role

	// Create a role
	role := suite.params.RoleInitial
	expectRoleMeta := role.Metadata
	expectRoleSpec := &role.Spec
	stepsBuilder.CreateOrUpdateRoleV1Step("Create a role", suite.Client.AuthorizationV1, role,
		steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
			Metadata:      expectRoleMeta,
			Spec:          expectRoleSpec,
			ResourceState: schema.ResourceStateCreating,
		},
	)

	// Get the created role
	roleTRef := secapi.TenantReference{
		Tenant: secapi.TenantID(suite.Tenant),
		Name:   role.Metadata.Name,
	}
	role = stepsBuilder.GetRoleV1Step("Get the created role", suite.Client.AuthorizationV1, roleTRef,
		steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
			Metadata:      expectRoleMeta,
			Spec:          expectRoleSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Update the role
	role.Spec = suite.params.RoleUpdated.Spec
	expectRoleSpec = &role.Spec
	stepsBuilder.CreateOrUpdateRoleV1Step("Update the role", suite.Client.AuthorizationV1, role,
		steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
			Metadata:      expectRoleMeta,
			Spec:          expectRoleSpec,
			ResourceState: schema.ResourceStateUpdating,
		},
	)

	// Get the updated role
	role = stepsBuilder.GetRoleV1Step("Get the updated role", suite.Client.AuthorizationV1, roleTRef,
		steps.ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
			Metadata:      expectRoleMeta,
			Spec:          expectRoleSpec,
			ResourceState: schema.ResourceStateActive,
		},
	)

	// Resources deletion
	stepsBuilder.DeleteRoleV1Step("Delete the role", suite.Client.AuthorizationV1, role)
	stepsBuilder.GetRoleWithErrorV1Step("Get the deleted role", suite.Client.AuthorizationV1, roleTRef, secapi.ErrResourceNotFound)

	suite.FinishScenario()
}

func (suite *RoleLifeCycleV1TestSuite) AfterAll(t provider.T) {
	suite.ResetAllScenarios()
}
