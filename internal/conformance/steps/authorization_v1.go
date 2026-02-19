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

func (configurator *StepsConfigurator) CreateOrUpdateRoleV1Step(stepName string, api secapi.AuthorizationV1, resource *schema.Role,
	responseExpects StepResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	configurator.logStepName(stepName)
	createOrUpdateTenantResourceStep(configurator.t, configurator.suite,
		createOrUpdateTenantResourceParams[schema.Role, schema.GlobalTenantResourceMetadata, schema.RoleSpec, schema.Status]{
			createOrUpdateResourceParams: createOrUpdateResourceParams[schema.Role, schema.GlobalTenantResourceMetadata, schema.RoleSpec, schema.Status]{
				resource: resource,
				createOrUpdateFunc: func(context.Context, *schema.Role) (
					*createOrUpdateStepFuncResponse[schema.Role, schema.GlobalTenantResourceMetadata, schema.RoleSpec, schema.Status], error,
				) {
					if resp, err := api.CreateOrUpdateRole(configurator.t.Context(), resource); err == nil {
						return newCreateOrUpdateStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), nil
					} else {
						return nil, err
					}
				},
				expectedMetadata:      responseExpects.Metadata,
				verifyMetadataFunc:    configurator.suite.VerifyGlobalTenantResourceMetadataStep,
				expectedSpec:          responseExpects.Spec,
				verifySpecFunc:        configurator.suite.VerifyRoleSpecStep,
				expectedResourceState: responseExpects.ResourceState,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetAuthorizationV1StepParams,
			operationName:  "CreateOrUpdateRole",
		},
	)
}

func (configurator *StepsConfigurator) GetRoleV1Step(stepName string, api secapi.AuthorizationV1, tref secapi.TenantReference,
	responseExpects StepResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec],
) *schema.Role {
	responseExpects.Metadata.Verb = http.MethodGet
	configurator.logStepName(stepName)
	return getTenantResourceStep(configurator.t, configurator.suite,
		getTenantResourceParams[schema.Role, schema.GlobalTenantResourceMetadata, schema.RoleSpec, schema.Status]{
			getResourceWithObserverParams: getResourceWithObserverParams[schema.Role, schema.GlobalTenantResourceMetadata, schema.RoleSpec, schema.Status, secapi.TenantReference, schema.ResourceState]{
				reference: tref,
				getFunc: func(ctx context.Context, tref secapi.TenantReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (
					*getStepFuncResponse[schema.Role, schema.GlobalTenantResourceMetadata, schema.RoleSpec, schema.Status], error,
				) {
					if resp, err := api.GetRoleUntilState(ctx, tref, config); err == nil {
						return newGetStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), nil
					} else {
						return nil, err
					}
				},
				expectedMetadata:      responseExpects.Metadata,
				verifyMetadataFunc:    configurator.suite.VerifyGlobalTenantResourceMetadataStep,
				expectedSpec:          responseExpects.Spec,
				verifySpecFunc:        configurator.suite.VerifyRoleSpecStep,
				expectedResourceState: responseExpects.ResourceState,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetAuthorizationV1StepParams,
			operationName:  "GetRole",
		},
	)
}

func (configurator *StepsConfigurator) ListRoleV1Step(stepName string,
	api secapi.AuthorizationV1,
	tref secapi.TenantReference,
	opts *secapi.ListOptions,
) {
	configurator.logStepName(stepName)
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetAuthorizationV1StepParams(sCtx, "GetListRole")

		var iter *secapi.Iterator[schema.Role]
		var err error
		if opts != nil {
			iter, err = api.ListRolesWithFilters(configurator.t.Context(), tref.Tenant, opts)
		} else {
			iter, err = api.ListRoles(configurator.t.Context(), tref.Tenant)
		}
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, configurator.t, *iter)
	})
}

func (configurator *StepsConfigurator) GetRoleWithErrorV1Step(stepName string, api secapi.AuthorizationV1, tref secapi.TenantReference, expectedError error) {
	configurator.logStepName(stepName)
	getTenantResourceWithErrorStep(configurator.t,
		getTenantResourceWithErrorParams{
			getResourceWithErrorParams: getResourceWithErrorParams[secapi.TenantReference]{
				reference: tref,
				getFunc: func(ctx context.Context, tref secapi.TenantReference) error {
					_, err := api.GetRole(ctx, tref)
					return err
				},
				expectedError: expectedError,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetAuthorizationV1StepParams,
			operationName:  "GetRole",
		},
	)
}

func (configurator *StepsConfigurator) DeleteRoleV1Step(stepName string, api secapi.AuthorizationV1, resource *schema.Role) {
	configurator.logStepName(stepName)
	deleteTenantResourceStep(configurator.t,
		deleteTenantResourceParams[schema.Role]{
			deleteResourceParams: deleteResourceParams[schema.Role]{
				resource: resource,
				deleteFunc: func(ctx context.Context, r *schema.Role) error {
					return api.DeleteRole(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetAuthorizationV1StepParams,
			operationName:  "DeleteRole",
		},
	)
}

// Role Assignment

func (configurator *StepsConfigurator) CreateOrUpdateRoleAssignmentV1Step(stepName string, api secapi.AuthorizationV1, resource *schema.RoleAssignment,
	responseExpects StepResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	configurator.logStepName(stepName)
	createOrUpdateTenantResourceStep(configurator.t, configurator.suite,
		createOrUpdateTenantResourceParams[schema.RoleAssignment, schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec, schema.Status]{
			createOrUpdateResourceParams: createOrUpdateResourceParams[schema.RoleAssignment, schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec, schema.Status]{
				resource: resource,
				createOrUpdateFunc: func(context.Context, *schema.RoleAssignment) (
					*createOrUpdateStepFuncResponse[schema.RoleAssignment, schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec, schema.Status], error,
				) {
					if resp, err := api.CreateOrUpdateRoleAssignment(configurator.t.Context(), resource); err == nil {
						return newCreateOrUpdateStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), nil
					} else {
						return nil, err
					}
				},
				expectedMetadata:      responseExpects.Metadata,
				verifyMetadataFunc:    configurator.suite.VerifyGlobalTenantResourceMetadataStep,
				expectedSpec:          responseExpects.Spec,
				verifySpecFunc:        configurator.suite.VerifyRoleAssignmentSpecStep,
				expectedResourceState: responseExpects.ResourceState,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetAuthorizationV1StepParams,
			operationName:  "CreateOrUpdateRoleAssignment",
		},
	)
}

func (configurator *StepsConfigurator) GetRoleAssignmentV1Step(stepName string, api secapi.AuthorizationV1, tref secapi.TenantReference,
	responseExpects StepResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec],
) *schema.RoleAssignment {
	responseExpects.Metadata.Verb = http.MethodGet
	configurator.logStepName(stepName)
	return getTenantResourceStep(configurator.t, configurator.suite,
		getTenantResourceParams[schema.RoleAssignment, schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec, schema.Status]{
			getResourceWithObserverParams: getResourceWithObserverParams[schema.RoleAssignment, schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec, schema.Status, secapi.TenantReference, schema.ResourceState]{
				reference: tref,
				getFunc: func(ctx context.Context, tref secapi.TenantReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (
					*getStepFuncResponse[schema.RoleAssignment, schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec, schema.Status], error,
				) {
					if resp, err := api.GetRoleAssignmentUntilState(ctx, tref, config); err == nil {
						return newGetStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), nil
					} else {
						return nil, err
					}
				},
				expectedMetadata:      responseExpects.Metadata,
				verifyMetadataFunc:    configurator.suite.VerifyGlobalTenantResourceMetadataStep,
				expectedSpec:          responseExpects.Spec,
				verifySpecFunc:        configurator.suite.VerifyRoleAssignmentSpecStep,
				expectedResourceState: responseExpects.ResourceState,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetAuthorizationV1StepParams,
			operationName:  "GetRoleAssignment",
		},
	)
}

func (configurator *StepsConfigurator) ListRoleAssignmentsV1(
	stepName string, api secapi.AuthorizationV1, tid secapi.TenantID, opts *secapi.ListOptions,
) []*schema.RoleAssignment {
	configurator.logStepName(stepName)
	return listTenantResourcesStep(configurator.t,
		listTenantResourcesParams[schema.RoleAssignment, schema.GlobalTenantResourceMetadata]{
			listResourcesParams: listResourcesParams[schema.RoleAssignment, schema.GlobalTenantResourceMetadata, secapi.TenantID]{
				reference: tid,
				listFunc: func(ctx context.Context, tref secapi.TenantID, opts *secapi.ListOptions) (
					*secapi.Iterator[schema.RoleAssignment], error,
				) {
					if opts != nil {
						return api.ListRoleAssignmentsWithFilters(ctx, tid, opts)
					} else {
						return api.ListRoleAssignments(ctx, tid)
					}
				},
				listOptions: opts,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetAuthorizationV1StepParams,
			operationName:  "ListRoleAssignments",
		},
	)
}

func (configurator *StepsConfigurator) GetRoleAssignmentWithErrorV1Step(stepName string, api secapi.AuthorizationV1, tref secapi.TenantReference, expectedError error) {
	configurator.logStepName(stepName)
	getTenantResourceWithErrorStep(configurator.t,
		getTenantResourceWithErrorParams{
			getResourceWithErrorParams: getResourceWithErrorParams[secapi.TenantReference]{
				reference: tref,
				getFunc: func(ctx context.Context, tref secapi.TenantReference) error {
					_, err := api.GetRoleAssignment(ctx, tref)
					return err
				},
				expectedError: expectedError,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetAuthorizationV1StepParams,
			operationName:  "GetRoleAssignment",
		},
	)
}

func (configurator *StepsConfigurator) DeleteRoleAssignmentV1Step(stepName string, api secapi.AuthorizationV1, resource *schema.RoleAssignment) {
	configurator.logStepName(stepName)
	deleteTenantResourceStep(configurator.t,
		deleteTenantResourceParams[schema.RoleAssignment]{
			deleteResourceParams: deleteResourceParams[schema.RoleAssignment]{
				resource: resource,
				deleteFunc: func(ctx context.Context, r *schema.RoleAssignment) error {
					return api.DeleteRoleAssignment(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetAuthorizationV1StepParams,
			operationName:  "DeleteRoleAssignment",
		},
	)
}
