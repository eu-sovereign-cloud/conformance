package secatest

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestAuthorizationV1Suite(t *testing.T) {
	testSuite := &AuthorizationV1TestSuite{
		globalTestSuite: globalTestSuite{
			testSuite: testSuite{
				tenant:        config.clientTenant,
				authToken:     config.clientAuthToken,
				mockEnabled:   config.mockEnabled,
				mockServerURL: &config.mockServerURL,
				scenarioName:  authorizationV1LifeCycleSuiteName,
			},
			client: clients.globalClient,
		},
		users: config.scenariosUsers,
	}

	if testSuite.canRun(config.scenariosRegexp) {
		suite.RunSuite(t, testSuite)
	}
}
