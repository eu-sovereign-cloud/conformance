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

func (builder *Builder) CreateOrUpdateInstanceV1Step(stepName string, api *secapi.ComputeV1, resource *schema.Instance,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec],
) {
	responseExpects.Metadata.Verb = http.MethodPut
	createOrUpdateWorkspaceResourceStep(builder.t, builder.suite,
		createOrUpdateWorkspaceResourceParams[schema.Instance, schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			stepName:       stepName,
			stepParamsFunc: builder.suite.SetComputeV1StepParams,
			operationName:  "CreateOrUpdateInstance",
			workspace:      resource.Metadata.Workspace,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.Instance) (*stepFuncResponse[schema.Instance, schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec], error) {
				resp, err := api.CreateOrUpdateInstance(builder.t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    builder.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        builder.suite.VerifyInstanceSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (builder *Builder) GetInstanceV1Step(stepName string, api *secapi.ComputeV1, wref secapi.WorkspaceReference,
	responseExpects ResponseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec],
) *schema.Instance {
	responseExpects.Metadata.Verb = http.MethodGet
	return getWorkspaceResourceStep(builder.t, builder.suite,
		getWorkspaceResourceParams[schema.Instance, schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			stepName:       stepName,
			stepParamsFunc: builder.suite.SetComputeV1StepParams,
			operationName:  "GetInstance",
			wref:           wref,
			getFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[schema.Instance, schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec], error) {
				resp, err := api.GetInstanceUntilState(builder.t.Context(), wref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.Metadata,
			verifyMetadataFunc:    builder.suite.VerifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.Spec,
			verifySpecFunc:        builder.suite.VerifyInstanceSpecStep,
			expectedResourceState: responseExpects.ResourceState,
		},
	)
}

func (builder *Builder) GetInstanceWithErrorV1Step(stepName string, api *secapi.ComputeV1, wref secapi.WorkspaceReference, expectedError error) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetComputeV1StepParams(sCtx, "GetInstance", string(wref.Workspace))

		_, err := api.GetInstance(builder.t.Context(), wref)
		requireError(sCtx, err, expectedError)
	})
}

func (builder *Builder) StartInstanceV1Step(stepName string, api *secapi.ComputeV1, resource *schema.Instance) {
	var err error

	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetComputeV1StepParams(sCtx, "StartInstance", resource.Metadata.Workspace)

		err = api.StartInstance(builder.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

func (builder *Builder) StopInstanceV1Step(stepName string, api *secapi.ComputeV1, resource *schema.Instance) {
	var err error

	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetComputeV1StepParams(sCtx, "StopInstance", resource.Metadata.Workspace)

		err = api.StopInstance(builder.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

func (builder *Builder) RestartInstanceV1Step(stepName string, api *secapi.ComputeV1, resource *schema.Instance) {
	var err error

	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetComputeV1StepParams(sCtx, "RestartInstance", resource.Metadata.Workspace)

		err = api.RestartInstance(builder.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

func (builder *Builder) DeleteInstanceV1Step(stepName string, api *secapi.ComputeV1, resource *schema.Instance) {
	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetComputeV1StepParams(sCtx, "DeleteInstance", resource.Metadata.Workspace)

		err := api.DeleteInstance(builder.t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

func (builder *Builder) GetListInstanceV1Step(
	stepName string,
	api *secapi.ComputeV1,
	wref secapi.WorkspaceReference,
	opts *secapi.ListOptions,
) []*schema.Instance {
	var resp []*schema.Instance

	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetComputeV1StepParams(sCtx, "ListInstances with parameters", wref.Name)
		var iter *secapi.Iterator[schema.Instance]
		var err error
		if opts != nil {
			iter, err = api.ListInstancesWithFilters(context.Background(), wref.Tenant, wref.Workspace, opts)
		} else {
			iter, err = api.ListInstances(context.Background(), wref.Tenant, wref.Workspace)
		}
		requireNoError(sCtx, err)
		resp, err := iter.All(builder.t.Context())
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)
		requireLenResponse(sCtx, len(resp))
	})
	return resp
}

func (builder *Builder) GetListSkusV1Step(
	stepName string,
	api *secapi.ComputeV1,
	tref secapi.TenantReference,
	opts *secapi.ListOptions,
) []*schema.InstanceSku {
	var resp []*schema.InstanceSku

	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetComputeV1StepParams(sCtx, "ListSkus", tref.Name)

		var iter *secapi.Iterator[schema.InstanceSku]
		var err error
		if opts != nil {
			iter, err = api.ListSkusWithFilters(builder.t.Context(), tref.Tenant, opts)
		} else {
			iter, err = api.ListSkus(builder.t.Context(), tref.Tenant)
		}
		requireNoError(sCtx, err)

		// Iterate through all items
		resp, err := iter.All(builder.t.Context())
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)
		requireLenResponse(sCtx, len(resp))
	})
	return resp
}
