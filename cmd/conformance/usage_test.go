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
	regionalTestSuite := suites.CreateRegionalTestSuite(config.Parameters, config.Clients)

	// Foundation Providers Suite
	foundationProvidersSuite := usage.CreateFoundationProvidersV1TestSuite(mixedTestSuite,
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
	if foundationProvidersSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, foundationProvidersSuite)
	}

	// Multi-Workspace Isolation Suite
	multiWorkspaceIsolationSuite := usage.CreateMultiWorkspaceIsolationV1TestSuite(mixedTestSuite,
		&usage.MultiWorkspaceIsolationV1Config{
			Users:          config.Parameters.ScenariosUsers,
			NetworkCidr:    config.Parameters.ScenariosCidr,
			PublicIpsRange: config.Parameters.ScenariosPublicIps,
			RegionZones:    config.Clients.RegionZones,
			InstanceSkus:   config.Clients.InstanceSkus,
			StorageSkus:    config.Clients.StorageSkus,
			NetworkSkus:    config.Clients.NetworkSkus,
		},
	)
	if multiWorkspaceIsolationSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, multiWorkspaceIsolationSuite)
	}

	// HA Multi-Zone Suite
	haMultiZoneSuite := usage.CreateHaMultiZoneV1TestSuite(regionalTestSuite,
		&usage.HaMultiZoneV1Config{
			RegionZones:  config.Clients.RegionZones,
			InstanceSkus: config.Clients.InstanceSkus,
			StorageSkus:  config.Clients.StorageSkus,
		},
	)
	if haMultiZoneSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, haMultiZoneSuite)
	}

	// Private Secure Workspace Suite
	privateSecureWorkspaceSuite := usage.CreatePrivateSecureWorkspaceV1TestSuite(regionalTestSuite,
		&usage.PrivateSecureWorkspaceV1Config{
			NetworkCidr:  config.Parameters.ScenariosCidr,
			RegionZones:  config.Clients.RegionZones,
			InstanceSkus: config.Clients.InstanceSkus,
			StorageSkus:  config.Clients.StorageSkus,
			NetworkSkus:  config.Clients.NetworkSkus,
		},
	)
	if privateSecureWorkspaceSuite.CanRun(config.Parameters.ScenariosRegexp) {
		suite.RunSuite(t, privateSecureWorkspaceSuite)
	}
}
