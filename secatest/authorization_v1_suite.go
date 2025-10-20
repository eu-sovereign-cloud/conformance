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
	ctx := context.Background()
	var err error
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

	// Role
	var roleResp *schema.Role
	var expectedRoleMeta *schema.GlobalTenantResourceMetadata
	var expectedRoleSpec *schema.RoleSpec

	t.WithNewStep("Create role", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "CreateOrUpdateRole")

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
		roleResp, err = suite.client.AuthorizationV1.CreateOrUpdateRole(ctx, role)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, roleResp)

		expectedRoleMeta = secalib.NewGlobalTenantResourceMetadata(roleName, secalib.AuthorizationProviderV1, roleResource, secalib.ApiVersion1, secalib.RoleKind, suite.tenant)
		expectedRoleMeta.Verb = http.MethodPut
		suite.verifyGlobalTenantResourceMetadataStep(sCtx, expectedRoleMeta, roleResp.Metadata)

		expectedRoleSpec = &schema.RoleSpec{
			Permissions: []schema.Permission{
				{
					Provider:  secalib.StorageProviderV1,
					Resources: []string{imageResource},
					Verb:      []string{http.MethodGet},
				},
			},
		}
		suite.verifyRoleSpecStep(sCtx, expectedRoleSpec, &roleResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.CreatingStatusState, *roleResp.Status.State)
	})

	t.WithNewStep("Get created role", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetRole")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   roleName,
		}
		roleResp, err = suite.client.AuthorizationV1.GetRole(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, roleResp)

		expectedRoleMeta.Verb = http.MethodGet
		suite.verifyGlobalTenantResourceMetadataStep(sCtx, expectedRoleMeta, roleResp.Metadata)

		suite.verifyRoleSpecStep(sCtx, expectedRoleSpec, &roleResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveStatusState, *roleResp.Status.State)
	})

	t.WithNewStep("Update role", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "CreateOrUpdateRole")

		roleResp, err = suite.client.AuthorizationV1.CreateOrUpdateRole(ctx, roleResp)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, roleResp)

		expectedRoleMeta.Verb = http.MethodPut
		suite.verifyGlobalTenantResourceMetadataStep(sCtx, expectedRoleMeta, roleResp.Metadata)

		expectedRoleSpec.Permissions[0].Verb = []string{http.MethodGet, http.MethodPut}
		suite.verifyRoleSpecStep(sCtx, expectedRoleSpec, &roleResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.UpdatingStatusState, *roleResp.Status.State)
	})

	t.WithNewStep("Get updated role", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetRole")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   roleName,
		}
		roleResp, err = suite.client.AuthorizationV1.GetRole(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, roleResp)

		expectedRoleMeta.Verb = http.MethodGet
		suite.verifyGlobalTenantResourceMetadataStep(sCtx, expectedRoleMeta, roleResp.Metadata)

		suite.verifyRoleSpecStep(sCtx, expectedRoleSpec, &roleResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveStatusState, *roleResp.Status.State)
	})

	// Role assignment
	var assignResp *schema.RoleAssignment
	var expectedAssignMeta *schema.GlobalTenantResourceMetadata
	var expectedAssignSpec *schema.RoleAssignmentSpec

	t.WithNewStep("Create role assignment", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "CreateOrUpdateRoleAssignment")

		assign := &schema.RoleAssignment{
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
		assignResp, err = suite.client.AuthorizationV1.CreateOrUpdateRoleAssignment(ctx, assign)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, assignResp)

		expectedAssignMeta = secalib.NewGlobalTenantResourceMetadata(roleAssignmentName, secalib.AuthorizationProviderV1, roleAssignmentResource, secalib.ApiVersion1, secalib.RoleAssignmentKind, suite.tenant)
		expectedAssignMeta.Verb = http.MethodPut
		suite.verifyGlobalTenantResourceMetadataStep(sCtx, expectedAssignMeta, assignResp.Metadata)

		expectedAssignSpec = &schema.RoleAssignmentSpec{
			Roles:  []string{roleName},
			Subs:   []string{roleAssignmentSub1},
			Scopes: []schema.RoleAssignmentScope{{Tenants: &[]string{suite.tenant}}},
		}
		suite.verifyRoleAssignmentSpecStep(sCtx, expectedAssignSpec, &assignResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.CreatingStatusState, *assignResp.Status.State)
	})

	t.WithNewStep("Get created role assignment", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetRoleAssignment")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   roleAssignmentName,
		}
		assignResp, err = suite.client.AuthorizationV1.GetRoleAssignment(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, assignResp)

		expectedAssignMeta.Verb = http.MethodGet
		suite.verifyGlobalTenantResourceMetadataStep(sCtx, expectedAssignMeta, assignResp.Metadata)

		suite.verifyRoleAssignmentSpecStep(sCtx, expectedAssignSpec, &assignResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveStatusState, *assignResp.Status.State)
	})

	t.WithNewStep("Update role assignment", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "CreateOrUpdateRoleAssignment")

		assignResp, err = suite.client.AuthorizationV1.CreateOrUpdateRoleAssignment(ctx, assignResp)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, assignResp)

		expectedAssignMeta.Verb = http.MethodPut
		suite.verifyGlobalTenantResourceMetadataStep(sCtx, expectedAssignMeta, assignResp.Metadata)

		expectedAssignSpec.Subs = []string{roleAssignmentSub1, roleAssignmentSub2}
		suite.verifyRoleAssignmentSpecStep(sCtx, expectedAssignSpec, &assignResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.UpdatingStatusState, *assignResp.Status.State)
	})

	t.WithNewStep("Get updated role assignment", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetRoleAssignment")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   roleAssignmentName,
		}
		assignResp, err = suite.client.AuthorizationV1.GetRoleAssignment(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, assignResp)

		expectedAssignMeta.Verb = http.MethodGet
		suite.verifyGlobalTenantResourceMetadataStep(sCtx, expectedAssignMeta, assignResp.Metadata)

		suite.verifyRoleAssignmentSpecStep(sCtx, expectedAssignSpec, &assignResp.Spec)

		suite.verifyStatusStep(sCtx, secalib.ActiveStatusState, *assignResp.Status.State)
	})

	t.WithNewStep("Delete role assignment", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "DeleteRoleAssignment")

		err = suite.client.AuthorizationV1.DeleteRoleAssignment(ctx, assignResp)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get deleted role assignment", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetRoleAssignment")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   roleAssignmentName,
		}
		_, err = suite.client.AuthorizationV1.GetRoleAssignment(ctx, tref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	t.WithNewStep("Delete role", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "DeleteRole")

		err = suite.client.AuthorizationV1.DeleteRole(ctx, roleResp)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get deleted role", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetRole")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   roleAssignmentName,
		}
		_, err = suite.client.AuthorizationV1.GetRole(ctx, tref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	slog.Info("Finishing " + suite.scenarioName)
}

func (suite *AuthorizationV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}
