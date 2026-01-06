package compute

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestComputeV1Suite(t *testing.T) {
	regionalTestSuite := suites.RegionalTestSuite{
		TestSuite: suites.TestSuite{
			Tenant:        config.Parameters.ClientTenant,
			AuthToken:     config.Parameters.ClientAuthToken,
			MockEnabled:   config.Parameters.MockEnabled,
			MockServerURL: &config.Parameters.MockServerURL,
			ScenarioName:  constants.ComputeV1LifeCycleSuiteName,
			BaseDelay:     config.Parameters.BaseDelay,
			BaseInterval:  config.Parameters.BaseInterval,
			MaxAttempts:   config.Parameters.MaxAttempts,
		},
		Region: config.Parameters.ClientRegion,
		Client: config.Clients.RegionalClient,
	}

	// LifeCycle Suite
	testLifeCycleSuite := &ComputeV1LifeCycleTestSuite{
		RegionalTestSuite: regionalTestSuite,
		AvailableZones:    config.Clients.RegionZones,
		InstanceSkus:      config.Clients.InstanceSkus,
		StorageSkus:       config.Clients.StorageSkus,
	}

	if testLifeCycleSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, testLifeCycleSuite)
	}

	// List Suite
	testListSuite := &ComputeV1ListTestSuite{
		RegionalTestSuite: regionalTestSuite,
		AvailableZones:    config.Clients.RegionZones,
		InstanceSkus:      config.Clients.InstanceSkus,
		StorageSkus:       config.Clients.StorageSkus,
	}
	testListSuite.RegionalTestSuite.ScenarioName = constants.ComputeV1ListSuiteName

	if testListSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, testListSuite)
	}
}
