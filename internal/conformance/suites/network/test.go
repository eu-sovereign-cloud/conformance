package network

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestNetworkLifecycleV1Suite(t *testing.T) {
	regionalTest := suites.RegionalTestSuite{
		TestSuite: suites.TestSuite{
			Tenant:        config.Parameters.ClientTenant,
			AuthToken:     config.Parameters.ClientAuthToken,
			MockEnabled:   config.Parameters.MockEnabled,
			MockServerURL: &config.Parameters.MockServerURL,
			ScenarioName:  constants.NetworkV1LifeCycleSuiteName,
			BaseDelay:     config.Parameters.BaseDelay,
			BaseInterval:  config.Parameters.BaseInterval,
			MaxAttempts:   config.Parameters.MaxAttempts,
		},
		Region: config.Parameters.ClientRegion,
		Client: config.Clients.RegionalClient,
	}

	// LifeCycle Suite
	testLifeCycleSuite := &NetworkLifeCycleV1TestSuite{
		RegionalTestSuite: regionalTest,
		NetworkCidr:       config.Parameters.ScenariosCidr,
		publicIpsRange:    config.Parameters.ScenariosPublicIps,
		RegionZones:       config.Clients.RegionZones,
		InstanceSkus:      config.Clients.InstanceSkus,
		StorageSkus:       config.Clients.StorageSkus,
		NetworkSkus:       config.Clients.NetworkSkus,
	}

	if testLifeCycleSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, testLifeCycleSuite)
	}

	// List Suite
	testListSuite := &NetworkListV1TestSuite{
		RegionalTestSuite: regionalTest,
		NetworkCidr:       config.Parameters.ScenariosCidr,
		PublicIpsRange:    config.Parameters.ScenariosPublicIps,
		RegionZones:       config.Clients.RegionZones,
		InstanceSkus:      config.Clients.InstanceSkus,
		StorageSkus:       config.Clients.StorageSkus,
		NetworkSkus:       config.Clients.NetworkSkus,
	}
	testListSuite.RegionalTestSuite.TestSuite.ScenarioName = constants.NetworkV1ListSuiteName

	if testListSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, testListSuite)
	}
}
