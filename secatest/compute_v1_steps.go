package secatest

import (
	"context"
	"net/http"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// Instance

func (suite *testSuite) createOrUpdateInstanceV1Step(stepName string, t provider.T, api *secapi.ComputeV1, resource *schema.Instance,
	responseExpects responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec],
) {
	responseExpects.metadata.Verb = http.MethodPut
	createOrUpdateWorkspaceResourceStep(t, suite,
		createOrUpdateWorkspaceResourceParams[schema.Instance, schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			stepName:       stepName,
			stepParamsFunc: suite.setComputeV1StepParams,
			operationName:  "CreateOrUpdateInstance",
			workspace:      resource.Metadata.Workspace,
			resource:       resource,
			createOrUpdateFunc: func(context.Context, *schema.Instance) (*stepFuncResponse[schema.Instance, schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec], error) {
				resp, err := api.CreateOrUpdateInstance(t.Context(), resource)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.metadata,
			verifyMetadataFunc:    suite.verifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.spec,
			verifySpecFunc:        suite.verifyInstanceSpecStep,
			expectedResourceState: responseExpects.resourceState,
		},
	)
}

func (suite *testSuite) getInstanceV1Step(stepName string, t provider.T, api *secapi.ComputeV1, wref secapi.WorkspaceReference,
	responseExpects responseExpects[schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec],
) *schema.Instance {
	responseExpects.metadata.Verb = http.MethodGet
	return getWorkspaceResourceStep(t, suite,
		getWorkspaceResourceParams[schema.Instance, schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec]{
			stepName:       stepName,
			stepParamsFunc: suite.setComputeV1StepParams,
			operationName:  "GetInstance",
			wref:           wref,
			getFunc: func(ctx context.Context, wref secapi.WorkspaceReference, config secapi.ResourceObserverConfig[schema.ResourceState]) (*stepFuncResponse[schema.Instance, schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec], error) {
				resp, err := api.GetInstanceUntilState(t.Context(), wref, config)
				return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
			},
			expectedMetadata:      responseExpects.metadata,
			verifyMetadataFunc:    suite.verifyRegionalWorkspaceResourceMetadataStep,
			expectedSpec:          responseExpects.spec,
			verifySpecFunc:        suite.verifyInstanceSpecStep,
			expectedResourceState: responseExpects.resourceState,
		},
	)
}

func (suite *testSuite) getInstanceWithErrorV1Step(stepName string, t provider.T, api *secapi.ComputeV1, wref secapi.WorkspaceReference, expectedError error) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "GetInstance", string(wref.Workspace))

		_, err := api.GetInstance(t.Context(), wref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) startInstanceV1Step(stepName string, t provider.T, api *secapi.ComputeV1, resource *schema.Instance) {
	var err error

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "StartInstance", resource.Metadata.Workspace)

		err = api.StartInstance(t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

func (suite *testSuite) stopInstanceV1Step(stepName string, t provider.T, api *secapi.ComputeV1, resource *schema.Instance) {
	var err error

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "StopInstance", resource.Metadata.Workspace)

		err = api.StopInstance(t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

func (suite *testSuite) restartInstanceV1Step(stepName string, t provider.T, api *secapi.ComputeV1, resource *schema.Instance) {
	var err error

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "RestartInstance", resource.Metadata.Workspace)

		err = api.RestartInstance(t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

func (suite *testSuite) deleteInstanceV1Step(stepName string, t provider.T, api *secapi.ComputeV1, resource *schema.Instance) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "DeleteInstance", resource.Metadata.Workspace)

		err := api.DeleteInstance(t.Context(), resource)
		requireNoError(sCtx, err)
	})
}

func (suite *testSuite) getListInstanceV1Step(
	stepName string,
	t provider.T,
	api *secapi.ComputeV1,
	wref secapi.WorkspaceReference,
	opts *secapi.ListOptions,
) []*schema.Instance {

	var resp []*schema.Instance

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "ListInstances with parameters", wref.Name)
		var iter *secapi.Iterator[schema.Instance]
		var err error
		if opts != nil {
			iter, err = api.ListInstancesWithFilters(context.Background(), wref.Tenant, wref.Workspace, opts)
		} else {
			iter, err = api.ListInstances(context.Background(), wref.Tenant, wref.Workspace)
		}
		requireNoError(sCtx, err)
		resp, err := iter.All(t.Context())
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)
		requireLenResponse(sCtx, len(resp))
	})
	return resp
}

func (suite *testSuite) getListSkusV1Step(
	stepName string,
	t provider.T,
	api *secapi.ComputeV1,
	tref secapi.TenantReference,
	opts *secapi.ListOptions,
) []*schema.InstanceSku {
	var resp []*schema.InstanceSku

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "ListSkus", tref.Name)

		var iter *secapi.Iterator[schema.InstanceSku]
		var err error
		if opts != nil {
			iter, err = api.ListSkusWithFilters(t.Context(), tref.Tenant, opts)
		} else {
			iter, err = api.ListSkus(t.Context(), tref.Tenant)
		}
		requireNoError(sCtx, err)

		// Iterate through all items
		resp, err := iter.All(t.Context())
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)
		requireLenResponse(sCtx, len(resp))
	})
	return resp
}
