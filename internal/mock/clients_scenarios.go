package mock

import (
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/wiremock/go-wiremock"
)

func ConfigureClientsInitScenario(params *ClientsInitParams) (*wiremock.Client, error) {
	slog.Info("Configuring mock to ClientsInit scenario")

	wm, err := newMockClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	url := generators.GenerateRegionURL(regionProviderV1, params.Region)

	headers := map[string]string{
		authorizationHttpHeaderKey: authorizationHttpHeaderValuePrefix + params.AuthToken,
	}

	spec := &schema.RegionSpec{
		AvailableZones: []string{zoneA, zoneB},
		Providers:      BuildProviderSpec(),
	}

	response, err := builders.NewRegionBuilder().
		Name(params.Region).
		Provider(regionProviderV1).ApiVersion(apiVersion1).
		Spec(spec).
		Build()
	if err != nil {
		return nil, err
	}

	if err := configureGetStub(wm, "ClientsInit",
		&stubConfig{url: url, params: params.getBaseParams(), headers: headers, responseBody: response}); err != nil {
		return nil, err
	}

	return wm, nil
}
