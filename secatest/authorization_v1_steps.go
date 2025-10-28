package secatest

import (
	"context"
	"net/http"
	"time"

	"github.com/eu-sovereign-cloud/conformance/secalib"
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
	expectedStatusState string,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "CreateOrUpdateRole")

		resp, err := api.CreateOrUpdateRole(ctx, resource)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		if expectedMeta != nil {
			expectedMeta.Verb = http.MethodPut
			suite.verifyGlobalTenantResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)
		}

		if expectedSpec != nil {
			suite.verifyRoleSpecStep(sCtx, expectedSpec, &resp.Spec)
		}

		suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
	})
}

func (suite *testSuite) getRoleV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.AuthorizationV1,
	tref secapi.TenantReference,
	expectedMeta *schema.GlobalTenantResourceMetadata,
	expectedSpec *schema.RoleSpec,
	expectedStatusState string,
) *schema.Role {
	var resp *schema.Role
	var err error

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetRole")
		time.Sleep(time.Duration(suite.baseDelay) * time.Second)
		for attempt := 1; attempt <= suite.maxAttempts; attempt++ {
			resp, err = api.GetRole(ctx, tref)
			requireNoError(sCtx, err)
			requireNotNilResponse(sCtx, resp)
			if resp.Status.State != nil && *resp.Status.State == *secalib.SetResourceState(expectedStatusState) {

				if expectedMeta != nil {
					expectedMeta.Verb = http.MethodGet
					suite.verifyGlobalTenantResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)
				}

				if expectedSpec != nil {
					suite.verifyRoleSpecStep(sCtx, expectedSpec, &resp.Spec)
				}

				suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
				return
			} else {
				time.Sleep(time.Duration(suite.baseInterval) * time.Second)
			}
			suite.verifyMaxAttempts(sCtx, attempt, "GetRole", expectedStatusState)
		}
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
	expectedStatusState string,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "CreateOrUpdateRoleAssignment")

		resp, err := api.CreateOrUpdateRoleAssignment(ctx, resource)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		expectedMeta.Verb = http.MethodPut
		suite.verifyGlobalTenantResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)

		suite.verifyRoleAssignmentSpecStep(sCtx, expectedSpec, &resp.Spec)

		suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
	})
}

func (suite *testSuite) getRoleAssignmentV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.AuthorizationV1,
	tref secapi.TenantReference,
	expectedMeta *schema.GlobalTenantResourceMetadata,
	expectedSpec *schema.RoleAssignmentSpec,
	expectedStatusState string,
) *schema.RoleAssignment {
	var resp *schema.RoleAssignment

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetRoleAssignment")
		retry := newStepRetry(
			suite.baseDelay,
			suite.baseInterval,
			suite.maxAttempts,
			func() schema.ResourceState {
				var err error
				resp, err = api.GetRoleAssignment(ctx, tref)
				requireNoError(sCtx, err)
				requireNotNilResponse(sCtx, resp)

				suite.requireNotNilStatus(sCtx, resp.Status)
				return *resp.Status.State
			},
			func() {
				expectedMeta.Verb = http.MethodGet
				suite.verifyGlobalTenantResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)

				suite.verifyRoleAssignmentSpecStep(sCtx, expectedSpec, &resp.Spec)

				suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
			},
		)
		retry.run(sCtx, "GetRoleAssignment", expectedStatusState)
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
