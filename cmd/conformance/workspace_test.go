package main

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites/workspace"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestWorkspaceV1Suites(t *testing.T) {
	regionalTestSuite := createRegionalTestSuite(config.Parameters, config.Clients)

	// LifeCycle Suite
	lifeCycleTestSuite := &workspace.WorkspaceLifeCycleV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
	}
	lifeCycleTestSuite.ScenarioName = constants.WorkspaceV1LifeCycleSuiteName
	if lifeCycleTestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, lifeCycleTestSuite)
	}

	// List Suite
	listTestSuite := &workspace.WorkspaceListV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
	}
	listTestSuite.ScenarioName = constants.WorkspaceV1ListSuiteName
	if listTestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, listTestSuite)
	}
}
