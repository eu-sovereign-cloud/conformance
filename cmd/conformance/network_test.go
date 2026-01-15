package main

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites/network"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestNetworkV1Suites(t *testing.T) {
	regionalTestSuite := createRegionalTestSuite(config.Parameters, config.Clients)

	// LifeCycle Suite
	lifeCycleTestSuite := &network.NetworkLifeCycleV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		NetworkCidr:       config.Parameters.ScenariosCidr,
		PublicIpsRange:    config.Parameters.ScenariosPublicIps,
		RegionZones:       config.Clients.RegionZones,
		InstanceSkus:      config.Clients.InstanceSkus,
		StorageSkus:       config.Clients.StorageSkus,
		NetworkSkus:       config.Clients.NetworkSkus,
	}
	lifeCycleTestSuite.ScenarioName = constants.NetworkV1LifeCycleSuiteName
	if lifeCycleTestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, lifeCycleTestSuite)
	}

	// List Suite
	listTestSuite := &network.NetworkListV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		NetworkCidr:       config.Parameters.ScenariosCidr,
		PublicIpsRange:    config.Parameters.ScenariosPublicIps,
		RegionZones:       config.Clients.RegionZones,
		InstanceSkus:      config.Clients.InstanceSkus,
		StorageSkus:       config.Clients.StorageSkus,
		NetworkSkus:       config.Clients.NetworkSkus,
	}
	listTestSuite.ScenarioName = constants.NetworkV1ListSuiteName
	if listTestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, listTestSuite)
	}
}
