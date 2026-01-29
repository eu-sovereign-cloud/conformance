package main

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites/workspace"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestWorkspaceV1Suites(t *testing.T) {
	regionalTestSuite := suites.NewRegionalTestSuite(config.Parameters, config.Clients)

	// LifeCycle Suite
	lifeCycleTestSuite := workspace.NewLifeCycleV1TestSuite(regionalTestSuite)
	if lifeCycleTestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, lifeCycleTestSuite)
	}

	// List Suite
	listTestSuite := workspace.NewListV1TestSuite(regionalTestSuite)
	if listTestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, listTestSuite)
	}

	// Create Workspace Suite
	createWorkspaceTestSuite := workspace.NewCreateWorkspaceV1TestSuite(regionalTestSuite)
	if createWorkspaceTestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, createWorkspaceTestSuite)
	}
}
