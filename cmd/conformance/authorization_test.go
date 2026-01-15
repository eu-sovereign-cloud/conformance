package main

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites/authorization"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestAuthorizationV1Suites(t *testing.T) {
	globalTestSuite := createGlobalTestSuite(config.Parameters, config.Clients)

	// LifeCycle Suite
	lifeCycleTestSuite := &authorization.AuthorizationLifeCycleV1TestSuite{
		GlobalTestSuite: globalTestSuite,
		Users:           config.Parameters.ScenariosUsers,
	}
	lifeCycleTestSuite.ScenarioName = constants.AuthorizationV1LifeCycleSuiteName
	if lifeCycleTestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, lifeCycleTestSuite)
	}

	// List Suite
	listTestSuite := &authorization.AuthorizationListV1TestSuite{
		GlobalTestSuite: globalTestSuite,
		Users:           config.Parameters.ScenariosUsers,
	}
	listTestSuite.ScenarioName = constants.AuthorizationV1ListSuiteName
	if listTestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, listTestSuite)
	}
}
