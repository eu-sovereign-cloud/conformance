package mockclients

import (
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/conformance/pkg/builders"
	"github.com/eu-sovereign-cloud/conformance/pkg/generators"
	sdkconsts "github.com/eu-sovereign-cloud/go-sdk/pkg/constants"
	"github.com/eu-sovereign-cloud/go-sdk/pkg/spec/schema"
)

func ConfigureInitScenarioV1(scenario *mockscenarios.Scenario, params *params.ClientsInitParams) error {
	configurator, err := scenario.StartConfiguration()
	if err != nil {
		return err
	}

	url := generators.GenerateRegionURL(sdkconsts.RegionProviderV1Name, params.Region)

	spec := &schema.RegionSpec{
		AvailableZones: []string{constants.ZoneA, constants.ZoneB},
		Providers:      mock.BuildProviderSpecV1(),
	}

	response, err := builders.NewRegionBuilder().
		Name(params.Region).
		Provider(sdkconsts.RegionProviderV1Name).ApiVersion(sdkconsts.ApiVersion1).
		Spec(spec).
		Build()
	if err != nil {
		return err
	}

	if err := configurator.ConfigureClientsInitStub(response, url, scenario.MockParams); err != nil {
		return err
	}

	if err := scenario.FinishConfiguration(configurator); err != nil {
		return err
	}
	return nil
}
