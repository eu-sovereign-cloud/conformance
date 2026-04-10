package mockregion

import (
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"

	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	region "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.region.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

func ConfigureProviderQueriesV1(scenario *mockscenarios.Scenario, params params.RegionProviderQueriesV1Params) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	var regionsMetadata = params.RegionsMetadata

	// Generate resource
	regionsResource := generators.GenerateRegionListResource()

	// Generate URLs
	regionsUrl := generators.GenerateRegionListURL(sdkconsts.RegionProviderV1Name)
	regionUrl := generators.GenerateRegionURL(sdkconsts.RegionProviderV1Name, regionsMetadata[0].Name)

	// Build the regions
	var regions []schema.Region
	for _, regionMeta := range regionsMetadata {

		resource, err := builders.NewRegionBuilder().
			Name(regionMeta.Name).
			Provider(regionMeta.Provider).ApiVersion(regionMeta.ApiVersion).
			Spec(&schema.RegionSpec{
				AvailableZones: []string{constants.ZoneA, constants.ZoneB},
				Providers:      mock.BuildProviderSpecV1(params.MockProviders),
			}).
			Build()
		if err != nil {
			return err
		}

		regions = append(regions, *resource)
	}

	// Build the response
	regionsResponse := &region.RegionIterator{
		Metadata: schema.ResponseMetadata{
			Provider: sdkconsts.RegionProviderV1Name,
			Resource: regionsResource,
			Verb:     http.MethodGet,
		},
	}
	regionsResponse.Items = regions

	// List regions
	if err := configurator.ConfigureListRegionStub(regionsResponse, regionsUrl, scenario.MockParams, nil); err != nil {
		return err
	}

	// Get region
	regionResponse := regions[0]
	if err := configurator.ConfigureGetRegionStub(&regionResponse, regionUrl, scenario.MockParams); err != nil {
		return err
	}

	if err := scenario.FinishConfiguration(configurator); err != nil {
		return err
	}
	return nil
}
