package mock

import (
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	regionV1 "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.region.v1"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/wiremock/go-wiremock"
)

func ConfigRegionLifecycleScenarioV1(scenario string, params *RegionParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}
	// Generate resource
	regionsResource := secalib.GenerateRegionResource("Regions")
	regionResource := secalib.GenerateRegionResource(params.Regions[0].Name)

	// Generate Url
	regionUrl := secalib.GenerateRegionURL(params.Regions[0].Name)

	// Build headers
	headerParams := map[string]string{
		authorizationHttpHeaderKey: authorizationHttpHeaderValuePrefix + params.AuthToken,
	}

	regionsResponse := &regionV1.RegionIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.RegionProviderV1,
			Resource: regionsResource,
			Verb:     http.MethodGet,
		},
	}
	var regionsList []schema.Region

	// Create Regions to be listed
	for _, region := range params.Regions {

		regionResource := secalib.GenerateRegionResource(region.Name)
		regionResponse := newRegionResponse(region.Name, secalib.RegionProviderV1, regionResource, secalib.ApiVersion1,
			region.InitialSpec)
		regionResponse.Metadata.Verb = http.MethodGet

		regionsList = append(regionsList, *regionResponse)

	}

	regionsResponse.Items = regionsList

	// 1 - Create ListRegions stub
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.RegionsURLV1, params: params, headers: headerParams, responseBody: regionsResponse, currentState: "", nextState: ""}); err != nil {
		return nil, err
	}

	// 2 - Create GetRegion stubs
	singleRegionResponse := &schema.Region{
		Metadata: &schema.GlobalResourceMetadata{
			Name:       params.Regions[0].Name,
			Provider:   secalib.RegionProviderV1,
			Resource:   regionResource,
			ApiVersion: secalib.ApiVersion1,
			Kind:       secalib.RegionKind,
			Verb:       http.MethodGet,
		},
	}

	if err := configureGetStub(wm, scenario,
		&stubConfig{url: regionUrl, params: params, headers: headerParams, responseBody: singleRegionResponse, currentState: "", nextState: ""}); err != nil {
		return nil, err
	}

	return wm, nil
}
