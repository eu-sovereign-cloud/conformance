package secatest

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestGlobalSuites(t *testing.T) {
	suite.RunNamedSuite(t, "Authorization V1", &AuthorizationV1TestSuite{
		globalTestSuite: globalTestSuite{
			testSuite: testSuite{
				tenant:        config.clientTenant,
				authToken:     config.clientAuthToken,
				mockEnabled:   config.mockEnabled,
				mockServerURL: &config.mockServerURL,
			},
			client: clients.globalClient,
		},
		users: config.scenarioUsers,
	})
}
