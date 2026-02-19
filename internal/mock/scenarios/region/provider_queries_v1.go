package mockregion

import (
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
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
	regionResource := generators.GenerateRegionResource(regions[0].Metadata.Name)

	// Generate URLs
	regionsUrl := generators.GenerateRegionListURL(constants.RegionProviderV1)
	regionUrl := generators.GenerateRegionURL(constants.RegionProviderV1, regions[0].Metadata.Name)
	regionsResponse := &region.RegionIterator{
		Metadata: schema.ResponseMetadata{
			Provider: constants.RegionProviderV1,
			Resource: regionsResource,
			Verb:     http.MethodGet,
		},
	}
	var regionsList []schema.Region

	// Create Regions to be listed
	for _, region := range regions {

		regionResponse, err := builders.NewRegionBuilder().
			Name(region.Metadata.Name).
			Provider(constants.RegionProviderV1).ApiVersion(constants.ApiVersion1).
			Spec(&region.Spec).
			Build()
		if err != nil {
			return err
		}

		regionsList = append(regionsList, *regionResponse)
	}

	regionsResponse.Items = regionsList

	// 1 - Create ListRegions stub
	if err := configurator.ConfigureGetListRegionStub(regionsResponse, regionsUrl, scenario.MockParams, nil); err != nil {
		return err
	}

	// 2 - Create GetRegion stubs
	region := regions[0]
	singleRegionResponse := &schema.Region{
		Metadata: &schema.GlobalResourceMetadata{
			Name:       region.Metadata.Name,
			Provider:   constants.RegionProviderV1,
			Resource:   regionResource,
			ApiVersion: constants.ApiVersion1,
			Kind:       schema.GlobalResourceMetadataKindResourceKindRegion,
			Verb:       http.MethodGet,
		},
		Spec: region.Spec,
	}

	if err := configurator.ConfigureGetRegionStub(singleRegionResponse, regionUrl, scenario.MockParams); err != nil {
		return err
	}

	if err := scenario.FinishConfiguration(configurator); err != nil {
		return err
	}
	return nil
}
