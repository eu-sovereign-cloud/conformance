package main

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites/workspace"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestWorkspaceV1Suites(t *testing.T) {
	regionalTest := suites.RegionalTestSuite{
		TestSuite: suites.TestSuite{
			Tenant:        config.Parameters.ClientTenant,
			AuthToken:     config.Parameters.ClientAuthToken,
			MockEnabled:   config.Parameters.MockEnabled,
			MockServerURL: &config.Parameters.MockServerURL,
			ScenarioName:  constants.WorkspaceV1LifeCycleSuiteName,
			BaseDelay:     config.Parameters.BaseDelay,
			BaseInterval:  config.Parameters.BaseInterval,
			MaxAttempts:   config.Parameters.MaxAttempts,
		},
		Region: config.Parameters.ClientRegion,
		Client: config.Clients.RegionalClient,
	}

	// LifeCycle Suite
	testLifeCycleSuite := &workspace.WorkspaceV1LifeCycleTestSuite{
		RegionalTestSuite: regionalTest,
	}
	testLifeCycleSuite.ScenarioName = constants.WorkspaceV1LifeCycleSuiteName
	if testLifeCycleSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, testLifeCycleSuite)
	}

	// List Suite
	testListSuite := &workspace.WorkspaceV1ListTestSuite{
		RegionalTestSuite: regionalTest,
	}
	testListSuite.ScenarioName = constants.WorkspaceV1ListSuiteName
	if testListSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, testListSuite)
	}
}
