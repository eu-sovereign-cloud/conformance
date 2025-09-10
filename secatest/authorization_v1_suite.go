package secatest

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/secalib"
	authorization "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.authorization.v1"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/google/uuid"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type AuthorizationV1TestSuite struct {
	globalTestSuite

	users []string
}

func (suite *AuthorizationV1TestSuite) generateLifecycleParams() *secalib.AuthorizationLifeCycleParamsV1 {
	roleName := secalib.GenerateRoleName()
	roleAssignmentName := secalib.GenerateRoleAssignmentName()
	imageName := secalib.GenerateImageName()

	imageResource := secalib.GenerateImageResource(suite.tenant, imageName)

	// Random data
	roleAssignmentSub1 := suite.users[rand.Intn(len(suite.users))]
	roleAssignmentSub2 := suite.users[rand.Intn(len(suite.users))]

	return &secalib.AuthorizationLifeCycleParamsV1{
		Role: &secalib.ResourceParams[secalib.RoleSpecV1]{
			Name: roleName,
			InitialSpec: &secalib.RoleSpecV1{
				Permissions: []*secalib.RoleSpecPermissionV1{
					{
						Provider:  secalib.StorageProviderV1,
						Resources: []string{imageResource},
						Verb:      []string{http.MethodGet},
					},
				},
			},
			UpdatedSpec: &secalib.RoleSpecV1{
				Permissions: []*secalib.RoleSpecPermissionV1{
					{
						Provider:  secalib.StorageProviderV1,
						Resources: []string{imageResource},
						Verb:      []string{http.MethodGet, http.MethodPut},
					},
				},
			},
		},
		RoleAssignment: &secalib.ResourceParams[secalib.RoleAssignmentSpecV1]{
			Name: roleAssignmentName,
			InitialSpec: &secalib.RoleAssignmentSpecV1{
				Roles:  []string{roleName},
				Subs:   []string{roleAssignmentSub1},
				Scopes: []*secalib.RoleAssignmentSpecScopeV1{{Tenants: []string{suite.tenant}}},
			},
			UpdatedSpec: &secalib.RoleAssignmentSpecV1{
				Roles:  []string{roleName},
				Subs:   []string{roleAssignmentSub1, roleAssignmentSub2},
				Scopes: []*secalib.RoleAssignmentSpecScopeV1{{Tenants: []string{suite.tenant}}},
			},
		},
	}
}

func (suite *AuthorizationV1TestSuite) TestAuthorizationLifeCycleV1(t provider.T) {
	slog.Info("Starting Authorization Lifecycle Test")

	t.Title("Authorization Lifecycle Test")
	configureTags(t, secalib.AuthorizationProviderV1, secalib.RoleKind, secalib.RoleAssignmentKind)

	params := suite.generateLifecycleParams()
	roleResource := secalib.GenerateRoleResource(suite.tenant, params.Role.Name)
	roleAssignmentResource := secalib.GenerateRoleAssignmentResource(suite.tenant, params.RoleAssignment.Name)

	// Setup mock, if configured to use
	if suite.isMockEnabled() {
		scenarios := mock.NewAuthorizationScenariosV1(suite.authToken, suite.tenant, suite.mockServerURL)

		id := uuid.New().String()
		wm, err := scenarios.ConfigureLifecycleScenario(id, params)
		if err != nil {
			t.Fatalf("Failed to configure mock scenario: %v", err)
		}
		suite.mockClient = wm
	}

	ctx := context.Background()
	var err error

	// Role
	var roleResp *authorization.Role
	var expectedRoleMeta *secalib.Metadata

	t.WithNewStep("Create role", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "CreateOrUpdateRole")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   params.Role.Name,
		}
		role := &authorization.Role{
			Spec: authorization.RoleSpec{
				Permissions: []authorization.Permission{},
			},
		}
		for _, perm := range params.Role.InitialSpec.Permissions {
			role.Spec.Permissions = append(role.Spec.Permissions, authorization.Permission{
				Provider:  perm.Provider,
				Resources: perm.Resources,
				Verb:      perm.Verb,
			})
		}
		roleResp, err = suite.client.AuthorizationV1.CreateOrUpdateRole(ctx, tref, role, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, roleResp)

		expectedRoleMeta = &secalib.Metadata{
			Name:       params.Role.Name,
			Provider:   secalib.AuthorizationProviderV1,
			Resource:   roleResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.RoleKind,
			Tenant:     suite.tenant,
		}
		verifyAuthorizationMetadataStep(sCtx, expectedRoleMeta, roleResp.Metadata)

		verifyRoleSpecStep(sCtx, params.Role.InitialSpec, &roleResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*roleResp.Status.State)},
		)
	})

	t.WithNewStep("Get created role", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetRole")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   params.Role.Name,
		}
		roleResp, err = suite.client.AuthorizationV1.GetRole(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, roleResp)

		expectedRoleMeta.Verb = http.MethodGet
		verifyAuthorizationMetadataStep(sCtx, expectedRoleMeta, roleResp.Metadata)

		verifyRoleSpecStep(sCtx, params.Role.InitialSpec, &roleResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*roleResp.Status.State)},
		)
	})

	t.WithNewStep("Update role", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "CreateOrUpdateRole")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   params.Role.Name,
		}
		roleResp, err = suite.client.AuthorizationV1.CreateOrUpdateRole(ctx, tref, roleResp, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, roleResp)

		expectedRoleMeta.Verb = http.MethodPut
		verifyAuthorizationMetadataStep(sCtx, expectedRoleMeta, roleResp.Metadata)

		verifyRoleSpecStep(sCtx, params.Role.UpdatedSpec, &roleResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*roleResp.Status.State)},
		)
	})

	t.WithNewStep("Get updated role", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetRole")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   params.Role.Name,
		}
		roleResp, err = suite.client.AuthorizationV1.GetRole(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, roleResp)

		expectedRoleMeta.Verb = http.MethodGet
		verifyAuthorizationMetadataStep(sCtx, expectedRoleMeta, roleResp.Metadata)

		verifyRoleSpecStep(sCtx, params.Role.UpdatedSpec, &roleResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*roleResp.Status.State)},
		)
	})

	// Role assignment
	var assignResp *authorization.RoleAssignment
	var expectedAssignMeta *secalib.Metadata

	t.WithNewStep("Create role assignment", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "CreateOrUpdateRoleAssignment")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   params.RoleAssignment.Name,
		}
		assign := &authorization.RoleAssignment{
			Spec: authorization.RoleAssignmentSpec{
				Roles:  params.RoleAssignment.InitialSpec.Roles,
				Subs:   params.RoleAssignment.InitialSpec.Subs,
			},
		}
		for _, scope := range params.RoleAssignment.InitialSpec.Scopes {
			assign.Spec.Scopes = append(assign.Spec.Scopes, authorization.RoleAssignmentScope{
				Tenants:    &scope.Tenants,
				Workspaces: &scope.Workspaces,
				Regions:    &scope.Regions,
			})
		}
		assignResp, err = suite.client.AuthorizationV1.CreateOrUpdateRoleAssignment(ctx, tref, assign, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, assignResp)

		expectedAssignMeta = &secalib.Metadata{
			Name:       params.RoleAssignment.Name,
			Provider:   secalib.AuthorizationProviderV1,
			Resource:   roleAssignmentResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.RoleAssignmentKind,
			Tenant:     suite.tenant,
		}
		verifyAuthorizationMetadataStep(sCtx, expectedAssignMeta, assignResp.Metadata)

		verifyRoleAssignmentSpecStep(sCtx, params.RoleAssignment.InitialSpec, &assignResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*assignResp.Status.State)},
		)
	})

	t.WithNewStep("Get created role assignment", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetRoleAssignment")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   params.RoleAssignment.Name,
		}
		assignResp, err = suite.client.AuthorizationV1.GetRoleAssignment(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, assignResp)

		expectedAssignMeta.Verb = http.MethodGet
		verifyAuthorizationMetadataStep(sCtx, expectedAssignMeta, assignResp.Metadata)

		verifyRoleAssignmentSpecStep(sCtx, params.RoleAssignment.InitialSpec, &assignResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*assignResp.Status.State)},
		)
	})

	t.WithNewStep("Update role assignment", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "CreateOrUpdateRoleAssignment")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   params.RoleAssignment.Name,
		}
		assignResp, err = suite.client.AuthorizationV1.CreateOrUpdateRoleAssignment(ctx, tref, assignResp, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, assignResp)

		expectedAssignMeta.Verb = http.MethodPut
		verifyAuthorizationMetadataStep(sCtx, expectedAssignMeta, assignResp.Metadata)

		verifyRoleAssignmentSpecStep(sCtx, params.RoleAssignment.UpdatedSpec, &assignResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*assignResp.Status.State)},
		)
	})

	t.WithNewStep("Get updated role assignment", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetRoleAssignment")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   params.RoleAssignment.Name,
		}
		assignResp, err = suite.client.AuthorizationV1.GetRoleAssignment(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, assignResp)

		expectedAssignMeta.Verb = http.MethodGet
		verifyAuthorizationMetadataStep(sCtx, expectedAssignMeta, assignResp.Metadata)

		verifyRoleAssignmentSpecStep(sCtx, params.RoleAssignment.UpdatedSpec, &assignResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*assignResp.Status.State)},
		)
	})

	t.WithNewStep("Delete role assignment", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "DeleteRoleAssignment")

		err = suite.client.AuthorizationV1.DeleteRoleAssignment(ctx, assignResp, nil)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get deleted role assignment", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetRoleAssignment")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   params.RoleAssignment.Name,
		}
		_, err = suite.client.AuthorizationV1.GetRoleAssignment(ctx, tref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	t.WithNewStep("Delete role", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "DeleteRole")

		err = suite.client.AuthorizationV1.DeleteRole(ctx, roleResp, nil)
		requireNoError(sCtx, err)
	})

	t.WithNewStep("Get deleted role", func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetRole")

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   params.RoleAssignment.Name,
		}
		_, err = suite.client.AuthorizationV1.GetRole(ctx, tref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	slog.Info("Finishing Authorization Lifecycle Test")
}

func (suite *AuthorizationV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}

// TODO Create a helper to perform this copy using reflection
func verifyAuthorizationMetadataStep(ctx provider.StepCtx, expected *secalib.Metadata, actual *authorization.GlobalResourceMetadata) {
	actualMetadata := &secalib.Metadata{
		Name:       actual.Name,
		Provider:   actual.Provider,
		Verb:       actual.Verb,
		Resource:   actual.Resource,
		ApiVersion: actual.ApiVersion,
		Kind:       string(actual.Kind),
		Tenant:     actual.Tenant,
	}
	verifyGlobalMetadataStep(ctx, expected, actualMetadata)
}

// TODO Create a helper to perform these asserts using reflection
func verifyRoleSpecStep(ctx provider.StepCtx, expected *secalib.RoleSpecV1, actual *authorization.RoleSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(len(expected.Permissions), len(actual.Permissions),
			"Permissions list length should match expected")

		for i := 0; i < len(expected.Permissions); i++ {
			expectedPerm := expected.Permissions[i]
			actualPerm := actual.Permissions[i]

			stepCtx.Require().Equal(expectedPerm.Provider, actualPerm.Provider,
				fmt.Sprintf("Permission [%d] provider should match expected", i))

			stepCtx.Require().Equal(expectedPerm.Resources, actualPerm.Resources,
				fmt.Sprintf("Permission [%d] resources should match expected", i))

			stepCtx.Require().Equal(expectedPerm.Verb, actualPerm.Verb,
				fmt.Sprintf("Permission [%d] verb should match expected", i))
		}
	})
}

func verifyRoleAssignmentSpecStep(ctx provider.StepCtx, expected *secalib.RoleAssignmentSpecV1, actual *authorization.RoleAssignmentSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {
		stepCtx.Require().Equal(expected.Roles, actual.Roles,
			"Roles provider should match expected")
		stepCtx.Require().Equal(expected.Subs, actual.Subs,
			"Subs should match expected")
		stepCtx.Require().Equal(len(expected.Scopes), len(actual.Scopes),
			"Scope list length should match expected")

		for i := 0; i < len(expected.Scopes); i++ {
			expectedScope := expected.Scopes[i]
			actualScope := actual.Scopes[i]

			if len(*actualScope.Tenants) > 0 {
				stepCtx.Require().Equal(expectedScope.Tenants, *actualScope.Tenants,
					fmt.Sprintf("Scope [%d] tenants should match expected", i))
			}
			if len(*actualScope.Regions) > 0 {
				stepCtx.Require().Equal(expectedScope.Regions, *actualScope.Regions,
					fmt.Sprintf("Scope [%d] regions should match expected", i))
			}
			if len(*actualScope.Workspaces) > 0 {
				stepCtx.Require().Equal(expectedScope.Workspaces, *actualScope.Workspaces,
					fmt.Sprintf("Scope [%d] workspaces should match expected", i))
			}
		}
	})
}
