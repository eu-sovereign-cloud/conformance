package secatest

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestFoundationUsageV1Suite(t *testing.T) {
	testSuite := &FoundationUsageV1TestSuite{
		mixedTestSuite: mixedTestSuite{
			testSuite: testSuite{
				tenant:        config.clientTenant,
				authToken:     config.clientAuthToken,
				mockEnabled:   config.mockEnabled,
				mockServerURL: &config.mockServerURL,
				scenarioName:  foundationV1UsageSuiteName,
				initialDelay:  config.initialDelay,
				baseInterval:  config.baseInterval,
				maxAttempts:   config.maxAttempts,
			},
			globalClient:   clients.globalClient,
			regionalClient: clients.regionalClient,
			region:         config.clientRegion,
		},
		users:          config.scenariosUsers,
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
