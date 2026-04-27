package stubs

import (
	"github.com/eu-sovereign-cloud/conformance/pkg/wiremock"
)

type Configurator struct {
	*wiremock.Configurator
}

func NewConfigurator(scenarioName string, params wiremock.MockParams) (*Configurator, error) {
	configurator, err := wiremock.NewConfigurator(scenarioName, params)
	if err != nil {
		return nil, err
	}

	return &Configurator{
		Configurator: configurator,
	}, nil
}
