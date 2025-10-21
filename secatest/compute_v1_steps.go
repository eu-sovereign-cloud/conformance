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

func (suite *testSuite) createOrUpdateInstanceV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.ComputeV1, instance *schema.Instance,
	expectedMeta *schema.RegionalWorkspaceResourceMetadata, expectedSpec *schema.InstanceSpec, expectedStatusState string,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "CreateOrUpdateInstance")

		resp, err := api.CreateOrUpdateInstance(ctx, instance)
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

func (suite *testSuite) getInstanceV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.ComputeV1, wref secapi.WorkspaceReference,
	expectedMeta *schema.RegionalWorkspaceResourceMetadata, expectedSpec *schema.InstanceSpec, expectedStatusState string,
) *schema.Instance {
	var resp *schema.Instance
	var err error

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetInstance")

		resp, err = api.GetInstance(ctx, wref)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		if expectedMeta != nil {
			expectedMeta.Verb = http.MethodGet
			suite.verifyRegionalWorkspaceResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)
		}

		if expectedSpec != nil {
			suite.verifyInstanceSpecStep(sCtx, expectedSpec, &resp.Spec)
		}

		suite.verifyStatusStep(sCtx, *secalib.SetResourceState(expectedStatusState), *resp.Status.State)
	})
	return resp
}

func (suite *testSuite) getInstanceWithErrorV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.ComputeV1, wref secapi.WorkspaceReference,
	expectedError error,
) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "GetInstance")

		_, err := api.GetInstance(ctx, wref)
		requireError(sCtx, err, expectedError)
	})
}

func (suite *testSuite) startInstanceV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.ComputeV1, instance *schema.Instance) {
	var err error

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "StartInstance")

		err = api.StartInstance(ctx, instance)
		requireNoError(sCtx, err)
	})
}

func (suite *testSuite) stopInstanceV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.ComputeV1, instance *schema.Instance) {
	var err error

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "StopInstance")

		err = api.StopInstance(ctx, instance)
		requireNoError(sCtx, err)
	})
}

func (suite *testSuite) restartInstanceV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.ComputeV1, instance *schema.Instance) {
	var err error

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "RestartInstance")

		err = api.RestartInstance(ctx, instance)
		requireNoError(sCtx, err)
	})
}

func (suite *testSuite) deleteInstanceV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.ComputeV1, instance *schema.Instance) {
	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setAuthorizationV1StepParams(sCtx, "DeleteInstance")

		err := api.DeleteInstance(ctx, instance)
		requireNoError(sCtx, err)
	})
}
