//nolint:dupl
package main

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites/network"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
)

func TestNetworkV1Suites(t *testing.T) {
	regionalTestSuite := suites.CreateRegionalTestSuite(config.Parameters, config.Clients)

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
	lifeCycleTestSuite.RunSuite(t, config.Parameters.ScenariosRegexp,
		func() { lifeCycleTestSuite.ScenarioName = constants.NetworkV1LifeCycleSuiteName },
	)

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
	listTestSuite.RunSuite(t, config.Parameters.ScenariosRegexp,
		func() { listTestSuite.ScenarioName = constants.NetworkV1ListSuiteName },
	)
}
