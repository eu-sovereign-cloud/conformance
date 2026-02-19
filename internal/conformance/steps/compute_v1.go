//nolint:dupl
package steps

import (
	"context"
	"net/http"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// Instance

func (configurator *StepsConfigurator) CreateOrUpdateInstanceV1Step(stepName string, api secapi.ComputeV1, resource *schema.Instance,
	responseExpects StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	configurator.logStepName(stepName)
	createOrUpdateWorkspaceResourceStep(configurator.t, configurator.suite,
		createOrUpdateWorkspaceResourceParams[schema.Instance, schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus]{
			createOrUpdateResourceParams: createOrUpdateResourceParams[schema.Instance, schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus]{
				resource: resource,
				createOrUpdateFunc: func(context.Context, *schema.Instance) (
					*createOrUpdateStepFuncResponse[schema.Instance, schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus], error,
				) {
					if resp, err := api.CreateOrUpdateInstance(configurator.t.Context(), resource); err == nil {
						return newCreateOrUpdateStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), nil
					} else {
						return nil, err
					}
				},
				expectedMetadata:      responseExpects.Metadata,
				verifyMetadataFunc:    configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
				expectedSpec:          responseExpects.Spec,
				verifySpecFunc:        configurator.suite.VerifyInstanceSpecStep,
				expectedResourceState: responseExpects.ResourceState,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetComputeV1StepParams,
			operationName:  "CreateOrUpdateInstance",
			workspace:      resource.Metadata.Workspace,
		},
	)
}

func (configurator *StepsConfigurator) GetInstanceV1Step(stepName string, api secapi.ComputeV1, wref secapi.WorkspaceReference,
	responseExpects StepResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec],
) *schema.Instance {
	responseExpects.Metadata.Verb = http.MethodGet
	configurator.logStepName(stepName)
	return getWorkspaceResourceStep(configurator.t, configurator.suite,
		getWorkspaceResourceParams[schema.Instance, schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus]{
			getResourceWithObserverParams: getResourceWithObserverParams[schema.Instance, schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus, secapi.WorkspaceReference, schema.ResourceState]{
				reference: wref,
				getFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (
					*getStepFuncResponse[schema.Instance, schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec, schema.InstanceStatus], error,
				) {
					if resp, err := api.GetInstanceUntilState(configurator.t.Context(), wref, config); err == nil {
						return newGetStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status), nil
					} else {
						return nil, err
					}
				},
				expectedMetadata:      responseExpects.Metadata,
				verifyMetadataFunc:    configurator.suite.VerifyRegionalWorkspaceResourceMetadataStep,
				expectedSpec:          responseExpects.Spec,
				verifySpecFunc:        configurator.suite.VerifyInstanceSpecStep,
				expectedResourceState: responseExpects.ResourceState,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetComputeV1StepParams,
			operationName:  "GetInstance",
		},
	)
}

func (configurator *StepsConfigurator) GetInstanceWithErrorV1Step(stepName string, api secapi.ComputeV1, wref secapi.WorkspaceReference, expectedError error) {
	configurator.logStepName(stepName)
	getWorkspaceResourceWithErrorStep(configurator.t,
		getWorkspaceResourceWithErrorParams{
			getResourceWithErrorParams: getResourceWithErrorParams[secapi.WorkspaceReference]{
				reference: wref,
				getFunc: func(ctx context.Context, wref secapi.WorkspaceReference) error {
					_, err := api.GetInstance(ctx, wref)
					return err
				},
				expectedError: expectedError,
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetComputeV1StepParams,
			operationName:  "GetInstance",
		},
	)
}

func (configurator *StepsConfigurator) StartInstanceV1Step(stepName string, api secapi.ComputeV1, resource *schema.Instance) {
	configurator.logStepName(stepName)
	actionWorkspaceResourceStep(configurator.t,
		actionWorkspaceResourceParams[schema.Instance]{
			actionResourceParams: actionResourceParams[schema.Instance]{
				resource: resource,
				actionFunc: func(ctx context.Context, r *schema.Instance) error {
					return api.StartInstance(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetComputeV1StepParams,
			operationName:  "StartInstance",
			workspace:      resource.Metadata.Workspace,
		},
	)
}

func (configurator *StepsConfigurator) StopInstanceV1Step(stepName string, api secapi.ComputeV1, resource *schema.Instance) {
	configurator.logStepName(stepName)
	actionWorkspaceResourceStep(configurator.t,
		actionWorkspaceResourceParams[schema.Instance]{
			actionResourceParams: actionResourceParams[schema.Instance]{
				resource: resource,
				actionFunc: func(ctx context.Context, r *schema.Instance) error {
					return api.StopInstance(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetComputeV1StepParams,
			operationName:  "StopInstance",
			workspace:      resource.Metadata.Workspace,
		},
	)
}

func (configurator *StepsConfigurator) RestartInstanceV1Step(stepName string, api secapi.ComputeV1, resource *schema.Instance) {
	configurator.logStepName(stepName)
	actionWorkspaceResourceStep(configurator.t,
		actionWorkspaceResourceParams[schema.Instance]{
			actionResourceParams: actionResourceParams[schema.Instance]{
				resource: resource,
				actionFunc: func(ctx context.Context, r *schema.Instance) error {
					return api.RestartInstance(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetComputeV1StepParams,
			operationName:  "RestartInstance",
			workspace:      resource.Metadata.Workspace,
		},
	)
}

func (configurator *StepsConfigurator) DeleteInstanceV1Step(stepName string, api secapi.ComputeV1, resource *schema.Instance) {
	configurator.logStepName(stepName)
	deleteWorkspaceResourceStep(configurator.t,
		deleteWorkspaceResourceParams[schema.Instance]{
			deleteResourceParams: deleteResourceParams[schema.Instance]{
				resource: resource,
				deleteFunc: func(ctx context.Context, r *schema.Instance) error {
					return api.DeleteInstance(ctx, r)
				},
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetComputeV1StepParams,
			operationName:  "DeleteInstance",
			workspace:      resource.Metadata.Workspace,
		},
	)
}

func (configurator *StepsConfigurator) ListInstanceV1Step(stepName string, api secapi.ComputeV1, wref secapi.WorkspaceReference, opts *secapi.ListOptions) []*schema.Instance {
	var resp []*schema.Instance
	configurator.logStepName(stepName)
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
		requireNotEmptyResponse(sCtx, resp)
	})
	return resp
}

func (configurator *StepsConfigurator) ListSkusV1Step(
	stepName string,
	api secapi.ComputeV1,
	tref secapi.TenantReference,
	opts *secapi.ListOptions,
) []*schema.InstanceSku {
	var resp []*schema.InstanceSku
	configurator.logStepName(stepName)
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
		requireNotEmptyResponse(sCtx, resp)
	})
	return resp
}
