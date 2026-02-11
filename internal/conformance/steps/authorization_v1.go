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

func (configurator *StepsConfigurator) CreateOrUpdateRoleV1Step(stepName string, api *secapi.AuthorizationV1, resource *schema.Role,
	responseExpects ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	createOrUpdateTenantResourceStep(configurator,
		createOrUpdateTenantResourceParams[schema.Role, schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetAuthorizationV1StepParams,
			operationName:  "CreateOrUpdateRole",
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.Role) (*stepFuncResponse[schema.Role, schema.GlobalTenantResourceMetadata, schema.RoleSpec], error) {
				resp, err := api.CreateOrUpdateRole(configurator.t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    configurator.suite.VerifyGlobalTenantResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        configurator.suite.VerifyRoleSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (configurator *StepsConfigurator) GetRoleV1Step(stepName string, api *secapi.AuthorizationV1, tref secapi.TenantReference,
	responseExpects ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec],
) *schema.Role {
	responseExpects.Metadata.Verb = http.MethodGet
	return getTenantResourceStep(configurator,
		getTenantResourceParams[schema.Role, schema.GlobalTenantResourceMetadata, schema.RoleSpec]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetAuthorizationV1StepParams,
			operationName:  "GetRole",
			tref:           tref,
			getFunc: func(ctx context.Context, tref secapi.TenantReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[schema.Role, schema.GlobalTenantResourceMetadata, schema.RoleSpec], error) {
				resp, err := api.GetRoleUntilState(ctx, tref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    configurator.suite.VerifyGlobalTenantResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        configurator.suite.VerifyRoleSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (configurator *StepsConfigurator) GetListRoleV1Step(stepName string,
	api *secapi.AuthorizationV1,
	tref secapi.TenantReference,
	opts *secapi.ListOptions,
) {
	configurator.withStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetAuthorizationV1StepParams(sCtx, "GetListRole")
		var iter *secapi.Iterator[schema.Role]
		var err error
		if opts != nil {
			iter, err = api.ListRolesWithFilters(configurator.t.Context(), tref.Tenant, opts)
		} else {
			iter, err = api.ListRoles(configurator.t.Context(), tref.Tenant)
		}
		requireNoError(sCtx, err)
		iterResp := verifyIterListStep(sCtx, configurator.t, *iter)

		if iterResp != nil {
			configurator.suite.ReportResponseStep(sCtx, iterResp)
		}
	})
}

func (configurator *StepsConfigurator) GetRoleWithErrorV1Step(stepName string, api *secapi.AuthorizationV1, tref secapi.TenantReference, expectedError error) {
	configurator.withStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetAuthorizationV1StepParams(sCtx, "GetRole")

		_, err := api.GetRole(configurator.t.Context(), tref)
		requireError(sCtx, err, expectedError)
	})
}

func (configurator *StepsConfigurator) DeleteRoleV1Step(stepName string, api *secapi.AuthorizationV1, resource *schema.Role) {
	configurator.withStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetAuthorizationV1StepParams(sCtx, "DeleteRole")

		err := api.DeleteRole(configurator.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

// Role Assignment

func (configurator *StepsConfigurator) CreateOrUpdateRoleAssignmentV1Step(stepName string, api *secapi.AuthorizationV1, resource *schema.RoleAssignment,
	responseExpects ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	createOrUpdateTenantResourceStep(configurator,
		createOrUpdateTenantResourceParams[schema.RoleAssignment, schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetAuthorizationV1StepParams,
			operationName:  "CreateOrUpdateRoleAssignment",
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.RoleAssignment) (*stepFuncResponse[schema.RoleAssignment, schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec], error) {
				resp, err := api.CreateOrUpdateRoleAssignment(configurator.t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    configurator.suite.VerifyGlobalTenantResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        configurator.suite.VerifyRoleAssignmentSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (configurator *StepsConfigurator) GetRoleAssignmentV1Step(stepName string, api *secapi.AuthorizationV1, tref secapi.TenantReference,
	responseExpects ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec],
) *schema.RoleAssignment {
	responseExpects.Metadata.Verb = http.MethodGet
	return getTenantResourceStep(configurator,
		getTenantResourceParams[schema.RoleAssignment, schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetAuthorizationV1StepParams,
			operationName:  "GetRoleAssignment",
			tref:           tref,
			getFunc: func(ctx context.Context, tref secapi.TenantReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[schema.RoleAssignment, schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec], error) {
				resp, err := api.GetRoleAssignmentUntilState(ctx, tref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    configurator.suite.VerifyGlobalTenantResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        configurator.suite.VerifyRoleAssignmentSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (configurator *StepsConfigurator) GetListRoleAssignmentsV1(stepName string,
	api *secapi.AuthorizationV1,
	tref secapi.TenantReference,
	opts *secapi.ListOptions,
) {
	configurator.withStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetAuthorizationV1StepParams(sCtx, "GetListRoleAssignment")

		var iter *secapi.Iterator[schema.RoleAssignment]
		var err error
		if opts != nil {
			iter, err = api.ListRoleAssignmentsWithFilters(configurator.t.Context(), tref.Tenant, opts)
		} else {
			iter, err = api.ListRoleAssignments(configurator.t.Context(), tref.Tenant)
		}
		requireNoError(sCtx, err)

		iterResp := verifyIterListStep(sCtx, configurator.t, *iter)

		if iterResp != nil {
			configurator.suite.ReportResponseStep(sCtx, iterResp)
		}
	})
}

func (configurator *StepsConfigurator) GetRoleAssignmentWithErrorV1Step(stepName string, api *secapi.AuthorizationV1, tref secapi.TenantReference, expectedError error) {
	configurator.withStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetAuthorizationV1StepParams(sCtx, "GetRoleAssignment")

		_, err := api.GetRoleAssignment(configurator.t.Context(), tref)
		requireError(sCtx, err, expectedError)
	})
}

func (configurator *StepsConfigurator) DeleteRoleAssignmentV1Step(stepName string, api *secapi.AuthorizationV1, resource *schema.RoleAssignment) {
	configurator.withStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetAuthorizationV1StepParams(sCtx, "DeleteRoleAssignment")

		err := api.DeleteRoleAssignment(configurator.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}
