//nolint:dupl
package steps

import (
	"context"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/pkg/wrappers"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
)

func (configurator *StepsConfigurator) GetRegionV1Step(stepName string, ctx context.Context, api secapi.RegionV1, regionName string,
	responseExpects ResponseExpects[schema.GlobalResourceMetadata, schema.RegionSpec],
) *schema.Region {
	responseExpects.Metadata.Verb = http.MethodGet
	return getGlobalResourceStep(configurator.t,
		getGlobalResourceParams[schema.Region, schema.GlobalResourceMetadata, schema.RegionSpec]{
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetRegionV1StepParams,
			operationName:  constants.GetRegionOperation,
			resourceName:   regionName,
			getFunc: func(ctx context.Context, resourceName string) (wrappers.GlobalResourceWrapper[schema.Region, schema.GlobalResourceMetadata, schema.RegionSpec], error) {
				resp, err := api.GetRegion(ctx, resourceName)
				return wrappers.NewRegionWrapper(resp), err
			},
			expectedMetadata:   responseExpects.Metadata,
			verifyMetadataFunc: configurator.suite.VerifyGlobalResourceMetadataStep,
			expectedSpec:       responseExpects.Spec,
			verifySpecFunc:     configurator.suite.VerifyRegionSpecStep,
		},
	)
}

func (configurator *StepsConfigurator) ListRegionsV1Step(stepName string, ctx context.Context, api secapi.RegionV1, opts *secapi.ListOptions) []*schema.Region {
	return listGlobalResourcesStep(configurator.t, configurator.suite, stepName,
		listGlobalResourcesParams[schema.Region, schema.GlobalResourceMetadata]{
			listOptions: opts,
			listFunc: func(ctx context.Context, options *secapi.ListOptions) (*secapi.Iterator[schema.Region], error) {
				return api.ListRegionsWithOptions(ctx, options)
			},
			stepName:       stepName,
			stepParamsFunc: configurator.suite.SetRegionV1StepParams,
			operationName:  constants.ListRegionsOperation,
		},
	)
}
