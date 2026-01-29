package main

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites/network"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestNetworkV1Suites(t *testing.T) {
	regionalTestSuite := suites.NewRegionalTestSuite(config.Parameters, config.Clients)

	// LifeCycle Suite
	lifeCycleTestSuite := network.NewLifeCycleV1TestSuite(regionalTestSuite,
		&network.NetworkLifeCycleV1Config{
			NetworkCidr:    config.Parameters.ScenariosCidr,
			PublicIpsRange: config.Parameters.ScenariosPublicIps,
			RegionZones:    config.Clients.RegionZones,
			InstanceSkus:   config.Clients.InstanceSkus,
			StorageSkus:    config.Clients.StorageSkus,
			NetworkSkus:    config.Clients.NetworkSkus,
		},
	)
	if lifeCycleTestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, lifeCycleTestSuite)
	}

	// List Suite
	listTestSuite := network.NewListV1TestSuite(regionalTestSuite,
		&network.NetworkListV1Config{
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
