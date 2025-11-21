package secatest

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/eu-sovereign-cloud/go-sdk/secapi/builders"

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

		resp, err = api.GetRole(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		if expectedMeta != nil {
			expectedMeta.Verb = http.MethodGet
			suite.verifyGlobalTenantResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)
		}

		if expectedSpec != nil {
			suite.verifyRoleSpecStep(sCtx, expectedSpec, &resp.Spec)
		}

		suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
	})
	return resp
}

func (suite *testSuite) getListRoleV1Step(stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.AuthorizationV1,
	tref secapi.TenantReference,
	opts *builders.ListOptions,
) []*schema.Role {
	var respNext []*schema.Role
	var respAll []*schema.Role

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetListRole")

		var iter *secapi.Iterator[schema.Role]

		var err error
		if opts != nil {
			iter, err = api.ListRolesWithFilters(ctx, tref.Tenant, opts)
		} else {
			iter, err = api.ListRoles(ctx, tref.Tenant)
		}
		requireNoError(sCtx, err)

		for {
			item, err := iter.Next(context.Background())
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				break
			}
			respNext = append(respNext, item)
		}
		requireNotNilResponse(sCtx, respNext)
		requireLenResponse(sCtx, len(respNext))
		/*
			respAll, err = iterAll.All(ctx)
			requireNoError(sCtx, err)
			requireNotNilResponse(sCtx, respAll)
			requireLenResponse(sCtx, len(respAll))

			compareIteratorsResponse(sCtx, len(respNext), len(respAll))
		*/
	})
	return respAll
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
	var err error

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetRoleAssignment")

		resp, err = api.GetRoleAssignment(ctx, tref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		expectedMeta.Verb = http.MethodGet
		suite.verifyGlobalTenantResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)

		suite.verifyRoleAssignmentSpecStep(sCtx, expectedSpec, &resp.Spec)

		suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
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

func (suite *testSuite) getListRoleAssignmentsV1(stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.AuthorizationV1,
	tref secapi.TenantReference,
	opts *builders.ListOptions,
) []*schema.RoleAssignment {
	var respNext []*schema.RoleAssignment
	var respAll []*schema.RoleAssignment

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetListRoleAssignment")

		var iter *secapi.Iterator[schema.RoleAssignment]
		var err error
		if opts != nil {
			iter, err = api.ListRoleAssignmentsWithFilters(ctx, tref.Tenant, opts)
		} else {
			iter, err = api.ListRoleAssignments(ctx, tref.Tenant)
		}
		requireNoError(sCtx, err)
		for {
			item, err := iter.Next(context.Background())
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				break
			}
			respNext = append(respNext, item)
		}
		requireNotNilResponse(sCtx, respNext)
		requireLenResponse(sCtx, len(respNext))
		/*
			respAll, err = iter.All(ctx)
			requireNoError(sCtx, err)
			requireNotNilResponse(sCtx, respAll)
			requireLenResponse(sCtx, len(respAll))

			compareIteratorsResponse(sCtx, len(respNext), len(respAll))
		*/
	})
	return respAll
}
