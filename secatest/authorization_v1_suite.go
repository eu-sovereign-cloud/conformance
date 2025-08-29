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

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

type AuthorizationV1TestSuite struct {
	globalTestSuite
}

func (suite *AuthorizationV1TestSuite) TestAuthorizationV1(t provider.T) {
	slog.Info("Starting Authorization Lifecycle Test")

	t.Title("Authorization Lifecycle Test")
	configureTags(t, secalib.AuthorizationProviderV1, secalib.RoleKind, secalib.RoleAssignmentKind)

	// TODO Export to configuration
	roleAssignmentSubs := []string{"user1@secalib.com", "user2@secalib.com"}

	// Select subs
	roleAssignmentSub1 := roleAssignmentSubs[rand.Intn(len(roleAssignmentSubs))]
	roleAssignmentSub2 := roleAssignmentSubs[rand.Intn(len(roleAssignmentSubs))]

	// Generate scenario data
	roleName := secalib.GenerateRoleName()
	roleResource := secalib.GenerateRoleResource(suite.tenant, roleName)

	roleAssignmentName := secalib.GenerateRoleAssignmentName()
	roleAssignmentResource := secalib.GenerateRoleAssignmentResource(suite.tenant, roleAssignmentName)

	imageName := secalib.GenerateImageName()
	imageResource := secalib.GenerateImageResource(suite.tenant, imageName)

	// Setup mock, if configured to use
	if suite.isMockEnabled() {
		wm, err := mock.CreateAuthorizationLifecycleScenarioV1("Authorization Lifecycle",
			mock.AuthorizationParamsV1{
				Params: &mock.Params{
					MockURL:   suite.mockServerURL,
					AuthToken: suite.authToken,
					Tenant:    suite.tenant,
				},
				Role: &mock.ResourceParams[secalib.RoleSpecV1]{
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
				RoleAssignment: &mock.ResourceParams[secalib.RoleAssignmentSpecV1]{
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
			})
		if err != nil {
			slog.Error("Failed to create storage scenario", "error", err)
			return
		}
		suite.mockClient = wm
	}

	ctx := context.Background()
	var roleResp *authorization.Role
	var assignResp *authorization.RoleAssignment
	var err error

	var expectedRoleMetadata *secalib.Metadata
	var expectedRoleSpec *secalib.RoleSpecV1

	// Step 1
	t.WithNewStep("Create role", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateRole",
			tenantStepParameter, suite.tenant,
		)

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   roleName,
		}
		role := &authorization.Role{
			Spec: authorization.RoleSpec{
				Permissions: []authorization.Permission{
					{
						Provider:  secalib.StorageProviderV1,
						Resources: []string{imageResource},
						Verb:      []string{http.MethodGet},
					},
				},
			},
		}
		roleResp, err = suite.client.AuthorizationV1.CreateOrUpdateRole(ctx, tref, role, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, roleResp)

		expectedRoleMetadata = &secalib.Metadata{
			Name:       roleName,
			Provider:   secalib.AuthorizationProviderV1,
			Resource:   roleResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.RoleKind,
			Tenant:     suite.tenant,
		}
		verifyAuthorizationMetadataStep(sCtx, expectedRoleMetadata, roleResp.Metadata)

		expectedRoleSpec = &secalib.RoleSpecV1{
			Permissions: []*secalib.RoleSpecPermissionV1{
				{
					Provider:  secalib.StorageProviderV1,
					Resources: []string{imageResource},
					Verb:      []string{http.MethodGet},
				},
			},
		}
		verifyRoleSpecStep(sCtx, expectedRoleSpec, &roleResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*roleResp.Status.State)},
		)
	})

	// Step 2
	t.WithNewStep("Get created role", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetRole",
			tenantStepParameter, suite.tenant,
		)

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   roleName,
		}
		roleResp, err = suite.client.AuthorizationV1.GetRole(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, roleResp)

		expectedRoleMetadata.Verb = http.MethodGet
		verifyAuthorizationMetadataStep(sCtx, expectedRoleMetadata, roleResp.Metadata)

		verifyRoleSpecStep(sCtx, expectedRoleSpec, &roleResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*roleResp.Status.State)},
		)
	})

	// Step 3
	t.WithNewStep("Update role", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateRole",
			tenantStepParameter, suite.tenant,
		)

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   roleName,
		}
		roleResp, err = suite.client.AuthorizationV1.CreateOrUpdateRole(ctx, tref, roleResp, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, roleResp)

		expectedRoleMetadata.Verb = http.MethodPut
		verifyAuthorizationMetadataStep(sCtx, expectedRoleMetadata, roleResp.Metadata)

		expectedRoleSpec.Permissions[0].Verb = []string{http.MethodGet, http.MethodPut}
		verifyRoleSpecStep(sCtx, expectedRoleSpec, &roleResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*roleResp.Status.State)},
		)
	})

	// Step 4
	t.WithNewStep("Get updated role", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetRole",
			tenantStepParameter, suite.tenant,
		)

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   roleName,
		}
		roleResp, err = suite.client.AuthorizationV1.GetRole(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, roleResp)

		expectedRoleMetadata.Verb = http.MethodGet
		verifyAuthorizationMetadataStep(sCtx, expectedRoleMetadata, roleResp.Metadata)

		verifyRoleSpecStep(sCtx, expectedRoleSpec, &roleResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*roleResp.Status.State)},
		)
	})

	var expectedRoleAssignmentMetadata *secalib.Metadata
	var expectedRoleAssignmentSpec *secalib.RoleAssignmentSpecV1

	// Step 5
	t.WithNewStep("Create role assignment", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateRoleAssignment",
			tenantStepParameter, suite.tenant,
		)

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   roleAssignmentName,
		}
		assign := &authorization.RoleAssignment{
			Spec: authorization.RoleAssignmentSpec{
				Roles:  []string{roleName},
				Subs:   []string{roleAssignmentSub1},
				Scopes: []authorization.RoleAssignmentScope{{Tenants: &[]string{suite.tenant}}},
			},
		}
		assignResp, err = suite.client.AuthorizationV1.CreateOrUpdateRoleAssignment(ctx, tref, assign, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, assignResp)

		expectedRoleAssignmentMetadata = &secalib.Metadata{
			Name:       roleAssignmentName,
			Provider:   secalib.AuthorizationProviderV1,
			Resource:   roleAssignmentResource,
			Verb:       http.MethodPut,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.RoleAssignmentKind,
			Tenant:     suite.tenant,
		}
		verifyAuthorizationMetadataStep(sCtx, expectedRoleAssignmentMetadata, assignResp.Metadata)

		expectedRoleAssignmentSpec = &secalib.RoleAssignmentSpecV1{
			Roles:  []string{roleName},
			Subs:   []string{roleAssignmentSub1},
			Scopes: []*secalib.RoleAssignmentSpecScopeV1{{Tenants: []string{suite.tenant}}},
		}
		verifyRoleAssignmentSpecStep(sCtx, expectedRoleAssignmentSpec, &assignResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.CreatingStatusState},
			&secalib.Status{State: string(*assignResp.Status.State)},
		)
	})

	// Step 6
	t.WithNewStep("Get created role assignment", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetRoleAssignment",
			tenantStepParameter, suite.tenant,
		)

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   roleAssignmentName,
		}
		assignResp, err = suite.client.AuthorizationV1.GetRoleAssignment(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, assignResp)

		expectedRoleAssignmentMetadata.Verb = http.MethodGet
		verifyAuthorizationMetadataStep(sCtx, expectedRoleAssignmentMetadata, assignResp.Metadata)

		verifyRoleAssignmentSpecStep(sCtx, expectedRoleAssignmentSpec, &assignResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*assignResp.Status.State)},
		)
	})

	// Step 7
	t.WithNewStep("Update role assignment", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateRoleAssignment",
			tenantStepParameter, suite.tenant,
		)

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   roleAssignmentName,
		}
		assignResp, err = suite.client.AuthorizationV1.CreateOrUpdateRoleAssignment(ctx, tref, assignResp, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, assignResp)

		expectedRoleAssignmentMetadata.Verb = http.MethodPut
		verifyAuthorizationMetadataStep(sCtx, expectedRoleAssignmentMetadata, assignResp.Metadata)

		expectedRoleAssignmentSpec.Subs = []string{roleAssignmentSub1, roleAssignmentSub2}
		verifyRoleAssignmentSpecStep(sCtx, expectedRoleAssignmentSpec, &assignResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.UpdatingStatusState},
			&secalib.Status{State: string(*assignResp.Status.State)},
		)
	})

	// Step 8
	t.WithNewStep("Get updated role assignment", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetRoleAssignment",
			tenantStepParameter, suite.tenant,
		)

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   roleAssignmentName,
		}
		assignResp, err = suite.client.AuthorizationV1.GetRoleAssignment(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, assignResp)

		expectedRoleAssignmentMetadata.Verb = http.MethodGet
		verifyAuthorizationMetadataStep(sCtx, expectedRoleAssignmentMetadata, assignResp.Metadata)

		verifyRoleAssignmentSpecStep(sCtx, expectedRoleAssignmentSpec, &assignResp.Spec)

		verifyStatusStep(sCtx,
			&secalib.Status{State: secalib.ActiveStatusState},
			&secalib.Status{State: string(*assignResp.Status.State)},
		)
	})

	// Step 9
	t.WithNewStep("Delete role assignment", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "DeleteRoleAssignment",
			tenantStepParameter, suite.tenant,
		)

		err = suite.client.AuthorizationV1.DeleteRoleAssignment(ctx, assignResp, nil)
		requireNoError(sCtx, err)
	})

	// Step 10
	t.WithNewStep("Get deleted role assignment", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetRoleAssignment",
			tenantStepParameter, suite.tenant,
		)

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   roleAssignmentName,
		}
		_, err = suite.client.AuthorizationV1.GetRoleAssignment(ctx, tref)
		requireError(sCtx, err, secapi.ErrResourceNotFound)
	})

	// Step 11
	t.WithNewStep("Delete role", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "DeleteRole",
			tenantStepParameter, suite.tenant,
		)

		err = suite.client.AuthorizationV1.DeleteRole(ctx, roleResp, nil)
		requireNoError(sCtx, err)
	})

	// Step 12
	t.WithNewStep("Get deleted role", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "GetRole",
			tenantStepParameter, suite.tenant,
		)

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   roleAssignmentName,
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
		stepCtx.WithNewParameters(
			"expected_permissions_length", len(expected.Permissions),
			"actual_permissions_length", len(actual.Permissions),
		)
		stepCtx.Require().Equal(len(expected.Permissions), len(actual.Permissions),
			"Permissions list length should match expected")

		for i := 0; i < len(expected.Permissions); i++ {
			expectedPerm := expected.Permissions[i]
			actualPerm := actual.Permissions[i]
			stepCtx.WithNewParameters(
				fmt.Sprintf("expected_permission[%d]_provider", i), expectedPerm.Provider,
				fmt.Sprintf("actual_permission[%d]_provider", i), actualPerm.Provider,

				fmt.Sprintf("expected_permission[%d]_resources", i), expectedPerm.Resources,
				fmt.Sprintf("actual_permission[%d]_resources", i), actualPerm.Resources,

				fmt.Sprintf("expected_permission[%d]_verb", i), expectedPerm.Verb,
				fmt.Sprintf("actual_permission[%d]_verb", i), actualPerm.Verb,
			)
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
		stepCtx.WithNewParameters(
			"expected_roles", expected.Roles,
			"actual_roles", actual.Roles,

			"expected_subs", expected.Subs,
			"actual_subs", actual.Subs,

			"expected_scopes_length", len(expected.Scopes),
			"actual_scopes_length", len(actual.Scopes),
		)
		stepCtx.Require().Equal(expected.Roles, actual.Roles,
			"Roles provider should match expected")
		stepCtx.Require().Equal(expected.Subs, actual.Subs,
			"Subs should match expected")
		stepCtx.Require().Equal(len(expected.Scopes), len(actual.Scopes),
			"Scopes list length should match expected")

		for i := 0; i < len(expected.Scopes); i++ {
			expectedScope := expected.Scopes[i]
			actualScope := actual.Scopes[i]
			stepCtx.WithNewParameters(
				fmt.Sprintf("expected_scope[%d]_tenants", i), expectedScope.Tenants,
				fmt.Sprintf("actual_scope[%d]_tenants", i), actualScope.Tenants,

				fmt.Sprintf("expected_scope[%d]_regions", i), expectedScope.Regions,
				fmt.Sprintf("actual_scope[%d]_regions", i), actualScope.Regions,

				fmt.Sprintf("expected_scope[%d]_workspaces", i), expectedScope.Workspaces,
				fmt.Sprintf("actual_scope[%d]_workspaces", i), actualScope.Workspaces,
			)
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
