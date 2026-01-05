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

type scenarioConfigurator struct {
	Client       *wiremock.Client
	scenarioName string
	stateID      int
	currentState string
}

func NewScenarioConfigurator(scenarioName, mockURL string) (*scenarioConfigurator, error) {
	client, err := newMockClient(mockURL)
	if err != nil {
		return nil, err
	}

	return &scenarioConfigurator{
		Client:       client,
		scenarioName: scenarioName,
		stateID:      0,
		currentState: startedScenarioState,
	}, nil
}

func (builder *scenarioConfigurator) ConfigureStub(
	stubFunc func(*wiremock.Client, string, *stubConfig) error,
	url string,
	params *mock.BaseParams,
	pathParams map[string]string,
	responseBody any,
) error {
	// Calculte next state
	stateID := builder.stateID + 1
	nextState := statePrefix + strconv.Itoa(stateID)

	if err := stubFunc(builder.Client, builder.scenarioName,
		&stubConfig{url: url, params: params, pathParams: pathParams, responseBody: responseBody, currentState: builder.currentState, nextState: nextState}); err != nil {
		return err
	}

	builder.stateID = stateID
	builder.currentState = nextState
	return nil
}

func (builder *scenarioConfigurator) ConfigurePutStub(url string, params *mock.BaseParams, setMetadataVerbFunc func(string), responseBody any) error {
	if setMetadataVerbFunc != nil {
		setMetadataVerbFunc(http.MethodPut)
	}
	if err := builder.ConfigureStub(ConfigurePutStub, url, params, nil, responseBody); err != nil {
		return err
	}
	return nil
}

func (builder *scenarioConfigurator) ConfigurePostStub(url string, params *mock.BaseParams, setMetadataVerbFunc func(string), responseBody any) error {
	if setMetadataVerbFunc != nil {
		setMetadataVerbFunc(http.MethodPost)
	}
	if err := builder.ConfigureStub(configurePostStub, url, params, nil, responseBody); err != nil {
		return err
	}
	return nil
}

func (builder *scenarioConfigurator) ConfigureGetStub(url string, params *mock.BaseParams, setMetadataVerbFunc func(string), responseBody any) error {
	if setMetadataVerbFunc != nil {
		setMetadataVerbFunc(http.MethodGet)
	}
	if err := builder.ConfigureStub(ConfigureGetStub, url, params, nil, responseBody); err != nil {
		return err
	}
	return nil
}

func (builder *scenarioConfigurator) ConfigureGetListStub(url string, params *mock.BaseParams, pathParams map[string]string, setMetadataVerbFunc func(string), responseBody any) error {
	if setMetadataVerbFunc != nil {
		setMetadataVerbFunc(http.MethodGet)
	}
	if err := builder.ConfigureStub(ConfigureGetStub, url, params, pathParams, responseBody); err != nil {
		return err
	}
	return nil
}

func (builder *scenarioConfigurator) ConfigureGetNotFoundStub(url string, params *mock.BaseParams) error {
	if err := builder.ConfigureStub(configureGetNotFoundStub, url, params, nil, nil); err != nil {
		return err
	}
	return nil
}

func (builder *scenarioConfigurator) ConfigureDeleteStub(url string, params *mock.BaseParams) error {
	if err := builder.ConfigureStub(configureDeleteStub, url, params, nil, nil); err != nil {
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
