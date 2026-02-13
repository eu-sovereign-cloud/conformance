//nolint:dupl
package steps

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// Instance

func (configurator *StepsConfigurator) CreateOrUpdateInstanceV1Step(stepName string, api *secapi.ComputeV1, resource *schema.Instance,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	createOrUpdateWorkspaceResourceStep(configurator.t, configurator.suite,
		createOrUpdateWorkspaceResourceParams[schema.Instance, schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetComputeV1StepParams,
			operationName:  "CreateOrUpdateInstance",
			workspace:      resource.Metadata.Workspace,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.Instance) (
				*stepFuncResponse[schema.Instance, schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus], error,
			) {
				resp, err := api.CreateOrUpdateInstance(configurator.t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        configurator.suite.VerifyInstanceSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (configurator *StepsConfigurator) GetInstanceV1Step(stepName string, api *secapi.ComputeV1, wref secapi.WorkspaceReference,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec],
) *schema.Instance {
	responseExpects.Metadata.Verb = http.MethodGet
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	return getWorkspaceResourceStep(configurator.t, configurator.suite,
		getWorkspaceResourceParams[schema.Instance, schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetComputeV1StepParams,
			operationName:  "GetInstance",
			wref:           wref,
			getFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (
				*stepFuncResponse[schema.Instance, schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus], error,
			) {
				resp, err := api.GetInstanceUntilState(configurator.t.Context(), wref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        configurator.suite.VerifyInstanceSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (configurator *StepsConfigurator) GetInstanceWithErrorV1Step(stepName string, api *secapi.ComputeV1, wref secapi.WorkspaceReference, expectedError error) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetComputeV1StepParams(sCtx, "GetInstance", string(wref.Workspace))

		_, err := api.GetInstance(configurator.t.Context(), wref)
		requireError(sCtx, err, expectedError)
	})
}

func (configurator *StepsConfigurator) StartInstanceV1Step(stepName string, api *secapi.ComputeV1, resource *schema.Instance) {
	var err error
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetComputeV1StepParams(sCtx, "StartInstance", resource.Metadata.Workspace)

		err = api.StartInstance(configurator.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

func (configurator *StepsConfigurator) StopInstanceV1Step(stepName string, api *secapi.ComputeV1, resource *schema.Instance) {
	var err error
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetComputeV1StepParams(sCtx, "StopInstance", resource.Metadata.Workspace)

		err = api.StopInstance(configurator.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

func (configurator *StepsConfigurator) RestartInstanceV1Step(stepName string, api *secapi.ComputeV1, resource *schema.Instance) {
	var err error
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetComputeV1StepParams(sCtx, "RestartInstance", resource.Metadata.Workspace)

		err = api.RestartInstance(configurator.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

func (configurator *StepsConfigurator) DeleteInstanceV1Step(stepName string, api *secapi.ComputeV1, resource *schema.Instance) {
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetComputeV1StepParams(sCtx, "DeleteInstance", resource.Metadata.Workspace)

		err := api.DeleteInstance(configurator.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

func (configurator *StepsConfigurator) GetListInstanceV1Step(
	stepName string,
	api *secapi.ComputeV1,
	wref secapi.WorkspaceReference,
	opts *secapi.ListOptions,
) []*schema.Instance {
	var resp []*schema.Instance
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetComputeV1StepParams(sCtx, "ListInstances with parameters", wref.Name)
		var iter *secapi.Iterator[schema.Instance]
		var err error
		if opts != nil {
			iter, err = api.ListInstancesWithFilters(context.Background(), wref.Tenant, wref.Workspace, opts)
		} else {
			iter, err = api.ListInstances(context.Background(), wref.Tenant, wref.Workspace)
		}
		requireNoError(sCtx, err)
		resp, err := iter.All(configurator.t.Context())
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)
		requireLenResponse(sCtx, len(resp))
	})
	return resp
}

func (configurator *StepsConfigurator) GetListSkusV1Step(
	stepName string,
	api *secapi.ComputeV1,
	tref secapi.TenantReference,
	opts *secapi.ListOptions,
) []*schema.InstanceSku {
	var resp []*schema.InstanceSku
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetComputeV1StepParams(sCtx, "ListSkus", tref.Name)

		var iter *secapi.Iterator[schema.InstanceSku]
		var err error
		if opts != nil {
			iter, err = api.ListSkusWithFilters(configurator.t.Context(), tref.Tenant, opts)
		} else {
			iter, err = api.ListSkus(configurator.t.Context(), tref.Tenant)
		}
		requireNoError(sCtx, err)

		// Iterate through all items
		resp, err := iter.All(configurator.t.Context())
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)
		requireLenResponse(sCtx, len(resp))
	})
	return resp
}
