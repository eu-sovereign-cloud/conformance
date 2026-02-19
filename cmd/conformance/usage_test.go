package main

import (
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/config"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites"
	"github.com/eu-sovereign-cloud/conformance/internal/conformance/suites/usage"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestUsageV1Suites(t *testing.T) {
	mixedTestSuite := suites.CreateMixedTestSuite(config.Parameters, config.Clients)

	// Foundation Suite
	foundationTestSuite := usage.CreateFoundationProvidersV1TestSuite(mixedTestSuite,
		&usage.FoundationProvidersV1Config{
			Users:          config.Parameters.ScenariosUsers,
			NetworkCidr:    config.Parameters.ScenariosCidr,
			PublicIpsRange: config.Parameters.ScenariosPublicIps,
			RegionZones:    config.Clients.RegionZones,
			InstanceSkus:   config.Clients.InstanceSkus,
			StorageSkus:    config.Clients.StorageSkus,
			NetworkSkus:    config.Clients.NetworkSkus,
		},
	)
	if foundationTestSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, foundationTestSuite)
	}
}
