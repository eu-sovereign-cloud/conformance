package region

import (
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	region "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.region.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"

	"github.com/wiremock/go-wiremock"
)

func ConfigureListScenarioV1(scenario string, params *mock.RegionListParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	configurator, err := stubs.NewScenarioConfigurator(scenario, params.MockURL)
	if err != nil {
		return nil, err
	}

	// Generate resource
	// TODO Create a region list resource function
	regionsResource := generators.GenerateRegionResource("Regions")
	regionResource := generators.GenerateRegionResource(params.Regions[0].Name)

	// Generate URLs
	regionsUrl := generators.GenerateRegionListURL(constants.RegionProviderV1)
	regionUrl := generators.GenerateRegionURL(constants.RegionProviderV1, params.Regions[0].Name)

	regionsResponse := &region.RegionIterator{
		Metadata: schema.ResponseMetadata{
			Provider: constants.RegionProviderV1,
			Resource: regionsResource,
			Verb:     http.MethodGet,
		},
	}
	var regionsList []schema.Region

	// Create Regions to be listed
	for _, region := range params.Regions {

		regionResponse, err := builders.NewRegionBuilder().
			Name(region.Name).
			Provider(constants.RegionProviderV1).ApiVersion(constants.ApiVersion1).
			Spec(region.InitialSpec).
			Build()
		if err != nil {
			return nil, err
		}

		regionsList = append(regionsList, *regionResponse)
	}

	regionsResponse.Items = regionsList

	// 1 - Create ListRegions stub
	if err := configurator.ConfigureGetListRegionStub(regionsResponse, regionsUrl, params.GetBaseParams(), nil); err != nil {
		return nil, err
	}

	// 2 - Create GetRegion stubs
	singleRegionResponse := &schema.Region{
		Metadata: &schema.GlobalResourceMetadata{
			Name:       params.Regions[0].Name,
			Provider:   constants.RegionProviderV1,
			Resource:   regionResource,
			ApiVersion: constants.ApiVersion1,
			Kind:       schema.GlobalResourceMetadataKindResourceKindRegion,
			Verb:       http.MethodGet,
		},
		Spec: regionsResponse.Items[0].Spec,
	}

	if err := configurator.ConfigureGetRegionStub(singleRegionResponse, regionUrl, params.GetBaseParams()); err != nil {
		return nil, err
	}

	return configurator.Client, nil
}
