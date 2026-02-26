package mockregion

import (
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"

	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	region "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.region.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

func ConfigureProviderQueriesV1(scenario *mockscenarios.Scenario, params *params.RegionProviderQueriesV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	regions := params.Regions

	// Generate resource
	regionsResource := generators.GenerateRegionListResource()

	// Generate URLs
	regionsUrl := generators.GenerateRegionListURL(sdkconsts.RegionProviderV1Name)
	regionUrl := generators.GenerateRegionURL(sdkconsts.RegionProviderV1Name, regions[0].Metadata.Name)
	regionsResponse := &region.RegionIterator{
		Metadata: schema.ResponseMetadata{
			Provider: sdkconsts.RegionProviderV1Name,
			Resource: regionsResource,
			Verb:     http.MethodGet,
		},
	}
	regionsResponse.Items = regions

	// 1 - Create ListRegions stub
	if err := configurator.ConfigureGetListRegionStub(regionsResponse, regionsUrl, scenario.MockParams, nil); err != nil {
		return err
	}

	// 2 - Create GetRegion stubs
	regionResponse := regions[0]
	if err := configurator.ConfigureGetRegionStub(&regionResponse, regionUrl, scenario.MockParams); err != nil {
		return err
	}

	if err := scenario.FinishConfiguration(configurator); err != nil {
		return err
	}
	return nil
}
