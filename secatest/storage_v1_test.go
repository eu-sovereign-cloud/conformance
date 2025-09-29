package secatest

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestStorageV1Suite(t *testing.T) {
	testSuite := &StorageV1TestSuite{
		regionalTestSuite: regionalTestSuite{
			testSuite: testSuite{
				tenant:        config.clientTenant,
				authToken:     config.clientAuthToken,
				mockEnabled:   config.mockEnabled,
				mockServerURL: &config.mockServerURL,
				scenarioName:  storageV1LifeCycleSuiteName,
			},
			region: config.clientRegion,
			client: clients.regionalClient,
		},
		storageSkus: clients.storageSkus,
	}

	if testSuite.canRun(config.scenariosRegexp) {
		suite.RunSuite(t, testSuite)
	}
}
