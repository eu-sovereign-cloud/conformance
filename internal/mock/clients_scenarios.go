package mock

import (
	"log/slog"

	"github.com/eu-sovereign-cloud/go-sdk/pkg/secalib/builders"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/secalib/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/wiremock/go-wiremock"
)

func ConfigClientsInitScenario(params *ClientsInitParams) (*wiremock.Client, error) {
	slog.Info("Configuring mock to ClientsInit scenario")

	wm, err := newClient(params.MockURL)
	if err != nil {
		return nil, err
	}

	resource := generators.GenerateRegionResource(params.Region)
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
		Provider(regionProviderV1).
		Resource(resource).
		ApiVersion(apiVersion1).
		Spec(spec).
		BuildResponse()
	if err != nil {
		return nil, err
	}

	if err := configureGetStub(wm, "ClientsInit",
		&stubConfig{url: url, params: params, headers: headers, responseBody: response}); err != nil {
		return nil, err
	}

	return wm, nil
}
