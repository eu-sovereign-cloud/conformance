package secatest

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestComputeV1Suite(t *testing.T) {
	testSuite := &ComputeV1TestSuite{
		regionalTestSuite: regionalTestSuite{
			testSuite: testSuite{
				tenant:        config.clientTenant,
				authToken:     config.clientAuthToken,
				mockEnabled:   config.mockEnabled,
				mockServerURL: &config.mockServerURL,
				scenarioName:  computeV1LifeCycleSuiteName,
				initialDelay:  config.initialDelay,
				baseInterval:  config.baseInterval,
				maxAttempts:   config.maxAttempts,
			},
			region: config.clientRegion,
			client: clients.regionalClient,
		},
		availableZones: clients.regionZones,
		instanceSkus:   clients.instanceSkus,
		storageSkus:    clients.storageSkus,
	}

	if testSuite.canRun(config.scenariosRegexp) {
		suite.RunSuite(t, testSuite)
	}
}
