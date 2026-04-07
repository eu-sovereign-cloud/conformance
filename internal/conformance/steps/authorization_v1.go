//nolint:dupl
package steps

import (
	"context"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/pkg/wrappers"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
)

// Role

func (configurator *StepsConfigurator) ListRoleV1Step(stepName string, api secapi.AuthorizationV1, tpath secapi.TenantPath, opts *secapi.ListOptions) {
	listTenantResourcesStep(configurator.t, configurator.suite,
		listTenantResourcesParams[schema.Role, schema.GlobalTenantResourceMetadata]{
			listResourcesParams: listResourcesParams[schema.Role, schema.GlobalTenantResourceMetadata, secapi.TenantPath]{
				path: tpath, listOptions: opts,
				listFunc: func(ctx context.Context, path secapi.TenantPath, options *secapi.ListOptions) (*secapi.Iterator[schema.Role], error) {
					return api.ListRolesWithOptions(ctx, path, options)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetAuthorizationV1StepParams,
			operationName:  constants.ListRolesOperation,
		},
	)
}

func (configurator *StepsConfigurator) GetRoleV1Step(stepName string, api secapi.AuthorizationV1, tref secapi.TenantReference,
	responseExpects StepResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec],
) *schema.Role {
	responseExpects.Metadata.Verb = http.MethodGet
	return getTenantResourceStep(configurator.t, configurator.suite,
		getTenantResourceParams[schema.Role, schema.GlobalTenantResourceMetadata, schema.RoleSpec, schema.Status]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetAuthorizationV1StepParams,
			operationName:  constants.GetRoleOperation,
			tref:           tref,
			getValueFunc: func(ctx context.Context, tref secapi.TenantReference, config secapi.ResourceObserverUntilValueConfig[schema.ResourceState]) (
				wrappers.ResourceWrapper[schema.Role, schema.GlobalTenantResourceMetadata, schema.RoleSpec, schema.Status], error,
			) {
				resp, err := api.GetRoleUntilState(ctx, tref, config)
				return wrappers.NewRoleWrapper(resp), err
			},
			expectedMetadata:       responseExpects.Metadata,
			verifyMetadataFunc:     configurator.suite.VerifyGlobalTenantResourceMetadataStep,
			expectedSpec:           responseExpects.Spec,
			verifySpecFunc:         configurator.suite.VerifyRoleSpecStep,
			expectedResourceStates: responseExpects.ResourceStates,
		},
	)
}

func (configurator *StepsConfigurator) WatchRoleUntilDeletedV1Step(stepName string, stepCreator StepCreator, api secapi.AuthorizationV1, tref secapi.TenantReference) {
	watchTenantResourceUntilDeletedStep(configurator.t.Context(), configurator.suite, stepCreator,
		watchTenantResourceUntilDeletedParams{
			watchResourceUntilDeletedParams: watchResourceUntilDeletedParams[secapi.TenantReference]{
				reference: tref,
				getErrorFunc: func(ctx context.Context, tref secapi.TenantReference, config secapi.ResourceObserverConfig) error {
					return api.WatchRoleUntilDeleted(configurator.t.Context(), tref, config)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetWorkspaceV1StepParams,
			operationName:  constants.GetRoleOperation,
		},
	)
}

func (configurator *StepsConfigurator) CreateOrUpdateRoleV1Step(stepName string, stepCreator StepCreator, api secapi.AuthorizationV1, resource *schema.Role,
	responseExpects StepResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	createOrUpdateTenantResourceStep(configurator.t.Context(), configurator.suite, stepCreator,
		createOrUpdateTenantResourceParams[schema.Role, schema.GlobalTenantResourceMetadata, schema.RoleSpec, schema.Status]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetAuthorizationV1StepParams,
			operationName:  constants.CreateOrUpdateRoleOperation,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.Role) (
				wrappers.ResourceWrapper[schema.Role, schema.GlobalTenantResourceMetadata, schema.RoleSpec, schema.Status], error,
			) {
				resp, err := api.CreateOrUpdateRole(configurator.t.Context(), resource)
				return wrappers.NewRoleWrapper(resp), err
			},
			expectedMetadata:       responseExpects.Metadata,
			verifyMetadataFunc:     configurator.suite.VerifyGlobalTenantResourceMetadataStep,
			expectedSpec:           responseExpects.Spec,
			verifySpecFunc:         configurator.suite.VerifyRoleSpecStep,
			expectedResourceStates: responseExpects.ResourceStates,
		},
	)
}

func (configurator *StepsConfigurator) DeleteRoleV1Step(stepName string, stepCreator StepCreator, api secapi.AuthorizationV1, resource *schema.Role) {
	deleteTenantResourceStep(configurator.t.Context(), configurator.suite, stepCreator,
		deleteTenantResourceParams[schema.Role]{
			deleteResourceParams: deleteResourceParams[schema.Role]{
				resource: resource,
				deleteFunc: func(ctx context.Context, r *schema.Role) error {
					return api.DeleteRole(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetAuthorizationV1StepParams,
			operationName:  constants.DeleteRoleOperation,
		},
	)
}

// Role Assignment

func (configurator *StepsConfigurator) ListRoleAssignmentsV1(stepName string, api secapi.AuthorizationV1, tpath secapi.TenantPath, opts *secapi.ListOptions) {
	listTenantResourcesStep(configurator.t, configurator.suite,
		listTenantResourcesParams[schema.RoleAssignment, schema.GlobalTenantResourceMetadata]{
			listResourcesParams: listResourcesParams[schema.RoleAssignment, schema.GlobalTenantResourceMetadata, secapi.TenantPath]{
				path: tpath, listOptions: opts,
				listFunc: func(ctx context.Context, path secapi.TenantPath, options *secapi.ListOptions) (*secapi.Iterator[schema.RoleAssignment], error) {
					return api.ListRoleAssignmentsWithOptions(ctx, path, options)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetAuthorizationV1StepParams,
			operationName:  constants.ListRoleAssignmentsOperation,
		},
	)
}

func (configurator *StepsConfigurator) GetRoleAssignmentV1Step(stepName string, api secapi.AuthorizationV1, tref secapi.TenantReference,
	responseExpects StepResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec],
) *schema.RoleAssignment {
	responseExpects.Metadata.Verb = http.MethodGet
	return getTenantResourceStep(configurator.t, configurator.suite,
		getTenantResourceParams[schema.RoleAssignment, schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec, schema.Status]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetAuthorizationV1StepParams,
			operationName:  constants.GetRoleAssignmentOperation,
			tref:           tref,
			getValueFunc: func(ctx context.Context, tref secapi.TenantReference, config secapi.ResourceObserverUntilValueConfig[schema.ResourceState]) (
				wrappers.ResourceWrapper[schema.RoleAssignment, schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec, schema.Status], error,
			) {
				resp, err := api.GetRoleAssignmentUntilState(ctx, tref, config)
				return wrappers.NewRoleAssignmentWrapper(resp), err
			},
			expectedMetadata:       responseExpects.Metadata,
			verifyMetadataFunc:     configurator.suite.VerifyGlobalTenantResourceMetadataStep,
			expectedSpec:           responseExpects.Spec,
			verifySpecFunc:         configurator.suite.VerifyRoleAssignmentSpecStep,
			expectedResourceStates: responseExpects.ResourceStates,
		},
	)
}

func (configurator *StepsConfigurator) WatchRoleAssignmentUntilDeletedV1Step(stepName string, stepCreator StepCreator, api secapi.AuthorizationV1, tref secapi.TenantReference) {
	watchTenantResourceUntilDeletedStep(configurator.t.Context(), configurator.suite, stepCreator,
		watchTenantResourceUntilDeletedParams{
			watchResourceUntilDeletedParams: watchResourceUntilDeletedParams[secapi.TenantReference]{
				reference: tref,
				getErrorFunc: func(ctx context.Context, tref secapi.TenantReference, config secapi.ResourceObserverConfig) error {
					return api.WatchRoleAssignmentUntilDeleted(configurator.t.Context(), tref, config)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetWorkspaceV1StepParams,
			operationName:  constants.GetRoleAssignmentOperation,
		},
	)
}

func (configurator *StepsConfigurator) CreateOrUpdateRoleAssignmentV1Step(stepName string, stepCreator StepCreator, api secapi.AuthorizationV1, resource *schema.RoleAssignment,
	responseExpects StepResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	createOrUpdateTenantResourceStep(configurator.t.Context(), configurator.suite, stepCreator,
		createOrUpdateTenantResourceParams[schema.RoleAssignment, schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec, schema.Status]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetAuthorizationV1StepParams,
			operationName:  constants.CreateOrUpdateRoleAssignmentOperation,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.RoleAssignment) (
				wrappers.ResourceWrapper[schema.RoleAssignment, schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec, schema.Status], error,
			) {
				resp, err := api.CreateOrUpdateRoleAssignment(configurator.t.Context(), resource)
				return wrappers.NewRoleAssignmentWrapper(resp), err
			},
			expectedMetadata:       responseExpects.Metadata,
			verifyMetadataFunc:     configurator.suite.VerifyGlobalTenantResourceMetadataStep,
			expectedSpec:           responseExpects.Spec,
			verifySpecFunc:         configurator.suite.VerifyRoleAssignmentSpecStep,
			expectedResourceStates: responseExpects.ResourceStates,
		},
	)
}

func (configurator *StepsConfigurator) DeleteRoleAssignmentV1Step(stepName string, stepCreator StepCreator, api secapi.AuthorizationV1, resource *schema.RoleAssignment) {
	deleteTenantResourceStep(configurator.t.Context(), configurator.suite, stepCreator,
		deleteTenantResourceParams[schema.RoleAssignment]{
			deleteResourceParams: deleteResourceParams[schema.RoleAssignment]{
				resource: resource,
				deleteFunc: func(ctx context.Context, r *schema.RoleAssignment) error {
					return api.DeleteRoleAssignment(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetAuthorizationV1StepParams,
			operationName:  constants.DeleteRoleAssignmentOperation,
		},
	)
}
