//nolint:dupl
package steps

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (configurator *StepsConfigurator) GetRegionV1Step(stepName string, ctx context.Context, api *secapi.RegionV1, expectedMeta *schema.GlobalResourceMetadata) *schema.Region {
	var resp *schema.Region
	var err error
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetRegionV1StepParams(sCtx, "GetRegion")

		resp, err = api.GetRegion(ctx, expectedMeta.Name)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)

		expectedMeta.Verb = http.MethodGet
		configurator.suite.VerifyGlobalResourceMetadataStep(sCtx, expectedMeta, resp.Metadata)

		configurator.suite.VerifyRegionSpecStep(sCtx, &resp.Spec)
	})
	return resp
}

func (configurator *StepsConfigurator) ListRegionsV1Step(stepName string, ctx context.Context, api *secapi.RegionV1) []*schema.Region {
	var resp []*schema.Region
	slog.Info(fmt.Sprintf("[%s] %s", configurator.suite.ScenarioName, stepName))
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetRegionV1StepParams(sCtx, "ListRegions")

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
