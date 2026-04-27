package wiremock

import (
	"net/http"
	"strconv"

	"github.com/wiremock/go-wiremock"
)

const (
	statePrefix = "State."
)

type Configurator struct {
	client       *MockClient
	scenarioName string
	stateID      int
	currentState string
}

func NewConfigurator(scenarioName string, params MockParams) (*Configurator, error) {
	client, err := newMockClient(params.ServerURL)
	if err != nil {
		return nil, err
	}

	return &Configurator{
		client:       client,
		scenarioName: scenarioName,
		stateID:      0,
		currentState: startedScenarioState,
	}, nil
}

func (configurator *Configurator) configureStub(
	stubFunc func(*wiremock.Client, string, *stubConfig) error,
	url string,
	mockParams MockParams,
	pathParams map[string]string,
	responseBody any,
) error {
	// Calculte next state
	stateID := configurator.stateID + 1
	nextState := statePrefix + strconv.Itoa(stateID)

	if err := stubFunc(configurator.client.wm, configurator.scenarioName,
		&stubConfig{url: url, params: mockParams, pathParams: pathParams, responseBody: responseBody, currentState: configurator.currentState, nextState: nextState}); err != nil {
		return err
	}

	configurator.stateID = stateID
	configurator.currentState = nextState
	return nil
}

func (configurator *Configurator) ConfigurePutStub(url string, params MockParams, setMetadataVerbFunc func(string), responseBody any) error {
	if setMetadataVerbFunc != nil {
		setMetadataVerbFunc(http.MethodPut)
	}
	if err := configurator.configureStub(ConfigurePutStub, url, params, nil, responseBody); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigurePostStub(url string, params MockParams, setMetadataVerbFunc func(string), responseBody any) error {
	if setMetadataVerbFunc != nil {
		setMetadataVerbFunc(http.MethodPost)
	}
	if err := configurator.configureStub(configurePostStub, url, params, nil, responseBody); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetWithPathStub(url string, params MockParams, pathParams map[string]string, setMetadataVerbFunc func(string), responseBody any) error {
	if setMetadataVerbFunc != nil {
		setMetadataVerbFunc(http.MethodGet)
	}
	if err := configurator.configureStub(ConfigureGetStub, url, params, pathParams, responseBody); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetStub(url string, params MockParams, setMetadataVerbFunc func(string), responseBody any) error {
	if setMetadataVerbFunc != nil {
		setMetadataVerbFunc(http.MethodGet)
	}
	if err := configurator.configureStub(ConfigureGetStub, url, params, nil, responseBody); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetNotFoundStub(url string, params MockParams) error {
	if err := configurator.configureStub(configureGetNotFoundStub, url, params, nil, nil); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureDeleteStub(url string, params MockParams) error {
	if err := configurator.configureStub(configureDeleteStub, url, params, nil, nil); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) Finish() (*MockClient, error) {
	if err := configureDefaultStub(configurator.client.wm); err != nil {
		return nil, err
	}

	return configurator.client, nil
}
