package mock

import (
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/secalib/builders"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/secalib/generators"
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
	regionsResource := generators.GenerateRegionResource("Regions")
	regionResource := generators.GenerateRegionResource(params.Regions[0].Name)

	// Generate URLs
	regionsUrl := generators.GenerateRegionsURL(regionProviderV1)
	regionUrl := generators.GenerateRegionURL(regionProviderV1, params.Regions[0].Name)

	// Build headers
	headerParams := map[string]string{
		authorizationHttpHeaderKey: authorizationHttpHeaderValuePrefix + params.AuthToken,
	}

	regionsResponse := &regionV1.RegionIterator{
		Metadata: schema.ResponseMetadata{
			Provider: regionProviderV1,
			Resource: regionsResource,
			Verb:     http.MethodGet,
		},
	}
	var regionsList []schema.Region

	// Create Regions to be listed
	for _, region := range params.Regions {

		regionResource := generators.GenerateRegionResource(region.Name)
		regionResponse, err := builders.NewRegionBuilder().
			Name(region.Name).Resource(regionResource).
			Provider(regionProviderV1).ApiVersion(apiVersion1).
			Spec(region.InitialSpec).
			BuildResponse()
		if err != nil {
			return nil, err
		}

		regionsList = append(regionsList, *regionResponse)
	}

	regionsResponse.Items = regionsList

	// 1 - Create ListRegions stub
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: regionsUrl, params: params, headers: headerParams, responseBody: regionsResponse, currentState: "", nextState: ""}); err != nil {
		return nil, err
	}

	// 2 - Create GetRegion stubs
	singleRegionResponse := &schema.Region{
		Metadata: &schema.GlobalResourceMetadata{
			Name:       params.Regions[0].Name,
			Provider:   regionProviderV1,
			Resource:   regionResource,
			ApiVersion: apiVersion1,
			Kind:       schema.GlobalResourceMetadataKindResourceKindRegion,
			Verb:       http.MethodGet,
		},
		Spec: regionsResponse.Items[0].Spec,
	}

	if err := configureGetStub(wm, scenario,
		&stubConfig{url: regionUrl, params: params, headers: headerParams, responseBody: singleRegionResponse, currentState: "", nextState: ""}); err != nil {
		return nil, err
	}

	return wm, nil
}
