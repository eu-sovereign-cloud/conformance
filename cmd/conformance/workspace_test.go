//nolint:dupl
package main

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites/workspace"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
)

func TestWorkspaceV1Suites(t *testing.T) {
	regionalTestSuite := suites.CreateRegionalTestSuite(config.Parameters, config.Clients)

	// LifeCycle Suite
	lifeCycleTestSuite := &workspace.WorkspaceLifeCycleV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
	}
	lifeCycleTestSuite.RunSuite(t, config.Parameters.ScenariosRegexp,
		func() { lifeCycleTestSuite.ScenarioName = constants.WorkspaceV1LifeCycleSuiteName },
	)

	// List Suite
	listTestSuite := &workspace.WorkspaceListV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
	}
	listTestSuite.RunSuite(t, config.Parameters.ScenariosRegexp,
		func() { listTestSuite.ScenarioName = constants.WorkspaceV1ListSuiteName },
	)
}
