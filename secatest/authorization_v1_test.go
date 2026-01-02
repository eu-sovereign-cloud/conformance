package secatest

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestAuthorizationV1Suite(t *testing.T) {
	globalTestSuite := globalTestSuite{
		testSuite: testSuite{
			tenant:        config.clientTenant,
			authToken:     config.clientAuthToken,
			mockEnabled:   config.mockEnabled,
			mockServerURL: &config.mockServerURL,
			scenarioName:  authorizationV1LifeCycleSuiteName,
			baseDelay:     config.baseDelay,
			baseInterval:  config.baseInterval,
			maxAttempts:   config.maxAttempts,
		},
		client: clients.globalClient,
	}

	// LifeCycle Suite
	testLifeCycleSuite := &AuthorizationV1LifeCycleTestSuite{
		globalTestSuite: globalTestSuite,
		users:           config.scenariosUsers,
	}

	if testLifeCycleSuite.canRun(config.scenariosRegexp) {
		suite.RunSuite(t, testLifeCycleSuite)
	}

	// List Suite
	globalTestSuite.testSuite.scenarioName = authorizationV1ListSuiteName
	testListSuite := &AuthorizationV1ListTestSuite{
		globalTestSuite: globalTestSuite,
		users:           config.scenariosUsers,
	}

	if testListSuite.canRun(config.scenariosRegexp) {
		suite.RunSuite(t, testListSuite)
	}
}
