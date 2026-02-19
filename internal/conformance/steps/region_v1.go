//nolint:dupl
package steps

import (
	"context"
	"net/http"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/ozontech/allure-go/pkg/framework/provider"
)

func (configurator *StepsConfigurator) GetRegionV1Step(stepName string, ctx context.Context, api secapi.RegionV1, regionName string,
	responseExpects StepResponseExpects[schema.GlobalResourceMetadata, schema.RegionSpec],
) *schema.Region {
	responseExpects.Metadata.Verb = http.MethodGet
	configurator.logStepName(stepName)
	return getGlobalResourceStep(configurator.t,
		getGlobalResourceParams[schema.Region, schema.GlobalResourceMetadata, schema.RegionSpec]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetRegionV1StepParams,
			operationName:  "GetRegion",
			resourceName:   regionName,
			getFunc: func(ctx context.Context, name string) (*globalStepFuncResponse[schema.Region, schema.GlobalResourceMetadata, schema.RegionSpec], error) {
				if resp, err := api.GetRegion(ctx, name); err == nil {
					return newGlobalStepFuncResponse(resp, resp.Metadata, resp.Spec), nil
				} else {
					return nil, err
				}
			},
			expectedMetadata:   responseExpects.Metadata,
			verifyMetadataFunc: configurator.suite.VerifyGlobalResourceMetadataStep,
			expectedSpec:       responseExpects.Spec,
			verifySpecFunc:     configurator.suite.VerifyRegionSpecStep,
		},
	)
}

func (configurator *StepsConfigurator) ListRegionsV1Step(stepName string, ctx context.Context, api secapi.RegionV1) []*schema.Region {
	var resp []*schema.Region
	configurator.logStepName(stepName)
	configurator.t.WithNewStep(stepName, func(sCtx provider.StepCtx) {
		configurator.suite.SetRegionV1StepParams(sCtx, "ListRegions")

		iter, err := api.ListRegions(ctx)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, iter)

		resp, err = iter.All(ctx)
		requireNoError(sCtx, err)
		requireNotNilResponse(sCtx, resp)
		requireNotEmptyResponse(sCtx, resp)
	})
	return resp
}
