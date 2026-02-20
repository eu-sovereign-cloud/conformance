package mockscenarios

import (
	"fmt"
	"log/slog"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/conformance/internal/mock/stubs"

	"github.com/wiremock/go-wiremock"
)

type Scenario struct {
	Name       string
	MockParams *mock.MockParams
	client     *wiremock.Client
}

func NewScenario(name string,
	params *mock.MockParams,
) *Scenario {
	return &Scenario{
		Name:       name,
		MockParams: params,
	}
}

func (scenario *Scenario) StartConfiguration() (*stubs.Configurator, error) {
	slog.Info("Mocking of scenario [" + scenario.Name + "]")

	configurator, err := stubs.NewConfigurator(scenario.Name, scenario.MockParams)
	if err != nil {
		return nil, err
	}

	return configurator, err
}

func (scenario *Scenario) FinishConfiguration(stubsConfigurator *stubs.Configurator) error {
	if client, err := stubsConfigurator.Finish(); err != nil {
		return err
	} else {
		scenario.client = client
	}

	return nil
}

func (scenario *Scenario) ResetScenario() error {
	if scenario.client != nil {
		if err := scenario.client.ResetAllScenarios(); err != nil {
			return fmt.Errorf("Failed to reset scenario: %w", err)
		}
	}
	return nil
}

func LogScenarioMocking(scenarioName string) {
	slog.Info("Mocking of scenario " + scenarioName)
}
