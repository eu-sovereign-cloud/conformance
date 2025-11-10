package mock

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/wiremock/go-wiremock"
)

func parseResponseBody(configResponse any) (string, error) {
	responseBytes, err := json.Marshal(configResponse)
	if err != nil {
		slog.Error("Error parsing response to JSON", "error", err)
		return "", err
	}
	return string(responseBytes), nil
}

func configureStub(wm *wiremock.Client, scenarioName string, method string, httpStatus int, stubConfig *stubConfig) error {
	// Build the stubResponse
	stubResponse := wiremock.NewResponse().WithStatus(int64(httpStatus))

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
	switch method {
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

	for key, value := range stubConfig.headers {
		matcher := wiremock.Matching(value)
		stubRule.WithHeader(key, matcher)
	}
	// Create a stub with scenario state if currentState it's defined
	if stubConfig.currentState != "" {
		if err := wm.StubFor(stubRule.
			InScenario(scenarioName).
			WhenScenarioStateIs(stubConfig.currentState).
			WillSetStateTo(stubConfig.nextState).
			WillReturnResponse(stubResponse).
			AtPriority(int64(priority))); err != nil {
			slog.Error("Error configuring stub", "method", method, "error", err)
			return err
		}
	} else {
		if err := wm.StubFor(stubRule.
			InScenario(scenarioName).
			WillReturnResponse(stubResponse).
			AtPriority(int64(priority))); err != nil {
			slog.Error("Error configuring stub", "method", method, "error", err)
			return err
		}
	}
	return nil
}

func configurePutStub(wm *wiremock.Client, scenarioName string, stubConfig *stubConfig) error {
	return configureStub(wm, scenarioName, http.MethodPut, http.StatusCreated, stubConfig)
}

func configurePostStub(wm *wiremock.Client, scenarioName string, stubConfig *stubConfig) error {
	return configureStub(wm, scenarioName, http.MethodPost, http.StatusAccepted, stubConfig)
}

func configureGetStub(wm *wiremock.Client, scenarioName string, stubConfig *stubConfig) error {
	return configureStub(wm, scenarioName, http.MethodGet, http.StatusOK, stubConfig)
}

func configureGetStubWithStatus(wm *wiremock.Client, scenarioName string, httpStatus int, stubConfig *stubConfig) error {
	return configureStub(wm, scenarioName, http.MethodGet, httpStatus, stubConfig)
}

func configureDeleteStub(wm *wiremock.Client, scenarioName string, stubConfig *stubConfig) error {
	return configureStub(wm, scenarioName, http.MethodDelete, http.StatusAccepted, stubConfig)
}
