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

	// Build headers
	headerParams := map[string]string{
		authorizationHttpHeaderKey: authorizationHttpHeaderValuePrefix + params.AuthToken,
	}

	regionsResponse := &regionV1.RegionIterator{
		Metadata: schema.ResponseMetadata{
			Provider: secalib.RegionProviderV1,
			Resource: secalib.ApiVersion1,
			Verb:     http.MethodGet,
		},
	}
	var regionsList []*schema.Region

	// Create Regions to be listed
	for _, region := range params.Regions {
		regionResource := secalib.GenerateRegionResource(region.Name)
		regionResponse := newRegionResponse(region.Name, secalib.RegionProviderV1, regionResource, secalib.ApiVersion1,
			region.InitialSpec)
		regionResponse.Metadata.Verb = http.MethodGet
		if err := configureGetStub(wm, scenario,
			&stubConfig{url: secalib.RegionsURLV1, params: params, headers: headerParams, responseBody: regionsResponse, currentState: "", nextState: ""}); err != nil {
			return nil, err
		}
		regionsList = append(regionsList, regionResponse)

	}

	// List Region

	// Convert []*schema.Region to []schema.Region
	var regionsValueList []schema.Region
	for _, r := range regionsList {
		if r != nil {
			regionsValueList = append(regionsValueList, *r)
		}
	}
	regionsResponse.Items = regionsValueList
	regionsResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.RegionsURLV1, params: params, headers: headerParams, responseBody: regionsResponse, currentState: "", nextState: ""}); err != nil {
		return nil, err
	}

	// List Region with labels
	// Limit = 1
	if len(regionsResponse.Items) > 0 {
		regionsResponse.Items = []schema.Region{regionsResponse.Items[0]}
	}

	headerParams = map[string]string{
		authorizationHttpHeaderKey: authorizationHttpHeaderValuePrefix + params.AuthToken,
		limitHeaderKey:             "1",
	}
	regionsResponse.Metadata.Verb = http.MethodGet
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.RegionsURLV1, params: params, headers: headerParams, responseBody: regionsResponse, currentState: "", nextState: ""}); err != nil {
		return nil, err
	}
	return wm, nil
}
