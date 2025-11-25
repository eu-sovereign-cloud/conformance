package secatest

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/eu-sovereign-cloud/go-sdk/secapi/builders"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// Role

func (suite *testSuite) createOrUpdateRoleV1Step(stepName string, t provider.T, api *secapi.AuthorizationV1, resource *schema.Role,
	responseExpects responseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec],
) {
	responseExpects.metadata.Verb = http.MethodPut
	createOrUpdateTenantResourceStep(t, suite,
		createOrUpdateTenantResourceParams[schema.Role, schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
			stepName:       stepName,
			stepParamsFunc: suite.setAuthorizationV1StepParams,
			operationName:  "CreateOrUpdateRole",
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.Role) (*stepFuncResponse[schema.Role, schema.GlobalTenantResourceMetadata, schema.RoleSpec], error) {
				resp, err := api.CreateOrUpdateRole(t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.metadata,
			verifyMetadataFunc:    suite.verifyGlobalTenantResourceMetadataStep,
			expectedSpec:          responseExpects.spec,
			verifySpecFunc:        suite.verifyRoleSpecStep,
			expectedResourceState: responseExpects.resourceState,
		},
	)
}

func (suite *testSuite) getRoleV1Step(stepName string, t provider.T, api *secapi.AuthorizationV1, tref secapi.TenantReference,
	responseExpects responseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec],
) *schema.Role {
	responseExpects.metadata.Verb = http.MethodGet
	return getTenantResourceStep(t, suite,
		getTenantResourceParams[schema.Role, schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
			stepName:       stepName,
			stepParamsFunc: suite.setAuthorizationV1StepParams,
			operationName:  "GetRole",
			tref:           tref,
			getFunc: func(ctx context.Context, tref secapi.TenantReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[schema.Role, schema.GlobalTenantResourceMetadata, schema.RoleSpec], error) {
				resp, err := api.GetRoleUntilState(ctx, tref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.metadata,
			verifyMetadataFunc:    suite.verifyGlobalTenantResourceMetadataStep,
			expectedSpec:          responseExpects.spec,
			verifySpecFunc:        suite.verifyRoleSpecStep,
			expectedResourceState: responseExpects.resourceState,
		},
	)
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
func (suite *testSuite) getRoleWithErrorV1Step(stepName string, t provider.T, api *secapi.AuthorizationV1, tref secapi.TenantReference, expectedError error) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetRole")

		_, err := api.GetRole(t.Context(), tref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deleteRoleV1Step(stepName string, t provider.T, api *secapi.AuthorizationV1, resource *schema.Role) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "DeleteRole")

		err := api.DeleteRole(t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

// Role Assignment

func (suite *testSuite) createOrUpdateRoleAssignmentV1Step(stepName string, t provider.T, api *secapi.AuthorizationV1, resource *schema.RoleAssignment,
	responseExpects responseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec],
) {
	responseExpects.metadata.Verb = http.MethodPut
	createOrUpdateTenantResourceStep(t, suite,
		createOrUpdateTenantResourceParams[schema.RoleAssignment, schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
			stepName:       stepName,
			stepParamsFunc: suite.setAuthorizationV1StepParams,
			operationName:  "CreateOrUpdateRoleAssignment",
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.RoleAssignment) (*stepFuncResponse[schema.RoleAssignment, schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec], error) {
				resp, err := api.CreateOrUpdateRoleAssignment(t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.metadata,
			verifyMetadataFunc:    suite.verifyGlobalTenantResourceMetadataStep,
			expectedSpec:          responseExpects.spec,
			verifySpecFunc:        suite.verifyRoleAssignmentSpecStep,
			expectedResourceState: responseExpects.resourceState,
		},
	)
}

func (suite *testSuite) getRoleAssignmentV1Step(stepName string, t provider.T, api *secapi.AuthorizationV1, tref secapi.TenantReference,
	responseExpects responseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec],
) *schema.RoleAssignment {
	responseExpects.metadata.Verb = http.MethodGet
	return getTenantResourceStep(t, suite,
		getTenantResourceParams[schema.RoleAssignment, schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
			stepName:       stepName,
			stepParamsFunc: suite.setAuthorizationV1StepParams,
			operationName:  "GetRoleAssignment",
			tref:           tref,
			getFunc: func(ctx context.Context, tref secapi.TenantReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[schema.RoleAssignment, schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec], error) {
				resp, err := api.GetRoleAssignmentUntilState(ctx, tref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.metadata,
			verifyMetadataFunc:    suite.verifyGlobalTenantResourceMetadataStep,
			expectedSpec:          responseExpects.spec,
			verifySpecFunc:        suite.verifyRoleAssignmentSpecStep,
			expectedResourceState: responseExpects.resourceState,
		},
	)
}

func (suite *testSuite) getRoleAssignmentWithErrorV1Step(stepName string, t provider.T, api *secapi.AuthorizationV1, tref secapi.TenantReference, expectedError error) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetRoleAssignment")

		_, err := api.GetRoleAssignment(t.Context(), tref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) deleteRoleAssignmentV1Step(stepName string, t provider.T, api *secapi.AuthorizationV1, resource *schema.RoleAssignment) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "DeleteRoleAssignment")

		err := api.DeleteRoleAssignment(t.Context(), resource)
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
