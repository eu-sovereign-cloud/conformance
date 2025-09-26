package secatest

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestStorageV1Suite(t *testing.T) {
	suite.RunNamedSuite(t, "Storage V1", &StorageV1TestSuite{
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
		storageSkus: clients.storageSkus,
	})
}
