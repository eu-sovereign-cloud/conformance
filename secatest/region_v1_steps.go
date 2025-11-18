package secatest

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (suite *testSuite) getRegionV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.RegionV1, expectedMeta *schema.GlobalResourceMetadata) *schema.Region {
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

func (suite *testSuite) listRegionsV1Step(stepName string, t provider.T, ctx context.Context, api *secapi.RegionV1) []*schema.Region {
	var respNext []*schema.Region
	var respAll []*schema.Region

	t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		suite.setRegionV1StepParams(sCtx, "ListRegions")

		iter, err := api.ListRegions(ctx)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, iter)

		for {
			item, err := iter.Next(ctx)
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

		iterAll, err := api.ListRegions(ctx)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, iterAll)

		respAll, err = iterAll.All(ctx)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, respAll)
		requireLenResponse(sCtx, len(respAll))

		compareIteratorsResponse(sCtx, len(respNext), len(respAll))
	})
	return respAll
}
