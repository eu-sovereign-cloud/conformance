package secatest

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestFoundationUsageV1Suite(t *testing.T) {
	suite.RunNamedSuite(t, "Foundation Usage V1", &FoundationUsageV1TestSuite{
		mixedTestSuite: mixedTestSuite{
			testSuite: testSuite{
				tenant:        config.clientTenant,
				authToken:     config.clientAuthToken,
				mockEnabled:   config.mockEnabled,
				mockServerURL: &config.mockServerURL,
			},
			globalClient:   clients.globalClient,
			regionalClient: clients.regionalClient,
			region:         config.clientRegion,
		},
		users:          config.scenarioUsers,
		networkCidr:    config.scenarioCidr,
		publicIpsRange: config.scenarioPublicIps,
		regionZones:    clients.regionZones,
		instanceSkus:   clients.instanceSkus,
		storageSkus:    clients.storageSkus,
		networkSkus:    clients.networkSkus,
	})
}
