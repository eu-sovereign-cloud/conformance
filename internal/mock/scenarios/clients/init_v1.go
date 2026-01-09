package mockclients

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"

	"github.com/wiremock/go-wiremock"
)

func ConfigureInitScenarioV1(params *params.ClientsInitParams) (*wiremock.Client, error) {
	scenarios.LogScenarioMocking(constants.ClientsInitSuiteName)

	configurator, err := stubs.NewStubConfigurator(constants.ClientsInitSuiteName, params.MockParams)
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

	if err := configurator.ConfigureClientsInitStub(response, url, params.MockParams); err != nil {
		return nil, err
	}

	return configurator.Client, nil
}
