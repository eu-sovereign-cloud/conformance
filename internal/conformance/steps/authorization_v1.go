//nolint:dupl
package steps

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/pkg/wrappers"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// Role

func (configurator *StepsConfigurator) CreateOrUpdateRoleV1Step(stepName string, api secapi.AuthorizationV1, resource *schema.Role,
	responseExpects ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	createOrUpdateTenantResourceStep(configurator.t, configurator.suite,
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

func (configurator *StepsConfigurator) GetRoleV1Step(stepName string, api secapi.AuthorizationV1, tref secapi.TenantReference,
	responseExpects ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleSpec],
) *schema.Role {
	responseExpects.Metadata.Verb = http.MethodGet
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
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

func (configurator *StepsConfigurator) ListRoleV1Step(
	stepName string, api secapi.AuthorizationV1, tref secapi.TenantReference, opts *secapi.FilterOptions,
) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetAuthorizationV1StepParams(sCtx, "ListRole")

		iter, err := api.ListRoles(configurator.t.Context(), secapi.TenantFilter{Tenant: tref.Tenant, Options: opts})
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, configurator.t, *iter)
	})
}

func (configurator *StepsConfigurator) WatchRoleUntilDeletedV1Step(stepName string, api secapi.AuthorizationV1, tref secapi.TenantReference) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		watchTenantResourceUntilDeletedStep(configurator.t, configurator.suite,
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
	})
}

func (configurator *StepsConfigurator) DeleteRoleV1Step(stepName string, api secapi.AuthorizationV1, resource *schema.Role) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetAuthorizationV1StepParams(sCtx, "DeleteRole")

		err := api.DeleteRole(configurator.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

// Role Assignment

func (configurator *StepsConfigurator) CreateOrUpdateRoleAssignmentV1Step(stepName string, api secapi.AuthorizationV1, resource *schema.RoleAssignment,
	responseExpects ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	createOrUpdateTenantResourceStep(configurator.t, configurator.suite,
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

func (configurator *StepsConfigurator) GetRoleAssignmentV1Step(stepName string, api secapi.AuthorizationV1, tref secapi.TenantReference,
	responseExpects ResponseExpects[schema.GlobalTenantResourceMetadata, schema.RoleAssignmentSpec],
) *schema.RoleAssignment {
	responseExpects.Metadata.Verb = http.MethodGet
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
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

func (configurator *StepsConfigurator) ListRoleAssignmentsV1(
	stepName string, api secapi.AuthorizationV1, tref secapi.TenantReference, opts *secapi.FilterOptions,
) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetAuthorizationV1StepParams(sCtx, "ListRoleAssignment")

		iter, err := api.ListRoleAssignments(configurator.t.Context(), secapi.TenantFilter{Tenant: tref.Tenant, Options: opts})
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, configurator.t, *iter)
	})
}

func (configurator *StepsConfigurator) WatchRoleAssignmentUntilDeletedV1Step(stepName string, api secapi.AuthorizationV1, tref secapi.TenantReference) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		watchTenantResourceUntilDeletedStep(configurator.t, configurator.suite,
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
	})
}

func (configurator *StepsConfigurator) DeleteRoleAssignmentV1Step(stepName string, api secapi.AuthorizationV1, resource *schema.RoleAssignment) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetAuthorizationV1StepParams(sCtx, "DeleteRoleAssignment")

		err := api.DeleteRoleAssignment(configurator.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}
