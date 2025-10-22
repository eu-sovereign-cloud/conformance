package secatest

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestNetworkV1Suite(t *testing.T) {
	testSuite := &NetworkV1TestSuite{
		regionalTestSuite: regionalTestSuite{
			testSuite: testSuite{
				tenant:        config.clientTenant,
				authToken:     config.clientAuthToken,
				mockEnabled:   config.mockEnabled,
				mockServerURL: &config.mockServerURL,
				scenarioName:  networkV1LifeCycleSuiteName,
			},
			region: config.clientRegion,
			client: clients.regionalClient,
		},
		networkCidr:    config.scenariosCidr,
		publicIpsRange: config.scenariosPublicIps,
		regionZones:    clients.regionZones,
		instanceSkus:   clients.instanceSkus,
		storageSkus:    clients.storageSkus,
		networkSkus:    clients.networkSkus,
	}

	if testSuite.canRun(config.scenariosRegexp) {
		suite.RunSuite(t, testSuite)
	}
}
