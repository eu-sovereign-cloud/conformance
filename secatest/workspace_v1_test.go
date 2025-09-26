package secatest

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestWorkspaceV1Suite(t *testing.T) {
	suite.RunNamedSuite(t, "Workspace V1", &WorkspaceV1TestSuite{
		regionalTestSuite: regionalTestSuite{
			testSuite: testSuite{
				tenant:        config.clientTenant,
				authToken:     config.clientAuthToken,
				mockEnabled:   config.mockEnabled,
				mockServerURL: &config.mockServerURL,
			},
			region: config.clientRegion,
			client: clients.regionalClient,
		},
	})
}
