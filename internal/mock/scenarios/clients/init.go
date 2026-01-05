package clients

import (
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/wiremock/go-wiremock"
)

func ConfigureInitScenario(params *mock.ClientsInitParams) (*wiremock.Client, error) {
	slog.Info("Configuring mock to ClientsInit scenario")

	configurator, err := stubs.NewScenarioConfigurator("ClientsInit", params.MockURL)
	if err != nil {
		return nil, err
	}

	url := generators.GenerateRegionURL(mock.RegionProviderV1, params.Region)

	spec := &schema.RegionSpec{
		AvailableZones: []string{mock.ZoneA, mock.ZoneB},
		Providers:      mock.BuildProviderSpecV1(),
	}

	response, err := builders.NewRegionBuilder().
		Name(params.Region).
		Provider(mock.RegionProviderV1).ApiVersion(mock.ApiVersion1).
		Spec(spec).
		Build()
	if err != nil {
		return nil, err
	}

	if err := configurator.ConfigureClientsInitStub(response, url, params.GetBaseParams()); err != nil {
		return nil, err
	}

	return configurator.Client, nil
}
