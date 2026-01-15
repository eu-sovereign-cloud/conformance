package main

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites/compute"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestComputeV1Suites(t *testing.T) {
	regionalTestSuite := createRegionalTestSuite(config.Parameters, config.Clients)

	// LifeCycle Suite
	lifeCycleTestSuite := &compute.ComputeLifeCycleV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		AvailableZones:    config.Clients.RegionZones,
		InstanceSkus:      config.Clients.InstanceSkus,
		StorageSkus:       config.Clients.StorageSkus,
	}
	lifeCycleTestSuite.ScenarioName = constants.ComputeV1LifeCycleSuiteName
	if lifeCycleTestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, lifeCycleTestSuite)
	}

	// List Suite
	listTestSuite := &compute.ComputeListV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		AvailableZones:    config.Clients.RegionZones,
		InstanceSkus:      config.Clients.InstanceSkus,
		StorageSkus:       config.Clients.StorageSkus,
	}
	listTestSuite.ScenarioName = constants.ComputeV1ListSuiteName
	if listTestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, listTestSuite)
	}
}
