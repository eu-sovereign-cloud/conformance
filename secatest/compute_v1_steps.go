package secatest

import (
	"context"
	"net/http"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
)

// Instance

func (suite *testSuite) createOrUpdateInstanceV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.ComputeV1,
	resource *schema.Instance,
	expectedMeta *schema.RegionalWorkspaceResourceMetadata,
	expectedSpec *schema.InstanceSpec,
	expectedState schema.ResourceState,
) {
	expectedMeta.Verb = http.MethodPut
	createOrUpdateWorkspaceResourceStep(
		t,
		ctx,
		suite,
		stepName,
		suite.setComputeV1StepParams,
		"CreateOrUpdateInstance",
		resource.Metadata.Workspace,
		resource,
		func(context.Context, *schema.Instance) (*stepFuncResponse[schema.Instance, schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec], error) {
			resp, err := api.CreateOrUpdateInstance(ctx, resource)
			return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
		},
		nil,
		expectedMeta,
		suite.verifyRegionalWorkspaceResourceMetadataStep,
		expectedSpec,
		suite.verifyInstanceSpecStep,
		expectedState,
	)
}

func (suite *testSuite) getInstanceV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.ComputeV1,
	wref secapi.WorkspaceReference,
	expectedMeta *schema.RegionalWorkspaceResourceMetadata,
	expectedSpec *schema.InstanceSpec,
	expectedState schema.ResourceState,
) *schema.Instance {
	expectedMeta.Verb = http.MethodGet
	return getWorkspaceResourceStep(
		t,
		ctx,
		suite,
		stepName,
		suite.setComputeV1StepParams,
		"GetInstance",
		wref,
		func(context.Context, secapi.WorkspaceReference) (*stepFuncResponse[schema.Instance, schema.RegionalWorkspaceResourceMetadata, schema.InstanceSpec], error) {
			resp, err := api.GetInstance(ctx, wref)
			return newStepFuncResponse(resp, resp.Labels, resp.Metadata, resp.Spec, resp.Status.State), err
		},
		nil,
		expectedMeta,
		suite.verifyRegionalWorkspaceResourceMetadataStep,
		expectedSpec,
		suite.verifyInstanceSpecStep,
		expectedState,
	)
}

func (suite *testSuite) getInstanceWithErrorV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.ComputeV1,
	wref secapi.WorkspaceReference,
	expectedError error,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "GetInstance", string(wref.Workspace))

		_, err := api.GetInstance(ctx, wref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) startInstanceV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.ComputeV1, resource *schema.Instance) {
	var err error

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "StartInstance", resource.Metadata.Workspace)

		err = api.StartInstance(ctx, resource)
		requireNoError(sCtx, err)
	})
}

func (suite *testSuite) stopInstanceV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.ComputeV1, resource *schema.Instance) {
	var err error

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "StopInstance", resource.Metadata.Workspace)

		err = api.StopInstance(ctx, resource)
		requireNoError(sCtx, err)
	})
}

func (suite *testSuite) restartInstanceV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.ComputeV1, resource *schema.Instance) {
	var err error

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "RestartInstance", resource.Metadata.Workspace)

		err = api.RestartInstance(ctx, resource)
		requireNoError(sCtx, err)
	})
}

func (suite *testSuite) deleteInstanceV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.ComputeV1, resource *schema.Instance) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "DeleteInstance", resource.Metadata.Workspace)

		err := api.DeleteInstance(ctx, resource)
		requireNoError(sCtx, err)
	})
}
