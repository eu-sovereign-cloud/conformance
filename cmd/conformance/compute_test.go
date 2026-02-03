package main

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites/compute"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestComputeV1Suites(t *testing.T) {
	regionalTestSuite := suites.NewRegionalTestSuite(config.Parameters, config.Clients)

	// LifeCycle Suite
	lifeCycleTestSuite := compute.NewLifeCycleV1TestSuite(regionalTestSuite,
		&compute.LifeCycleV1Config{
			AvailableZones: config.Clients.RegionZones,
			InstanceSkus:   config.Clients.InstanceSkus,
			StorageSkus:    config.Clients.StorageSkus,
		},
	)
	if lifeCycleTestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, lifeCycleTestSuite)
	}

	// List Suite
	listTestSuite := compute.NewListV1TestSuite(regionalTestSuite,
		&compute.ListV1Config{
			AvailableZones: config.Clients.RegionZones,
			InstanceSkus:   config.Clients.InstanceSkus,
			StorageSkus:    config.Clients.StorageSkus,
		},
	)
	if listTestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, listTestSuite)
	}
}
