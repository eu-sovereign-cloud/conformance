package secatest

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestComputeV1Suite(t *testing.T) {
	regionalTestSuite := regionalTestSuite{
		testSuite: testSuite{
			tenant:        config.clientTenant,
			authToken:     config.clientAuthToken,
			mockEnabled:   config.mockEnabled,
			mockServerURL: &config.mockServerURL,
			scenarioName:  computeV1LifeCycleSuiteName,
			baseDelay:     config.baseDelay,
			baseInterval:  config.baseInterval,
			maxAttempts:   config.maxAttempts,
		},
		region: config.clientRegion,
		client: clients.regionalClient,
	}

	// LifeCycle Suite
	testLifeCycleSuite := &ComputeV1LifeCycleTestSuite{
		regionalTestSuite: regionalTestSuite,
		availableZones:    clients.regionZones,
		instanceSkus:      clients.instanceSkus,
		storageSkus:       clients.storageSkus,
	}

	if testLifeCycleSuite.canRun(config.scenariosRegexp) {
		suite.RunSuite(t, testLifeCycleSuite)
	}

	// List Suite
	testListSuite := &ComputeV1ListTestSuite{
		regionalTestSuite: regionalTestSuite,
		availableZones:    clients.regionZones,
		instanceSkus:      clients.instanceSkus,
		storageSkus:       clients.storageSkus,
	}

	testListSuite.regionalTestSuite.scenarioName = computeV1ListSuiteName

	if testListSuite.canRun(config.scenariosRegexp) {
		suite.RunSuite(t, testListSuite)
	}
}
