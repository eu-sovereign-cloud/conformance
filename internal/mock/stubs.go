package mock

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/wiremock/go-wiremock"
)

func configureStub(wm *wiremock.Client, scenarioName string, method string, httpStatus int, stubConfig *stubConfig) error {
	// Build the response
	response := wiremock.NewResponse().WithStatus(int64(httpStatus))

	// Response with body
	if stubConfig.response != nil {
		responseBytes, err := json.Marshal(stubConfig.response)
		if err != nil {
			slog.Error("Error parsing response to JSON", "error", err)
			return err
		}
		responseJSON := string(responseBytes)

		response = response.
			WithHeader(contentTypeHttpHeaderKey, contentTypeHttpHeaderValue).
			WithJSONBody(responseJSON)
	}

	params := stubConfig.params.getParams()

	// Request matchers
	urlMatcher := wiremock.URLPathMatching(stubConfig.url)
	headerMatcher := wiremock.Matching(authorizationHttpHeaderValuePrefix + params.AuthToken)

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

	if err := wm.StubFor(stubRule.
		WithHeader(authorizationHttpHeaderKey, headerMatcher).
		InScenario(scenarioName).
		WhenScenarioStateIs(stubConfig.currentState).
		WillSetStateTo(stubConfig.nextState).
		WillReturnResponse(response).
		AtPriority(int64(priority))); err != nil {
		slog.Error("Error configuring stub", "method", method, "error", err)
		return err
	}
	return nil
}

func configurePutStub(wm *wiremock.Client, scenarioName string, httpStatus int, stubConfig *stubConfig) error {
	return configureStub(wm, scenarioName, http.MethodPut, httpStatus, stubConfig)
}

func configurePutSuccessStub(wm *wiremock.Client, scenarioName string, stubConfig *stubConfig) error {
	return configureStub(wm, scenarioName, http.MethodPut, http.StatusCreated, stubConfig)
}

func configurePostStub(wm *wiremock.Client, scenarioName string, httpStatus int, stubConfig *stubConfig) error {
	return configureStub(wm, scenarioName, http.MethodPost, httpStatus, stubConfig)
}

func configurePostSuccessStub(wm *wiremock.Client, scenarioName string, stubConfig *stubConfig) error {
	return configureStub(wm, scenarioName, http.MethodPost, http.StatusAccepted, stubConfig)
}

func configureGetStub(wm *wiremock.Client, scenarioName string, httpStatus int, stubConfig *stubConfig) error {
	return configureStub(wm, scenarioName, http.MethodGet, httpStatus, stubConfig)
}

func configureGetSuccessStub(wm *wiremock.Client, scenarioName string, stubConfig *stubConfig) error {
	return configureStub(wm, scenarioName, http.MethodGet, http.StatusOK, stubConfig)
}

func configureDeleteStub(wm *wiremock.Client, scenarioName string, httpStatus int, stubConfig *stubConfig) error {
	return configureStub(wm, scenarioName, http.MethodDelete, httpStatus, stubConfig)
}

func configureDeleteSuccessStub(wm *wiremock.Client, scenarioName string, stubConfig *stubConfig) error {
	return configureStub(wm, scenarioName, http.MethodDelete, http.StatusAccepted, stubConfig)
}
