package secatest

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestComputeV1Suite(t *testing.T) {
	suite.RunNamedSuite(t, "Compute V1", &ComputeV1TestSuite{
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
		availableZones: clients.regionZones,
		instanceSkus:   clients.instanceSkus,
		storageSkus:    clients.storageSkus,
	})
}
