package mock

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"text/template"

	"github.com/wiremock/go-wiremock"
)

func configureStub(wm *wiremock.Client, method string, name string, config stubConfig) error {
	// Build the response
	response := wiremock.NewResponse().WithStatus(int64(config.httpStatus))
	if config.template != "" {
		// Response with body
		processTemplate, err := processTemplate(config.template, config.response)
		if err != nil {
			return err
		}
		response = response.
			WithHeader(contentTypeHttpHeaderKey, contentTypeHttpHeaderValue).
			WithJSONBody(processTemplate)
	}

	params := config.params.getParams()

	// Request matchers
	urlMatcher := wiremock.URLPathMatching(config.url)
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
	if config.priority != 0 {
		priority = config.priority
	}

	if err := wm.StubFor(stubRule.
		WithHeader(authorizationHttpHeaderKey, headerMatcher).
		InScenario(name).
		WhenScenarioStateIs(config.currentState).
		WillSetStateTo(config.nextState).
		WillReturnResponse(response).
		AtPriority(int64(priority))); err != nil {
		slog.Error("Error configuring stub", "method", method, "error", err)
		return err
	}
	return nil
}

func configureWithoutStateStub(wm *wiremock.Client, method string, name string, config stubConfig) error {
	// Build the response
	response := wiremock.NewResponse().WithStatus(int64(config.httpStatus))
	if config.template != "" {
		// Response with body
		processTemplate, err := processTemplate(config.template, config.response)
		if err != nil {
			return err
		}
		response = response.
			WithHeader(contentTypeHttpHeaderKey, contentTypeHttpHeaderValue).
			WithJSONBody(processTemplate)
	}

	params := config.params.getParams()

	// Request matchers
	urlMatcher := wiremock.URLPathMatching(config.url)
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
	if config.priority != 0 {
		priority = config.priority
	}

	if err := wm.StubFor(stubRule.
		WithHeader(authorizationHttpHeaderKey, headerMatcher).
		InScenario(name).
		WillReturnResponse(response).
		AtPriority(int64(priority))); err != nil {
		slog.Error("Error configuring stub", "method", method, "error", err)
		return err
	}
	return nil
}

func configurePutStub(wm *wiremock.Client, scenarioName string, scenarioConfig stubConfig) error {
	return configureStub(wm, http.MethodPut, scenarioName, scenarioConfig)
}

func configurePostStub(wm *wiremock.Client, scenarioName string, scenarioConfig stubConfig) error {
	return configureStub(wm, http.MethodPost, scenarioName, scenarioConfig)
}

func configureGetStub(wm *wiremock.Client, scenarioName string, scenarioConfig stubConfig) error {
	return configureStub(wm, http.MethodGet, scenarioName, scenarioConfig)
}

func configureGetWithoutStateStub(wm *wiremock.Client, scenarioName string, scenarioConfig stubConfig) error {
	return configureWithoutStateStub(wm, http.MethodGet, scenarioName, scenarioConfig)
}

func configureDeleteStub(wm *wiremock.Client, scenarioName string, scenarioConfig stubConfig) error {
	return configureStub(wm, http.MethodDelete, scenarioName, scenarioConfig)
}

func processTemplate(templ string, data any) (any, error) {
	tmpl := template.Must(template.New("response").Delims("[[", "]]").Parse(templ))

	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, data); err != nil {
		return nil, err
	}
	var dataJsonMap any
	err := json.Unmarshal(buffer.Bytes(), &dataJsonMap)
	if err != nil {
		return nil, err
	}

	return dataJsonMap, nil
}
