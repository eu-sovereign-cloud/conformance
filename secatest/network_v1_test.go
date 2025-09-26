package secatest

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestNetworkV1Suite(t *testing.T) {
	suite.RunNamedSuite(t, "Network V1", &NetworkV1TestSuite{
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
		networkCidr:    config.scenarioCidr,
		publicIpsRange: config.scenarioPublicIps,
		regionZones:    clients.regionZones,
		instanceSkus:   clients.instanceSkus,
		storageSkus:    clients.storageSkus,
		networkSkus:    clients.networkSkus,
	})
}
