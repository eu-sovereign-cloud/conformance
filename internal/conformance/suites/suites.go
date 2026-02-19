package suites

import (
	"log/slog"
	"regexp"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/wiremock/go-wiremock"
)

// Test Suite

type TestSuite struct {
	suite.Suite
	Tenant        string
	AuthToken     string
	MockEnabled   bool
	MockServerURL *string

	MockClient   *wiremock.Client
	MockScenario *mockscenarios.Scenario
	ScenarioName string

	BaseDelay    int
	BaseInterval int
	MaxAttempts  int
}

func createTestSuite(params *config.ParametersHolder) *TestSuite {
	return &TestSuite{
		Tenant:        params.ClientTenant,
		AuthToken:     params.ClientAuthToken,
		MockEnabled:   params.MockEnabled,
		MockServerURL: &params.MockServerURL,
		BaseDelay:     params.BaseDelay,
		BaseInterval:  params.BaseInterval,
		MaxAttempts:   params.MaxAttempts,
	}
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

func (suite *TestSuite) ConfigureTags(t provider.T, provider string, kinds ...string) {
	t.Tags("provider:" + provider)
	for _, kind := range kinds {
		t.Tags("kind:" + kind)
	}
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

func SetupMockIfEnabledV2[P any](suite *TestSuite, configFunc func(*mockscenarios.Scenario, *P) error, suiteParams *P) error {
	// Setup mock, if configured to use
	if suite.MockEnabled {
		mockParams := &mock.MockParams{
			ServerURL: *suite.MockServerURL,
			AuthToken: suite.AuthToken,
		}

		suite.MockScenario = mockscenarios.NewScenario(suite.ScenarioName, mockParams)

		err := configFunc(suite.MockScenario, suiteParams)
		if err != nil {
			return err
		}
	}
	return nil
}

// Global Test Suite

type GlobalTestSuite struct {
	*TestSuite
	Client *secapi.GlobalClient
}

func CreateGlobalTestSuite(params *config.ParametersHolder, clients *config.ClientsHolder) GlobalTestSuite {
	return GlobalTestSuite{
		TestSuite: createTestSuite(params),
		Client:    clients.GlobalClient,
	}
}

// Regional Test Suite

type RegionalTestSuite struct {
	*TestSuite
	Region string
	Client *secapi.RegionalClient
}

func CreateRegionalTestSuite(params *config.ParametersHolder, clients *config.ClientsHolder) RegionalTestSuite {
	return RegionalTestSuite{
		TestSuite: createTestSuite(params),
		Region:    params.ClientRegion,
		Client:    clients.RegionalClient,
	}
}

// Mixed Test Suite

type MixedTestSuite struct {
	*TestSuite

	GlobalClient *secapi.GlobalClient

	Region         string
	RegionalClient *secapi.RegionalClient
}

func CreateMixedTestSuite(params *config.ParametersHolder, clients *config.ClientsHolder) MixedTestSuite {
	return MixedTestSuite{
		TestSuite:      createTestSuite(params),
		Region:         params.ClientRegion,
		GlobalClient:   clients.GlobalClient,
		RegionalClient: clients.RegionalClient,
	}
}
