package main

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites/network"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestNetworkV1Suites(t *testing.T) {
	globalTestSuite := suites.RegionalTestSuite{
		TestSuite: suites.TestSuite{
			Tenant:        config.Parameters.ClientTenant,
			AuthToken:     config.Parameters.ClientAuthToken,
			MockEnabled:   config.Parameters.MockEnabled,
			MockServerURL: &config.Parameters.MockServerURL,
			BaseDelay:     config.Parameters.BaseDelay,
			BaseInterval:  config.Parameters.BaseInterval,
			MaxAttempts:   config.Parameters.MaxAttempts,
		},
		Region: config.Parameters.ClientRegion,
		Client: config.Clients.RegionalClient,
	}

	// LifeCycle Suite
	testLifeCycleSuite := &network.LifeCycleV1TestSuite{
		RegionalTestSuite: globalTestSuite,
		NetworkCidr:       config.Parameters.ScenariosCidr,
		PublicIpsRange:    config.Parameters.ScenariosPublicIps,
		RegionZones:       config.Clients.RegionZones,
		InstanceSkus:      config.Clients.InstanceSkus,
		StorageSkus:       config.Clients.StorageSkus,
		NetworkSkus:       config.Clients.NetworkSkus,
	}
	testLifeCycleSuite.ScenarioName = constants.NetworkV1LifeCycleSuiteName
	if testLifeCycleSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, testLifeCycleSuite)
	}

	// List Suite
	testListSuite := &network.ListV1TestSuite{
		RegionalTestSuite: globalTestSuite,
		NetworkCidr:       config.Parameters.ScenariosCidr,
		PublicIpsRange:    config.Parameters.ScenariosPublicIps,
		RegionZones:       config.Clients.RegionZones,
		InstanceSkus:      config.Clients.InstanceSkus,
		StorageSkus:       config.Clients.StorageSkus,
		NetworkSkus:       config.Clients.NetworkSkus,
	}
	testListSuite.ScenarioName = constants.NetworkV1ListSuiteName
	if testListSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, testListSuite)
	}
}
