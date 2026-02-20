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
	providerLifeCycleSuite := network.CreateProviderLifeCycleV1TestSuite(regionalTestSuite,
		&network.ProviderLifeCycleV1Config{
			NetworkCidr:    config.Parameters.ScenariosCidr,
			PublicIpsRange: config.Parameters.ScenariosPublicIps,
			RegionZones:    config.Clients.RegionZones,
			InstanceSkus:   config.Clients.InstanceSkus,
			StorageSkus:    config.Clients.StorageSkus,
			NetworkSkus:    config.Clients.NetworkSkus,
		},
	)
	if providerLifeCycleSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, providerLifeCycleSuite)
	}

	// Provider Queries Suite
	providerQueriesSuite := network.CreateProviderQueriesV1TestSuite(regionalTestSuite,
		&network.ProviderQueriesV1Config{
			NetworkCidr:    config.Parameters.ScenariosCidr,
			PublicIpsRange: config.Parameters.ScenariosPublicIps,
			RegionZones:    config.Clients.RegionZones,
			InstanceSkus:   config.Clients.InstanceSkus,
			StorageSkus:    config.Clients.StorageSkus,
			NetworkSkus:    config.Clients.NetworkSkus,
		},
	)
	if providerQueriesSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, providerQueriesSuite)
	}
}
