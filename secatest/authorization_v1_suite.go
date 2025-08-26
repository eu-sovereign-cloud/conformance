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

type verifyRoleSpecStepParams struct {
	permissions []verifyRolePermissionStepParams
}
type verifyRolePermissionStepParams struct {
	provider  string
	resources []string
	verb      []string
}

type verifyRoleAssignmentSpecStepParams struct {
	roles  []string
	subs   []string
	scopes []verifyRoleAssignmentScopeStepParams
}
type verifyRoleAssignmentScopeStepParams struct {
	tenants    []string
	regions    []string
	workspaces []string
}

func (suite *AuthorizationV1TestSuite) TestAuthorizationV1(t provider.T) {
	t.Title("Authorization Lifecycle Test")
	configureTags(t, secalib.AuthorizationProviderV1, secalib.RoleKind, secalib.RoleAssignmentKind)

	// TODO Export configuration
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
				Params: mock.Params{
					MockURL:   suite.mockServerURL,
					AuthToken: suite.authToken,
					Tenant:    suite.tenant,
				},
				Role: mock.RoleParamsV1{
					Name: roleName,
					Permissions: []mock.RolePermissionParamsV1{
						{
							Provider:    secalib.StorageProviderV1,
							Resources:   []string{imageResource},
							VerbInitial: []string{http.MethodGet},
							VerbUpdated: []string{http.MethodGet, http.MethodPut},
						},
					},
				},
				RoleAssignment: mock.RoleAssignmentParamsV1{
					Name:        roleAssignmentName,
					Roles:       []string{roleName},
					SubsInitial: []string{roleAssignmentSub1},
					SubsUpdated: []string{roleAssignmentSub1, roleAssignmentSub2},
					Scopes:      []mock.RoleAssignmentScopeParamsV1{{Tenants: []string{suite.tenant}}},
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

	var expectedRoleMetadata verifyGlobalMetadataStepParams
	var expectedRoleSpec verifyRoleSpecStepParams

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

		expectedRoleMetadata = verifyGlobalMetadataStepParams{
			name:       roleName,
			provider:   secalib.AuthorizationProviderV1,
			resource:   roleResource,
			verb:       http.MethodPut,
			apiVersion: secalib.ApiVersion1,
			kind:       secalib.RoleKind,
			tenant:     suite.tenant,
		}
		verifyAuthorizationMetadataStep(sCtx, expectedRoleMetadata, roleResp.Metadata)

		expectedRoleSpec = verifyRoleSpecStepParams{
			permissions: []verifyRolePermissionStepParams{
				{
					provider:  secalib.StorageProviderV1,
					resources: []string{imageResource},
					verb:      []string{http.MethodGet},
				},
			},
		}
		verifyRoleSpecStep(sCtx, expectedRoleSpec, &roleResp.Spec)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: secalib.CreatingStatusState,
			actualState:   string(*roleResp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
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

		expectedRoleMetadata.verb = http.MethodGet
		verifyAuthorizationMetadataStep(sCtx, expectedRoleMetadata, roleResp.Metadata)

		verifyRoleSpecStep(sCtx, expectedRoleSpec, &roleResp.Spec)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: secalib.ActiveStatusState,
			actualState:   string(*roleResp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
	})

	// Step 3
	t.WithNewStep("Update role", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateRole",
			tenantStepParameter, suite.tenant,
		)

		roleResp.Spec.Permissions[0].Verb = []string{http.MethodGet, http.MethodPut}

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   roleName,
		}
		roleResp, err = suite.client.AuthorizationV1.CreateOrUpdateRole(ctx, tref, roleResp, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, roleResp)

		expectedRoleMetadata.verb = http.MethodPut
		verifyAuthorizationMetadataStep(sCtx, expectedRoleMetadata, roleResp.Metadata)

		expectedRoleSpec.permissions[0].verb = []string{http.MethodGet, http.MethodPut}
		verifyRoleSpecStep(sCtx, expectedRoleSpec, &roleResp.Spec)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: secalib.UpdatingStatusState,
			actualState:   string(*roleResp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
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

		expectedRoleMetadata.verb = http.MethodGet
		verifyAuthorizationMetadataStep(sCtx, expectedRoleMetadata, roleResp.Metadata)

		verifyRoleSpecStep(sCtx, expectedRoleSpec, &roleResp.Spec)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: secalib.ActiveStatusState,
			actualState:   string(*roleResp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
	})

	var expectedRoleAssignmentMetadata verifyGlobalMetadataStepParams
	var expectedRoleAssignmentSpec verifyRoleAssignmentSpecStepParams

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

		expectedRoleAssignmentMetadata = verifyGlobalMetadataStepParams{
			name:       roleAssignmentName,
			provider:   secalib.AuthorizationProviderV1,
			resource:   roleAssignmentResource,
			verb:       http.MethodPut,
			apiVersion: secalib.ApiVersion1,
			kind:       secalib.RoleAssignmentKind,
			tenant:     suite.tenant,
		}
		verifyAuthorizationMetadataStep(sCtx, expectedRoleAssignmentMetadata, assignResp.Metadata)

		expectedRoleAssignmentSpec = verifyRoleAssignmentSpecStepParams{
			roles:  []string{roleName},
			subs:   []string{roleAssignmentSub1},
			scopes: []verifyRoleAssignmentScopeStepParams{{tenants: []string{suite.tenant}}},
		}
		verifyRoleAssignmentSpecStep(sCtx, expectedRoleAssignmentSpec, &assignResp.Spec)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: secalib.CreatingStatusState,
			actualState:   string(*assignResp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
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

		expectedRoleAssignmentMetadata.verb = http.MethodGet
		verifyAuthorizationMetadataStep(sCtx, expectedRoleAssignmentMetadata, assignResp.Metadata)

		verifyRoleAssignmentSpecStep(sCtx, expectedRoleAssignmentSpec, &assignResp.Spec)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: secalib.ActiveStatusState,
			actualState:   string(*assignResp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
	})

	// Step 7
	t.WithNewStep("Update role assignment", func(sCtx provider.StepCtx) {
		sCtx.WithNewParameters(
			operationStepParameter, "CreateOrUpdateRoleAssignment",
			tenantStepParameter, suite.tenant,
		)

		assignResp.Spec.Subs = []string{roleAssignmentSub1, roleAssignmentSub2}

		tref := secapi.TenantReference{
			Tenant: secapi.TenantID(suite.tenant),
			Name:   roleAssignmentName,
		}
		assignResp, err = suite.client.AuthorizationV1.CreateOrUpdateRoleAssignment(ctx, tref, assignResp, nil)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, assignResp)

		expectedRoleAssignmentMetadata.verb = http.MethodPut
		verifyAuthorizationMetadataStep(sCtx, expectedRoleAssignmentMetadata, assignResp.Metadata)

		expectedRoleAssignmentSpec.subs = []string{roleAssignmentSub1, roleAssignmentSub2}
		verifyRoleAssignmentSpecStep(sCtx, expectedRoleAssignmentSpec, &assignResp.Spec)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: secalib.UpdatingStatusState,
			actualState:   string(*assignResp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
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

		expectedRoleAssignmentMetadata.verb = http.MethodGet
		verifyAuthorizationMetadataStep(sCtx, expectedRoleAssignmentMetadata, assignResp.Metadata)

		verifyRoleAssignmentSpecStep(sCtx, expectedRoleAssignmentSpec, &assignResp.Spec)

		verifyStatusParams := verifyStatusStepParams{
			expectedState: secalib.ActiveStatusState,
			actualState:   string(*assignResp.Status.State),
		}
		verifyStatusStep(sCtx, verifyStatusParams)
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
}

func (suite *AuthorizationV1TestSuite) AfterEach(t provider.T) {
	suite.resetAllScenarios()
}

func verifyAuthorizationMetadataStep(ctx provider.StepCtx, expected verifyGlobalMetadataStepParams, metadata *authorization.GlobalResourceMetadata) {
	actualMetadata := verifyGlobalMetadataStepParams{
		name:       metadata.Name,
		provider:   metadata.Provider,
		verb:       metadata.Verb,
		resource:   metadata.Resource,
		apiVersion: metadata.ApiVersion,
		kind:       string(metadata.Kind),
		tenant:     metadata.Tenant,
	}
	verifyGlobalMetadataStep(ctx, expected, actualMetadata)
}

func verifyRoleSpecStep(ctx provider.StepCtx, expected verifyRoleSpecStepParams, actual *authorization.RoleSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {

		stepCtx.WithNewParameters(
			"expected_permissions_length", len(expected.permissions),
			"actual_permissions_length", len(actual.Permissions),
		)
		stepCtx.Require().Equal(len(expected.permissions), len(actual.Permissions),
			"Permissions list length should match expected")

		for i := 0; i < len(expected.permissions); i++ {
			expectedPerm := expected.permissions[i]
			actualPerm := actual.Permissions[i]
			stepCtx.WithNewParameters(
				fmt.Sprintf("expected_permission[%d]_provider", i), expectedPerm.provider,
				fmt.Sprintf("actual_permission[%d]_provider", i), actualPerm.Provider,

				fmt.Sprintf("expected_permission[%d]_resources", i), expectedPerm.resources,
				fmt.Sprintf("actual_permission[%d]_resources", i), actualPerm.Resources,

				fmt.Sprintf("expected_permission[%d]_verb", i), expectedPerm.verb,
				fmt.Sprintf("actual_permission[%d]_verb", i), actualPerm.Verb,
			)
			stepCtx.Require().Equal(expectedPerm.provider, actualPerm.Provider,
				fmt.Sprintf("Permission [%d] provider should match expected", i))
			stepCtx.Require().Equal(expectedPerm.resources, actualPerm.Resources,
				fmt.Sprintf("Permission [%d] resources should match expected", i))
			stepCtx.Require().Equal(expectedPerm.verb, actualPerm.Verb,
				fmt.Sprintf("Permission [%d] verb should match expected", i))
		}
	})
}

func verifyRoleAssignmentSpecStep(ctx provider.StepCtx, expected verifyRoleAssignmentSpecStepParams, actual *authorization.RoleAssignmentSpec) {
	ctx.WithNewStep("Verify spec", func(stepCtx provider.StepCtx) {

		stepCtx.WithNewParameters(
			"expected_roles", expected.roles,
			"actual_roles", actual.Roles,

			"expected_subs", expected.subs,
			"actual_subs", actual.Subs,

			"expected_scopes_length", len(expected.scopes),
			"actual_scopoes_length", len(actual.Scopes),
		)
		stepCtx.Require().Equal(expected.roles, actual.Roles,
			"Roles provider should match expected")
		stepCtx.Require().Equal(expected.subs, actual.Subs,
			"Subs should match expected")
		stepCtx.Require().Equal(len(expected.scopes), len(actual.Scopes),
			"Scopes list length should match expected")

		for i := 0; i < len(expected.scopes); i++ {
			expectedScope := expected.scopes[i]
			actualScope := actual.Scopes[i]
			stepCtx.WithNewParameters(
				fmt.Sprintf("expected_scope[%d]_tenants", i), expectedScope.tenants,
				fmt.Sprintf("actual_scope[%d]_tenants", i), actualScope.Tenants,

				fmt.Sprintf("expected_scope[%d]_regions", i), expectedScope.regions,
				fmt.Sprintf("actual_scope[%d]_regions", i), actualScope.Regions,

				fmt.Sprintf("expected_scope[%d]_workspaces", i), expectedScope.workspaces,
				fmt.Sprintf("actual_scope[%d]_workspaces", i), actualScope.Workspaces,
			)
			if len(*actualScope.Tenants) > 0 {
				stepCtx.Require().Equal(expectedScope.tenants, *actualScope.Tenants,
					fmt.Sprintf("Scope [%d] tenants should match expected", i))
			}
			if len(*actualScope.Regions) > 0 {
				stepCtx.Require().Equal(expectedScope.regions, *actualScope.Regions,
					fmt.Sprintf("Scope [%d] regions should match expected", i))
			}
			if len(*actualScope.Workspaces) > 0 {
				stepCtx.Require().Equal(expectedScope.workspaces, *actualScope.Workspaces,
					fmt.Sprintf("Scope [%d] workspaces should match expected", i))
			}
		}
	})
}
