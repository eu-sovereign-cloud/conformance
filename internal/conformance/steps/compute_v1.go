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

// Instance

func (configurator *StepsConfigurator) CreateOrUpdateInstanceV1Step(stepName string, api secapi.ComputeV1, resource *schema.Instance,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	createOrUpdateWorkspaceResourceStep(configurator.t, configurator.suite,
		createOrUpdateWorkspaceResourceParams[schema.Instance, schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetComputeV1StepParams,
			operationName:  constants.CreateOrUpdateInstanceOperation,
			workspace:      resource.Metadata.Workspace,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.Instance) (
				wrappers.ResourceWrapper[schema.Instance, schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus], error,
			) {
				resp, err := api.CreateOrUpdateInstance(configurator.t.Context(), resource)
				return wrappers.NewInstanceWrapper(resp), err
			},
			expectedMetadata:       responseExpects.Metadata,
			verifyMetadataFunc:     configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:           responseExpects.Spec,
			verifySpecFunc:         configurator.suite.VerifyInstanceSpecStep,
			expectedResourceStates: responseExpects.ResourceStates,
		},
	)
}

func (configurator *StepsConfigurator) GetInstanceV1Step(stepName string, api secapi.ComputeV1, wref secapi.WorkspaceReference,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec],
) *schema.Instance {
	responseExpects.Metadata.Verb = http.MethodGet
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	return getWorkspaceResourceStep(configurator.t, configurator.suite,
		getWorkspaceResourceParams[schema.Instance, schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetComputeV1StepParams,
			operationName:  constants.GetInstanceOperation,
			wref:           wref,
			getValueFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverUntilValueConfig[schema.ResourceState]) (
				wrappers.ResourceWrapper[schema.Instance, schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus], error,
			) {
				resp, err := api.GetInstanceUntilState(configurator.t.Context(), wref, config)
				return wrappers.NewInstanceWrapper(resp), err
			},
			expectedMetadata:       responseExpects.Metadata,
			verifyMetadataFunc:     configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:           responseExpects.Spec,
			verifySpecFunc:         configurator.suite.VerifyInstanceSpecStep,
			expectedResourceStates: responseExpects.ResourceStates,
		},
	)
}

func (configurator *StepsConfigurator) WatchInstanceUntilDeletedV1Step(stepName string, api secapi.ComputeV1, tref secapi.WorkspaceReference) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		watchWorkspaceResourceUntilDeletedStep(configurator.t, configurator.suite,
			watchWorkspaceResourceUntilDeletedParams{
				watchResourceUntilDeletedParams: watchResourceUntilDeletedParams[secapi.WorkspaceReference]{
					reference: tref,
					getErrorFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig) error {
						return api.WatchInstanceUntilDeleted(configurator.t.Context(), wref, config)
					},
				},
				stepName:       stepName,
				stepParamsFunc: configurator.suite.SetComputeV1StepParams,
				operationName:  constants.GetInstanceOperation,
			},
		)
	})
}

func (configurator *StepsConfigurator) StartInstanceV1Step(stepName string, api secapi.ComputeV1, resource *schema.Instance) {
	var err error
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetComputeV1StepParams(sCtx, "StartInstance", resource.Metadata.Workspace)

		err = api.StartInstance(configurator.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

func (configurator *StepsConfigurator) StopInstanceV1Step(stepName string, api secapi.ComputeV1, resource *schema.Instance) {
	var err error
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetComputeV1StepParams(sCtx, "StopInstance", resource.Metadata.Workspace)

		err = api.StopInstance(configurator.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

func (configurator *StepsConfigurator) RestartInstanceV1Step(stepName string, api secapi.ComputeV1, resource *schema.Instance) {
	var err error
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetComputeV1StepParams(sCtx, "RestartInstance", resource.Metadata.Workspace)

		err = api.RestartInstance(configurator.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

func (configurator *StepsConfigurator) DeleteInstanceV1Step(stepName string, api secapi.ComputeV1, resource *schema.Instance) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetComputeV1StepParams(sCtx, "DeleteInstance", resource.Metadata.Workspace)

		err := api.DeleteInstance(configurator.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

func (configurator *StepsConfigurator) ListInstanceV1Step(
	stepName string, api secapi.ComputeV1, wref secapi.WorkspaceReference, opts *secapi.FilterOptions,
) []*schema.Instance {
	var resp []*schema.Instance
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetComputeV1StepParams(sCtx, "ListInstances with parameters", wref.Name)

		iter, err := api.ListInstances(context.Background(), secapi.WorkspaceFilter{Tenant: wref.Tenant, Workspace: wref.Workspace, Options: opts})
		requireNoError(sCtx, err)

		// Iterate through all items
		resp, err := iter.All(configurator.t.Context())
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)
		requireLenResponse(sCtx, len(resp))
	})
	return resp
}

func (configurator *StepsConfigurator) ListSkusV1Step(
	stepName string, api secapi.ComputeV1, tref secapi.TenantReference, opts *secapi.FilterOptions,
) []*schema.InstanceSku {
	var resp []*schema.InstanceSku
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetComputeV1StepParams(sCtx, "ListSkus", tref.Name)

		iter, err := api.ListSkus(configurator.t.Context(), secapi.TenantFilter{Tenant: tref.Tenant, Options: opts})
		requireNoError(sCtx, err)

		// Iterate through all items
		resp, err := iter.All(configurator.t.Context())
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)
		requireLenResponse(sCtx, len(resp))
	})
	return resp
}
