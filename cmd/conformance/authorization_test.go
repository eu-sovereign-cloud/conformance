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
	providerLifeCycleSuite := authorization.CreateProviderLifeCycleV1TestSuite(globalTestSuite, config.Parameters.ScenariosUsers)
	if providerLifeCycleSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, providerLifeCycleSuite)
	}

	// Provider Queries Suite
	providerQueriesSuite := authorization.CreateProviderQueriesV1TestSuite(globalTestSuite, config.Parameters.ScenariosUsers)
	if providerQueriesSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, providerQueriesSuite)
	}

	// Role LifeCycle Suite
	roleLifeCycleSuite := authorization.CreateRoleLifeCycleV1TestSuite(globalTestSuite)
	if roleLifeCycleSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, roleLifeCycleSuite)
	}

	// Role Assignment LifeCycle Suite
	roleAssignmentLifeCycleSuite := authorization.CreateRoleAssignmentLifeCycleV1TestSuite(globalTestSuite, config.Parameters.ScenariosUsers)
	if roleAssignmentLifeCycleSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, roleAssignmentLifeCycleSuite)
	}
}
