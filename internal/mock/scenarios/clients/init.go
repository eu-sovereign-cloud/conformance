package clients

import (
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
	"github.com/wiremock/go-wiremock"
)

func ConfigureInitScenario(params *mock.ClientsInitParams) (*wiremock.Client, error) {
	slog.Info("Configuring mock to ClientsInit scenario")

	configurator, err := stubs.NewStubConfigurator("ClientsInit", params.MockURL)
	if err != nil {
		return nil, err
	}

	url := generators.GenerateRegionURL(constants.RegionProviderV1, params.Region)

	spec := &schema.RegionSpec{
		AvailableZones: []string{constants.ZoneA, constants.ZoneB},
		Providers:      mock.BuildProviderSpecV1(),
	}

	response, err := builders.NewRegionBuilder().
		Name(params.Region).
		Provider(constants.RegionProviderV1).ApiVersion(constants.ApiVersion1).
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
