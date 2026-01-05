package workspace

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestWorkspaceV1LifeCycleSuite(t *testing.T) {
	regionalTest := suites.RegionalTestSuite{
		TestSuite: suites.TestSuite{
			Tenant:        conformance.Config.ClientTenant,
			AuthToken:     conformance.Config.ClientAuthToken,
			MockEnabled:   conformance.Config.MockEnabled,
			MockServerURL: &conformance.Config.MockServerURL,
			ScenarioName:  conformance.WorkspaceV1LifeCycleSuiteName,
			BaseDelay:     conformance.Config.BaseDelay,
			BaseInterval:  conformance.Config.BaseInterval,
			MaxAttempts:   conformance.Config.MaxAttempts,
		},
		Region: conformance.Config.ClientRegion,
		Client: conformance.Clients.RegionalClient,
	}

	// LifeCycle Suite
	testLifeCycleSuite := &WorkspaceV1LifeCycleTestSuite{
		RegionalTestSuite: regionalTest,
	}

	if testLifeCycleSuite.CanRun(conformance.Config.ScenariosRegexp) {
		suite.RunSuite(t, testLifeCycleSuite)
	}

	// List Suite
	testListSuite := &WorkspaceV1ListTestSuite{
		RegionalTestSuite: regionalTest,
	}
	testListSuite.RegionalTestSuite.TestSuite.ScenarioName = conformance.WorkspaceV1ListSuiteName
	if testListSuite.CanRun(conformance.Config.ScenariosRegexp) {
		suite.RunSuite(t, testListSuite)
	}
}
