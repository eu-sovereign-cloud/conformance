package secatest

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestNetworkLifecycleV1Suite(t *testing.T) {
	regionalTest := regionalTestSuite{
		testSuite: testSuite{
			tenant:        config.clientTenant,
			authToken:     config.clientAuthToken,
			mockEnabled:   config.mockEnabled,
			mockServerURL: &config.mockServerURL,
			scenarioName:  networkV1LifeCycleSuiteName,
			baseDelay:     config.baseDelay,
			baseInterval:  config.baseInterval,
			maxAttempts:   config.maxAttempts,
		},
		region: config.clientRegion,
		client: clients.regionalClient,
	}

	// LifeCycle Suite
	testLifeCycleSuite := &NetworkLifeCycleV1TestSuite{
		regionalTestSuite: regionalTest,
		networkCidr:       config.scenariosCidr,
		publicIpsRange:    config.scenariosPublicIps,
		regionZones:       clients.regionZones,
		instanceSkus:      clients.instanceSkus,
		storageSkus:       clients.storageSkus,
		networkSkus:       clients.networkSkus,
	}

	if testLifeCycleSuite.canRun(config.scenariosRegexp) {
		suite.RunSuite(t, testLifeCycleSuite)
	}

	// List Suite
	testListSuite := &NetworkListV1TestSuite{
		regionalTestSuite: regionalTest,
		networkCidr:       config.scenariosCidr,
		publicIpsRange:    config.scenariosPublicIps,
		regionZones:       clients.regionZones,
		instanceSkus:      clients.instanceSkus,
		storageSkus:       clients.storageSkus,
		networkSkus:       clients.networkSkus,
	}
	testListSuite.regionalTestSuite.testSuite.scenarioName = networkV1ListSuiteName
	if testListSuite.canRun(config.scenariosRegexp) {
		suite.RunSuite(t, testListSuite)
	}
}
