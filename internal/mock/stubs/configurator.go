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

type stubConfigurator struct {
	client       *wiremock.Client
	scenarioName string
	stateID      int
	currentState string
}

func NewStubConfigurator(scenarioName string, params *mock.MockParams) (*stubConfigurator, error) {
	client, err := newMockClient(params.ServerURL)
	if err != nil {
		return nil, err
	}

	return &stubConfigurator{
		client:       client,
		scenarioName: scenarioName,
		stateID:      0,
		currentState: startedScenarioState,
	}, nil
}

func (configurator *stubConfigurator) configureStub(
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

func (configurator *stubConfigurator) ConfigurePutStub(url string, params *mock.MockParams, setMetadataVerbFunc func(string), responseBody any) error {
	if setMetadataVerbFunc != nil {
		setMetadataVerbFunc(http.MethodPut)
	}
	if err := configurator.configureStub(ConfigurePutStub, url, params, nil, responseBody); err != nil {
		return err
	}
	return nil
}

func (configurator *stubConfigurator) ConfigurePostStub(url string, params *mock.MockParams, setMetadataVerbFunc func(string), responseBody any) error {
	if setMetadataVerbFunc != nil {
		setMetadataVerbFunc(http.MethodPost)
	}
	if err := configurator.configureStub(configurePostStub, url, params, nil, responseBody); err != nil {
		return err
	}
	return nil
}

func (configurator *stubConfigurator) ConfigureGetStub(url string, params *mock.MockParams, setMetadataVerbFunc func(string), responseBody any) error {
	if setMetadataVerbFunc != nil {
		setMetadataVerbFunc(http.MethodGet)
	}
	if err := configurator.configureStub(ConfigureGetStub, url, params, nil, responseBody); err != nil {
		return err
	}
	return nil
}

func (configurator *stubConfigurator) ConfigureGetListStub(url string, params *mock.MockParams, pathParams map[string]string, setMetadataVerbFunc func(string), responseBody any) error {
	if setMetadataVerbFunc != nil {
		setMetadataVerbFunc(http.MethodGet)
	}
	if err := configurator.configureStub(ConfigureGetStub, url, params, pathParams, responseBody); err != nil {
		return err
	}
	return nil
}

func (configurator *stubConfigurator) ConfigureGetNotFoundStub(url string, params *mock.MockParams) error {
	if err := configurator.configureStub(configureGetNotFoundStub, url, params, nil, nil); err != nil {
		return err
	}
	return nil
}

func (configurator *stubConfigurator) ConfigureDeleteStub(url string, params *mock.MockParams) error {
	if err := configurator.configureStub(configureDeleteStub, url, params, nil, nil); err != nil {
		return err
	}
	return nil
}

func (configurator *stubConfigurator) Finish() (*wiremock.Client, error) {
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
