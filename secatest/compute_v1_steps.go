package secatest

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/eu-sovereign-cloud/go-sdk/secapi/builders"

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
	var err error

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "GetInstance", string(wref.Workspace))

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
		suite.setComputeV1StepParams(sCtx, "StopInstance", string(resource.Metadata.Workspace))

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

func (suite *testSuite) getListInstanceV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.ComputeV1,
	tref secapi.TenantReference,
	wref secapi.WorkspaceReference,
	opts *builders.ListOptions,
) []*schema.Instance {
	var respNext []*schema.Instance
	var respAll []*schema.Instance

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "ListInstances with parameters", string(wref.Workspace))
		var iter *secapi.Iterator[schema.Instance]
		var err error
		if opts != nil {
			iter, err = api.ListInstancesWithFilters(ctx, tref.Tenant, wref.Workspace, opts)
		} else {
			iter, err = api.ListInstances(ctx, tref.Tenant, wref.Workspace)
		}
		requireNoError(sCtx, err)
		for {
			item, err := iter.Next(context.Background())
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				break
			}
			respNext = append(respNext, item)
		}

		requireNotNilResponse(sCtx, respNext)
		requireLenResponse(sCtx, len(respNext))

		/*
			respAll, err = iter.All(ctx)
			requireNoError(sCtx, err)
			requireNotNilResponse(sCtx, respAll)
			requireLenResponse(sCtx, len(respAll))

			compareIteratorsResponse(sCtx, len(respNext), len(respAll))
		*/
	})
	return respAll
}

func (suite *testSuite) getListSkusV1Step(
	stepName string,
	t provider.T,
	ctx context.Context,
	api *secapi.ComputeV1,
	tref secapi.TenantReference,
	opts *builders.ListOptions,
) []*schema.InstanceSku {
	var respNext []*schema.InstanceSku
	var respAll []*schema.InstanceSku

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setComputeV1StepParams(sCtx, "ListSkus", string(tref.Name))

		var iter *secapi.Iterator[schema.InstanceSku]
		var err error
		if opts != nil {
			iter, err = api.ListSkusWithFilters(ctx, tref.Tenant, opts)
		} else {
			iter, err = api.ListSkus(ctx, tref.Tenant)
		}
		requireNoError(sCtx, err)
		for {
			item, err := iter.Next(context.Background())
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				break
			}
			respNext = append(respNext, item)
		}

		requireNotNilResponse(sCtx, respNext)
		requireLenResponse(sCtx, len(respNext))
		/*
			respAll, err = iter.All(ctx)
			requireNoError(sCtx, err)
			requireNotNilResponse(sCtx, respAll)
			requireLenResponse(sCtx, len(respAll))

			compareIteratorsResponse(sCtx, len(respNext), len(respAll))
		*/
	})
	return respAll
}
