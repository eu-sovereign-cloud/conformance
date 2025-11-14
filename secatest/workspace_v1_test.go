package secatest

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestWorkspaceV1Suite(t *testing.T) {
	testSuite := &WorkspaceV1TestSuite{
		regionalTestSuite: regionalTestSuite{
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
		},
	}

	if testSuite.canRun(config.scenariosRegexp) {
		suite.RunSuite(t, testSuite)
	}
}
