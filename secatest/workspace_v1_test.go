package secatest

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestWorkspaceV1LifeCycleSuite(t *testing.T) {
	regionalTest := regionalTestSuite{
		testSuite: testSuite{
			tenant:        config.clientTenant,
			authToken:     config.clientAuthToken,
			mockEnabled:   config.mockEnabled,
			mockServerURL: &config.mockServerURL,
			scenarioName:  workspaceV1LifeCycleSuiteName,
			baseDelay:     config.baseDelay,
			baseInterval:  config.baseInterval,
			maxAttempts:   config.maxAttempts,
		},
		region: config.clientRegion,
		client: clients.regionalClient,
	}

	// LifeCycle Suite
	testLifeCycleSuite := &WorkspaceV1LifeCycleTestSuite{
		regionalTestSuite: regionalTest,
	}

	if testLifeCycleSuite.canRun(config.scenariosRegexp) {
		suite.RunSuite(t, testLifeCycleSuite)
	}

	// List Suite
	testListSuite := &WorkspaceV1ListTestSuite{
		regionalTestSuite: regionalTest,
	}
	testListSuite.regionalTestSuite.testSuite.scenarioName = workspaceV1ListSuiteName
	if testListSuite.canRun(config.scenariosRegexp) {
		suite.RunSuite(t, testListSuite)
	}
}
