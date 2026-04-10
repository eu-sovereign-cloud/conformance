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

// Sku

func (configurator *StepsConfigurator) ListSkusV1Step(stepName string, api secapi.ComputeV1, tpath secapi.TenantPath, opts *secapi.ListOptions) {
	listTenantResourcesStep(configurator.t, configurator.suite,
		listTenantResourcesParams[schema.InstanceSku, schema.SkuResourceMetadata]{
			listResourcesParams: listResourcesParams[schema.InstanceSku, schema.SkuResourceMetadata, secapi.TenantPath]{
				path: tpath, listOptions: opts,
				listFunc: func(ctx context.Context, path secapi.TenantPath, options *secapi.ListOptions) (*secapi.Iterator[schema.InstanceSku], error) {
					return api.ListSkusWithOptions(ctx, path, options)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetComputeSkuV1StepParams,
			operationName:  constants.ListSkusOperation,
		},
	)
}

// Instance

func (configurator *StepsConfigurator) CreateOrUpdateInstanceV1Step(stepName string, stepCreator StepCreator, api secapi.ComputeV1, resource *schema.Instance,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	createOrUpdateWorkspaceResourceStep(configurator.t.Context(), configurator.suite, stepCreator,
		createOrUpdateWorkspaceResourceParams[schema.Instance, schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetComputeV1StepParams,
			operationName:  constants.CreateOrUpdateInstanceOperation,
			workspace:      secapi.WorkspaceID(resource.Metadata.Workspace),
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.Instance) (
				wrappers.ResourceWrapper[schema.Instance, schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus], error,
			) {
				resp, err := api.CreateOrUpdateInstance(configurator.t.Context(), resource)
				return wrappers.NewInstanceWrapper(resp), err
			},
			expectedMetadata:       responseExpects.Metadata,
			verifyMetadataFunc:     configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedLabels:         responseExpects.Labels,
			expectedAnnotations:    responseExpects.Annotations,
			expectedExtensions:     responseExpects.Extensions,
			expectedSpec:           responseExpects.Spec,
			verifySpecFunc:         configurator.suite.VerifyInstanceSpecStep,
			expectedResourceStates: responseExpects.ResourceStates,
		},
	)
}

func (configurator *StepsConfigurator) ListInstanceV1Step(stepName string, api secapi.ComputeV1, wpath secapi.WorkspacePath, opts *secapi.ListOptions) {
	listWorkspaceResourcesStep(configurator.t, configurator.suite,
		listWorkspaceResourcesParams[schema.Instance, schema.RegionalWorkspaceResourceMetadata]{
			listResourcesParams: listResourcesParams[schema.Instance, schema.RegionalWorkspaceResourceMetadata, secapi.WorkspacePath]{
				path: wpath, listOptions: opts,
				listFunc: func(ctx context.Context, path secapi.WorkspacePath, options *secapi.ListOptions) (*secapi.Iterator[schema.Instance], error) {
					return api.ListInstancesWithOptions(ctx, path, options)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetComputeV1StepParams,
			operationName:  constants.ListInstancesOperation,
			workspace:      wpath.Workspace,
		},
	)
}

func (configurator *StepsConfigurator) GetInstanceV1Step(stepName string, api secapi.ComputeV1, wref secapi.WorkspaceReference,
	responseExpects ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus],
) *schema.Instance {
	responseExpects.Metadata.Verb = http.MethodGet
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
			expectedResourceStatus: responseExpects.ResourceStatus,
		},
	)
}

func (configurator *StepsConfigurator) GetInstancePowerStateV1Step(stepName string, api secapi.ComputeV1, wref secapi.WorkspaceReference,
	responseExpects ResponseExpectsWithCondition[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus],
) *schema.Instance {
	responseExpects.Metadata.Verb = http.MethodGet
	return getWorkspaceInstanceResourceStep(configurator.t, configurator.suite,
		getWorkspaceInstanceResourceParams[schema.Instance, schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetComputeV1StepParams,
			operationName:  constants.GetInstanceOperation,
			wref:           wref,
			getValueFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverUntilValueConfig[schema.InstanceStatusPowerState]) (
				wrappers.ResourceWrapper[schema.Instance, schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus], error,
			) {
				resp, err := api.GetInstanceUntilPowerState(configurator.t.Context(), wref, config)
				return wrappers.NewInstanceWrapper(resp), err
			},
			expectedMetadata:   responseExpects.Metadata,
			verifyMetadataFunc: configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:       responseExpects.Spec,
			verifySpecFunc:     configurator.suite.VerifyInstanceSpecStep,
			expectedPowerState: responseExpects.ResourceStatus.PowerState,
		},
	)
}

func (configurator *StepsConfigurator) WatchInstanceUntilDeletedV1Step(stepName string, stepCreator StepCreator, api secapi.ComputeV1, wref secapi.WorkspaceReference) {
	watchWorkspaceResourceUntilDeletedStep(configurator.t.Context(), configurator.suite, stepCreator,
		watchWorkspaceResourceUntilDeletedParams{
			watchResourceUntilDeletedParams: watchResourceUntilDeletedParams[secapi.WorkspaceReference]{
				reference: wref,
				getErrorFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig) error {
					return api.WatchInstanceUntilDeleted(configurator.t.Context(), wref, config)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetComputeV1StepParams,
			operationName:  constants.GetInstanceOperation,
		},
	)
}

func (configurator *StepsConfigurator) StartInstanceV1Step(stepName string, api secapi.ComputeV1, resource *schema.Instance) {
	actionWorkspaceResourceStep(configurator.t, configurator.suite,
		actionWorkspaceResourceParams[schema.Instance]{
			actionResourceParams: actionResourceParams[schema.Instance]{
				resource: resource,
				actionFunc: func(ctx context.Context, r *schema.Instance) error {
					return api.StartInstance(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetComputeV1StepParams,
			operationName:  constants.StartInstanceOperation,
			workspace:      secapi.WorkspaceID(resource.Metadata.Workspace),
		},
	)
}

func (configurator *StepsConfigurator) StopInstanceV1Step(stepName string, api secapi.ComputeV1, resource *schema.Instance) {
	actionWorkspaceResourceStep(configurator.t, configurator.suite,
		actionWorkspaceResourceParams[schema.Instance]{
			actionResourceParams: actionResourceParams[schema.Instance]{
				resource: resource,
				actionFunc: func(ctx context.Context, r *schema.Instance) error {
					return api.StopInstance(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetComputeV1StepParams,
			operationName:  constants.StopInstanceOperation,
			workspace:      secapi.WorkspaceID(resource.Metadata.Workspace),
		},
	)
}

func (configurator *StepsConfigurator) RestartInstanceV1Step(stepName string, api secapi.ComputeV1, resource *schema.Instance) {
	actionWorkspaceResourceStep(configurator.t, configurator.suite,
		actionWorkspaceResourceParams[schema.Instance]{
			actionResourceParams: actionResourceParams[schema.Instance]{
				resource: resource,
				actionFunc: func(ctx context.Context, r *schema.Instance) error {
					return api.RestartInstance(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetComputeV1StepParams,
			operationName:  constants.RestartInstanceOperation,
			workspace:      secapi.WorkspaceID(resource.Metadata.Workspace),
		},
	)
}

func (configurator *StepsConfigurator) DeleteInstanceV1Step(stepName string, stepCreator StepCreator, api secapi.ComputeV1, resource *schema.Instance) {
	deleteWorkspaceResourceStep(configurator.t.Context(), configurator.suite, stepCreator,
		deleteWorkspaceResourceParams[schema.Instance]{
			deleteResourceParams: deleteResourceParams[schema.Instance]{
				resource: resource,
				deleteFunc: func(ctx context.Context, r *schema.Instance) error {
					return api.DeleteInstance(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetComputeV1StepParams,
			operationName:  constants.DeleteInstanceOperation,
			workspace:      secapi.WorkspaceID(resource.Metadata.Workspace),
		},
	)
}
