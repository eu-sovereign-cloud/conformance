//nolint:dupl
package steps

import (
	"context"
	"net/http"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (builder *Builder) GetRegionV1Step(stepName string, ctx context.Context, api *secapi.RegionV1, expectedMeta *schema.GlobalResourceMetadata) *schema.Region {
	var resp *schema.Region
	var err error

	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetRegionV1StepParams(sCtx, "GetRegion")

		resp, err = api.GetRegion(ctx, expectedMeta.Name)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		expectedMeta.Verb = http.MethodGet
		builder.suite.VerifyGlobalResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)

		builder.suite.VerifyRegionSpecStep(sCtx, &resp.Spec)
	})
	return resp
}

func (builder *Builder) ListRegionsV1Step(stepName string, ctx context.Context, api *secapi.RegionV1) []*schema.Region {
	var resp []*schema.Region

	builder.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		builder.suite.SetRegionV1StepParams(sCtx, "ListRegions")

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
