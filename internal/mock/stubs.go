package mock

import (
	"log/slog"

	"github.com/wiremock/go-wiremock"
)

func configurePutStub(wm *wiremock.Client, scenarioConfig scenarioConfig) error {
	processTemplate, err := processTemplate(scenarioConfig.template, scenarioConfig.response)
	if err != nil {
		return err
	}

	err = wm.StubFor(wiremock.Put(wiremock.URLPathMatching(scenarioConfig.params.MockURL)).
		WithHeader(authorizationHttpHeaderKey, wiremock.Matching(authorizationHttpHeaderValuePrefix+scenarioConfig.params.AuthToken)).
		InScenario(scenarioConfig.name).
		WhenScenarioStateIs(scenarioConfig.currentState).
		WillSetStateTo(scenarioConfig.nextState).
		WillReturnResponse(
			wiremock.NewResponse().
				WithStatus(int64(scenarioConfig.httpStatus)).
				WithHeader(contentTypeHttpHeaderKey, contentTypeHttpHeaderValue).
				WithJSONBody(processTemplate),
		).
		AtPriority(int64(scenarioConfig.priority)))
	if err != nil {
		slog.Error("Error configuring put stub", "error", err)
		return err
	}
	return nil
}

func configureGetStub(wm *wiremock.Client, scenarioConfig scenarioConfig) error {
	var response wiremock.Response

	if scenarioConfig.template != "" {
		processTemplate, err := processTemplate(scenarioConfig.template, scenarioConfig.response)
		if err != nil {
			return err
		}

		response = wiremock.NewResponse().
			WithStatus(int64(scenarioConfig.httpStatus)).
			WithHeader(contentTypeHttpHeaderKey, contentTypeHttpHeaderValue).
			WithJSONBody(processTemplate)
	} else {
		response = wiremock.NewResponse().
			WithStatus(int64(scenarioConfig.httpStatus))
	}

	err := wm.StubFor(wiremock.Get(wiremock.URLPathMatching(scenarioConfig.params.MockURL)).
		WithHeader(authorizationHttpHeaderKey, wiremock.Matching(authorizationHttpHeaderValuePrefix+scenarioConfig.params.AuthToken)).
		InScenario(scenarioConfig.name).
		WhenScenarioStateIs(scenarioConfig.currentState).
		WillSetStateTo(scenarioConfig.nextState).
		WillReturnResponse(response).
		AtPriority(int64(scenarioConfig.priority)))
	if err != nil {
		slog.Error("Error configuring get stub", "error", err)
		return err
	}
	return nil
}

func configureDeleteStub(wm *wiremock.Client, scenarioConfig scenarioConfig) error {
	err := wm.StubFor(wiremock.Delete(wiremock.URLPathMatching(scenarioConfig.params.MockURL)).
		WithHeader(authorizationHttpHeaderKey, wiremock.Matching(authorizationHttpHeaderValuePrefix+scenarioConfig.params.AuthToken)).
		InScenario(scenarioConfig.name).
		WhenScenarioStateIs(scenarioConfig.currentState).
		WillSetStateTo(scenarioConfig.nextState).
		WillReturnResponse(
			wiremock.NewResponse().
				WithStatus(int64(scenarioConfig.httpStatus)),
		).
		AtPriority(int64(scenarioConfig.priority)))
	if err != nil {
		slog.Error("Error configuring delete stub", "error", err)
		return err
	}
	return nil
}
