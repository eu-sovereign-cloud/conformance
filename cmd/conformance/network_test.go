package main

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites/network"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestNetworkV1Suites(t *testing.T) {
	regionalTestSuite := suites.CreateRegionalTestSuite(config.Parameters, config.Clients)

	// Provider LifeCycle Suite
	providerLifeCycleTestSuite := network.CreateProviderLifeCycleV1TestSuite(regionalTestSuite,
		&network.ProviderLifeCycleV1Config{
			NetworkCidr:    config.Parameters.ScenariosCidr,
			PublicIpsRange: config.Parameters.ScenariosPublicIps,
			RegionZones:    config.Clients.RegionZones,
			InstanceSkus:   config.Clients.InstanceSkus,
			StorageSkus:    config.Clients.StorageSkus,
			NetworkSkus:    config.Clients.NetworkSkus,
		},
	)
	if providerLifeCycleTestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, providerLifeCycleTestSuite)
	}

	// List Suite
	listTestSuite := network.CreateProviderQueriesV1TestSuite(regionalTestSuite,
		&network.ProviderQueriesV1Config{
			NetworkCidr:    config.Parameters.ScenariosCidr,
			PublicIpsRange: config.Parameters.ScenariosPublicIps,
			RegionZones:    config.Clients.RegionZones,
			InstanceSkus:   config.Clients.InstanceSkus,
			StorageSkus:    config.Clients.StorageSkus,
			NetworkSkus:    config.Clients.NetworkSkus,
		},
	)
	if listTestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, listTestSuite)
	}
}
