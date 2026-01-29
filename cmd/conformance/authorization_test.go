package main

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites/authorization"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestAuthorizationV1Suites(t *testing.T) {
	globalTestSuite := suites.NewGlobalTestSuite(config.Parameters, config.Clients)

	// LifeCycle Suite
	lifeCycleTestSuite := authorization.NewLifeCycleV1TestSuite(globalTestSuite, config.Parameters.ScenariosUsers)
	if lifeCycleTestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, lifeCycleTestSuite)
	}

	// List Suite
	listTestSuite := authorization.NewListV1TestSuite(globalTestSuite, config.Parameters.ScenariosUsers)
	if listTestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, listTestSuite)
	}
}
