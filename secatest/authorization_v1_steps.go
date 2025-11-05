package secatest

import (
	"context"
	"net/http"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// Role

func (suite *testSuite) createOrUpdateRoleV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.AuthorizationV1,
	resource *schema.Role,
	expectedMeta *schema.GlobalTenantResourceMetadata,
	expectedSpec *schema.RoleSpec,
	expectedState schema.ResourceState,
) {
	expectedMeta.Verb = http.MethodPut
	createOrUpdateResourceStep(
		t,
		ctx,
		suite,
		stepName,
		suite.setAuthorizationV1StepParams,
		"CreateOrUpdateRole",
		resource,
		func(context.Context, *schema.Role) (*stepFuncResponse[schema.Role, schema.GlobalTenantResourceMetadata, schema.RoleSpec], error) {
			resp, err := api.CreateOrUpdateRole(ctx, resource)
			return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
		},
		nil,
		expectedMeta,
		suite.verifyGlobalTenantResourceMetadataStep,
		expectedSpec,
		suite.verifyRoleSpecStep,
		expectedState,
	)
}

func (suite *testSuite) getRoleV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.AuthorizationV1,
	tref secapi.TenantReference,
	expectedMeta *schema.GlobalTenantResourceMetadata,
	expectedSpec *schema.RoleSpec,
	expectedState schema.ResourceState,
) *schema.Role {
	var resp *schema.Role

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetRole")
		retry := newStepResourceStateRetry(
			suite.baseDelay,
			suite.baseInterval,
			suite.maxAttempts,
			func() (schema.ResourceState, error) {
				var err error
				resp, err = api.GetRole(ctx, tref)
				requireNoError(sCtx, err)
				requireNotNilResponse(sCtx, resp)

				suite.requireNotNilStatus(sCtx, resp.Status)
				return *resp.Status.State, nil
			},
			func() {
				expectedMeta.Verb = http.MethodGet
				suite.verifyGlobalTenantResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)

				suite.verifyRoleSpecStep(sCtx, expectedSpec, &resp.Spec)

				suite.verifyStatusStep(sCtx, expectedState, *resp.Status.State)
			},
		)
		retry.run(sCtx, "GetRole", expectedState)
	})
	return resp
}

func (suite *testSuite) getRoleWithErrorV1Step(stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.AuthorizationV1,
	tref secapi.TenantReference,
	expectedError error,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetRole")

		_, err := api.GetRole(ctx, tref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deleteRoleV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.AuthorizationV1, resource *schema.Role) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "DeleteRole")

		err := api.DeleteRole(ctx, resource)
		requireNoError(sCtx, err)
	})
}

// Role Assignment

func (suite *testSuite) createOrUpdateRoleAssignmentV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.AuthorizationV1,
	resource *schema.RoleAssignment,
	expectedMeta *schema.GlobalTenantResourceMetadata,
	expectedSpec *schema.RoleAssignmentSpec,
	expectedState schema.ResourceState,
) {
	expectedMeta.Verb = http.MethodPut
	createOrUpdateResourceStep(
		t,
		ctx,
		suite,
		stepName,
		suite.setAuthorizationV1StepParams,
		"CreateOrUpdateRoleAssignment",
		resource,
		func(context.Context, *schema.RoleAssignment) (*stepFuncResponse[schema.RoleAssignment, schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec], error) {
			resp, err := api.CreateOrUpdateRoleAssignment(ctx, resource)
			return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
		},
		nil,
		expectedMeta,
		suite.verifyGlobalTenantResourceMetadataStep,
		expectedSpec,
		suite.verifyRoleAssignmentSpecStep,
		expectedState,
	)
}

func (suite *testSuite) getRoleAssignmentV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.AuthorizationV1,
	tref secapi.TenantReference,
	expectedMeta *schema.GlobalTenantResourceMetadata,
	expectedSpec *schema.RoleAssignmentSpec,
	expectedState schema.ResourceState,
) *schema.RoleAssignment {
	var resp *schema.RoleAssignment

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetRoleAssignment")
		retry := newStepResourceStateRetry(
			suite.baseDelay,
			suite.baseInterval,
			suite.maxAttempts,
			func() (schema.ResourceState, error) {
				var err error
				resp, err = api.GetRoleAssignment(ctx, tref)
				requireNoError(sCtx, err)
				requireNotNilResponse(sCtx, resp)

				suite.requireNotNilStatus(sCtx, resp.Status)
				return *resp.Status.State, nil
			},
			func() {
				expectedMeta.Verb = http.MethodGet
				suite.verifyGlobalTenantResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)

				suite.verifyRoleAssignmentSpecStep(sCtx, expectedSpec, &resp.Spec)

				suite.verifyStatusStep(sCtx, expectedState, *resp.Status.State)
			},
		)
		retry.run(sCtx, "GetRoleAssignment", expectedState)
	})
	return resp
}

func (suite *testSuite) getRoleAssignmentWithErrorV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.AuthorizationV1,
	tref secapi.TenantReference,
	expectedError error,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetRoleAssignment")

		_, err := api.GetRoleAssignment(ctx, tref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deleteRoleAssignmentV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.AuthorizationV1, resource *schema.RoleAssignment) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "DeleteRoleAssignment")

		err := api.DeleteRoleAssignment(ctx, resource)
		requireNoError(sCtx, err)
	})
}
