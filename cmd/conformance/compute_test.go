package main

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites/compute"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestComputeV1Suites(t *testing.T) {
	regionalTestSuite := suites.CreateRegionalTestSuite(config.Parameters, config.Clients)

	// Provider LifeCycle Suite
	providerLifeCycleTestSuite := compute.CreateProviderLifeCycleV1TestSuite(regionalTestSuite,
		&compute.ProviderLifeCycleV1Config{
			AvailableZones: config.Clients.RegionZones,
			InstanceSkus:   config.Clients.InstanceSkus,
			StorageSkus:    config.Clients.StorageSkus,
		},
	)
	if providerLifeCycleTestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, providerLifeCycleTestSuite)
	}

	// List Suite
	listTestSuite := compute.CreateProviderQueriesV1TestSuite(regionalTestSuite,
		&compute.ProviderQueriesV1Config{
			AvailableZones: config.Clients.RegionZones,
			InstanceSkus:   config.Clients.InstanceSkus,
			StorageSkus:    config.Clients.StorageSkus,
		},
	)
	if listTestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, listTestSuite)
	}
}
