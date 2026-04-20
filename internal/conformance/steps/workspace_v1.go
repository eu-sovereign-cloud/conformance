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

func (configurator *StepsConfigurator) CreateOrUpdateWorkspaceV1Step(stepName string, api secapi.WorkspaceV1, resource *schema.Workspace,
	responseExpects ResponseExpects[schema.RegionalResourceMetadata, schema.WorkspaceSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	createOrUpdateTenantResourceStep(configurator.t, configurator.suite,
		createOrUpdateTenantResourceParams[schema.Workspace, schema.RegionalResourceMetadata, schema.WorkspaceSpec, schema.WorkspaceStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetWorkspaceV1StepParams,
			operationName:  constants.CreateOrUpdateWorkspaceOperation,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.Workspace) (
				wrappers.ResourceWrapper[schema.Workspace, schema.RegionalResourceMetadata, schema.WorkspaceSpec, schema.WorkspaceStatus], error,
			) {
				resp, err := api.CreateOrUpdateWorkspace(configurator.t.Context(), resource)
				return wrappers.NewWorkspaceWrapper(resp), err
			},
			expectedLabels:         responseExpects.Labels,
			expectedAnnotations:    responseExpects.Annotations,
			expectedExtensions:     responseExpects.Extensions,
			expectedMetadata:       responseExpects.Metadata,
			verifyMetadataFunc:     configurator.suite.VerifyRegionalResourceMetadataStep,
			expectedResourceStates: responseExpects.ResourceStates,
		},
	)
}

func (configurator *StepsConfigurator) GetWorkspaceV1Step(stepName string, api secapi.WorkspaceV1, tref secapi.TenantReference,
	responseExpects ResponseExpectsWithCondition[schema.RegionalResourceMetadata, schema.WorkspaceSpec, schema.WorkspaceStatus],
) *schema.Workspace {
	responseExpects.Metadata.Verb = http.MethodGet
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	return getTenantResourceStep(configurator.t, configurator.suite,
		getTenantResourceParams[schema.Workspace, schema.RegionalResourceMetadata, schema.WorkspaceSpec, schema.WorkspaceStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetWorkspaceV1StepParams,
			operationName:  constants.GetWorkspaceOperation,
			tref:           tref,
			getValueFunc: func(ctx context.Context, tref secapi.TenantReference, config secapi.ResourceObserverUntilValueConfig[schema.ResourceState]) (
				wrappers.ResourceWrapper[schema.Workspace, schema.RegionalResourceMetadata, schema.WorkspaceSpec, schema.WorkspaceStatus], error,
			) {
				resp, err := api.GetWorkspaceUntilState(ctx, tref, config)
				return wrappers.NewWorkspaceWrapper(resp), err
			},
			expectedLabels:         responseExpects.Labels,
			expectedMetadata:       responseExpects.Metadata,
			verifyMetadataFunc:     configurator.suite.VerifyRegionalResourceMetadataStep,
			expectedResourceStatus: responseExpects.ResourceStatus,
		},
	)
}

func (configurator *StepsConfigurator) ListWorkspaceV1Step(
	stepName string, api secapi.WorkspaceV1, tref secapi.TenantReference, opts *secapi.ListOptions,
) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetStorageWorkspaceV1StepParams(sCtx, "ListWorkspace", string(tref.Tenant))
		iter, err := api.ListWorkspacesWithOptions(configurator.t.Context(), secapi.TenantPath{Tenant: tref.Tenant}, opts)
		requireNoError(sCtx, err)

		verifyIterListStep(sCtx, configurator.t, *iter)
	})
}

func (configurator *StepsConfigurator) WatchWorkspaceUntilDeletedV1Step(stepName string, api secapi.WorkspaceV1, tref secapi.TenantReference) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		watchTenantResourceUntilDeletedStep(configurator.t, configurator.suite,
			watchTenantResourceUntilDeletedParams{
				watchResourceUntilDeletedParams: watchResourceUntilDeletedParams[secapi.TenantReference]{
					reference: tref,
					getErrorFunc: func(ctx context.Context, tref secapi.TenantReference, config secapi.ResourceObserverConfig) error {
						return api.WatchWorkspaceUntilDeleted(configurator.t.Context(), tref, config)
					},
				},
				stepName:       stepName,
				stepParamsFunc: configurator.suite.SetWorkspaceV1StepParams,
				operationName:  constants.GetWorkspaceOperation,
			},
		)
	})
}

func (configurator *StepsConfigurator) DeleteWorkspaceV1Step(stepName string, api secapi.WorkspaceV1, resource *schema.Workspace) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetWorkspaceV1StepParams(sCtx, "DeleteWorkspace")

		err := api.DeleteWorkspace(configurator.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

func (configurator *StepsConfigurator) CreateOrUpdateWorkspaceExpectViolationV1Step(stepName string, api secapi.WorkspaceV1, resource *schema.Workspace) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetWorkspaceV1StepParams(sCtx, constants.CreateOrUpdateWorkspaceOperation)

		_, err := api.CreateOrUpdateWorkspace(configurator.t.Context(), resource)
		requireError(sCtx, err)
	})
}
