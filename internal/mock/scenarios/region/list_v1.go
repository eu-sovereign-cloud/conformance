package mockregion

import (
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	region "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.region.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"

	"github.com/wiremock/go-wiremock"
)

func ConfigureListScenarioV1(scenario string, mockParams *mock.MockParams, suiteParams *params.RegionListV1Params) (*wiremock.Client, error) {
	scenarios.LogScenarioMocking(scenario)

	regions := suiteParams.Regions

	configurator, err := stubs.NewStubConfigurator(scenario, mockParams)
	if err != nil {
		return nil, err
	}

	// Generate resource
	regionsResource := generators.GenerateRegionListResource()

	// Generate URLs
	regionsUrl := generators.GenerateRegionListURL(constants.RegionProviderV1)
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
			return nil, err
		}

		regionsList = append(regionsList, *regionResponse)
	}

	regionsResponse.Items = regionsList

	// Configure stub
	if err := configurator.ConfigureGetListRegionStub(regionsResponse, regionsUrl, mockParams, nil); err != nil {
		return nil, err
	}

	// Finish the stubs configuration
	if client, err := configurator.Finish(); err != nil {
		return nil, err
	} else {
		return client, nil
	}
}
