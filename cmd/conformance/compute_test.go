//nolint:dupl
package main

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites/compute"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
)

func TestComputeV1Suites(t *testing.T) {
	regionalTestSuite := suites.CreateRegionalTestSuite(config.Parameters, config.Clients)

	// LifeCycle Suite
	lifeCycleTestSuite := &compute.ComputeLifeCycleV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		AvailableZones:    config.Clients.RegionZones,
		InstanceSkus:      config.Clients.InstanceSkus,
		StorageSkus:       config.Clients.StorageSkus,
	}
	lifeCycleTestSuite.RunSuite(t, config.Parameters.ScenariosRegexp,
		func() { lifeCycleTestSuite.ScenarioName = constants.ComputeV1LifeCycleSuiteName },
	)

	// List Suite
	listTestSuite := &compute.ComputeListV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		AvailableZones:    config.Clients.RegionZones,
		InstanceSkus:      config.Clients.InstanceSkus,
		StorageSkus:       config.Clients.StorageSkus,
	}
	listTestSuite.RunSuite(t, config.Parameters.ScenariosRegexp,
		func() { listTestSuite.ScenarioName = constants.ComputeV1ListSuiteName },
	)
}
