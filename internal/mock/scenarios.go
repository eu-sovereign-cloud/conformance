package mock

import (
	"net/http"
	"strconv"

	"github.com/wiremock/go-wiremock"
)

const (
	statePrefix = "State."
)

type scenarioConfigurator struct {
	client       *wiremock.Client
	scenarioName string
	stateID      int
	currentState string
}

func newScenarioConfigurator(scenarioName, mockURL string) (*scenarioConfigurator, error) {
	client, err := newMockClient(mockURL)
	if err != nil {
		return nil, err
	}

	return &scenarioConfigurator{
		client:       client,
		scenarioName: scenarioName,
		stateID:      0,
		currentState: startedScenarioState,
	}, nil
}

func (builder *scenarioConfigurator) configureStub(
	stubFunc func(*wiremock.Client, string, *stubConfig) error,
	url string,
	params HasParams,
	pathParams map[string]string,
	responseBody any,
) error {
	// Calculte next state
	stateID := builder.stateID + 1
	nextState := statePrefix + strconv.Itoa(stateID)

	if err := stubFunc(builder.client, builder.scenarioName,
		&stubConfig{url: url, params: params, pathParams: pathParams, responseBody: responseBody, currentState: builder.currentState, nextState: nextState}); err != nil {
		return err
	}

	builder.stateID = stateID
	builder.currentState = nextState
	return nil
}

func (builder *scenarioConfigurator) configurePutStub(url string, params HasParams, setMetadataVerbFunc func(string), responseBody any) error {
	if setMetadataVerbFunc != nil {
		setMetadataVerbFunc(http.MethodPut)
	}
	if err := builder.configureStub(configurePutStub, url, params, nil, responseBody); err != nil {
		return err
	}
	return nil
}

func (builder *scenarioConfigurator) configurePostStub(url string, params HasParams, setMetadataVerbFunc func(string), responseBody any) error {
	if setMetadataVerbFunc != nil {
		setMetadataVerbFunc(http.MethodPost)
	}
	if err := builder.configureStub(configurePostStub, url, params, nil, responseBody); err != nil {
		return err
	}
	return nil
}

func (builder *scenarioConfigurator) configureGetStub(url string, params HasParams, setMetadataVerbFunc func(string), responseBody any) error {
	if setMetadataVerbFunc != nil {
		setMetadataVerbFunc(http.MethodGet)
	}
	if err := builder.configureStub(configureGetStub, url, params, nil, responseBody); err != nil {
		return err
	}
	return nil
}

func (builder *scenarioConfigurator) configureGetListStub(url string, params HasParams, pathParams map[string]string, setMetadataVerbFunc func(string), responseBody any) error {
	if setMetadataVerbFunc != nil {
		setMetadataVerbFunc(http.MethodGet)
	}
	if err := builder.configureStub(configureGetStub, url, params, pathParams, responseBody); err != nil {
		return err
	}
	return nil
}

func (builder *scenarioConfigurator) configureGetNotFoundStub(url string, params HasParams) error {
	if err := builder.configureStub(configureGetNotFoundStub, url, params, nil, nil); err != nil {
		return err
	}
	return nil
}

func (builder *scenarioConfigurator) configureDeleteStub(url string, params HasParams) error {
	if err := builder.configureStub(configureDeleteStub, url, params, nil, nil); err != nil {
		return err
	}
	return nil
}

func newMockClient(mockURL string) (*wiremock.Client, error) {
	wm := wiremock.NewClient(mockURL)
	if err := wm.ResetAllScenarios(); err != nil {
		return nil, err
	}
	return wm, nil
}
