package suites

import (
	"log/slog"
	"regexp"
	"strings"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/wiremock/go-wiremock"
)

type TestSuite struct {
	suite.Suite
	Tenant        string
	AuthToken     string
	MockEnabled   bool
	MockServerURL *string

	MockClient   *wiremock.Client
	ScenarioName string

	BaseDelay    int
	BaseInterval int
	MaxAttempts  int
}

func (suite *TestSuite) CanRun(regexp *regexp.Regexp) bool {
	if regexp == nil {
		return true
	} else {
		return regexp.MatchString(suite.ScenarioName)
	}
}

func (suite *TestSuite) StartScenario(t provider.T) {
	slog.Info("Starting execution of scenario " + suite.ScenarioName)
	t.Title(suite.ScenarioName)
}

func (suite *TestSuite) FinishScenario() {
	slog.Info("Finishing execution of scenario " + suite.ScenarioName)
}

type MixedTestSuite struct {
	TestSuite

	GlobalClient *secapi.GlobalClient

	Region         string
	RegionalClient *secapi.RegionalClient
}

type GlobalTestSuite struct {
	TestSuite
	Client *secapi.GlobalClient
}

type RegionalTestSuite struct {
	TestSuite
	Region string
	Client *secapi.RegionalClient
}

func (suite *TestSuite) ConfigureTags(t provider.T, provider string, kinds ...string) {
	t.Tags(
		"provider:"+provider,
		"resources:"+strings.Join(kinds, ", "),
	)
}

func (suite *TestSuite) ResetAllScenarios() {
	// Cleanup configured mock scenarios
	if suite.MockClient != nil {
		if err := suite.MockClient.ResetAllScenarios(); err != nil {
			slog.Error("Failed to reset scenarios", "error", err)
		}
	}
}

func SetupMockIfEnabled[P any](suite *TestSuite, configFunc func(string, *mock.MockParams, *P) (*wiremock.Client, error), suiteParams *P) error {
	// Setup mock, if configured to use
	if suite.MockEnabled {
		mockParams := &mock.MockParams{
			ServerURL: *suite.MockServerURL,
			AuthToken: suite.AuthToken,
		}

		wm, err := configFunc(suite.ScenarioName, mockParams, suiteParams)
		if err != nil {
			return err
		}
		suite.MockClient = wm
	}
	return nil
}
