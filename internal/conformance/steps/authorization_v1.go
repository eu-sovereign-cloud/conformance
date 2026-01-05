//nolint:dupl
package steps

import (
	"context"
	"net/http"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// Role

func (builder *Builder) CreateOrUpdateRoleV1Step(stepName string, api *secapi.AuthorizationV1, resource *schema.Role,
	responseExpects ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	createOrUpdateTenantResourceStep(builder.t, builder.suite,
		createOrUpdateTenantResourceParams[schema.Role, schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
			stepName:       stepName,
			stepParamsFunc: builder.suite.SetAuthorizationV1StepParams,
			operationName:  "CreateOrUpdateRole",
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.Role) (*stepFuncResponse[schema.Role, schema.GlobalTenantResourceMetadata, schema.RoleSpec], error) {
				resp, err := api.CreateOrUpdateRole(builder.t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    builder.suite.VerifyGlobalTenantResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        builder.suite.VerifyRoleSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (builder *Builder) GetRoleV1Step(stepName string, api *secapi.AuthorizationV1, tref secapi.TenantReference,
	responseExpects ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec],
) *schema.Role {
	responseExpects.Metadata.Verb = http.MethodGet
	return getTenantResourceStep(builder.t, builder.suite,
		getTenantResourceParams[schema.Role, schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
			stepName:       stepName,
			stepParamsFunc: builder.suite.SetAuthorizationV1StepParams,
			operationName:  "GetRole",
			tref:           tref,
			getFunc: func(ctx context.Context, tref secapi.TenantReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[schema.Role, schema.GlobalTenantResourceMetadata, schema.RoleSpec], error) {
				resp, err := api.GetRoleUntilState(ctx, tref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    builder.suite.VerifyGlobalTenantResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        builder.suite.VerifyRoleSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (builder *Builder) GetListRoleV1Step(stepName string,
	api *secapi.AuthorizationV1,
	tref secapi.TenantReference,
	opts *secapi.ListOptions,
) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetAuthorizationV1StepParams(sCtx, "GetListRole")

		var iter *secapi.Iterator[schema.Role]

		var err error
		if opts != nil {
			iter, err = api.ListRolesWithFilters(builder.t.Context(), tref.Tenant, opts)
		} else {
			iter, err = api.ListRoles(builder.t.Context(), tref.Tenant)
		}
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, builder.t, *iter)
	})
}

func (builder *Builder) GetRoleWithErrorV1Step(stepName string, api *secapi.AuthorizationV1, tref secapi.TenantReference, expectedError error) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetAuthorizationV1StepParams(sCtx, "GetRole")

		_, err := api.GetRole(builder.t.Context(), tref)
		requireError(sCtx, err, expectedError)
	})
}

func (builder *Builder) DeleteRoleV1Step(stepName string, api *secapi.AuthorizationV1, resource *schema.Role) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetAuthorizationV1StepParams(sCtx, "DeleteRole")

		err := api.DeleteRole(builder.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

// Role Assignment

func (builder *Builder) CreateOrUpdateRoleAssignmentV1Step(stepName string, api *secapi.AuthorizationV1, resource *schema.RoleAssignment,
	responseExpects ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	createOrUpdateTenantResourceStep(builder.t, builder.suite,
		createOrUpdateTenantResourceParams[schema.RoleAssignment, schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
			stepName:       stepName,
			stepParamsFunc: builder.suite.SetAuthorizationV1StepParams,
			operationName:  "CreateOrUpdateRoleAssignment",
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.RoleAssignment) (*stepFuncResponse[schema.RoleAssignment, schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec], error) {
				resp, err := api.CreateOrUpdateRoleAssignment(builder.t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    builder.suite.VerifyGlobalTenantResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        builder.suite.VerifyRoleAssignmentSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (builder *Builder) GetRoleAssignmentV1Step(stepName string, api *secapi.AuthorizationV1, tref secapi.TenantReference,
	responseExpects ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec],
) *schema.RoleAssignment {
	responseExpects.Metadata.Verb = http.MethodGet
	return getTenantResourceStep(builder.t, builder.suite,
		getTenantResourceParams[schema.RoleAssignment, schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
			stepName:       stepName,
			stepParamsFunc: builder.suite.SetAuthorizationV1StepParams,
			operationName:  "GetRoleAssignment",
			tref:           tref,
			getFunc: func(ctx context.Context, tref secapi.TenantReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[schema.RoleAssignment, schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec], error) {
				resp, err := api.GetRoleAssignmentUntilState(ctx, tref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    builder.suite.VerifyGlobalTenantResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        builder.suite.VerifyRoleAssignmentSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (builder *Builder) GetListRoleAssignmentsV1(stepName string,
	api *secapi.AuthorizationV1,
	tref secapi.TenantReference,
	opts *secapi.ListOptions,
) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetAuthorizationV1StepParams(sCtx, "GetListRoleAssignment")

		var iter *secapi.Iterator[schema.RoleAssignment]
		var err error
		if opts != nil {
			iter, err = api.ListRoleAssignmentsWithFilters(builder.t.Context(), tref.Tenant, opts)
		} else {
			iter, err = api.ListRoleAssignments(builder.t.Context(), tref.Tenant)
		}
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, builder.t, *iter)
	})
}

func (builder *Builder) GetRoleAssignmentWithErrorV1Step(stepName string, api *secapi.AuthorizationV1, tref secapi.TenantReference, expectedError error) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetAuthorizationV1StepParams(sCtx, "GetRoleAssignment")

		_, err := api.GetRoleAssignment(builder.t.Context(), tref)
		requireError(sCtx, err, expectedError)
	})
}

func (builder *Builder) DeleteRoleAssignmentV1Step(stepName string, api *secapi.AuthorizationV1, resource *schema.RoleAssignment) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetAuthorizationV1StepParams(sCtx, "DeleteRoleAssignment")

		err := api.DeleteRoleAssignment(builder.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}
