package main

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites/authorization"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestAuthorizationV1Suites(t *testing.T) {
	globalTestSuite := suites.GlobalTestSuite{
		TestSuite: suites.TestSuite{
			Tenant:        config.Parameters.ClientTenant,
			AuthToken:     config.Parameters.ClientAuthToken,
			MockEnabled:   config.Parameters.MockEnabled,
			MockServerURL: &config.Parameters.MockServerURL,
			BaseDelay:     config.Parameters.BaseDelay,
			BaseInterval:  config.Parameters.BaseInterval,
			MaxAttempts:   config.Parameters.MaxAttempts,
		},
		Client: config.Clients.GlobalClient,
	}

	// LifeCycle Suite
	testLifeCycleSuite := &authorization.LifeCycleV1TestSuite{
		GlobalTestSuite: globalTestSuite,
		Users:           config.Parameters.ScenariosUsers,
	}
	testLifeCycleSuite.ScenarioName = constants.AuthorizationV1LifeCycleSuiteName
	if testLifeCycleSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, testLifeCycleSuite)
	}

	// List Suite
	testListSuite := &authorization.ListV1TestSuite{
		GlobalTestSuite: globalTestSuite,
		Users:           config.Parameters.ScenariosUsers,
	}
	testListSuite.ScenarioName = constants.AuthorizationV1ListSuiteName
	if testListSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, testListSuite)
	}
}
