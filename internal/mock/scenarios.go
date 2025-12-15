package mock

import (
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
	restart bool,
) error {
	// Calculte next state
	var stateID int
	var nextState string
	if restart {
		stateID = 0
		nextState = startedScenarioState
	} else {
		stateID = builder.stateID + 1
		nextState = statePrefix + strconv.Itoa(stateID)
	}

	if err := stubFunc(builder.client, builder.scenarioName,
		&stubConfig{url: url, params: params, pathParams: pathParams, responseBody: responseBody, currentState: builder.currentState, nextState: nextState}); err != nil {
		return err
	}

	builder.stateID = stateID
	builder.currentState = nextState
	return nil
}

func (builder *scenarioConfigurator) configurePutStub(url string, params HasParams, responseBody any, restart bool) error {
	if err := builder.configureStub(configurePutStub, url, params, nil, responseBody, restart); err != nil {
		return err
	}
	return nil
}

func (builder *scenarioConfigurator) configurePostStub(url string, params HasParams, responseBody any, restart bool) error {
	if err := builder.configureStub(configurePostStub, url, params, nil, responseBody, restart); err != nil {
		return err
	}
	return nil
}

func (builder *scenarioConfigurator) configureGetStub(url string, params HasParams, responseBody any, restart bool) error {
	if err := builder.configureStub(configureGetStub, url, params, nil, responseBody, restart); err != nil {
		return err
	}
	return nil
}
func (builder *scenarioConfigurator) configureGetListStub(url string, params HasParams, pathParams map[string]string, responseBody any, restart bool) error {
	if err := builder.configureStub(configureGetStub, url, params, pathParams, responseBody, restart); err != nil {
		return err
	}
	return nil
}
func (builder *scenarioConfigurator) configureGetStubWithHeaders(url string, params HasParams, headers map[string]string, responseBody any, restart bool) error {
	if err := builder.configureStub(configureGetStub, url, params, headers, responseBody, restart); err != nil {
		return err
	}
	return nil
}

func (builder *scenarioConfigurator) configureGetNotFoundStub(url string, params HasParams, restart bool) error {
	if err := builder.configureStub(configureGetNotFoundStub, url, params, nil, nil, restart); err != nil {
		return err
	}
	return nil
}

func (builder *scenarioConfigurator) configureDeleteStub(url string, params HasParams, restart bool) error {
	if err := builder.configureStub(configureDeleteStub, url, params, nil, nil, restart); err != nil {
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
