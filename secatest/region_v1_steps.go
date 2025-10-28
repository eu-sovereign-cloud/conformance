package secatest

import (
	"context"
	"net/http"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/eu-sovereign-cloud/go-sdk/secapi/builders"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (suite *testSuite) getRegion(stepName string, t provider.T, ctx context.Context, api *secapi.RegionV1,
	expectedMeta *schema.GlobalResourceMetadata,
) *schema.Region {
	var resp *schema.Region
	var err error

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setRegionV1StepParams(sCtx, "GetRegion")

		resp, err = api.GetRegion(ctx, expectedMeta.Name)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		expectedMeta.Verb = http.MethodGet
		suite.verifyGlobalResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)

		suite.verifyRegionSpecStep(sCtx, &resp.Spec)

	})
	return resp
}

func (suite *testSuite) getListRegion(stepName string, t provider.T, ctx context.Context, api *secapi.RegionV1,
) []*schema.Region {
	var resp []*schema.Region

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setRegionV1StepParams(sCtx, "GetListRegion")

		iter, err := api.ListRegions(ctx)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, iter)

		resp, err = iter.All(ctx)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)
		requireLenResponse(sCtx, len(resp))
	})
	return resp
}

func (suite *testSuite) getListRegionWithParameters(stepName string, t provider.T, ctx context.Context, api *secapi.RegionV1, opts *builders.ListOptions,
) []*schema.Region {
	var resp []*schema.Region

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setRegionV1StepParams(sCtx, "GetListRegion")

		iter, err := api.ListRegionsWithFilters(ctx, opts)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, iter)

		resp, err = iter.All(ctx)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)
		requireLenResponse(sCtx, len(resp))
	})
	return resp
}
