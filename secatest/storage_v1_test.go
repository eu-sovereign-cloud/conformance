package secatest

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestStorageV1LifeCycleSuite(t *testing.T) {
	regionalTestSuite := regionalTestSuite{
		testSuite: testSuite{
			tenant:        config.clientTenant,
			authToken:     config.clientAuthToken,
			mockEnabled:   config.mockEnabled,
			mockServerURL: &config.mockServerURL,
			scenarioName:  storageV1LifeCycleSuiteName,
			baseDelay:     config.baseDelay,
			baseInterval:  config.baseInterval,
			maxAttempts:   config.maxAttempts,
		},
		region: config.clientRegion,
		client: clients.regionalClient,
	}

	// LifeCycle Suite
	testLifeCycleSuite := &StorageV1LifeCycleTestSuite{
		regionalTestSuite: regionalTestSuite,
		storageSkus:       clients.storageSkus,
	}

	if testLifeCycleSuite.canRun(config.scenariosRegexp) {
		suite.RunSuite(t, testLifeCycleSuite)
	}

	// List Suite
	testListSuite := &StorageV1ListTestSuite{
		regionalTestSuite: regionalTestSuite,
		storageSkus:       clients.storageSkus,
	}
	testListSuite.regionalTestSuite.testSuite.scenarioName = storageV1ListSuiteName
	if testListSuite.canRun(config.scenariosRegexp) {
		suite.RunSuite(t, testListSuite)
	}
}
