package mock

import (
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/secalib"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"

	"github.com/wiremock/go-wiremock"
)

func CreateRegionLifecycleScenarioV1(scenario string, params RegionParamsV1) (*wiremock.Client, error) {
	slog.Info("Configuring mock to scenario " + scenario)

	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	var regionsResponse []*schema.Region

	for _, region := range params.Regions {
		regionResource := secalib.GenerateRegionResource(region.Name)
		regionsResponse = append(regionsResponse, newRegionResponse(region.Name, secalib.RegionProviderV1, regionResource, secalib.ApiVersion1,
			region.InitialSpec))
	}

	// List Region

	//setCreatedGlobalResourceMetadata(regionsResponse.Metadata)
	if err := configureGetStub(wm, scenario,
		&stubConfig{url: secalib.RegionsURLV1, params: params, responseBody: regionsResponse, currentState: startedScenarioState, nextState: "end"}); err != nil {
		return nil, err
	}

	return wm, nil
}
