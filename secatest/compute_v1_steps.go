package secatest

import (
	"context"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/secalib"
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
	expectedStatusState string,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "CreateOrUpdateInstance", resource.Metadata.Workspace)

		resp, err := api.CreateOrUpdateInstance(ctx, resource)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		if expectedMeta != nil {
			expectedMeta.Verb = http.MethodPut
			suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)
		}

		if expectedSpec != nil {
			suite.verifyInstanceSpecStep(sCtx, expectedSpec, &resp.Spec)
		}

		suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
	})
}

func (suite *testSuite) getInstanceV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.ComputeV1,
	wref secapi.WorkspaceReference,
	expectedMeta *schema.RegionalWorkspaceResourceMetadata,
	expectedSpec *schema.InstanceSpec,
	expectedStatusState string,
) *schema.Instance {
	var resp *schema.Instance

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "GetInstance", string(wref.Workspace))

		retry := newStepRetry(
			suite.baseDelay,
			suite.baseInterval,
			suite.maxAttempts,
			func() schema.ResourceState {
				var err error
				resp, err = api.GetInstance(ctx, wref)
				requireNoError(sCtx, err)
				requireNotNilResponse(sCtx, resp)

				suite.requireNotNilStatus(sCtx, resp.Status)
				return *resp.Status.State
			},
			func() {
				expectedMeta.Verb = http.MethodGet
				suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)

				suite.verifyInstanceSpecStep(sCtx, expectedSpec, &resp.Spec)

				suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
			},
		)
		retry.run(sCtx, "GetInstance", expectedStatusState)
	})
	return resp
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
