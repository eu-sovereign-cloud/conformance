package stubs

import (
	"net/http"
	"strconv"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"

	"github.com/wiremock/go-wiremock"
)

const (
	statePrefix = "State."
)

type Configurator struct {
	client       *wiremock.Client
	scenarioName string
	stateID      int
	currentState string
}

func NewConfigurator(scenarioName string, params *mock.MockParams) (*Configurator, error) {
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
	params *mock.MockParams,
	pathParams map[string]string,
	responseBody any,
) error {
	// Calculte next state
	stateID := configurator.stateID + 1
	nextState := statePrefix + strconv.Itoa(stateID)

	if err := stubFunc(configurator.client, configurator.scenarioName,
		&stubConfig{url: url, params: params, pathParams: pathParams, responseBody: responseBody, currentState: configurator.currentState, nextState: nextState}); err != nil {
		return err
	}

	configurator.stateID = stateID
	configurator.currentState = nextState
	return nil
}

func (configurator *Configurator) ConfigurePutStub(url string, params *mock.MockParams, setMetadataVerbFunc func(string), responseBody any) error {
	if setMetadataVerbFunc != nil {
		setMetadataVerbFunc(http.MethodPut)
	}
	if err := configurator.configureStub(ConfigurePutStub, url, params, nil, responseBody); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigurePostStub(url string, params *mock.MockParams, setMetadataVerbFunc func(string), responseBody any) error {
	if setMetadataVerbFunc != nil {
		setMetadataVerbFunc(http.MethodPost)
	}
	if err := configurator.configureStub(configurePostStub, url, params, nil, responseBody); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetStub(url string, params *mock.MockParams, setMetadataVerbFunc func(string), responseBody any) error {
	if setMetadataVerbFunc != nil {
		setMetadataVerbFunc(http.MethodGet)
	}
	if err := configurator.configureStub(ConfigureGetStub, url, params, nil, responseBody); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetListStub(url string, params *mock.MockParams, pathParams map[string]string, setMetadataVerbFunc func(string), responseBody any) error {
	if setMetadataVerbFunc != nil {
		setMetadataVerbFunc(http.MethodGet)
	}
	if err := configurator.configureStub(ConfigureGetStub, url, params, pathParams, responseBody); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureGetNotFoundStub(url string, params *mock.MockParams) error {
	if err := configurator.configureStub(configureGetNotFoundStub, url, params, nil, nil); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) ConfigureDeleteStub(url string, params *mock.MockParams) error {
	if err := configurator.configureStub(configureDeleteStub, url, params, nil, nil); err != nil {
		return err
	}
	return nil
}

func (configurator *Configurator) Finish() (*wiremock.Client, error) {
	if err := configureDefaultStub(configurator.client); err != nil {
		return nil, err
	}

	return configurator.client, nil
}

func newMockClient(mockURL string) (*wiremock.Client, error) {
	wm := wiremock.NewClient(mockURL)
	if err := wm.ResetAllScenarios(); err != nil {
		return nil, err
	}
	return wm, nil
}
