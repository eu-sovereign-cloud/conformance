package main

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites/authorization"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestAuthorizationV1Suites(t *testing.T) {
	globalTestSuite := suites.CreateGlobalTestSuite(config.Parameters, config.Clients)

	// Provider LifeCycle Suite
	providerLifeCycleTestSuite := authorization.CreateProviderLifeCycleV1TestSuite(globalTestSuite, config.Parameters.ScenariosUsers)
	if providerLifeCycleTestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, providerLifeCycleTestSuite)
	}

	// List Suite
	listTestSuite := authorization.CreateProviderQueriesV1TestSuite(globalTestSuite, config.Parameters.ScenariosUsers)
	if listTestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, listTestSuite)
	}
}
