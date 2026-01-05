package stubs

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"

	"github.com/wiremock/go-wiremock"
)

type stubConfig struct {
	url          string
	httpMethod   string
	httpStatus   int
	params       *mock.BaseParams
	pathParams   map[string]string
	responseBody any
	currentState string
	nextState    string
	priority     int
}

func configureStub(wm *wiremock.Client, scenarioName string, stubConfig *stubConfig) error {
	// Build the stubResponse
	stubResponse := wiremock.NewResponse().WithStatus(int64(stubConfig.httpStatus))

	// Parse the response body
	if stubConfig.responseBody != nil {
		responseBody, err := parseResponseBody(stubConfig.responseBody)
		if err != nil {
			return err
		}

		stubResponse = stubResponse.
			WithHeader(contentTypeHttpHeaderKey, contentTypeHttpHeaderValue).
			WithBody(responseBody)
	}

	// Request matchers
	urlMatcher := wiremock.URLPathMatching(stubConfig.url)

	// Configure the stub
	var stubRule *wiremock.StubRule
	switch stubConfig.httpMethod {
	case http.MethodPut:
		stubRule = wiremock.Put(urlMatcher)
	case http.MethodPost:
		stubRule = wiremock.Post(urlMatcher)
	case http.MethodGet:
		stubRule = wiremock.Get(urlMatcher)
	case http.MethodDelete:
		stubRule = wiremock.Delete(urlMatcher)
	}

	priority := defaultScenarioPriority
	if stubConfig.priority != 0 {
		priority = stubConfig.priority
	}

	stubRule.WithHeader(authorizationHttpHeaderKey, wiremock.Matching(authorizationHttpHeaderValuePrefix+stubConfig.params.AuthToken))

	for key, value := range stubConfig.pathParams {
		matcher := wiremock.Matching(value)
		stubRule.WithPathParam(key, matcher)
	}
	// Create a stub with scenario state if currentState it's defined
	if stubConfig.currentState != "" {
		if err := wm.StubFor(stubRule.
			InScenario(scenarioName).
			WhenScenarioStateIs(stubConfig.currentState).
			WillSetStateTo(stubConfig.nextState).
			WillReturnResponse(stubResponse).
			AtPriority(int64(priority))); err != nil {
			slog.Error("Error configuring stub", "method", stubConfig.httpMethod, "error", err)
			return err
		}
	} else {
		if err := wm.StubFor(stubRule.
			InScenario(scenarioName).
			WillReturnResponse(stubResponse).
			AtPriority(int64(priority))); err != nil {
			slog.Error("Error configuring stub", "method", stubConfig.httpMethod, "error", err)
			return err
		}
	}
	return nil
}

func ConfigurePutStub(wm *wiremock.Client, scenarioName string, stubConfig *stubConfig) error {
	stubConfig.httpMethod = http.MethodPut
	stubConfig.httpStatus = http.StatusCreated
	return configureStub(wm, scenarioName, stubConfig)
}

func configurePostStub(wm *wiremock.Client, scenarioName string, stubConfig *stubConfig) error {
	stubConfig.httpMethod = http.MethodPost
	stubConfig.httpStatus = http.StatusAccepted
	return configureStub(wm, scenarioName, stubConfig)
}

func ConfigureGetStub(wm *wiremock.Client, scenarioName string, stubConfig *stubConfig) error {
	stubConfig.httpMethod = http.MethodGet
	stubConfig.httpStatus = http.StatusOK
	return configureStub(wm, scenarioName, stubConfig)
}

func configureGetNotFoundStub(wm *wiremock.Client, scenarioName string, stubConfig *stubConfig) error {
	stubConfig.httpMethod = http.MethodGet
	stubConfig.httpStatus = http.StatusNotFound
	return configureStub(wm, scenarioName, stubConfig)
}

func configureDeleteStub(wm *wiremock.Client, scenarioName string, stubConfig *stubConfig) error {
	stubConfig.httpMethod = http.MethodDelete
	stubConfig.httpStatus = http.StatusAccepted
	return configureStub(wm, scenarioName, stubConfig)
}

func parseResponseBody(configResponse any) (string, error) {
	responseBytes, err := json.Marshal(configResponse)
	if err != nil {
		slog.Error("Error parsing response to JSON", "error", err)
		return "", err
	}
	return string(responseBytes), nil
}
