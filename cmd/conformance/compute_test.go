package main

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites/compute"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestComputeV1Suites(t *testing.T) {
	regionalTestSuite := suites.RegionalTestSuite{
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
	testLifeCycleSuite := &compute.LifeCycleV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		AvailableZones:    config.Clients.RegionZones,
		InstanceSkus:      config.Clients.InstanceSkus,
		StorageSkus:       config.Clients.StorageSkus,
	}
	testLifeCycleSuite.ScenarioName = constants.ComputeV1LifeCycleSuiteName
	if testLifeCycleSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, testLifeCycleSuite)
	}

	// List Suite
	testListSuite := &compute.ListV1TestSuite{
		RegionalTestSuite: regionalTestSuite,
		AvailableZones:    config.Clients.RegionZones,
		InstanceSkus:      config.Clients.InstanceSkus,
		StorageSkus:       config.Clients.StorageSkus,
	}
	testListSuite.ScenarioName = constants.ComputeV1ListSuiteName
	if testListSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, testListSuite)
	}
}
