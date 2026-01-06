package storage

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestStorageV1LifeCycleSuite(t *testing.T) {
	regionalTestSuite := suites.RegionalTestSuite{
		TestSuite: suites.TestSuite{
			Tenant:        config.Parameters.ClientTenant,
			AuthToken:     config.Parameters.ClientAuthToken,
			MockEnabled:   config.Parameters.MockEnabled,
			MockServerURL: &config.Parameters.MockServerURL,
			ScenarioName:  constants.StorageV1LifeCycleSuiteName,
			BaseDelay:     config.Parameters.BaseDelay,
			BaseInterval:  config.Parameters.BaseInterval,
			MaxAttempts:   config.Parameters.MaxAttempts,
		},
		Region: config.Parameters.ClientRegion,
		Client: config.Clients.RegionalClient,
	}

	// LifeCycle Suite
	testLifeCycleSuite := &StorageV1LifeCycleTestSuite{
		RegionalTestSuite: regionalTestSuite,
		StorageSkus:       config.Clients.StorageSkus,
	}

	if testLifeCycleSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, testLifeCycleSuite)
	}

	// List Suite
	testListSuite := &StorageV1ListTestSuite{
		RegionalTestSuite: regionalTestSuite,
		StorageSkus:       config.Clients.StorageSkus,
	}
	testListSuite.RegionalTestSuite.TestSuite.ScenarioName = constants.StorageV1ListSuiteName

	if testListSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, testListSuite)
	}
}
